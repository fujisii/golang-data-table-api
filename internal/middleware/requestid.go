package middleware

import (
	"net/http"

	chimiddleware "github.com/go-chi/chi/v5/middleware"
)

// SetRequestIDHeader は chi の RequestID ミドルウェアが生成したリクエスト ID を
// レスポンスヘッダー X-Request-Id に設定する（plugins/request-id.ts 相当）
func SetRequestIDHeader(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestID := chimiddleware.GetReqID(r.Context())
		if requestID != "" {
			w.Header().Set("X-Request-Id", requestID)
		}
		next.ServeHTTP(w, r)
	})
}
