# テストルール

## ファイル配置

- テストファイルは実装ファイルと同一パッケージ・同一ディレクトリに置く（`util.go` → `util_test.go`）
- `internal/` のテストは `package root`、`pkg/` のテストは `package pkg`

## テストの書き方

- テーブル駆動テスト（`[]struct{ name, input, want }`）を基本形とする
- サブテストは `t.Run` で命名する。テスト名は日本語可
- 具体的なエッジケースとエラーパスを列挙して検証する

```go
func TestCompPath(t *testing.T) {
    cases := []struct {
        name     string
        specPath string
        testPath string
        want     bool
    }{
        {"完全一致", "/v1/users", "/v1/users", true},
        {"パスパラメーター一致", "/v1/users/{id}", "/v1/users/123", true},
        {"セグメント数不一致", "/v1/users", "/v1/users/123", false},
    }
    for _, tc := range cases {
        t.Run(tc.name, func(t *testing.T) {
            got := CompPath(tc.specPath, tc.testPath)
            if got != tc.want {
                t.Errorf("CompPath(%q, %q) = %v, want %v", tc.specPath, tc.testPath, got, tc.want)
            }
        })
    }
}
```

## 外部依存のテスト方針

- **HTTP リクエスト（GetResponse）**: テストで外部ネットワークにアクセスしない。`net/http/httptest` でモックサーバーを立てるか、テストケースから除外してカバレッジコメントを残す
- **ファイル I/O（GenItem, AddParam）**: `testdata/` ディレクトリにフィクスチャを置いて実ファイルを使うテストを書く
- **ファイル書き込み（GenScenario）**: `os.CreateTemp` で一時ファイルに書き出してテスト後に削除する

```go
func TestGenItem(t *testing.T) {
    cases := []struct {
        name      string
        inputFile string
        cases     []string
        wantTitle string
        wantErr   bool
    }{
        {"正常系", "testdata/input.yml", []string{}, "Sample API", false},
        {"ファイルなし", "testdata/not_exist.yml", []string{}, "", true},
    }
    for _, tc := range cases {
        t.Run(tc.name, func(t *testing.T) {
            got, err := GenItem(tc.inputFile, tc.cases)
            if (err != nil) != tc.wantErr {
                t.Errorf("GenItem() error = %v, wantErr %v", err, tc.wantErr)
            }
            if !tc.wantErr && got.Title != tc.wantTitle {
                t.Errorf("got.Title = %q, want %q", got.Title, tc.wantTitle)
            }
        })
    }
}
```

## 禁止事項

- `time.Sleep` をテスト内で使わない
- テストで外部ネットワークにアクセスしない（CI が壊れる）
- `os.Exit` を呼ぶコードはテストできないため、`main` 以外では使わない

## カバレッジ

- `go test ./... -coverprofile=coverage.out` でカバレッジを計測する
- コアロジック（`internal/util.go`, `pkg/util.go`）は 80% 以上を目標とする
