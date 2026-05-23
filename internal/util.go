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

func parseQueryParams(params openapi3.Parameters) []paramSpec {
	var queries []paramSpec
	for _, q := range params {
		queries = append(queries, paramSpec{
			Name:    q.Value.Name,
			Type:    typeString(q.Value.Schema.Value.Type),
			Example: q.Value.Example,
		})
	}
	return queries
}

func parseBodyParams(requestBody *openapi3.RequestBodyRef) []paramSpec {
	var bodies []paramSpec
	if requestBody == nil {
		return bodies
	}
	content, ok := requestBody.Value.Content["application/json"]
	if !ok || content.Schema == nil {
		return bodies
	}
	for name, b := range content.Schema.Value.Properties {
		bodies = append(bodies, paramSpec{
			Name:    name,
			Type:    typeString(b.Value.Type),
			Example: b.Value.Example,
		})
	}
	return bodies
}

func parseResponses(responses *openapi3.Responses, cases []string) []responseSpec {
	var result []responseSpec
	if responses == nil {
		return result
	}
	for name, r := range responses.Map() {
		if len(cases) > 0 && !pkg.InArray(name, cases) {
			continue
		}
		if r.Value.Description == nil {
			continue
		}
		content, ok := r.Value.Content["application/json"]
		if !ok {
			continue
		}
		result = append(result, responseSpec{
			Name:        name,
			Description: *r.Value.Description,
			Example:     content.Example,
		})
	}
	return result
}

func GenItem(inputFileName string, cases []string) (*APISpec, error) {
	var apiSpec APISpec

	doc, err := openapi3.NewLoader().LoadFromFile(inputFileName)
	if err != nil {
		return &APISpec{}, err
	}

	apiSpec.Title = doc.Info.Title
	apiSpec.Description = doc.Info.Description
	apiSpec.Version = doc.Info.Version
	if doc.Servers != nil {
		apiSpec.BaseUrl = doc.Servers[0].URL
	} else {
		apiSpec.BaseUrl = "dummy URL"
	}

	var pathSpecs []pathSpec
	for _, path := range doc.Paths.InMatchingOrder() {
		obj := doc.Paths.Find(path).Operations()
		var baseSpecs []baseSpec
		for method, op := range obj {
			baseSpecs = append(baseSpecs, baseSpec{
				Method:   method,
				Summary:  op.Summary,
				Body:     parseBodyParams(op.RequestBody),
				Params:   parseQueryParams(op.Parameters),
				Response: parseResponses(op.Responses, cases),
			})
		}
		pathSpecs = append(pathSpecs, pathSpec{
			Path:    path,
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

	var currentStep step
	var req requestInfo
	var exp expectInfo
	var reqURL string

	scenario := Scenario{
		Title: apiSpec.Title,
		Vars:  map[string]string{"endpoint": apiSpec.BaseUrl},
	}

	for _, spec := range apiSpec.PathSpec {
		if param != nil {
			for path, p := range *param {
				if ok := pkg.CompPath(spec.Path, path); ok {
					req.Query = p.Query
					req.Url = "{{ vars.endpoint }}" + path
					reqURL = apiSpec.BaseUrl + path
					break
				}
				reqURL = apiSpec.BaseUrl + spec.Path
				req.Url = "{{ vars.endpoint }}" + spec.Path
			}
		} else {
			reqURL = apiSpec.BaseUrl + spec.Path
			req.Url = "{{ vars.endpoint }}" + spec.Path
		}

		for _, method := range spec.Methods {
			currentStep.Title = method.Summary
			currentStep.Protocol = "http"
			req.Method = method.Method
			req.Header = map[string]string{"x-api-key": "{{ env.API_KEY }}"}

			for _, r := range method.Response {
				httpMethod := strings.ToUpper(req.Method)
				res, err := GetResponse(reqURL, req.Query, httpMethod)
				if err != nil {
					return err
				}

				i, _ := strconv.Atoi(r.Name)
				exp.StatusCode = i
				exp.Body = res
				currentStep.Request = req
				currentStep.Expect = exp
				scenario.Step = append(scenario.Step, currentStep)
			}
		}
	}

	f, err := os.Create(outputFileName) // #nosec G304 -- path comes from CLI argument supplied by the user
	if err != nil {
		return errors.New("scenario file cannot create")
	}

	b, err := yaml.Marshal(scenario)
	if err != nil {
		_ = f.Close()
		return errors.New("scenario file cannot convert")
	}

	if _, err = f.Write(b); err != nil {
		_ = f.Close()
		return errors.New("scenario file cannot write")
	}

	if err = f.Close(); err != nil {
		return fmt.Errorf("scenario file cannot close: %w", err)
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

	reqURL := url
	if qs != "" {
		reqURL = fmt.Sprintf("%s?%s", url, qs)
	}

	req, err := http.NewRequest(method, reqURL, nil)
	if err != nil {
		return nil, fmt.Errorf("request create error: %w", err)
	}
	req.Header.Set("x-api-key", os.Getenv("API_KEY"))

	client := new(http.Client)
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("access error: %v", err)
	}
	defer func() { _ = resp.Body.Close() }()

	response, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("response read error: %v", err)
	}

	var convertResponse map[string]interface{}
	if err = json.Unmarshal(response, &convertResponse); err != nil {
		return nil, fmt.Errorf("convert error: %v", err)
	}
	return convertResponse, nil
}

func AddParam(inputFile string) (*map[string]addParam, error) {
	param := make(map[string]addParam)

	f, err := os.Open(inputFile) // #nosec G304 -- path comes from CLI argument supplied by the user
	if err != nil {
		return nil, err
	}
	defer func() { _ = f.Close() }()

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
