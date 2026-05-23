package root

import (
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/o-ga09/spec2scenarigo/pkg"
	yaml "gopkg.in/yaml.v2"
)

func typeString(t *openapi3.Types) string {
	if t == nil {
		return ""
	}
	return strings.Join(*t, ",")
}

func GenItem(inputFileName string, cases []string) (*APISpec, error) {
	// パス毎の構造体を格納するスライスを定義
	var apiSpec APISpec
	var pathSpecs []pathSpec

	// OpenAPIのYAMLファイルを読み込み
	doc, err := openapi3.NewLoader().LoadFromFile(inputFileName)
	if err != nil {
		fmt.Println("Error:", err)
		return &APISpec{}, err
	}

	// API Specのメタデータを格納
	apiSpec.Title = doc.Info.Title
	apiSpec.Description = doc.Info.Description
	apiSpec.Version = doc.Info.Version
	// デフォルトでは、一番最初のURLを取得する。オプションで指定したURLが優先される。
	if doc.Servers != nil {
		apiSpec.BaseUrl = doc.Servers[0].URL
	} else {
		apiSpec.BaseUrl = "dummy URL"
	}

	// パス毎に処理
	for _, path := range doc.Paths.InMatchingOrder() {

		// それぞれのパスに対するメソッドの一覧を取得
		obj := doc.Paths.Find(path).Operations()

		// メソッド毎の構造体を格納するスライスを定義
		var baseSpecs []baseSpec

		// メソッド毎に処理
		for method, op := range obj {

			// クエリとボディに当たるパラメータ構造体を格納するスライスを定義
			var queries []paramSpec
			var bodies []paramSpec
			var responses []responseSpec

			// 元データにクエリパラメータがある場合
			if op.Parameters != nil {
				for _, q := range op.Parameters {

					// クエリ毎にクエリパラメータ構造体を生成
					queries = append(queries, paramSpec{
						// フィールドの名前
						Name: q.Value.Name,
						// フィールドの型
						Type: typeString(q.Value.Schema.Value.Type),
						// フィールドのサンプル値
						Example: q.Value.Example,
					})
				}
			}

			// 元データにボディパラメータがある場合
			if op.RequestBody != nil {
				for name, b := range op.RequestBody.Value.Content["application/json"].Schema.Value.Properties {

					// ボディ毎にボディパラメータ構造体を生成
					bodies = append(bodies, paramSpec{
						// フィールドの名前
						Name: name,
						// フィールドの型
						Type: typeString(b.Value.Type),
						// フィールドのサンプル値
						Example: b.Value.Example,
					})
				}
			}

			if op.Responses != nil && len(cases) == 0 {
				for name, r := range op.Responses.Map() {
					responses = append(responses, responseSpec{
						Name:        name,
						Description: *r.Value.Description,
						Example:     r.Value.Content["application/json"].Example,
					})
				}
			} else if op.Responses != nil && len(cases) > 0 {
				for name, r := range op.Responses.Map() {
					if !pkg.InArray(name, cases) {
						continue
					}
					responses = append(responses, responseSpec{
						Name:        name,
						Description: *r.Value.Description,
						Example:     r.Value.Content["application/json"].Example,
					})
				}
			}

			// メソッド毎にメソッド構造体を生成して末尾に追加
			baseSpecs = append(baseSpecs, baseSpec{
				// メソッド名
				Method: method,
				// APIのサマリ
				Summary: op.Summary,
				// ボディパラメータ構造体のスライス
				Body: bodies,
				// クエリパラメータ構造体のスライス
				Params: queries,
				// レスポンス構造体のスライス
				Response: responses,
			})

		}

		// パス毎にパス構造体を生成して末尾に追加
		pathSpecs = append(pathSpecs, pathSpec{
			// パス
			Path: path,
			// メソッド構造体のスライス
			Methods: baseSpecs,
		})

	}

	apiSpec.PathSpec = pathSpecs
	return &apiSpec, nil
}

