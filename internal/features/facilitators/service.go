package facilitators

import (
	"errors"
	"fmt"
	"sort"
	"strings"

	"golang.org/x/text/collate"
	"golang.org/x/text/language"
)

// japaneseCollator は日本語ロケール対応の文字列比較器
// localeCompare("ja") 相当。毎回生成するとコストがかかるためパッケージ変数として保持する
var japaneseCollator = collate.New(language.Japanese)

// Service はファシリテーター一覧取得のビジネスロジック
type Service struct {
	repo Repository
}

// NewService は Service を生成する
func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

// List はクエリパラメータに基づいてファシリテーター一覧を返す
func (s *Service) List(params ListParams) (ListResponse, error) {
	if params.Page <= 0 {
		params.Page = 1
	}
	if params.Limit <= 0 {
		params.Limit = 20
	}
	if params.Order == "" {
		params.Order = "asc"
	}

	all, err := s.repo.FindAll()
	if err != nil {
		return ListResponse{}, fmt.Errorf("facilitators.FindAll: %w", err)
	}

	filtered := all
	if params.Search != "" {
		filtered = filterBySearch(all, params.Search)
	}

	// totalCount はページネーション前のフィルタ済み件数
	totalCount := len(filtered)

	if params.Sort != "" {
		sortFacilitators(filtered, params.Sort, params.Order)
	}

	data := paginate(filtered, params.Page, params.Limit)

	return ListResponse{
		Data:       data,
		TotalCount: totalCount,
	}, nil
}

// ErrInvalidSort は不正なソートキーを表す
var ErrInvalidSort = errors.New("invalid sort key")

func filterBySearch(items []Facilitator, search string) []Facilitator {
	s := strings.ToLower(search)
	result := make([]Facilitator, 0, len(items))
	for _, f := range items {
		if strings.Contains(strings.ToLower(f.Name), s) ||
			strings.Contains(strings.ToLower(f.LoginID), s) {
			result = append(result, f)
		}
	}
	return result
}

func sortFacilitators(items []Facilitator, sortKey, order string) {
	sort.SliceStable(items, func(i, j int) bool {
		var cmp int
		switch sortKey {
		case "loginId":
			cmp = japaneseCollator.CompareString(items[i].LoginID, items[j].LoginID)
		default: // "name"
			cmp = japaneseCollator.CompareString(items[i].Name, items[j].Name)
		}
		if order == "desc" {
			return cmp > 0
		}
		return cmp < 0
	})
}

// paginate はページネーション後のスライスを返す。
// 範囲外のページでは空スライスを返す（nil ではなく初期化済みのため JSON で [] になる）
func paginate(items []Facilitator, page, limit int) []Facilitator {
	start := (page - 1) * limit
	if start >= len(items) {
		return make([]Facilitator, 0)
	}
	end := start + limit
	if end > len(items) {
		end = len(items)
	}
	return items[start:end]
}
