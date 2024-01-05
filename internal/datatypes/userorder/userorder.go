package userorder

import "database/sql"

type UserOrder struct {
	Id                  uint           `db:"id"`
	OrderNumber         string         `db:"orderNumber"`
	CarOwnerName        string         `db:"carOwnerName"`
	CarOwnerMultiLvAddr string         `db:"carOwnerMultiLvAddr"`
	CarOwnerFullAddress string         `db:"carOwnerFullAddress"`
	PartnerStore        sql.NullString `db:"partnerStore"`
	PartnerStoreAddress string         `db:"partnerStoreAddress"`
	CarBrandName        string         `db:"carBrandName"`
	CarSeriesName       string         `db:"carSeriesName"`
	Comment             string         `db:"comment"`
	OrderStatus         uint8          `db:"orderStatus"`
	CreatedAt           string         `db:"createdAt"`
	UpdatedAt           string         `db:"updatedAt"`
	// other fields.
}

type UserOrderListItem struct {
	Id           uint           `db:"id" json:"id"`
	OrderNumber  string         `db:"orderNumber" json:"orderNumber"`
	PartnerStore sql.NullString `db:"partnerStore" json:"partnerStore"`
	Requirements string         `db:"requirements" json:"requirements"`
	OrderStatus  uint8          `db:"orderStatus" json:"orderStatus"`
	CreatedAt    string         `db:"createdAt" json:"createdAt"`
	UpdatedAt    string         `db:"updatedAt" json:"updatedAt"`
}
