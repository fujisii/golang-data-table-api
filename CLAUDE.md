# CLAUDE.md - プロジェクトガイド

## 基本ルール

- **日本語**でやり取りすること
- `go get` は自動実行しないこと。新規依存追加が必要な場合はユーザーに確認してから `go get` → `go mod tidy` を実行する
- コード識別子（変数名・関数名・型名）は**英語**、コメント・ドキュメントは**日本語**
- ファイル書き込み後は `go build ./...` と `go vet ./...` が自動実行される（フックによる）

## 移植元リポジトリ

| 項目 | 内容 |
|---|---|
| パス | `../react-data-table-api` |
| 技術スタック | Node.js 24 / Fastify / TypeScript |
| エントリーポイント | `../react-data-table-api/src/server.ts` |
| ビジネスロジック | `../react-data-table-api/src/features/facilitators/services/facilitator.service.ts` |
| テスト | `../react-data-table-api/src/features/facilitators/services/facilitator.service.test.ts` |
| 型定義 | `../react-data-table-api/src/features/facilitators/types.ts` |
| データファイル | `../react-data-table-api/data/facilitators.json` |

## 技術スタック

| カテゴリ | 採用技術 |
|---|---|
| 言語 | Go 1.23+ |
| ルーター | `github.com/go-chi/chi/v5` |
| CORS | `github.com/go-chi/cors` |
| OpenAPI | `github.com/swaggo/swag`（CLI）+ `github.com/swaggo/http-swagger` |
| テスト | 標準 `testing` + `github.com/stretchr/testify` |
| 日本語ソート | `golang.org/x/text/collate` |
| リクエストID | `github.com/google/uuid` |

## アーキテクチャ

### ディレクトリ構造

```
golang-data-table-api/
├── cmd/
│   └── server/
│       └── main.go                  # エントリーポイント（server.ts 相当）
├── internal/
│   ├── features/
│   │   └── facilitators/
│   │       ├── types.go             # types.ts 相当
│   │       ├── repository.go        # facilitator.repository.ts 相当
│   │       ├── service.go           # facilitator.service.ts 相当
│   │       ├── service_test.go      # facilitator.service.test.ts 相当
│   │       ├── handler.go           # facilitator.handler.ts 相当（swaggo アノテーション含む）
│   │       └── handler_test.go      # httptest 統合テスト
│   └── middleware/
│       ├── cors.go                  # plugins/cors.ts 相当
│       └── requestid.go             # plugins/request-id.ts 相当
├── docs/                            # swag init が自動生成（コミットしない）
├── data/
│   └── facilitators.json            # 移植元データ
├── go.mod
├── go.sum
└── CLAUDE.md
```

### レイヤーの責務

| レイヤー | ファイル | 責務 | Node.js 相当 |
|---|---|---|---|
| Handler | `handler.go` | HTTP リクエスト受信・レスポンス送信 | `facilitator.handler.ts` |
| Service | `service.go` | 検索・ソート・ページネーションのビジネスロジック | `facilitator.service.ts` |
| Repository | `repository.go` | JSON ファイルからのデータ取得 | `facilitator.repository.ts` |
| Types | `types.go` | 型定義 | `types.ts` |

**重要:** TypeScript は 4 ディレクトリ構成だが、Go ではパッケージが封装単位なので `facilitators/` 内に同居させている（クロスパッケージ循環インポートを回避）。

## API 仕様

### エンドポイント

`GET /api/facilitators`

### クエリパラメータ

| パラメータ | 型 | デフォルト | 制約 | 説明 |
|---|---|---|---|---|
| `page` | int | 1 | 1以上 | ページ番号 |
| `limit` | int | 20 | 1〜100 | 1ページあたりの件数 |
| `sort` | string | - | `name` または `loginId` | ソートキー |
| `order` | string | `asc` | `asc` または `desc` | ソート順 |
| `search` | string | - | - | name/loginId の部分一致（OR）大文字小文字無視 |

### レスポンス

```json
{
  "data": [
    { "id": 1, "name": "田中太郎", "loginId": "tanaka_taro" }
  ],
  "totalCount": 200
}
```

### ヘッダー

- リクエスト: `X-Request-Id` （chi の RequestID ミドルウェアが生成・引き継ぎ）
- レスポンス: `X-Request-Id` を `Access-Control-Expose-Headers` で公開
- CORS: `http://localhost:5173` を許可

## Go コーディング規約

- `gofmt` / `goimports` を使用する
- エラーは必ず呼び出し元でハンドリングする（`_` で捨てない）
- ログは標準ライブラリ `log/slog` を使用する
- インターフェースはモック可能性のため `repository.go` に定義する（`Repository interface`）
- エラーラップは `fmt.Errorf("contextmsg: %w", err)` の形式を使用する
- JSON フィールド名は移植元と同一にする（`json:"loginId"` 等）

## 日本語ソートの注意点

移植元の `localeCompare("ja")` は日本語ロケール対応ソート。
Go 標準の `strings.Compare` は非対応のため、**`golang.org/x/text/collate`** を使用する。

```go
import "golang.org/x/text/collate"
import "golang.org/x/text/language"

c := collate.New(language.Japanese)
// c.CompareString(a, b) で localeCompare("ja") 相当
```

## 開発コマンド

```bash
go run ./cmd/server                      # 開発サーバー起動 (port 3000)
go build -o bin/server ./cmd/server      # バイナリビルド（bin/ に出力）
go build ./...                           # ビルド確認（バイナリ出力なし）
go vet ./...                             # 静的解析
go test ./...                # 全テスト実行
go test -v -run TestXxx      # 特定テスト実行
go mod tidy                  # 依存関係整理
swag init -g cmd/server/main.go -o docs  # OpenAPI ドキュメント生成
```

## 自動フック（`.claude/settings.json`）

- **Write/Edit 後**: `go build ./... && go vet ./...` が自動実行される
- テストは自動実行されないため、適宜 `go test ./...` を手動実行すること

## Swagger UI

- UI: `http://localhost:3000/docs/index.html`
- JSON スキーマ: `http://localhost:3000/docs/doc.json`
- `docs/` は `.gitignore` に追加する（`swag init` で都度生成）

## コミットメッセージ規約

Conventional Commits に従い英語1行で記述する:
- `feat: add facilitators list endpoint`
- `fix: handle empty search parameter correctly`
- `test: add service unit tests for pagination`
