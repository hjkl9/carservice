package logic

import (
	"context"
	"fmt"
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

func (l *PingDbLogic) PingDb(req *types.PingDbReq) (resp *types.PingDbRep, err error) {
	err = l.svcCtx.DBC.Ping()
	if err != nil {
		return nil, errcode.New(http.StatusInternalServerError, "feature.", err.Error())
	}
	fmt.Println(req.AsResult)
	query := "SELECT ? AS `result`"
	var result string
	l.svcCtx.DBC.Get(&result, query, req.AsResult)
	return &types.PingDbRep{Result: result}, nil
}
