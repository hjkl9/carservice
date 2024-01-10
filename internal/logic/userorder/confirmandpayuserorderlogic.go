package userorder

import (
	"context"

	"carservice/internal/enum/userorder"
	"carservice/internal/pkg/common/errcode"
	"carservice/internal/pkg/jwt"
	"carservice/internal/pkg/wechat/payment"
	"carservice/internal/svc"
	"carservice/internal/types"

	"github.com/zeromicro/go-zero/core/logc"
	"github.com/zeromicro/go-zero/core/logx"
)

type ConfirmAndPayUserOrderLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewConfirmAndPayUserOrderLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ConfirmAndPayUserOrderLogic {
	return &ConfirmAndPayUserOrderLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// ConfirmUserOrder
//
// 用户确认订单
func (l *ConfirmAndPayUserOrderLogic) ConfirmAndPayUserOrder(req *types.ConfirmAndPayUserOrderReq) error {
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
	if orderStatus != userorder.Pending {
		return errcode.ConfirmedOrderErr.SetMessage("不符合确认和支付用户订单的条件")
	}

	var order struct {
		OrderNumber string  `db:"orderNumber"`
		Amount      float64 `db:"amount"`
		OpenId      string  `db:"openid"` // todo: relation query
	}
	query = "SELECT `user_orders`.`order_number` `orderNumber`, `user_orders`.`act_amount` `amount`, `user_binds`.`open_id` `openid` FROM `user_orders` JOIN `user_binds`.`user_id` = `user_orders`.`member_id` WHERE `user_orders`.`id` = ? LIMIT 1"
	stmt, err = l.svcCtx.DBC.PreparexContext(l.ctx, query)
	if err != nil {
		return errcode.DatabasePrepareErr.WithDetails(err.Error())
	}
	if err = stmt.GetContext(l.ctx, &order, req.Id); err != nil {
		return errcode.DatabaseGetErr.WithDetails(err.Error())
	}

	// 开始支付
	payload := payment.PaymentPayload{
		Description: "",
		OutTradeNo:  order.OrderNumber,
		Attach:      "",
		NotifyUrl:   "",
		Amount:      int64(order.Amount),
		OpenId:      order.OpenId,
	}
	payment.JsApiOrder(payment.PaymentConfig{ /* TODO */ }, payload)

	// 确认订单到待安装
	query = "UPDATE `user_orders` SET `order_status` = ? WHERE `id` = ?"
	stmt, err = l.svcCtx.DBC.PreparexContext(l.ctx, query)
	if err != nil {
		logc.Errorf(l.ctx, "预处理更改订单状态时发生错误, err: %s\n", err.Error())
		return errcode.DatabasePrepareErr
	}
	rs, err := stmt.ExecContext(l.ctx, userorder.PrepareToInstall, req.Id)
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
