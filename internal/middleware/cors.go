package middleware

import (
	"net/http"

	"github.com/go-chi/cors"
)

// NewCORS は CORS ミドルウェアを生成する（plugins/cors.ts 相当）
// 移植元と同様に http://localhost:5173 からのリクエストを許可し、
// X-Request-Id ヘッダーを公開する
func NewCORS() func(http.Handler) http.Handler {
	return cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://localhost:5173"},
		AllowedMethods:   []string{"GET", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Content-Type", "X-Request-Id"},
		ExposedHeaders:   []string{"X-Request-Id"},
		AllowCredentials: false,
		MaxAge:           300,
	})
}
