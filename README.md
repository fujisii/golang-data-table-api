# golang-data-table-api

[react-data-table-api](https://github.com/fujisii/react-data-table-api)（Node.js / Fastify / TypeScript）の Go 移植版。

## 技術スタック

- **言語**: Go 1.26+
- **ルーター**: [chi](https://github.com/go-chi/chi)
- **API ドキュメント**: [swaggo](https://github.com/swaggo/swag)
- **テスト**: 標準 `testing` + [testify](https://github.com/stretchr/testify)

## API

### `GET /api/facilitators`

ファシリテーター一覧を返す。検索・ソート・ページネーションに対応。

| パラメータ | 型 | デフォルト | 説明 |
|---|---|---|---|
| `page` | int | 1 | ページ番号（1以上） |
| `limit` | int | 20 | 1ページあたりの件数（1〜100） |
| `sort` | string | - | ソートキー: `name` または `loginId` |
| `order` | string | `asc` | ソート順: `asc` または `desc` |
| `search` | string | - | name/loginId の部分一致検索（OR・大文字小文字無視） |

**レスポンス例**

```json
{
  "data": [
    { "id": 1, "name": "田中太郎", "loginId": "tanaka_taro" }
  ],
  "totalCount": 200
}
```

## 開発

### 前提条件

- Go 1.26 以上
- [swag CLI](https://github.com/swaggo/swag)（Swagger ドキュメント生成用）

```bash
go install github.com/swaggo/swag/cmd/swag@latest
```

### セットアップ

```bash
git clone https://github.com/fujisii/golang-data-table-api.git
cd golang-data-table-api
go mod tidy
```

### サーバー起動

```bash
# 開発時（バイナリ生成なし）
go run ./cmd/server

# バイナリをビルドして起動（bin/ に出力）
go build -o bin/server ./cmd/server
./bin/server
# http://localhost:3000 で起動
```

### テスト

```bash
go test ./...
```

### Swagger ドキュメント再生成

ハンドラーのアノテーションを変更した場合に実行。

```bash
swag init -g cmd/server/main.go -o docs
```

Swagger UI: http://localhost:3000/docs/index.html

## ディレクトリ構成

```
.
├── cmd/server/          # エントリーポイント
├── data/                # モックデータ（JSON）
├── docs/                # swag init で自動生成される OpenAPI ドキュメント
└── internal/
    ├── features/
    │   └── facilitators/ # ファシリテーター機能（型定義・リポジトリ・サービス・ハンドラー）
    └── middleware/       # CORS・リクエストID ミドルウェア
```
