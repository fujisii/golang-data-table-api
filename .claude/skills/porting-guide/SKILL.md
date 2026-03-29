---
name: porting-guide
description: TypeScript/Node.js から Go への移植ガイド。概念対応表とコードパターンの変換方法を提供する。
---

# TypeScript → Go 移植ガイド

## 概念対応表

| TypeScript / Node.js | Go | 備考 |
|---|---|---|
| `interface Foo { ... }` | `type Foo struct { ... }` | Go の struct が TypeScript interface に相当 |
| `interface Repository` | `type Repository interface { ... }` | Go にも interface がある |
| `class Service { constructor(private repo) {} }` | `type Service struct { repo Repository }` + `func NewService(repo Repository) *Service` | コンストラクタは関数で表現 |
| `async function list(): Promise<T>` | `func List() (T, error)` | 非同期は不要、エラーを戻り値で返す |
| `await repo.findAll()` | `repo.FindAll()` | 同期呼び出し |
| `Promise.reject(new Error("msg"))` | `return T{}, fmt.Errorf("msg")` | panic は使わない |
| `try { ... } catch (e) { ... }` | `result, err := ...; if err != nil { ... }` | エラーは戻り値でハンドリング |
| `throw new Error("msg")` | `return fmt.Errorf("msg")` | |
| `arr.filter(x => ...)` | `for _, x := range arr { if ... { result = append(result, x) } }` | |
| `arr.sort((a, b) => a.localeCompare(b, "ja"))` | `collate.New(language.Japanese)` + `sort.Slice(...)` | |
| `arr.slice(start, end)` | `arr[start:end]`（end は `min(end, len(arr))` でガード） | |
| `str.toLowerCase().includes(search)` | `strings.Contains(strings.ToLower(str), strings.ToLower(search))` | |
| `JSON.stringify(obj)` | `json.Marshal(obj)` | |
| `JSON.parse(str)` | `json.Unmarshal([]byte(str), &obj)` | |
| `fs.readFileSync(path)` | `os.ReadFile(path)` または `//go:embed` | |
| `vi.fn().mockResolvedValue(x)` | スタブ構造体でインターフェースを実装 | |
| `expect(x).toBe(y)` | `assert.Equal(t, y, x)` | testify を使用 |
| `expect(arr).toHaveLength(n)` | `assert.Len(t, arr, n)` | |
| `fastify.inject({ method, url })` | `httptest.NewRequest(method, url, nil)` + `httptest.NewRecorder()` | |

## よく使うコードパターン

### スタブリポジトリ（Vitest の vi.fn() 相当）

```go
// TypeScript:
// repository = { findAll: vi.fn().mockResolvedValue([...testData]) }

// Go:
type stubRepository struct {
    data []Facilitator
    err  error
}

func (r *stubRepository) FindAll() ([]Facilitator, error) {
    return slices.Clone(r.data), r.err
}

// テスト内での使用:
repo := &stubRepository{data: testData}
svc := NewService(repo)
```

### 部分一致検索（大文字小文字無視）

```go
// TypeScript:
// items.filter(f =>
//   f.name.toLowerCase().includes(search.toLowerCase()) ||
//   f.loginId.toLowerCase().includes(search.toLowerCase())
// )

// Go:
func matchSearch(f Facilitator, search string) bool {
    s := strings.ToLower(search)
    return strings.Contains(strings.ToLower(f.Name), s) ||
        strings.Contains(strings.ToLower(f.LoginID), s)
}
```

### 日本語ロケール対応ソート

```go
// TypeScript:
// items.sort((a, b) => a.name.localeCompare(b.name, "ja"))

// Go:
import (
    "golang.org/x/text/collate"
    "golang.org/x/text/language"
)

c := collate.New(language.Japanese)
sort.Slice(items, func(i, j int) bool {
    return c.CompareString(items[i].Name, items[j].Name) < 0
})
```

### JSON レスポンスの送信

```go
// TypeScript（Fastify）:
// return reply.send({ data, totalCount })

// Go:
w.Header().Set("Content-Type", "application/json")
w.WriteHeader(http.StatusOK)
json.NewEncoder(w).Encode(ListResponse{Data: data, TotalCount: totalCount})
```

### ページネーションのスライス（範囲外ガード）

```go
// TypeScript:
// const start = (page - 1) * limit
// return items.slice(start, start + limit)

// Go:
start := (params.Page - 1) * params.Limit
end := start + params.Limit
if start >= len(items) {
    return []Facilitator{}, nil  // nil ではなく空スライスを返す（JSONで [] になる）
}
if end > len(items) {
    end = len(items)
}
return items[start:end], nil
```

### embed.FS でデータファイルを埋め込む

```go
// TypeScript:
// import data from "../../data/facilitators.json"

// Go（cmd/server/main.go または repository.go）:
import _ "embed"

//go:embed data/facilitators.json
var facilitatorsJSON []byte
```

## 注意点

1. **`nil` vs 空スライス**: Go で初期化していないスライス（`var s []T`）は JSON で `null` になる。`make([]T, 0)` または `[]T{}` で初期化すると `[]` になる。
2. **エラーの `panic` 禁止**: TypeScript の `throw` に相当する `panic` は、プログラムをクラッシュさせる。代わりに `error` を戻り値で返す。
3. **ゴルーチンは不要**: この API は単純な同期処理で十分。非同期処理は不要。
