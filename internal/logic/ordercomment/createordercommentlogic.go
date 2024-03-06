package ordercomment

import (
	"context"

	"carservice/internal/pkg/common/errcode"
	"carservice/internal/pkg/jwt"
	"carservice/internal/svc"
	"carservice/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateOrderCommentLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreateOrderCommentLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateOrderCommentLogic {
	return &CreateOrderCommentLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateOrderCommentLogic) CreateOrderComment(req *types.CreateOrderCommentReq) error {
	query := "INSERT INTO `order_comments` (`user_order_id`, `user_id`, `title`, `rate`, `content`, `client`, `published_at`, `created_at`, `updated_at`) VALUES(?, ?, ?, ?, ?, ?, NOW(), NOW(), NOW());"

	stmt, err := l.svcCtx.DBC.PrepareContext(l.ctx, query)
	if err != nil {
		return errcode.DatabasePrepareErr
	}

	userOrderId := req.UserOrderId
	userId := jwt.GetUserId(l.ctx)
	rate := req.Rate
	content := req.Content
	client := 1

	var title = ""
	if len(req.Title) == 0 {
		title = ""
	} else {
		title = req.Title
	}

	res, err := stmt.ExecContext(l.ctx, userOrderId, userId, title, rate, content, client)
	if err != nil {
		return errcode.DatabaseExecuteErr
	}
	newId, err := res.LastInsertId()
	if err != nil {
		return err
	}

	_ = newId

	return nil
}
