package userorder

import (
	"context"

	"carservice/internal/pkg/common/errcode"
	"carservice/internal/pkg/jwt"
	"carservice/internal/svc"
	"carservice/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteUserOrderLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDeleteUserOrderLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteUserOrderLogic {
	return &DeleteUserOrderLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// DeleteUserOrder
// todo: 被删除的条件:
// - 条件 1 (TODO)
// - 条件 2 (TODO)
func (l *DeleteUserOrderLogic) DeleteUserOrder(req *types.DeleteUserOrderReq) error {
	// 用户
	userId := jwt.GetUserId(l.ctx)
	// 检查用户订单是否存在
	var count int
	query := "SELECT COUNT(1) AS `count` FROM `user_orders` WHERE `member_id` = ? AND `id` = ? LIMIT 1"
	if err := l.svcCtx.DBC.Get(&count, query, userId, req.Id); err != nil {
		return errcode.NewDatabaseErrorx().GetError(err)
	}
	// 订单不存在
	if count == 0 {
		return errcode.NotFound.SetMsg("该用户订单不存在或已被删除")
	}
	// todo: 删除的条件
	// - ...

	// 软删除
	query = "UPDATE `user_orders` SET `deleted_at` = NOW() WHERE `member_id` = ? AND `id` = ?"
	rs, err := l.svcCtx.DBC.ExecContext(l.ctx, query, userId, req.Id)
	if err != nil {
		return errcode.NewDatabaseErrorx().DeleteError(err)
	}
	n, _ := rs.RowsAffected()
	if n == 0 {
		return errcode.NewDatabaseErrorx().DeleteError(err)
	}
	return nil
}
