package user

import (
	"net/http"

	"carservice/internal/logic/user"
	"carservice/internal/pkg/httper/api"
	"carservice/internal/svc"
	"carservice/internal/types"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func WechatAuthorizationHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.WechatAuthorizationReq
		if err := httpx.Parse(r, &req); err != nil {
			api.ResponseWithCtx(r.Context(), w, nil, err)
			// stdresponse.ResponseWithCtx(r.Context(), w, nil, errcode.New(http.StatusBadRequest, "feature.", err.Error()))
			return
		}

		l := user.NewWechatAuthorizationLogic(r.Context(), svcCtx)
		resp, err := l.WechatAuthorization(&req)
		api.ResponseWithCtx(r.Context(), w, resp, err)
		// stdresponse.Response(w, resp, err)
	}
}
