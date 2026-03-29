---
name: api-implementer
description: Go API の実装を進める際に使用するエージェント。移植元の TypeScript 実装を参照しながら、正しい順序でレイヤーを実装する。
---

# API 実装エージェント

## 実装順序チェックリスト

以下の順序で実装する。各ステップ完了後に `go build ./...` が通ることを確認してから次へ進む。

1. **types.go** — 型定義（`Facilitator`, `ListParams`, `ListResponse`）
2. **repository.go** — `Repository` インターフェース + `JSONRepository` 実装（`embed.FS` 使用）
3. **service.go** — `Service` 構造体 + `List()` メソッド（フィルタ・ソート・ページネーション）
4. **service_test.go** — サービス単体テスト（移植元のテストケースを全件移植）
5. **handler.go** — `Handler` 構造体 + `List()` HTTP ハンドラー（swaggo アノテーション付き）
6. **handler_test.go** — `httptest` を使った統合テスト
7. **middleware/cors.go** — CORS ミドルウェア設定
8. **middleware/requestid.go** — `X-Request-Id` レスポンスヘッダー設定
9. **cmd/server/main.go** — ルーター・ミドルウェア・ルート接続

## 移植元の参照先

| 実装内容 | 移植元ファイル |
|---|---|
| 型定義 | `../react-data-table-api/src/features/facilitators/types.ts` |
| ビジネスロジック | `../react-data-table-api/src/features/facilitators/services/facilitator.service.ts` |
| テストケース | `../react-data-table-api/src/features/facilitators/services/facilitator.service.test.ts` |
| HTTP ハンドラー | `../react-data-table-api/src/features/facilitators/handlers/facilitator.handler.ts` |
| CORS 設定 | `../react-data-table-api/src/plugins/cors.ts` |

## 各ステップの実装ガイド

### service.go の実装

`List()` メソッドの処理順序（移植元 `facilitator.service.ts` と同じ）:
1. `repo.FindAll()` で全件取得
2. `search` が空でなければ name/loginId の部分一致フィルタ（OR・大文字小文字無視）
3. フィルタ後の件数を `totalCount` に記録
4. `sort` が指定されていれば `golang.org/x/text/collate` で日本語ロケール対応ソート
5. `items[(page-1)*limit : min(page*limit, len(items))]` でページネーション

**注意:** Go の `strings.Compare` は日本語非対応。必ず `collate.New(language.Japanese)` を使用する。

### service_test.go の移植対応表

| TypeScript（Vitest） | Go（testing + testify） |
|---|---|
| `describe("...", () => {})` | `func TestList(t *testing.T) { ... }` |
| `it("...", async () => {})` | `t.Run("...", func(t *testing.T) { ... })` |
| `vi.fn().mockResolvedValue([...])` | スタブ構造体 `stubRepository` を実装 |
| `expect(x).toHaveLength(n)` | `assert.Len(t, x, n)` |
| `expect(x).toBe(n)` | `assert.Equal(t, n, x)` |
| `expect(x).toEqual(y)` | `assert.Equal(t, y, x)` |

### handler.go の swaggo アノテーション

```go
// List はファシリテーター一覧を返す
// @Summary      ファシリテーター一覧取得
// @Tags         facilitators
// @Produce      json
// @Param        page    query  int    false  "ページ番号（1以上）"  default(1)  minimum(1)
// @Param        limit   query  int    false  "件数（1-100）"        default(20) minimum(1) maximum(100)
// @Param        sort    query  string false  "ソートキー"           Enums(name,loginId)
// @Param        order   query  string false  "ソート順"             Enums(asc,desc)
// @Param        search  query  string false  "検索キーワード"
// @Success      200     {object} ListResponse
// @Header       200     {string} X-Request-Id  "リクエストID"
// @Router       /api/facilitators [get]
func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
```
