package userorder

import (
	"net/http"

	"carservice/internal/logic/userorder"
	"carservice/internal/svc"
	"carservice/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func PaymentOrderCallbackHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.PaymentOrderCallbackReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := userorder.NewPaymentOrderCallbackLogic(r.Context(), svcCtx)
		err := l.PaymentOrderCallback(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.Ok(w)
		}
	}
}
