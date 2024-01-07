package userorder

import (
	"context"
	"database/sql"
	"fmt"

	"carservice/internal/data/tables"
	"carservice/internal/enum/userorder"
	"carservice/internal/pkg/common/errcode"
	"carservice/internal/pkg/jwt"
	"carservice/internal/svc"
	"carservice/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetUserOrderLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetUserOrderLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserOrderLogic {
	return &GetUserOrderLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

type Order struct {
	Id                  uint           `db:"id"`
	OrderNumber         string         `db:"orderNumber"`
	CarOwnerName        string         `db:"carOwnerName"`
	CarOwnerMultiLvAddr string         `db:"carOwnerMultiLvAddr"`
	CarOwnerFullAddress string         `db:"carOwnerFullAddress"`
	PartnerStore        sql.NullString `db:"partnerStore"`
	PartnerStoreAddress sql.NullString `db:"partnerStoreAddress"`
	CarBrandName        string         `db:"carBrandName"`
	CarSeriesName       string         `db:"carSeriesName"`
	Comment             string         `db:"comment"`
	OrderStatus         uint8          `db:"orderStatus"`
	CreatedAt           string         `db:"createdAt"`
	UpdatedAt           string         `db:"updatedAt"`
	// other fields.
}

func (l *GetUserOrderLogic) GetUserOrder(req *types.GetUserOrderReq) (resp *types.GetUserOrderRep, err error) {
	userId := jwt.GetUserId(l.ctx)
	// 检查订单是否存在
	var count int = 0
	query := "SELECT COUNT(1) AS `count` FROM `%s` WHERE `id` = ? AND `member_id` = ? AND `deleted_at` IS NULL LIMIT 1"
	if err = l.svcCtx.DBC.Get(&count, fmt.Sprintf(query, tables.UserOrder), req.Id, userId); err != nil {
		return nil, errcode.DatabaseError.SetDetails(err.Error())
	}
	if count == 0 {
		return nil, errcode.NotFound.SetMsg("该用户订单不存在或已被删除")
	}
	// 通过 ID 查询订单
	var order Order
	query = "SELECT `uo`.`id` AS `id`, `uo`.`order_number` AS `orderNumber`, `coi`.`name` AS `carOwnerName`, `coi`.`multilevel_address` AS `carOwnerMultiLvAddr`, `coi`.`full_address` AS `carOwnerFullAddress` , `ps`.`title` AS `partnerStore`, `ps`.`full_address` AS `partnerStoreAddress`, `uo`.`comment` AS `comment`, `cb`.`brand_name` AS `carBrandName`, `cbs`.`series_name` AS `carSeriesName` , `uo`.`order_status` AS `orderStatus`, `uo`.`created_at` AS `createdAt`, `uo`.`updated_at` AS `updatedAt` FROM `user_orders` `uo` LEFT JOIN `partner_stores` `ps` ON `ps`.`id` = `uo`.`partner_store_id` JOIN `car_owner_infos` `coi` ON `coi`.`id` = `uo`.`car_owner_info_id` JOIN `car_brands` `cb` ON `cb`.`brand_id` = `uo`.`car_brand_id` JOIN `car_brand_series` `cbs` ON `cbs`.`series_id` = `uo`.`car_brand_series_id` WHERE `uo`.`id` = ? AND `uo`.`member_id` = ? LIMIT 1"

	stmt, err := l.svcCtx.DBC.PreparexContext(l.ctx, query)
	if err != nil {
		return nil, errcode.DatabaseError.SetDetails(err.Error())
	}
	if err = stmt.GetContext(l.ctx, &order, req.Id, userId); err != nil {
		return nil, errcode.DatabaseError.SetDetails(err.Error())
	}
	return &types.GetUserOrderRep{
		Id:                  order.Id,
		OrderNumber:         order.OrderNumber,
		CarOwnerName:        order.CarOwnerName,
		CarOwnerMultiLvAddr: order.CarOwnerMultiLvAddr,
		CarOwnerFullAddr:    order.CarOwnerFullAddress,
		CarBrandName:        order.CarBrandName,
		CarSeriesName:       order.CarSeriesName,
		PartnerStore: func() string {
			if !order.PartnerStore.Valid {
				return "未绑定合作门店"
			}
			return order.PartnerStore.String
		}(),
		PartnerStoreAddr: func() string {
			if !order.PartnerStoreAddress.Valid {
				return "未绑定合作门店"
			}
			return order.PartnerStoreAddress.String
		}(),
		Requirements: order.Comment,
		OrderStatus:  userorder.OrderStatusDesc(order.OrderStatus),
		CreatedAt:    order.CreatedAt,
		UpdatedAt:    order.UpdatedAt,
	}, nil
}
