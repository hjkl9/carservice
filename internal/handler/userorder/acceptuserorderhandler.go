package userorder

import (
	"net/http"

	"carservice/internal/logic/userorder"
	"carservice/internal/pkg/httper/api"
	"carservice/internal/svc"
	"carservice/internal/types"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func AcceptUserOrderHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.AcceptUserOrderReq
		if err := httpx.Parse(r, &req); err != nil {
			api.ResponseWithCtx(r.Context(), w, nil, err)
			return
		}

		l := userorder.NewAcceptUserOrderLogic(r.Context(), svcCtx)
		err := l.AcceptUserOrder(&req)
		api.ResponseWithCtx(r.Context(), w, nil, err)
	}
}
