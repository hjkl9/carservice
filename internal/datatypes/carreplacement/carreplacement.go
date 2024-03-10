package carreplacement

type Replacement struct {
	Id          uint    `db:"id"`
	ParentId    uint    `db:"parentId"`
	Title       string  `db:"title"`
	EstF32Price float32 `db:"estF32Price"`
	EstU64Price uint64  `db:"estU64Price"`
	Counter     uint8   `db:"counter"`
}
