package userorder

import (
	"context"

	"carservice/internal/enum/userorder"
	"carservice/internal/pkg/common/errcode"
	"carservice/internal/pkg/jwt"
	"carservice/internal/svc"
	"carservice/internal/types"

	"github.com/zeromicro/go-zero/core/logc"
	"github.com/zeromicro/go-zero/core/logx"
)

type ConfirmUserOrderLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewConfirmUserOrderLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ConfirmUserOrderLogic {
	return &ConfirmUserOrderLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// ConfirmUserOrder
//
// 用户确认订单
func (l *ConfirmUserOrderLogic) ConfirmUserOrder(req *types.ConfirmUserOrderReq) error {
	userId := jwt.GetUserId(l.ctx)

	hasOrder, err := l.svcCtx.Repo.
		UserOrderRelated().
		GetIfOrderExistsById(l.ctx, userId, uint(req.Id))
	if err != nil {
		logc.Errorf(l.ctx, "检查订单是否存在时发生错误, err: %s\n", err.Error())
		return errcode.DatabaseGetErr
	}
	// 用户订单不存在
	if !hasOrder {
		return errcode.OrderNotFoundErr
	}

	// 获取当前订单状态
	var orderStatus uint8
	query := "SELECT `order_status` `orderStatus` FROM `user_orders` WHERE `id` = ? LIMIT 1"
	stmt, err := l.svcCtx.DBC.PreparexContext(l.ctx, query)
	if err != nil {
		return errcode.DatabasePrepareErr
	}
	if err = stmt.GetContext(l.ctx, &orderStatus, req.Id); err != nil {
		return errcode.DatabaseGetErr
	}

	// 是否符合更改状态
	if orderStatus == userorder.ToBePaid {
		return errcode.DuplicateConfirmedOrderErr
	}

	// 确认订单到待支付状态
	query = "UPDATE `user_orders` SET `order_status` = ? WHERE `id` = ?"
	stmt, err = l.svcCtx.DBC.PreparexContext(l.ctx, query)
	if err != nil {
		logc.Errorf(l.ctx, "预处理更改订单状态时发生错误, err: %s\n", err.Error())
		return errcode.DatabasePrepareErr
	}
	rs, err := stmt.ExecContext(l.ctx, userorder.ToBePaid, req.Id)
	if err != nil {
		logc.Errorf(l.ctx, "更改订单状态时发生错误 1, err: %s\n", err.Error())
		return errcode.DatabaseUpdateErr
	}
	n, err := rs.RowsAffected()
	if err != nil {
		logc.Errorf(l.ctx, "更改订单状态时发生错误 2, err: %s\n", err.Error())
		return errcode.DatabaseUpdateErr
	}
	if n != 1 {
		logc.Errorf(l.ctx, "更改订单状态时发生错误 3, err: %s\n", err.Error())
		return errcode.DatabaseUpdateErr
	}
	return nil
}
