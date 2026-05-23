# Context 伝搬方針

## 現状

`spec2scenarigo` は CLI ツールであり、現時点のコードは `context.Context` を使用していない。HTTP リクエスト（`GetResponse`）も context なしで実装されている。

## context を追加するときのルール

将来的に context が必要になった場合（HTTP タイムアウト制御、テストのキャンセルなど）は以下のルールに従う。

- `context.Context` は**必ず第一引数**。変数名は `ctx`
- `context.Background()` を関数の途中で作らない
- CLI の場合は `context.Background()` を `main` または `Execute()` で一度だけ生成して伝搬する
- テストでは `t.Context()` を使う（Go 1.21+）

```go
// HTTP リクエストにタイムアウトを付ける例
func GetResponse(ctx context.Context, url string, query any, method string) (any, error) {
    req, err := http.NewRequestWithContext(ctx, method, reqUrl, nil)
    ...
}

// テスト
func TestGetResponse(t *testing.T) {
    ctx := t.Context()
    got, err := GetResponse(ctx, server.URL, nil, "GET")
    ...
}
```

## context を使わない理由を記録する

context を追加しない判断をした場合、その理由をコミットメッセージか PR 説明に残す。サイレントに省略しない。

## HTTP リクエストのタイムアウト

`GetResponse` で外部 API を呼ぶ際は、context にタイムアウトを設定するか `http.Client` の `Timeout` フィールドで制御する。

```go
// context でタイムアウト制御する場合
callCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
defer cancel()
req, _ := http.NewRequestWithContext(callCtx, method, reqUrl, nil)

// http.Client でタイムアウト制御する場合（context なしのシンプルなケース）
client := &http.Client{Timeout: 10 * time.Second}
```
