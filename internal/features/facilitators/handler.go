package facilitators

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"strconv"
)

// Handler は HTTP ハンドラー
type Handler struct {
	service *Service
}

// NewHandler は Handler を生成する
func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

//	@Summary		ファシリテーター一覧取得
//	@Tags			facilitators
//	@Produce		json
//	@Param			page	query		int		false	"ページ番号（1以上）"	default(1)	minimum(1)
//	@Param			limit	query		int		false	"件数（1-100）"			default(20)	minimum(1)	maximum(100)
//	@Param			sort	query		string	false	"ソートキー"			Enums(name,loginId)
//	@Param			order	query		string	false	"ソート順"				Enums(asc,desc)
//	@Param			search	query		string	false	"検索キーワード"
//	@Success		200		{object}	ListResponse
//	@Header			200		{string}	X-Request-Id	"リクエストID"
//	@Router			/api/facilitators [get]
func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	params, err := parseListParams(r)
	if err != nil {
		var he *httpError
		if errors.As(err, &he) {
			http.Error(w, he.message, he.code)
		} else {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
		return
	}

	result, err := h.service.List(params)
	if err != nil {
		slog.Error("ファシリテーター一覧取得に失敗", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(result); err != nil {
		slog.Error("レスポンスのエンコードに失敗", "error", err)
	}
}

func parseListParams(r *http.Request) (ListParams, error) {
	q := r.URL.Query()

	params := ListParams{
		Sort:   q.Get("sort"),
		Order:  q.Get("order"),
		Search: q.Get("search"),
	}

	if pageStr := q.Get("page"); pageStr != "" {
		page, err := strconv.Atoi(pageStr)
		if err != nil || page < 1 {
			return ListParams{}, &httpError{code: http.StatusBadRequest, message: "page は1以上の整数で指定してください"}
		}
		params.Page = page
	}

	if limitStr := q.Get("limit"); limitStr != "" {
		limit, err := strconv.Atoi(limitStr)
		if err != nil || limit < 1 || limit > 100 {
			return ListParams{}, &httpError{code: http.StatusBadRequest, message: "limit は1〜100の整数で指定してください"}
		}
		params.Limit = limit
	}

	if params.Sort != "" && params.Sort != "name" && params.Sort != "loginId" {
		return ListParams{}, &httpError{code: http.StatusBadRequest, message: "sort は name または loginId で指定してください"}
	}

	if params.Order != "" && params.Order != "asc" && params.Order != "desc" {
		return ListParams{}, &httpError{code: http.StatusBadRequest, message: "order は asc または desc で指定してください"}
	}

	return params, nil
}

type httpError struct {
	code    int
	message string
}

func (e *httpError) Error() string {
	return e.message
}
