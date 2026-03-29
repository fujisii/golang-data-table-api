---
name: reviewer
description: Go コードのレビューを行うエージェント。移植の正確性・Go らしさ・セキュリティの観点でチェックする。
---

# コードレビューエージェント

## レビューチェックリスト

### 必須確認項目

- [ ] エラーを `_` で捨てていないこと（`result, _ := ...` は禁止）
- [ ] `Repository` インターフェースが定義され、`Service` がインターフェース経由でリポジトリを使用していること（テストでモック可能）
- [ ] swaggo アノテーションがハンドラー関数に記載されていること
- [ ] `X-Request-Id` ミドルウェアが `main.go` でルーターに接続されていること
- [ ] CORS の `AllowedOrigins` が `http://localhost:5173` に設定されていること
- [ ] CORS の `ExposedHeaders` に `X-Request-Id` が含まれていること
- [ ] JSON レスポンスのフィールド名が移植元と一致していること（`data`, `totalCount`, `loginId` 等）

### ビジネスロジックの正確性

- [ ] 検索は name と loginId の**両方**に対して部分一致（OR）していること
- [ ] 検索は大文字小文字を無視していること（`strings.EqualFold` または `strings.ToLower`）
- [ ] ソートに `golang.org/x/text/collate` を使用していること（`strings.Compare` は不可）
- [ ] `totalCount` はページネーション前のフィルタ済み件数であること
- [ ] ページ範囲外では空配列（`[]Facilitator{}`）を返すこと（`nil` ではなく初期化済みスライス）

### Go コーディング規約

- [ ] `gofmt` / `goimports` でフォーマット済みであること
- [ ] エクスポートされた型・関数にコメントが付いていること
- [ ] ログに `log/slog` を使用していること（`fmt.Println` でのデバッグログは除去）
- [ ] エラーラップは `fmt.Errorf("...: %w", err)` 形式であること

### テスト

- [ ] 移植元テストの全ケース（10件）が網羅されていること
- [ ] サービステストでスタブリポジトリを使用していること（実ファイルI/O不使用）
- [ ] ハンドラーテストで `httptest.NewRecorder()` と `httptest.NewRequest()` を使用していること
- [ ] `go test ./... -count=1` がすべてパスすること

## よくある移植ミス

| TypeScript の振る舞い | Go での注意点 |
|---|---|
| `localeCompare("ja")` | `strings.Compare` は日本語非対応 → `collate.New(language.Japanese)` を使う |
| `Promise.reject(new Error(...))` | Go はエラーを戻り値で返す。`panic` は使わない |
| `undefined` チェック | Go の `string` のゼロ値は `""` なので `if params.Sort == ""` でチェック |
| 配列の末尾スライス越え | `items[start:end]` で `end > len(items)` の場合パニックする → `min(end, len(items))` を使う |
| `JSON.stringify` の null | Go で初期化されていないスライスは JSON で `null` になる → `make([]T, 0)` で初期化する |
