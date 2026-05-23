# パッケージ分割ルール

## 命名規則

- `internal/` のパッケージ名は `root`（`package root`）
- `pkg/` のパッケージ名は `pkg`（`package pkg`）
- パッケージ名は小文字・短く

## ディレクトリと役割

```
main.go              # main パッケージ。Execute() を呼ぶだけ
internal/
├── root.go          # Cobra CLI 定義・フラグ設定・コマンドハンドラ
├── util.go          # コアロジック（GenItem, GenScenario, GetResponse, AddParam）
├── types.go         # 構造体定義（APISpec, Scenario, step, requestInfo など）
└── util_test.go     # テスト
pkg/
└── util.go          # 共有ユーティリティ（CompPath, InArray）
example/
├── input.yml        # サンプル OpenAPI Spec
├── param.csv        # サンプル CSV
└── scenario.yml     # サンプル出力
```

## ファイル分割の基準

- 1 ファイルが 300 行を超えたら分割を検討する
- 構造体は `internal/types.go` に集約する。型が増えてきたらドメインごとにファイルを分けてよい
- 新しいユーティリティは既存ファイルへの追加を優先し、必要な場合のみ新ファイルを作る

## 新しいパッケージを作るとき

以下の順序で判断する:

1. 既存ファイルに追加できないか検討する
2. 複数パッケージから使う汎用ロジックなら `pkg/` に置く
3. CLI コマンドを増やすなら `internal/` に追加コマンドファイルを作る（例: `internal/validate.go`）
4. `cmd/` ディレクトリは不要。`main.go` のみで十分

## 禁止パターン

- `utils/`, `helpers/`, `common/` のような曖昧なパッケージ名は作らない
- `pkg/` が `internal/` を import しない（循環 import 禁止）
- グローバル変数でインスタンスを共有しない（`rootCmd` は Cobra の慣習として許容）
- `init()` 関数は `internal/root.go` のフラグ登録のみで使う
