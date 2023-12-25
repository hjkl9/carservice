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

// DeleteUserOrderHandler 删除用户订单
func DeleteUserOrderHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.DeleteUserOrderReq
		if err := httpx.Parse(r, &req); err != nil {
			stdresponse.ResponseWithCtx(r.Context(), w, errcode.InvalidParamsError.Lazy(err.Error()))
			return
		}

		l := userorder.NewDeleteUserOrderLogic(r.Context(), svcCtx)
		err := l.DeleteUserOrder(&req)
		stdresponse.Response(w, nil, err)
	}
}
