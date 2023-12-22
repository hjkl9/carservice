package main

import (
	"flag"
	"fmt"
	"net/http"

	"carservice/internal/config"
	"carservice/internal/handler"
	"carservice/internal/pkg/jwt"
	"carservice/internal/svc"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/logc"
	"github.com/zeromicro/go-zero/rest"
)

var configFile = flag.String("f", "etc/carservice.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)
	// setup logger
	logc.MustSetup(c.LogConf)

	server := rest.MustNewServer(
		c.RestConf,
		rest.WithUnauthorizedCallback(jwt.UnauthorizedCallback()),
	)
	defer server.Stop()

	ctx := svc.NewServiceContext(c)
	// register global middlewares.
	handler.RegisterGlobalMiddleware(server, ctx)
	handler.RegisterHandlers(server, ctx)

	fmt.Printf("Starting server at %s:%d...\n", c.Host, c.Port)
	server.PrintRoutes()

	server.Start()
}

func NotFoundHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
	}
}
