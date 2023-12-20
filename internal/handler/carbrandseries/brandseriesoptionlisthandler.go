package carbrandseries

import (
	"net/http"

	"carservice/internal/logic/carbrandseries"
	"carservice/internal/pkg/common/errcode"
	stdresponse "carservice/internal/pkg/httper/response"
	"carservice/internal/svc"
	"carservice/internal/types"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func BrandSeriesOptionListHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.BrandSeriesOptionListReq
		if err := httpx.Parse(r, &req); err != nil {
			stdresponse.ResponseWithCtx(r.Context(), w, errcode.New(http.StatusBadRequest, "feature.", err.Error()))
			return
		}
		l := carbrandseries.NewBrandSeriesOptionListLogic(r.Context(), svcCtx)
		resp, err := l.BrandSeriesOptionList(&req)
		stdresponse.Response(w, resp, err)
	}
}
