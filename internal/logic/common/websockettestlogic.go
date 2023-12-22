package common

import (
	"context"

	"carservice/internal/svc"
	"carservice/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type WebsocketTestLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewWebsocketTestLogic(ctx context.Context, svcCtx *svc.ServiceContext) *WebsocketTestLogic {
	return &WebsocketTestLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *WebsocketTestLogic) WebsocketTest(req *types.WebsocketTestReq) error {
	return nil
}
