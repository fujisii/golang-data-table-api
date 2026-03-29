// @title			golang-data-table-api
// @version		1.0.0
// @description	ファシリテーター一覧 API（react-data-table-api の Go 移植版）
// @host			localhost:3000
// @BasePath		/
package main

import (
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
	httpSwagger "github.com/swaggo/http-swagger"

	appdata "github.com/fujisii/golang-data-table-api/data"
	_ "github.com/fujisii/golang-data-table-api/docs"
	"github.com/fujisii/golang-data-table-api/internal/features/facilitators"
	"github.com/fujisii/golang-data-table-api/internal/middleware"
)

func main() {
	// リポジトリ・サービス・ハンドラーの初期化
	repo, err := facilitators.NewJSONRepository(appdata.FacilitatorsJSON)
	if err != nil {
		slog.Error("リポジトリの初期化に失敗", "error", err)
		return
	}
	svc := facilitators.NewService(repo)
	h := facilitators.NewHandler(svc)

	// ルーターの設定
	r := chi.NewRouter()

	// ミドルウェアの登録（plugins/ 相当）
	r.Use(chimiddleware.RequestID)    // リクエスト ID の生成
	r.Use(middleware.SetRequestIDHeader) // X-Request-Id レスポンスヘッダーの設定
	r.Use(middleware.NewCORS())       // CORS の設定
	r.Use(chimiddleware.Logger)       // アクセスログ
	r.Use(chimiddleware.Recoverer)    // パニックリカバリー

	// API ルートの登録
	r.Get("/api/facilitators", h.List)

	// Swagger UI の登録
	r.Get("/docs/*", httpSwagger.WrapHandler)

	slog.Info("サーバーを起動します", "addr", ":3000")
	if err := http.ListenAndServe(":3000", r); err != nil {
		slog.Error("サーバーの起動に失敗", "error", err)
	}
}
