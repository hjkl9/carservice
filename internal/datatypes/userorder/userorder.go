package userorder

import "database/sql"

type Facades_UserOrder struct {
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

type Native_UserOrder struct {
	Id                   uint64         `db:"id"`
	CarBrandId           uint64         `db:"car_brand_id"`
	CarBrandSeriesId     uint64         `db:"car_brand_series_id"`
	MemberId             uint64         `db:"member_id"`
	CarInfoId            uint64         `db:"car_info_id"`
	OrderNumber          string         `db:"order_number"`
	EstAmount            float32        `db:"est_amount"`
	ActAmount            float32        `db:"act_amount"`
	ExpiredAt            sql.NullString `db:"expired_at"`
	PaymentMethod        int8           `db:"payment_method"`
	PaidAt               sql.NullString `db:"paid_at"`
	OrderStatus          int8           `db:"order_status"`
	Comment              string         `db:"comment"`
	CreatedAt            string         `db:"created_at"`
	UpdatedAt            string         `db:"updated_at"`
	PartnerStoreId       uint64         `db:"partner_store_id"`
	DeletedAt            sql.NullString `db:"deleted_at"`
	InstallerPhoneNumber sql.NullString `db:"installer_phone_number"`
	InstallerName        sql.NullString `db:"installer_name"`
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

type UpdatePayload struct {
	MemberId         uint `db:"member_id"`
	CarBrandId       uint `db:"car_brand_id"`
	CarBrandSeriesId uint `db:"car_brand_series_id"`
	// CarOwnerInfoId   uint    `db:"car_owner_info_id"` // ! deprecated
	PartnerStoreId uint    `db:"partner_store_id"` // ! deprecated
	OrderNumber    string  `db:"order_number"`
	Comment        string  `db:"comment"`
	EstAmount      float64 `db:"est_amount"`
	ActAmount      float64 `db:"act_amount"`
	PaymentMethod  uint8   `db:"payment_method"`
	OrderStatus    uint8   `db:"order_status"`
	// CreatedAt        time.Duration `db:"created_at"`
	// UpdatedAt        time.Duration `db:"updated_at"`
}

// createUserOrderPayload 创建用户订单数据
// 内存对齐 OK
type CreatePayload struct {
	MemberId         uint `db:"member_id"`
	CarBrandId       uint `db:"car_brand_id"`
	CarBrandSeriesId uint `db:"car_brand_series_id"`
	// CarOwnerInfoId   uint    `db:"car_owner_info_id"` // ! deprecated
	PartnerStoreId uint    `db:"partner_store_id"` // ! deprecated
	OrderNumber    string  `db:"order_number"`
	Comment        string  `db:"comment"`
	EstAmount      float64 `db:"est_amount"`
	ActAmount      float64 `db:"act_amount"`
	PaymentMethod  uint8   `db:"payment_method"`
	OrderStatus    uint8   `db:"order_status"`
	// CreatedAt        time.Duration `db:"created_at"`
	// UpdatedAt        time.Duration `db:"updated_at"`
}
