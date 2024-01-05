package common

import (
	"context"

	"carservice/internal/svc"
	"carservice/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type WebsocketServicesLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewWebsocketServicesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *WebsocketServicesLogic {
	return &WebsocketServicesLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *WebsocketServicesLogic) WebsocketServices() (resp []types.WebsocketServiceListItem, err error) {
	// todo: add your logic here and delete this line

	return
}
