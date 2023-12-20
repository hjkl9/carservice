// Code generated by goctl. DO NOT EDIT.
package handler

import (
	"net/http"

	carbrand "carservice/internal/handler/carbrand"
	carbrandseries "carservice/internal/handler/carbrandseries"
	carownerinfo "carservice/internal/handler/carownerinfo"
	common "carservice/internal/handler/common"
	sms "carservice/internal/handler/sms"
	user "carservice/internal/handler/user"
	"carservice/internal/svc"

	"github.com/zeromicro/go-zero/rest"
)

func RegisterHandlers(server *rest.Server, serverCtx *svc.ServiceContext) {
	server.AddRoutes(
		[]rest.Route{
			{
				Method:  http.MethodGet,
				Path:    "/ping",
				Handler: PingHandler(serverCtx),
			},
			{
				Method:  http.MethodGet,
				Path:    "/pingDb",
				Handler: PingDbHandler(serverCtx),
			},
		},
	)

	server.AddRoutes(
		[]rest.Route{
			{
				Method:  http.MethodGet,
				Path:    "/wechat/mp/getUserPhoneNumber",
				Handler: common.GetUserPhoneNumberHandler(serverCtx),
			},
			{
				Method:  http.MethodPost,
				Path:    "/fs/uploadFile",
				Handler: common.UploadFileHandler(serverCtx),
			},
			{
				Method:  http.MethodPost,
				Path:    "/fs/uploadMultipleFiles",
				Handler: common.UploadMultipleFilesHandler(serverCtx),
			},
		},
		rest.WithPrefix("/common"),
	)

	server.AddRoutes(
		[]rest.Route{
			{
				Method:  http.MethodGet,
				Path:    "/common/sms/sendCaptcha",
				Handler: sms.SendCaptchaHandler(serverCtx),
			},
		},
	)

	server.AddRoutes(
		[]rest.Route{
			{
				Method:  http.MethodPost,
				Path:    "/user/login/phoneNumber",
				Handler: user.PhoneNumberLoginHandler(serverCtx),
			},
			{
				Method:  http.MethodGet,
				Path:    "/user/login/wechatAuthorization",
				Handler: user.WechatAuthorizationHandler(serverCtx),
			},
		},
		rest.WithPrefix("/v1"),
	)

	server.AddRoutes(
		[]rest.Route{
			{
				Method:  http.MethodGet,
				Path:    "/user/getUserByPhoneNumber",
				Handler: user.GetUserByPhoneNumberHandler(serverCtx),
			},
		},
		rest.WithJwt(serverCtx.Config.JwtConf.AccessSecret),
		rest.WithPrefix("/v1"),
	)

	server.AddRoutes(
		[]rest.Route{
			{
				Method:  http.MethodGet,
				Path:    "/carBrand/brandOptionList",
				Handler: carbrand.BrandOptionListHandler(serverCtx),
			},
		},
		rest.WithPrefix("/v1"),
	)

	server.AddRoutes(
		[]rest.Route{
			{
				Method:  http.MethodGet,
				Path:    "/carBrandSeries/brandSeriesOptionList",
				Handler: carbrandseries.BrandSeriesOptionListHandler(serverCtx),
			},
		},
		rest.WithPrefix("/v1"),
	)

	server.AddRoutes(
		[]rest.Route{
			{
				Method:  http.MethodGet,
				Path:    "/carOwnerInfo/checkEmptyList",
				Handler: carownerinfo.CheckEmptyListHandler(serverCtx),
			},
			{
				Method:  http.MethodPost,
				Path:    "/carOwnerInfo",
				Handler: carownerinfo.CreateCarOwnerInfoHandler(serverCtx),
			},
			{
				Method:  http.MethodPut,
				Path:    "/carOwnerInfo/:id",
				Handler: carownerinfo.UpdateCarOwnerInfoHandler(serverCtx),
			},
		},
		rest.WithJwt(serverCtx.Config.JwtConf.AccessSecret),
		rest.WithPrefix("/v1"),
	)
}
