package userorder

import (
	"context"
	"database/sql"
	"fmt"

	"carservice/internal/data/tables"
	uo_enum "carservice/internal/enum/userorder"
	"carservice/internal/pkg/common/errcode"
	"carservice/internal/pkg/jwt"
	"carservice/internal/svc"
	"carservice/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetrUserOrderListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetrUserOrderListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetrUserOrderListLogic {
	return &GetrUserOrderListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// todo: 应该使用单独文件存放
type OrderListItem struct {
	Id           uint           `db:"id" json:"id"`
	OrderNumber  string         `db:"orderNumber" json:"orderNumber"`
	PartnerStore sql.NullString `db:"partnerStore" json:"partnerStore"`
	Requirements string         `db:"requirements" json:"requirements"`
	OrderStatus  uint8          `db:"orderStatus" json:"orderStatus"`
	CreatedAt    string         `db:"createdAt" json:"createdAt"`
	UpdatedAt    string         `db:"updatedAt" json:"updatedAt"`
}

func (l *GetrUserOrderListLogic) GetrUserOrderList(req *types.GetUserOrderListReq) (resp []types.UserOrderListItem, err error) {
	// 用户
	userId := jwt.GetUserId(l.ctx)
	// 待导出数据
	var orders []*OrderListItem
	// 查询语句
	query := "SELECT `uo`.`id`, `uo`.`order_number` AS `orderNumber`, `ps`.`title` AS `partnerStore`, `uo`.`comment` AS `requirements`, `uo`.`order_status` AS `orderStatus`, `uo`.`created_at` AS `createdAt`, `uo`.`updated_at` AS `updatedAt` FROM `%s` AS `uo` LEFT JOIN `%s` AS `ps` ON `uo`.`partner_store_id` = `ps`.`id` WHERE `uo`.`member_id` = ? AND `uo`.`deleted_at` IS NULL"
	// 开始查询并扫描数据到变量 orders
	if err = l.svcCtx.DBC.SelectContext(
		l.ctx,
		&orders,
		// 匹配表名
		fmt.Sprintf(query, tables.UserOrder, tables.PartnerStore),
		userId,
	); err != nil {
		// 处理并抛出查询发生的错误
		return nil, errcode.NewDatabaseErrorx().GetError(err)
	}
	// 面向客户端的数据切片 data
	var data []types.UserOrderListItem
	// 处理空切片
	if len(orders) == 0 {
		return make([]types.UserOrderListItem, 0), nil
	}
	// 遍历导出数据到 data
	for _, v := range orders {
		data = append(data, types.UserOrderListItem{
			Id:          (*v).Id,
			OrderNumber: (*v).OrderNumber,
			PartnerStore: func() string {
				// 处理如果是 Nil 的字符串
				if (*v).PartnerStore.Valid {
					return (*v).PartnerStore.String
				}
				return "未绑定合作门店"
			}(),
			Requirements: (*v).Requirements,
			OrderStatus:  uo_enum.OrderStatusDesc((*v).OrderStatus),
			CreatedAt:    (*v).CreatedAt,
			UpdatedAt:    (*v).UpdatedAt,
		})
	}
	return data, nil
}
