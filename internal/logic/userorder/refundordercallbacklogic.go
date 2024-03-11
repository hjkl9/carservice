package userorder

import (
	"context"

	"carservice/internal/svc"
	"carservice/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type RefundOrderCallbackLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewRefundOrderCallbackLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RefundOrderCallbackLogic {
	return &RefundOrderCallbackLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RefundOrderCallbackLogic) RefundOrderCallback(req *types.RefundOrderCallbackReq) error {
	// todo: add your logic here and delete this line

	return nil
}
