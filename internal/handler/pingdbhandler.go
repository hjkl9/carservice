package handler

import (
	"net/http"

	"carservice/internal/logic"
	"carservice/internal/pkg/common/errcode"
	stdresponse "carservice/internal/pkg/httper/response"
	"carservice/internal/svc"
	"carservice/internal/types"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func PingDbHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.PingDbReq
		if err := httpx.Parse(r, &req); err != nil {
			stdresponse.ResponseWithCtx(r.Context(), w, errcode.New(http.StatusBadRequest, "feature.", err.Error()))
			return
		}

		l := logic.NewPingDbLogic(r.Context(), svcCtx)
		resp, err := l.PingDb(&req)
		stdresponse.Response(w, resp, err)
	}
}
