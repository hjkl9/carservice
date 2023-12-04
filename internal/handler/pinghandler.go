package handler

import (
	"net/http"

	"carservice/internal/svc"
	"carservice/internal/types"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func PingHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.PingReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		// l := logic.NewPingLogic(r.Context(), svcCtx)
		// resp, err := l.Ping(&req)
		httpx.OkJsonCtx(r.Context(), w, "Pong")
	}
}
