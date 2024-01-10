package userorder

import (
	"net/http"

	"carservice/internal/logic/userorder"
	"carservice/internal/pkg/httper/api"
	"carservice/internal/svc"
	"carservice/internal/types"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func ConfirmAndPayUserOrderHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.ConfirmAndPayUserOrderReq
		if err := httpx.Parse(r, &req); err != nil {
			api.ResponseWithCtx(r.Context(), w, nil, err)
			return
		}

		l := userorder.NewConfirmAndPayUserOrderLogic(r.Context(), svcCtx)
		err := l.ConfirmAndPayUserOrder(&req)
		api.ResponseWithCtx(r.Context(), w, nil, err)
	}
}
