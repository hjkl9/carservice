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

func UpdateUserOrderHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.UpdateUserOrderReq
		if err := httpx.Parse(r, &req); err != nil {
			api.ResponseWithCtx(r.Context(), w, nil, errcode.InvalidParametersErr)
			return
		}

		l := userorder.NewUpdateUserOrderLogic(r.Context(), svcCtx)
		err := l.UpdateUserOrder(&req)
		api.Response(w, nil, err)
	}
}
