// NON-USE.

package logic

import (
	"context"
	"net/http"

	"carservice/internal/pkg/common/errcode"
	"carservice/internal/svc"
	"carservice/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type PingLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewPingLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PingLogic {
	return &PingLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *PingLogic) Ping(req *types.ServerPingReq) (resp *types.ServerPingRep, err error) {
	switch req.HttpCode {
	case http.StatusOK:
		return &types.ServerPingRep{
			Result: "请求成功",
		}, nil
	case http.StatusInternalServerError:
		return nil, errcode.New(http.StatusInternalServerError, "呜呜呜", "No ok")
	case http.StatusBadRequest:
		return &types.ServerPingRep{}, errcode.New(http.StatusInternalServerError, "呜呜呜", "No ok")
	default:
		return &types.ServerPingRep{
			Result: "",
		}, errcode.New(http.StatusAccepted, "请求已接受", "莫名其妙")
	}
}
