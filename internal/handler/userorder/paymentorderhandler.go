package userorder

import (
	"fmt"
	"net/http"

	"carservice/internal/logic/userorder"
	"carservice/internal/pkg/common/errcode"
	"carservice/internal/pkg/httper/api"
	"carservice/internal/svc"
	"carservice/internal/types"

	"github.com/zeromicro/go-zero/rest/httpx"
)

// PaymentOrderHandler 处理订单支付
func PaymentOrderHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.PaymentOrderReq
		if err := httpx.Parse(r, &req); err != nil {
			fmt.Println(err.Error())
			api.ResponseWithCtx(r.Context(), w, nil, errcode.InvalidParametersErr)
			return
		}

		l := userorder.NewPaymentOrderLogic(r.Context(), svcCtx)
		resp, err := l.PaymentOrder(&req)
		api.ResponseWithCtx(r.Context(), w, resp, err)
	}
}
