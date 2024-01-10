package userorder

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"carservice/internal/enum/userorder"
	uo_enum "carservice/internal/enum/userorder"
	"carservice/internal/pkg/common/errcode"
	"carservice/internal/pkg/conv"
	"carservice/internal/pkg/jwt"
	"carservice/internal/svc"
	"carservice/internal/types"

	"github.com/zeromicro/go-zero/core/logc"
	"github.com/zeromicro/go-zero/core/logx"
)

type GetUserOrderListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetUserOrderListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserOrderListLogic {
	return &GetUserOrderListLogic{
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
	CreatedAt    time.Time      `db:"createdAt" json:"createdAt"`
	UpdatedAt    time.Time      `db:"updatedAt" json:"updatedAt"`
}

func (l *GetUserOrderListLogic) GetUserOrderList(req *types.GetUserOrderListReq) (resp []types.UserOrderListItem, err error) {
	// handler status of order list as string query.
	statusSubQuery := l.handleStatusSubQuery(req.Status)
	// 用户
	userId := jwt.GetUserId(l.ctx)
	// 先查询是否存在订单
	var hasList uint8
	query := "SELECT (count(1) > 0) AS `hasList` FROM `user_orders` WHERE `member_id` = ? AND `deleted_at` IS NULL"
	stmt, err := l.svcCtx.DBC.PreparexContext(l.ctx, query)
	if err != nil {
		logc.Errorf(l.ctx, "查询订单列表语句预处理发生错误, err: %s\n", err.Error())
		return []types.UserOrderListItem{}, errcode.DatabasePrepareErr
	}
	if err = stmt.GetContext(l.ctx, &hasList, userId); err != nil {
		logc.Errorf(l.ctx, "查询订单获取数据发生错误, err: %s\n", err.Error())
		return []types.UserOrderListItem{}, errcode.DatabaseGetErr
	}
	if hasList == 0 {
		return []types.UserOrderListItem{}, nil
	}
	// 是否可删除条件函数
	deletable := func(status uint8) bool { return status == uo_enum.Cancelled || status == uo_enum.Completed }
	// 待导出数据
	var orders []*OrderListItem
	// 查询语句
	query = "SELECT `uo`.`id`, `uo`.`order_number` AS `orderNumber`, `ps`.`title` AS `partnerStore`, `uo`.`comment` AS `requirements`, `uo`.`order_status` AS `orderStatus`, `uo`.`created_at` AS `createdAt`, `uo`.`updated_at` AS `updatedAt` FROM `user_orders` AS `uo` LEFT JOIN `partner_stores` AS `ps` ON `uo`.`partner_store_id` = `ps`.`id` WHERE 1=1 " + statusSubQuery + " AND `uo`.`member_id` = ? AND `uo`.`deleted_at` IS NULL"
	// 开始查询并扫描数据到变量 orders
	stmt, err = l.svcCtx.DBC.PreparexContext(l.ctx, query)
	if err != nil {
		return nil, errcode.DatabasePrepareErr
	}
	if err = stmt.SelectContext(
		l.ctx,
		&orders,
		// 匹配表名
		userId,
	); err != nil {
		fmt.Printf("here err: %s\n", err.Error())
		// 处理并抛出查询发生的错误
		return nil, errcode.DatabaseGetErr
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
			Deletable:   deletable((*v).OrderStatus),
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
			CreatedAt:    (*v).CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt:    (*v).UpdatedAt.Format("2006-01-02 15:04:05"),
		})
	}
	return data, nil
}

func (l *GetUserOrderListLogic) handleStatusSubQuery(status string) string {
	if status == "" {
		return ""
	}
	var in string
	// ? 在编译的时候确定
	switch status {
	case "0":
		// 全部数据
		return ""
	case "1":
		// 待处理和待确认和付款
		in = conv.ToStringWithSep([]uint8{
			userorder.Pending,
			userorder.ToBeConfirmedAndPay,
		}, ',')
	case "2":
		// 已付款和待安装
		in = conv.ToStringWithSep([]uint8{
			userorder.Paid,
			userorder.PrepareToInstall,
		}, ',')
	case "3":
		// 已完成
		in = conv.ToStringWithSep([]uint8{
			userorder.Completed,
		}, ',')
	case "4":
		// 已关闭
		in = conv.ToStringWithSep([]uint8{
			userorder.Closed,
		}, ',')
	}
	return fmt.Sprintf("AND `order_status` IN (%s)", in)
}
