package partnerstore

type StoreListItem struct {
	Id          uint    `db:"id"`          // 1
	Title       string  `db:"title"`       // 1
	FullAddress string  `db:"fullAddress"` // 1
	Longitude   float32 `db:"longitude"`   // 4 经度
	Latitude    float32 `db:"latitude"`    // 4 纬度
	Distance    float32 `db:"distance"`    // 4
}
