package user

import (
	"context"

	"carservice/internal/svc"
	"carservice/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type WechatAuthorizationLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewWechatAuthorizationLogic(ctx context.Context, svcCtx *svc.ServiceContext) *WechatAuthorizationLogic {
	return &WechatAuthorizationLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *WechatAuthorizationLogic) WechatAuthorization(req *types.WechatAuthorizationReq) (resp *types.WechatAuthorizationRep, err error) {
	// todo: add your logic here and delete this line

	return
}
