package errcode

import (
	"carservice/internal/pkg/httper/api"
	"net/http"
)

// basic errors.
var (
	Ok                   = New(http.StatusOK, "-", "欧克")
	NotFound             = New(http.StatusNotFound, "-", "找不到资源")
	InvalidParamsError   = New(http.StatusBadRequest, "-", "无效的参数")
	StatusForbiddenError = New(http.StatusForbidden, "-", "请求已被拒绝")
	UnauthorizedError    = New(http.StatusUnauthorized, "-", "未认证")
	InternalServerError  = New(http.StatusInternalServerError, "-", "服务器内部发生错误")
)

// 数据库错误
var (
	DatabaseError = New(http.StatusInternalServerError, "-", "操作数据库时发生错误")
)

// SMS-related errors.
var (
	InvalidPhoneNumberError = New(http.StatusBadRequest, "-", "无效的手机号码")
	SmsSdkCallingError      = New(http.StatusInternalServerError, "-", "短信服务调用 SDK 时出现错误")
	SmsSdkServiceException  = New(http.StatusServiceUnavailable, "-", "短信服务异常")
)

// New version //

var OK = api.OK

// System errors
var (
	SystemErr             = api.NewApiCode(http.StatusInternalServerError, ("10001"), "系统发生错误")
	ServiceUnavailableErr = api.NewApiCode(http.StatusServiceUnavailable, ("10002"), "服务不可用")
	RemoteServiceErr      = api.NewApiCode(http.StatusInternalServerError, ("10003"), "远程服务发生错误")
	RpcServiceErr         = api.NewApiCode(http.StatusInternalServerError, ("10004"), "RPC 服务发生错误")
	IllegalRequestErr     = api.NewApiCode(http.StatusInternalServerError, ("10005"), "非法请求")
	InvalidParametersErr  = api.NewApiCode(http.StatusBadRequest, ("10006"), "无效的参数")
	RouteNotFoundErr      = api.NewApiCode(http.StatusNotFound, ("10007"), "找不到该接口")
	MethodNotAllowedErr   = api.NewApiCode(http.StatusMethodNotAllowed, ("10008"), "错误的请求方式")
	RequestEntityTooLarge = api.NewApiCode(http.StatusRequestEntityTooLarge, ("10009"), "请求实体超出范围")
	// Database errors
	DatabaseCreateErr    = api.NewApiCode(http.StatusInternalServerError, ("10010"), "数据库创建数据时发生错误")
	DatabaseDeleteErr    = api.NewApiCode(http.StatusInternalServerError, ("10011"), "数据库删除数据时发生错误")
	DatabaseUpdateErr    = api.NewApiCode(http.StatusInternalServerError, ("10012"), "数据库更新数据时发生错误")
	DatabaseGetErr       = api.NewApiCode(http.StatusInternalServerError, ("10013"), "数据库获取数据时发生错误")
	DatabasePrepareErr   = api.NewApiCode(http.StatusInternalServerError, ("10014"), "数据库预处理时发生错误")
	DatabaseTrasationErr = api.NewApiCode(http.StatusInternalServerError, ("10015"), "数据库创建事务时发生错误")
	DatabaseRollbackErr  = api.NewApiCode(http.StatusInternalServerError, ("10016"), "数据库回滚事务时发生错误")
	DatabaseCommitErr    = api.NewApiCode(http.StatusInternalServerError, ("10017"), "数据库提交事务时发生错误")
	DatabaseExecuteErr   = api.NewApiCode(http.StatusInternalServerError, ("10018"), "数据库操作执行时发生错误")
	// Redis errors
	RedisGetErr = api.NewApiCode(http.StatusInternalServerError, ("10019"), "缓存读取时发生错误")
	RedisSetErr = api.NewApiCode(http.StatusInternalServerError, ("10019"), "缓存写入时发生错误")
	// JSON Marshal errors
	JSONMarshalErr   = api.NewApiCode(http.StatusInternalServerError, ("10020"), "JSON 序列化时发生错误")
	JSONUnmarshalErr = api.NewApiCode(http.StatusInternalServerError, ("10021"), "JSON 反序列化时发生错误")
)

// Each service errors //
// User service
var (
	UserNotFoundErr       = api.NewApiCode(http.StatusNotFound, ("20001"), "该用户不存在")
	UserIsLockedErr       = api.NewApiCode(http.StatusNotFound, ("20002"), "该用户已被锁定")
	UserUnauthorizedErr   = api.NewApiCode(http.StatusUnauthorized, ("20003"), "未登录")
	UserTokenExpiredErr   = api.NewApiCode(http.StatusUnauthorized, ("20004"), "登录过期")
	WechatCode2SessionErr = api.NewApiCode(http.StatusOK, ("20005"), "TODO")
	GenTokenErr           = api.NewApiCode(http.StatusInternalServerError, ("20006"), "生成 Token 时发生错误")
	InvalidWechatCodeErr  = api.NewApiCode(http.StatusBadRequest, ("20007"), "无效的 jscode")
	UserInvalidTokenErr   = api.NewApiCode(http.StatusUnauthorized, ("20008"), "无效的 Token")
)

// User order service
var (
	OrderNotFoundErr           = api.NewApiCode(http.StatusNotFound, ("20401"), "该用户订单不存在")
	OrderCannotBeCancelledErr  = api.NewApiCode(http.StatusBadRequest, ("20402"), "该用户订单无法被取消")
	DuplicateConfirmedOrderErr = api.NewApiCode(http.StatusBadRequest, "20403", "该用户订单已经重复确认")
	OrderUncomfirmedQuote      = api.NewApiCode(http.StatusBadRequest, "20404", "商家未接单")
	ConfirmedOrderErr          = api.NewApiCode(http.StatusBadRequest, "20405", "")
	OrderConfirmAndPayErr      = api.NewApiCode(http.StatusInternalServerError, "20406", "处理确认订单和支付发生错误")
	OrderPrepayErr             = api.NewApiCode(http.StatusInternalServerError, "20407", "订单预支付时发生错误")
)
