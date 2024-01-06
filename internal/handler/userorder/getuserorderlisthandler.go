package userorder

import (
	"net/http"

	"carservice/internal/logic/userorder"
	"carservice/internal/pkg/common/errcode"
	stdresponse "carservice/internal/pkg/httper/response"
	"carservice/internal/svc"
	"carservice/internal/types"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func GetUserOrderListHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.GetUserOrderListReq
		if err := httpx.Parse(r, &req); err != nil {
			stdresponse.ResponseWithCtx(r.Context(), w, nil, errcode.InvalidParamsError.Lazy(err.Error()))
			return
		}

		l := userorder.NewGetUserOrderListLogic(r.Context(), svcCtx)
		resp, err := l.GetUserOrderList(&req)
		stdresponse.ResponseWithCtx(r.Context(), w, resp, err)
	}
}
