package api

import "net/http"

var OK = NewApiCode(http.StatusOK, ("0"), "欧克")

// System errors
var (
	SystemErr             = NewApiCode(http.StatusInternalServerError, ("10001"), "系统发生错误")
	ServiceUnavailableErr = NewApiCode(http.StatusServiceUnavailable, ("10002"), "服务不可用")
	RemoteServiceErr      = NewApiCode(http.StatusInternalServerError, ("10003"), "远程服务发生错误")
	RpcServiceErr         = NewApiCode(http.StatusInternalServerError, ("10004"), "RPC 服务发生错误")
	IllegalRequestErr     = NewApiCode(http.StatusInternalServerError, ("10005"), "非法请求")
	InvalidParametersErr  = NewApiCode(http.StatusBadRequest, ("10006"), "无效的参数")
	RouteNotFoundErr      = NewApiCode(http.StatusNotFound, ("10007"), "找不到该接口")
	MethodNotAllowedErr   = NewApiCode(http.StatusMethodNotAllowed, ("10008"), "错误的请求方式")
	RequestEntityTooLarge = NewApiCode(http.StatusRequestEntityTooLarge, ("10009"), "请求实体超出范围")
	// Database errors
	DatabaseCreateErr    = NewApiCode(http.StatusInternalServerError, ("10010"), "数据库创建数据时发生错误")
	DatabaseDeleteErr    = NewApiCode(http.StatusInternalServerError, ("10011"), "数据库删除数据时发生错误")
	DatabaseUpdateErr    = NewApiCode(http.StatusInternalServerError, ("10012"), "数据库更新数据时发生错误")
	DatabaseGetErr       = NewApiCode(http.StatusInternalServerError, ("10013"), "数据库获取数据时发生错误")
	DatabasePrepareErr   = NewApiCode(http.StatusInternalServerError, ("10014"), "数据库预处理语句时发生错误")
	DatabaseTrasationErr = NewApiCode(http.StatusInternalServerError, ("10015"), "数据库创建事务时发生错误")
	DatabaseRollbackErr  = NewApiCode(http.StatusInternalServerError, ("10016"), "数据库回滚事务时发生错误")
	DatabaseCommitErr    = NewApiCode(http.StatusInternalServerError, ("10017"), "数据库提交事务时发生错误")
)

// Each service errors //
// User service
var (
	UserNotFoundErr       = NewApiCode(http.StatusNotFound, ("20001"), "该用户不存在")
	UserIsLockedErr       = NewApiCode(http.StatusNotFound, ("20002"), "该用户已被锁定")
	UserUnauthorizedErr   = NewApiCode(http.StatusUnauthorized, ("20003"), "未登录")
	UserTokenExpiredErr   = NewApiCode(http.StatusUnauthorized, ("20004"), "登录过期")
	WechatCode2SessionErr = NewApiCode(http.StatusOK, ("20005"), "TODO")
	GenTokenErr           = NewApiCode(http.StatusInternalServerError, ("20006"), "生成 Token 时发生错误")
)
