package facilitators

import (
	"encoding/json"
	"fmt"
)

// Repository はデータアクセスのインターフェース
type Repository interface {
	FindAll() ([]Facilitator, error)
}

// JSONRepository は JSON データをメモリにキャッシュするリポジトリ実装
type JSONRepository struct {
	cache []Facilitator
}

// NewJSONRepository は JSON バイト列からリポジトリを初期化する
func NewJSONRepository(data []byte) (*JSONRepository, error) {
	var facilitators []Facilitator
	if err := json.Unmarshal(data, &facilitators); err != nil {
		return nil, fmt.Errorf("facilitators.json の解析に失敗: %w", err)
	}
	return &JSONRepository{cache: facilitators}, nil
}

func (r *JSONRepository) FindAll() ([]Facilitator, error) {
	// 呼び出し元による変更がキャッシュに影響しないようコピーを返す
	result := make([]Facilitator, len(r.cache))
	copy(result, r.cache)
	return result, nil
}
