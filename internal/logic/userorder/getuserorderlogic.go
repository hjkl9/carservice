package userorder

import (
	"context"
	"fmt"

	"carservice/internal/data/tables"
	"carservice/internal/pkg/common/errcode"
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

func (l *GetUserOrderLogic) GetUserOrder(req *types.GetUserOrderReq) (resp *types.GetUserOrderRep, err error) {
	// 检查订单是否存在
	var count int = 0
	query := "SELECT COUNT(1) AS `count` FROM `%s` WHERE `id` = ? LIMIT 1"
	if err = l.svcCtx.DBC.Get(&count, fmt.Sprintf(query, tables.UserOrder), req.Id); err != nil {
		return nil, errcode.DatabaseError.SetDetails(err.Error())
	}
	// 通过 ID 查询订单
	query = "SELECT `id` AS `id`, '11111' AS `orderNumber`, `comment` AS `requirements`, `order_status` AS `orderStatus`, `created_at` AS `createdAt`, `updated_at` AS `updatedAt`, `partner_store_id` AS `partnerStore` FROM `%s` WHERE `id` = ? LIMIT 1"
	var data types.GetUserOrderRep
	if err = l.svcCtx.DBC.Get(&data, fmt.Sprintf(query, tables.UserOrder), req.Id); err != nil {
		return nil, errcode.DatabaseError.SetDetails(err.Error())
	}
	return &data, nil
}
