package carreplacement

import (
	"net/http"

	"carservice/internal/logic/carreplacement"
	"carservice/internal/pkg/httper/api"
	"carservice/internal/svc"
)

func ReplacementListHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := carreplacement.NewReplacementListLogic(r.Context(), svcCtx)
		resp, err := l.ReplacementList()
		api.ResponseWithCtx(r.Context(), w, resp, err)
	}
}
