// data パッケージは JSON データファイルを embed.FS で提供する
package data

import _ "embed"

// FacilitatorsJSON は data/facilitators.json の内容（embed.FS で埋め込み）
//
//go:embed facilitators.json
var FacilitatorsJSON []byte
