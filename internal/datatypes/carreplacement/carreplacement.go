package carreplacement

type Replacement struct {
	Id          uint    `db:"id"`
	ParentId    uint    `db:"parentId"`
	Title       string  `db:"title"`
	EstF32Price float32 `db:"estF32Price"`
	EstU64Price uint64  `db:"estU64Price"`
	Counter     uint8   `db:"counter"`
}


// // ItemId
// // 实现 contract.ItemPricer 接口
// func (r *Replacement) ItemId() uint64 {
// 	return uint64(r.Id)
// }

// // GetUint64
// // 实现 contract.ItemPricer 接口
// func (r *Replacement) GetUint64() uint64 {
// 	return r.EstU64Price
// }

// // GetFloat64
// // 实现 contract.ItemPricer 接口
// func (r *Replacement) GetFloat64() float64 {
// 	return float64(r.EstF32Price)
// }

// ItemId
// 实现 contract.ItemPricer 接口
func (r Replacement) ItemId() uint64 {
	return uint64(r.Id)
}

// GetUint64
// 实现 contract.ItemPricer 接口
func (r Replacement) GetUint64() uint64 {
	return r.EstU64Price
}

// GetFloat64
// 实现 contract.ItemPricer 接口
func (r Replacement) GetFloat64() float64 {
	return float64(r.EstF32Price)
}
