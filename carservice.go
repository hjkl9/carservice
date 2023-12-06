package main

import (
	"flag"
	"fmt"
	"net/http"

	"carservice/internal/config"
	"carservice/internal/handler"
	"carservice/internal/pkg/common/errcode"
	stdresponse "carservice/internal/pkg/httper/response"
	"carservice/internal/svc"

	"github.com/jmoiron/sqlx"
	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/rest"
)

var configFile = flag.String("f", "etc/carservice.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)

	server := rest.MustNewServer(c.RestConf)
	defer server.Stop()

	ctx := svc.NewServiceContext(c)
	checkDB := checkMysqlMiddleware(ctx.DBC)
	server.Use(checkDB)
	handler.RegisterHandlers(server, ctx)

	fmt.Printf("Starting server at %s:%d...\n", c.Host, c.Port)
	server.Start()
}

// check mysql connection middleware.
// It's possible to delete this middleware.
func checkMysqlMiddleware(db *sqlx.DB) rest.Middleware {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			err := db.Ping()
			if err != nil {
				stdresponse.ResponseWithCtx(r.Context(), w, errcode.New(http.StatusOK, "-", "MySQL 连接发生错误或配置不正确"))
				return
			}
			next(w, r)
		}
	}
}
