# アーキテクチャルール

## 概要

`spec2scenarigo` は OpenAPI Spec から scenarigo 用 E2E テストシナリオ YAML を生成する CLI ツール。

## ディレクトリ構成と役割

```
main.go          → エントリポイント。internal.Execute() を呼ぶだけ
internal/
  root.go        → Cobra CLI 定義とコマンドハンドラ
  util.go        → コアロジック（GenItem, GenScenario, GetResponse, AddParam）
  types.go       → 全構造体定義（APISpec, Scenario, step, requestInfo など）
  util_test.go   → テーブル駆動テスト
pkg/
  util.go        → 共有ヘルパー（CompPath, InArray）
example/
  input.yml      → サンプル OpenAPI Spec
  param.csv      → パス/クエリパラメーター上書き用サンプル CSV
```

## データフロー

```
GenItem(inputFile, cases) → *APISpec
  └─ OpenAPI YAML を kin-openapi でパース
  └─ パス/メソッド/パラメーター/レスポンスを APISpec に格納

GenScenario(apiSpec, outputFile, [param]) → scenario.yml
  └─ APISpec を走査
  └─ GetResponse() で実 API にリクエストしてテストデータを取得
  └─ Scenario 構造体を YAML にマーシャルして書き出し

AddParam(csvFile) → *map[string]addParam
  └─ ヘッダーなし CSV を読み込み（列: path?query=value, METHOD, body）
  └─ CompPath でパスパラメーターをワイルドカードマッチ
```

## 依存関係

- `internal/` は `pkg/` を import できる
- `pkg/` は `internal/` を import しない
- 循環 import 禁止

## 設定のルール

- 認証は `x-api-key` ヘッダーのみ。環境変数 `API_KEY` から取得する
- ベース URL は OpenAPI Spec の `servers[0].URL` をデフォルトとし、`--host` フラグで上書き可能
- 出力ファイル名は `--output-file` フラグで指定。デフォルトは `scenario.yml`

## 追加実装時のルール

- 新しいコアロジックは `internal/util.go` に追加する（300 行を超えたらファイル分割を検討）
- 複数パッケージで使う汎用ユーティリティは `pkg/` に置く
- CLI フラグの追加は `internal/root.go` の `init()` 内で行う
- 新しい構造体は `internal/types.go` に追加する
