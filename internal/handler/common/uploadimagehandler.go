package common

import (
	"net/http"

	"carservice/internal/logic/common"
	"carservice/internal/pkg/common/errcode"
	stdresponse "carservice/internal/pkg/httper/response"
	"carservice/internal/svc"
)

func UploadImageHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		file, fileHeader, err := r.FormFile("file")
		if err != nil {
			if err == http.ErrMissingFile {
				stdresponse.ResponseWithCtx(r.Context(), w, nil, errcode.New(http.StatusBadRequest, "feature.", err.Error()))
				return
			}
			stdresponse.ResponseWithCtx(r.Context(), w, nil, errcode.New(http.StatusAccepted, "feature.", err.Error()))
			return
		}
		defer file.Close()
		l := common.NewUploadImageLogic(r.Context(), svcCtx)
		resp, err := l.UploadImage(file, fileHeader)
		stdresponse.Response(w, resp, err)
	}
}
