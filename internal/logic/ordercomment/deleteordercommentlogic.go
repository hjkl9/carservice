package ordercomment

import (
	"context"

	"carservice/internal/pkg/common/errcode"
	"carservice/internal/svc"
	"carservice/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteOrderCommentLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDeleteOrderCommentLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteOrderCommentLogic {
	return &DeleteOrderCommentLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteOrderCommentLogic) DeleteOrderComment(req *types.DeleteOrderCommentReq) error {
	// 检查订单评价是否存在
	var count int8
	query := "SELECT COUNT(`id`) AS `count` FROM `order_comment` WHERE `id` = ? LIMIT 1"

	err := l.svcCtx.DBC.GetContext(l.ctx, &count, query, req.Id)
	if err != nil {
		return errcode.DatabaseGetErr
	}
	if count == 0 {
		return errcode.OrderCommentNotFoundErr
	}

	// 开始删除评价
	query = "DELETE FROM `order_comment` WHERE `id` = ? LIMIT 1"
	_, err = l.svcCtx.DBC.ExecContext(l.ctx, query, req.Id)
	if err != nil {
		return errcode.DatabaseDeleteErr
	}

	return nil
}
