package carbrand

import (
	"net/http"

	"carservice/internal/logic/carbrand"
	"carservice/internal/pkg/common/errcode"
	stdresponse "carservice/internal/pkg/httper/response"
	"carservice/internal/svc"
	"carservice/internal/types"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func BrandOptionListHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.CarBrandOptionListReq
		if err := httpx.Parse(r, &req); err != nil {
			stdresponse.ResponseWithCtx(r.Context(), w, errcode.New(http.StatusBadRequest, "feature.", err.Error()))
			return
		}
		l := carbrand.NewBrandOptionListLogic(r.Context(), svcCtx)
		resp, err := l.BrandOptionList()
		stdresponse.Response(w, resp, err)
	}
}
