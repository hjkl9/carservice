package bulletin

import (
	"net/http"

	"carservice/internal/logic/bulletin"
	"carservice/internal/pkg/common/errcode"
	"carservice/internal/pkg/httper/api"
	"carservice/internal/svc"
	"carservice/internal/types"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func GetBulletinListHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.GetBulletinListReq
		if err := httpx.Parse(r, &req); err != nil {
			api.ResponseWithCtx(r.Context(), w, nil, errcode.InvalidParametersErr)
			return
		}

		l := bulletin.NewGetBulletinListLogic(r.Context(), svcCtx)
		resp, err := l.GetBulletinList(&req)
		api.ResponseWithCtx(r.Context(), w, resp, err)
	}
}
