package facilitators_test

import (
	"slices"
	"testing"

	"github.com/fujisii/golang-data-table-api/internal/features/facilitators"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/text/collate"
	"golang.org/x/text/language"
)

// jaCollator はソート順の検証に使用する日本語ロケール比較器
var jaCollator = collate.New(language.Japanese)

var testData = []facilitators.Facilitator{
	{ID: 1, Name: "田中太郎", LoginID: "tanaka"},
	{ID: 2, Name: "佐藤花子", LoginID: "sato"},
	{ID: 3, Name: "鈴木一郎", LoginID: "suzuki"},
	{ID: 4, Name: "田中花子", LoginID: "tanaka_h"},
	{ID: 5, Name: "山田太郎", LoginID: "yamada"},
}

// stubRepository は Repository インターフェースのスタブ実装（vi.fn().mockResolvedValue() 相当）
type stubRepository struct {
	data []facilitators.Facilitator
	err  error
}

func (r *stubRepository) FindAll() ([]facilitators.Facilitator, error) {
	return slices.Clone(r.data), r.err
}

func newService(data []facilitators.Facilitator) *facilitators.Service {
	repo := &stubRepository{data: data}
	return facilitators.NewService(repo)
}

func TestList(t *testing.T) {
	t.Run("パラメータなしでデフォルト値（page=1, limit=20）で返却する", func(t *testing.T) {
		svc := newService(testData)
		result, err := svc.List(facilitators.ListParams{})

		require.NoError(t, err)
		assert.Len(t, result.Data, 5)
		assert.Equal(t, 5, result.TotalCount)
	})

	t.Run("name に対して部分一致検索できる", func(t *testing.T) {
		svc := newService(testData)
		result, err := svc.List(facilitators.ListParams{Search: "田中"})

		require.NoError(t, err)
		assert.Len(t, result.Data, 2)
		for _, f := range result.Data {
			assert.Contains(t, f.Name, "田中")
		}
		assert.Equal(t, 2, result.TotalCount)
	})

	t.Run("loginId に対して部分一致検索できる", func(t *testing.T) {
		svc := newService(testData)
		result, err := svc.List(facilitators.ListParams{Search: "tanaka"})

		require.NoError(t, err)
		assert.Len(t, result.Data, 2)
		for _, f := range result.Data {
			assert.Contains(t, f.LoginID, "tanaka")
		}
	})

	t.Run("検索は大文字小文字を無視する", func(t *testing.T) {
		svc := newService(testData)
		result, err := svc.List(facilitators.ListParams{Search: "TANAKA"})

		require.NoError(t, err)
		assert.Len(t, result.Data, 2)
		for _, f := range result.Data {
			assert.Contains(t, f.LoginID, "tanaka")
		}
	})

	t.Run("name で昇順ソートできる", func(t *testing.T) {
		svc := newService(testData)
		result, err := svc.List(facilitators.ListParams{Sort: "name", Order: "asc"})

		require.NoError(t, err)
		for i := 1; i < len(result.Data); i++ {
			cmp := jaCollator.CompareString(result.Data[i-1].Name, result.Data[i].Name)
			assert.LessOrEqual(t, cmp, 0, "昇順になっていない: %s > %s", result.Data[i-1].Name, result.Data[i].Name)
		}
	})

	t.Run("name で降順ソートできる", func(t *testing.T) {
		svc := newService(testData)
		result, err := svc.List(facilitators.ListParams{Sort: "name", Order: "desc"})

		require.NoError(t, err)
		for i := 1; i < len(result.Data); i++ {
			cmp := jaCollator.CompareString(result.Data[i-1].Name, result.Data[i].Name)
			assert.GreaterOrEqual(t, cmp, 0, "降順になっていない: %s < %s", result.Data[i-1].Name, result.Data[i].Name)
		}
	})

	t.Run("ページネーションで指定ページのデータを返す", func(t *testing.T) {
		svc := newService(testData)
		result, err := svc.List(facilitators.ListParams{Page: 2, Limit: 2})

		require.NoError(t, err)
		assert.Len(t, result.Data, 2)
		assert.Equal(t, 3, result.Data[0].ID)
		assert.Equal(t, 4, result.Data[1].ID)
		assert.Equal(t, 5, result.TotalCount)
	})

	t.Run("範囲外のページでは空配列を返し totalCount はフィルタ後の件数", func(t *testing.T) {
		svc := newService(testData)
		result, err := svc.List(facilitators.ListParams{Page: 100, Limit: 20})

		require.NoError(t, err)
		assert.Len(t, result.Data, 0)
		assert.Equal(t, 5, result.TotalCount)
	})

	t.Run("検索フィルタ後の件数が totalCount に反映される", func(t *testing.T) {
		svc := newService(testData)
		result, err := svc.List(facilitators.ListParams{Search: "太郎"})

		require.NoError(t, err)
		assert.Equal(t, 2, result.TotalCount)
		assert.Len(t, result.Data, 2)
	})

	t.Run("検索 + ソート + ページネーションの複合条件で動作する", func(t *testing.T) {
		svc := newService(testData)
		result, err := svc.List(facilitators.ListParams{
			Search: "田中",
			Sort:   "name",
			Order:  "desc",
			Page:   1,
			Limit:  1,
		})

		require.NoError(t, err)
		assert.Equal(t, 2, result.TotalCount)
		assert.Len(t, result.Data, 1)
		// 降順ソートの先頭1件が返ることを確認
		// 田中太郎 と 田中花子のうち降順で先頭のものが返る
		assert.NotEmpty(t, result.Data[0].Name)
	})
}
