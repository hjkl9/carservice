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
	userorder "carservice/internal/handler/userorder"
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
			{
				Method:  http.MethodGet,
				Path:    "/ws/test",
				Handler: common.WebsocketTestHandler(serverCtx),
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
				Method:  http.MethodGet,
				Path:    "/user/login/mock",
				Handler: user.MockLoginHandler(serverCtx),
			},
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
				Path:    "/carBrand/optionList",
				Handler: carbrand.OptionListHandler(serverCtx),
			},
		},
		rest.WithPrefix("/v1"),
	)

	server.AddRoutes(
		[]rest.Route{
			{
				Method:  http.MethodGet,
				Path:    "/carBrandSeries/optionList",
				Handler: carbrandseries.OptionListHandler(serverCtx),
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
			{
				Method:  http.MethodGet,
				Path:    "/carOwnerInfo/:id/list",
				Handler: carownerinfo.GetCarOwnerInfoListHandler(serverCtx),
			},
			{
				Method:  http.MethodGet,
				Path:    "/carOwnerInfo/:id",
				Handler: carownerinfo.GetCarOwnerInfoHandler(serverCtx),
			},
			{
				Method:  http.MethodDelete,
				Path:    "/carOwnerInfo/:id",
				Handler: carownerinfo.DeleteCarOwnerInfoHandler(serverCtx),
			},
		},
		rest.WithJwt(serverCtx.Config.JwtConf.AccessSecret),
		rest.WithPrefix("/v1"),
	)

	server.AddRoutes(
		[]rest.Route{
			{
				Method:  http.MethodPost,
				Path:    "/userOrder",
				Handler: userorder.CreateUserOrderHandler(serverCtx),
			},
			{
				Method:  http.MethodGet,
				Path:    "/userOrder/:id",
				Handler: userorder.GetUserOrderHandler(serverCtx),
			},
			{
				Method:  http.MethodGet,
				Path:    "/userOrder/list",
				Handler: userorder.GetrUserOrderListHandler(serverCtx),
			},
			{
				Method:  http.MethodDelete,
				Path:    "/userOrder/:id",
				Handler: userorder.DeleteUserOrderHandler(serverCtx),
			},
		},
		rest.WithJwt(serverCtx.Config.JwtConf.AccessSecret),
		rest.WithPrefix("/v1"),
	)
}
