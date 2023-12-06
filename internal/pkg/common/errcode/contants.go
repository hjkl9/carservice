package errcode

import "net/http"

// basic errors.
var (
	Ok                  = New(http.StatusOK, "-", "欧克")
	NotFound            = New(http.StatusNotFound, "-", "找不到资源")
	InvalidParamsError  = New(http.StatusBadRequest, "-", "无效的参数")
	UnauthorizedError   = New(http.StatusUnauthorized, "-", "未认证")
	InternalServerError = New(http.StatusInternalServerError, "-", "服务器内部发生错误")
)

// SMS-related errors.
var (
	InvalidPhoneNumberError = New(http.StatusBadRequest, "-", "无效的手机号码")
	SmsSdkCallingError      = New(http.StatusInternalServerError, "-", "短信服务调用 SDK 时出现错误")
)
