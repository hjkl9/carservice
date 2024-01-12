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

type CancelUserOrderLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCancelUserOrderLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CancelUserOrderLogic {
	return &CancelUserOrderLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CancelUserOrderLogic) CancelUserOrder(req *types.CancelUserOrderReq) error {
	userId := jwt.GetUserId(l.ctx)

	// 订单是否存在
	hasOrder, err := l.svcCtx.Repo.UserOrderRelated().GetIfOrderExistsById(l.ctx, userId, uint(req.Id))
	if err != nil {
		logc.Error(l.ctx, "查询订单是否存在时发生错误")
	}
	if !hasOrder {
		return errcode.NotFound.Lazy("该订单不存在")
	}

	// 取消订单的条件
	// userorder.ToBePaid
	// userorder.Pending
	// userorder.ToBeAcceptedByUser

	// CheckIfOrderCanBeCanceled
	// 检查是否满足取消订单的条件
	var cancellable uint8
	query := "SELECT (COUNT(1) = 1) AS `cancellable` FROM `user_orders` WHERE 1 = 1 AND `order_status` IN (?, ?) AND `id` = ? AND `member_id` = ? LIMIT 1"
	stmt, err := l.svcCtx.DBC.PreparexContext(l.ctx, query)
	if err != nil {
		logc.Error(l.ctx, "预处理查询是否满足取消订单时发生错误")
		return errcode.NewDatabaseErrorx().GetError(err)
	}
	if err = stmt.GetContext(
		l.ctx,
		&cancellable,
		userorder.Pending,
		userorder.AwaitingPayment,
		req.Id,
		userId,
	); err != nil {
		logc.Error(l.ctx, "开始查询是否满足取消订单时发生错误")
		return errcode.NewDatabaseErrorx().GetError(err)
	}
	if cancellable != 1 {
		return errcode.StatusForbiddenError.SetMsg("该订单无法被取消")
	}
	// 更新订单状态
	query = "UPDATE `user_orders` SET `order_status` = ?, `updated_at` = NOW() WHERE `id` = ? AND `member_id` = ?"
	stmt, err = l.svcCtx.DBC.PreparexContext(l.ctx, query)
	if err != nil {
		logc.Error(l.ctx, "预处理更新状态的 SQL 语句时发生错误")
		return errcode.NewDatabaseErrorx().UpdateError(err)
	}
	_, err = stmt.ExecContext(l.ctx, userorder.Cancelled, req.Id, userId)
	if err != nil {
		logc.Error(l.ctx, "执行更新状态的 SQL 语句时发生错误")
		return errcode.NewDatabaseErrorx().UpdateError(err)
	}
	return nil
}
