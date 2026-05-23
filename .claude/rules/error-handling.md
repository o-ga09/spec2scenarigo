# エラーハンドリングルール

## 基本原則

標準ライブラリの `errors` と `fmt` を使う。外部エラーパッケージは導入しない。

```go
// 新しいエラーを生成する
return errors.New("scenario file cannot create")

// 下位エラーを文脈付きでラップする
return fmt.Errorf("openapi load failed: %w", err)
```

## エラーの返し方

- 下位関数のエラーは `fmt.Errorf("操作の説明: %w", err)` でラップして文脈を付ける
- 新規エラーは `errors.New("メッセージ")` で生成する
- エラーメッセージは小文字で書く（Go の慣習）

```go
// 誤: エラーメッセージが大文字始まり
return errors.New("Scenario file cannot create")

// 正: 小文字始まり
return errors.New("scenario file cannot create")

// 誤: ラップせずに新規エラーを返す（元のエラー情報が失われる）
if err != nil {
    return errors.New("something went wrong")
}

// 正: %w でラップして元のエラーを保持する
if err != nil {
    return fmt.Errorf("openapi load failed: %w", err)
}
```

## nil チェック

- 関数冒頭でガード節を書き、早期 return する（ネストを深めない）

```go
// 正: 早期 return でネストを抑える
if err != nil {
    return fmt.Errorf("csv read failed: %w", err)
}
// 正常系処理を続ける
```

## CLI でのエラー出力

- `internal/root.go` のコマンドハンドラではエラーを `fmt.Println` で表示して return する
- `os.Exit(1)` は回復不能なエラー（CSV 読み込み失敗など）でのみ使う

```go
result, err := GenItem(input, cases)
if err != nil {
    fmt.Println("error:", err)
    return
}
```

## errors.Is / errors.As

- sentinel エラーとの比較は `errors.Is(err, target)` を使う
- 型アサーションが必要な場合は `errors.As(err, &target)` を使う
- `%w` でラップした場合、`errors.Is` はラップチェーン全体を検索するため多重ラップに注意する