func GenScenario(apiSpec *APISpec, outputFileName string, opts ...interface{}) error {
	var param *map[string]addParam
	if len(opts) > 0 {
		param = opts[0].(*map[string]addParam)
	}

	// シナリオの構造体
	var scenario Scenario

	// 各ステップの構造体
	var step step

	// メタデータの構造体
	var requestInfo requestInfo
	var expectInfo expectInfo

	// 想定結果レスポンス取得用URL変数
	var reqURL string

	// シナリオを作成
	scenario.Title = apiSpec.Title

	// 共通で使用する変数を定義
	scenario.Vars = map[string]string{"endpoint": apiSpec.BaseUrl}

	// ステップ毎にシナリオを作成
	for _, spec := range apiSpec.PathSpec {
		// SpecのURLとテスト対象のURLを比較して、パスパラメーターが含まれる場合、テスト対象のURLを置き換える
		if param != nil {
			for path, p := range *param {
				if ok := pkg.CompPath(spec.Path, path); ok {
					requestInfo.Query = p.Query
					requestInfo.Url = "{{ vars.endpoint }}" + path
					reqURL = apiSpec.BaseUrl + path
					break
				}
				reqURL = apiSpec.BaseUrl + spec.Path
				requestInfo.Url = "{{ vars.endpoint }}" + spec.Path
			}
		} else {
			reqURL = apiSpec.BaseUrl + spec.Path
			requestInfo.Url = "{{ vars.endpoint }}" + spec.Path
		}

		for _, method := range spec.Methods {
			step.Title = method.Summary
			step.Protocol = "http"
			requestInfo.Method = method.Method
			requestInfo.Header = map[string]string{"x-api-key": "{{ env.API_KEY }}"}

			for _, r := range method.Response {
				// APIにリクエストしてテストデータを取得する
				method := strings.ToUpper(requestInfo.Method)
				res, err := GetResponse(reqURL, requestInfo.Query, method)
				if err != nil {
					return err
				}

				// シナリオのステップ作成する
				i, _ := strconv.Atoi(r.Name)
				expectInfo.StatusCode = i
				expectInfo.Body = res
				step.Request = requestInfo
				step.Expect = expectInfo
				scenario.Step = append(scenario.Step, step)
			}
		}
	}

	// シナリオファイルを作成
	f, err := os.Create(outputFileName)
	if err != nil {
		return errors.New("Scenario file cannot create")
	}
	defer f.Close()

	// シナリオをyaml形式の変換
	b, err := yaml.Marshal(scenario)
	if err != nil {
		return errors.New("Scenario file cannot convert")
	}

	// シナリオをファイルに書き込み
	_, err = f.Write(b)
	if err != nil {
		return errors.New("Scenario file cannot write")
	}

	return nil
}

func GetResponse(url string, query any, method string) (any, error) {
	qs := ""

	if query != nil {
		q := query.(map[string]interface{})
		for k, v := range q {
			qs += fmt.Sprintf("%s=%s&", k, v)
		}
	}

	reqUrl := ""
	if qs == "" {
		reqUrl = url
	} else {
		reqUrl = fmt.Sprintf("%s?%s", url, qs)
	}
	req, _ := http.NewRequest(method, reqUrl, nil)
	req.Header.Set("x-api-key", os.Getenv("API_KEY"))
	client := new(http.Client)
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("access error: %v", err)
	}
	response, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("response read error: %v", err)
	}

	var ConvertResponse map[string]interface{}
	err = json.Unmarshal(response, &ConvertResponse)
	if err != nil {
		return nil, fmt.Errorf("convert error: %v", err)
	}
	return ConvertResponse, nil
}

func AddParam(intpufile string) (*map[string]addParam, error) {
	param := make(map[string]addParam)

	f, err := os.Open(intpufile)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	csvFile, err := csv.NewReader(f).ReadAll()
	if err != nil {
		return nil, err
	}
	for _, record := range csvFile {
		query := make(map[string]interface{})

		paths := strings.Split(record[0], "?")
		path := paths[0]
		if len(paths) > 1 {
			queries := strings.Split(paths[1], "&")
			for _, q := range queries {
				str := strings.Split(q, "=")
				p1 := str[0]
				p2 := str[1]
				query[p1] = p2
			}
		}

		param[path] = addParam{
			Method: record[1],
			Query:  query,
			Body:   record[2],
		}
	}

	return &param, nil
}
