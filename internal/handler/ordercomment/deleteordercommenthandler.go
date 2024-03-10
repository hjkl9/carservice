package ordercomment

import (
	"net/http"

	"carservice/internal/logic/ordercomment"
	"carservice/internal/pkg/common/errcode"
	"carservice/internal/pkg/httper/api"
	"carservice/internal/svc"
	"carservice/internal/types"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func DeleteOrderCommentHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.DeleteOrderCommentReq
		if err := httpx.Parse(r, &req); err != nil {
			api.ResponseWithCtx(r.Context(), w, nil, errcode.InvalidParametersErr)
			return
		}
		l := ordercomment.NewDeleteOrderCommentLogic(r.Context(), svcCtx)
		err := l.DeleteOrderComment(&req)
		api.ResponseWithCtx(r.Context(), w, nil, err)
	}
}
