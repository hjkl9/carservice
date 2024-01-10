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

func (l *ConfirmUserOrderLogic) ConfirmUserOrder(req *types.AcceptUserOrderReq) error {
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
	// FIXME: 当前任务是同意订单报价
	// 是否符合更改状态
	if orderStatus == userorder.Cancelled {
		return errcode.OrderCannotBeCancelledErr.SetMessage("该用户订单已经被取消")
	}
	var cancellabe = func() bool {
		switch orderStatus {
		case userorder.Pending:
		case userorder.ToBeConfirmed:
		case userorder.ToBePaid:
			return true
		default:
			return false
		}
		return false
	}()
	if !cancellabe {
		return errcode.OrderCannotBeCancelledErr.SetMessage("该用户订单不满足取消条件")
	}

	// 取消订单
	query = "UPDATE `user_orders` SET `order_status` = ? WHERE `id` = ?"
	stmt, err = l.svcCtx.DBC.PreparexContext(l.ctx, query)
	if err != nil {
		logc.Errorf(l.ctx, "预处理更改订单状态时发生错误, err: %s\n", err.Error())
		return errcode.DatabasePrepareErr
	}
	rs, err := stmt.ExecContext(l.ctx, userorder.Cancelled, req.Id)
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
