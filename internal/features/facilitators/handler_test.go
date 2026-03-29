package facilitators_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/fujisii/golang-data-table-api/internal/features/facilitators"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// newTestHandler はテスト用の Handler を生成する
func newTestHandler(data []facilitators.Facilitator) *facilitators.Handler {
	repo := &stubRepository{data: data}
	svc := facilitators.NewService(repo)
	return facilitators.NewHandler(svc)
}

func TestHandlerList(t *testing.T) {
	t.Run("デフォルトパラメータで200を返す", func(t *testing.T) {
		h := newTestHandler(testData)
		req := httptest.NewRequest(http.MethodGet, "/api/facilitators", nil)
		w := httptest.NewRecorder()

		h.List(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, "application/json", w.Header().Get("Content-Type"))

		var result facilitators.ListResponse
		require.NoError(t, json.Unmarshal(w.Body.Bytes(), &result))
		assert.Len(t, result.Data, 5)
		assert.Equal(t, 5, result.TotalCount)
	})

	t.Run("ページネーションパラメータが反映される", func(t *testing.T) {
		h := newTestHandler(testData)
		req := httptest.NewRequest(http.MethodGet, "/api/facilitators?page=2&limit=2", nil)
		w := httptest.NewRecorder()

		h.List(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var result facilitators.ListResponse
		require.NoError(t, json.Unmarshal(w.Body.Bytes(), &result))
		assert.Len(t, result.Data, 2)
		assert.Equal(t, 5, result.TotalCount)
	})

	t.Run("不正な page パラメータで400を返す", func(t *testing.T) {
		h := newTestHandler(testData)
		req := httptest.NewRequest(http.MethodGet, "/api/facilitators?page=0", nil)
		w := httptest.NewRecorder()

		h.List(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("不正な limit パラメータで400を返す", func(t *testing.T) {
		h := newTestHandler(testData)
		req := httptest.NewRequest(http.MethodGet, "/api/facilitators?limit=101", nil)
		w := httptest.NewRecorder()

		h.List(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("不正な sort パラメータで400を返す", func(t *testing.T) {
		h := newTestHandler(testData)
		req := httptest.NewRequest(http.MethodGet, "/api/facilitators?sort=invalid", nil)
		w := httptest.NewRecorder()

		h.List(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}
