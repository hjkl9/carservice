package userorder

import (
	"net/http"

	"carservice/internal/logic/userorder"
	"carservice/internal/pkg/httper/api"
	"carservice/internal/svc"
	"carservice/internal/types"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func GetUserOrderListHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.GetUserOrderListReq
		if err := httpx.Parse(r, &req); err != nil {
			api.ResponseWithCtx(r.Context(), w, nil, err)
			return
		}

		l := userorder.NewGetUserOrderListLogic(r.Context(), svcCtx)
		resp, err := l.GetUserOrderList(&req)
		_ = err
		api.ResponseWithCtx(r.Context(), w, resp, err)
		// stdresponse.ResponseWithCtx(r.Context(), w, resp, err)
	}
}
