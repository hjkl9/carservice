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
	Id           uint           `db:"id"`
	OrderNumber  string         `db:"orderNumber"`
	PartnerStore sql.NullString `db:"partnerStore"`
	Comment      string         `db:"comment"`
	OrderStatus  uint8          `db:"orderStatus"`
	CreatedAt    string         `db:"createdAt"`
	UpdatedAt    string         `db:"updatedAt"`
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
	query = "SELECT `o`.`id` AS `id`, `o`.`order_number` AS `orderNumber`, `ps`.`title` AS `partnerStore`, `o`.`comment` AS `comment`, `o`.`order_status` AS `orderStatus`, `o`.`created_at` AS `createdAt`, `o`.`updated_at` AS `updatedAt` FROM `%s` AS `o` LEFT JOIN `partner_stores` AS `ps` ON `ps`.`id` = `o`.`partner_store_id` WHERE `o`.`id` = ? AND `o`.`member_id` = ? LIMIT 1"
	if err = l.svcCtx.DBC.Get(&order, fmt.Sprintf(query, tables.UserOrder), req.Id, userId); err != nil {
		return nil, errcode.DatabaseError.SetDetails(err.Error())
	}
	fmt.Printf("%#v\n", order)
	return &types.GetUserOrderRep{
		Id:          order.Id,
		OrderNumber: order.OrderNumber,
		PartnerStore: func() string {
			if !order.PartnerStore.Valid {
				return "未绑定商家"
			}
			return order.PartnerStore.String
		}(),
		Requirements: order.Comment,
		OrderStatus:  userorder.OrderStatusDesc(order.OrderStatus),
		CreatedAt:    order.CreatedAt,
		UpdatedAt:    order.UpdatedAt,
	}, nil
}
