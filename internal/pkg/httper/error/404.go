package error

import (
	"carservice/internal/pkg/common/errcode"
	stdresponse "carservice/internal/pkg/httper/response"
	"net/http"
)

func NotFoundHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		stdresponse.Response(w, nil, errcode.NotFound.SetMsg("找不到该路由"))
	}
}
