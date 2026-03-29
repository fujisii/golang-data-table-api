package facilitators

// Facilitator はファシリテーターの型定義（types.ts の Facilitator 相当）
type Facilitator struct {
	ID      int    `json:"id"`
	Name    string `json:"name"`
	LoginID string `json:"loginId"`
}

// ListParams は GET /api/facilitators のクエリパラメータ（FacilitatorQueryParams 相当）
type ListParams struct {
	Page   int    // ページ番号（デフォルト: 1）
	Limit  int    // 1ページあたりの件数（デフォルト: 20、範囲: 1〜100）
	Sort   string // ソートキー: "name" | "loginId" | ""
	Order  string // ソート順: "asc" | "desc"（デフォルト: "asc"）
	Search string // name/loginId の部分一致検索キーワード（OR・大文字小文字無視）
}

// ListResponse は GET /api/facilitators のレスポンス型
type ListResponse struct {
	Data       []Facilitator `json:"data"`
	TotalCount int           `json:"totalCount"`
}
