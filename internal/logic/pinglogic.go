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

func (l *PingLogic) Ping(req *types.PingReq) (resp *types.PingRep, err error) {
	switch req.HttpCode {
	case http.StatusOK:
		return &types.PingRep{
			Result: "请求成功",
		}, nil
	case http.StatusInternalServerError:
		return nil, errcode.New(http.StatusInternalServerError, "呜呜呜", "No ok")
	case http.StatusBadRequest:
		return &types.PingRep{}, errcode.New(http.StatusInternalServerError, "呜呜呜", "No ok")
	default:
		return &types.PingRep{
			Result: "",
		}, errcode.New(http.StatusAccepted, "请求已接受", "莫名其妙")
	}
}
