package errcode

import "net/http"

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
