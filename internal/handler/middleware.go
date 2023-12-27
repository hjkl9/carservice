package handler

import (
	"carservice/internal/pkg/common/errcode"
	stdresponse "carservice/internal/pkg/httper/response"
	"carservice/internal/svc"
	"net/http"

	"github.com/zeromicro/go-zero/core/logc"
	"github.com/zeromicro/go-zero/rest"
)

func RegisterGlobalMiddleware(server *rest.Server, serverCtx *svc.ServiceContext) {
	server.Use(func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			err := serverCtx.DBC.Ping()
			if err != nil {
				logc.Errorf(r.Context(), "MySQL 连接发生错误或配置不正确, err: %#v\n", err)
				stdresponse.ResponseWithCtx(r.Context(), w, errcode.New(http.StatusInternalServerError, "-", "MySQL 连接发生错误或配置不正确"))
				return
			}
			next(w, r)
		}
	})
}
