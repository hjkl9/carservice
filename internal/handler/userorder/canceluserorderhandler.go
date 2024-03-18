package userorder

import (
	"net/http"

	"carservice/internal/logic/userorder"
	"carservice/internal/pkg/common/errcode"
	"carservice/internal/pkg/httper/api"
	"carservice/internal/svc"
	"carservice/internal/types"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func CancelUserOrderHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.CancelUserOrderReq
		if err := httpx.Parse(r, &req); err != nil {
			api.ResponseWithCtx(r.Context(), w, nil, errcode.InvalidParametersErr.SetMessage(err.Error()))
			return
		}

		l := userorder.NewCancelUserOrderLogic(r.Context(), svcCtx)
		err := l.CancelUserOrder(&req)
		api.ResponseWithCtx(r.Context(), w, nil, err)

	}
}
