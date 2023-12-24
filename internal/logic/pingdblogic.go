package logic

import (
	"context"
	"net/http"

	"carservice/internal/pkg/common/errcode"
	"carservice/internal/svc"
	"carservice/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type PingDbLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewPingDbLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PingDbLogic {
	return &PingDbLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *PingDbLogic) PingDb(req *types.ServerPingDbReq) (resp *types.ServerPingDbRep, err error) {
	err = l.svcCtx.DBC.Ping()
	if err != nil {
		return nil, errcode.New(http.StatusInternalServerError, "feature.", err.Error())
	}
	result, err := l.svcCtx.Repo.PingRelated().EchoAsResult(req.AsResult)
	if err != nil {
		return nil, errcode.New(http.StatusNotFound, "-", "找不到记录")
	}
	return &types.ServerPingDbRep{Result: result}, nil
}
