package georegeo

type AMapError struct {
	code string
	msg  string
}

func NewMapError(code, msg string) *AMapError {
	return &AMapError{code, msg}
}

func (e *AMapError) Error() string {
	return e.msg
}

func (e *AMapError) GetCode() string {
	return e.code
}

func (e *AMapError) GetMsg() string {
	return e.msg
}

var (
	ENGINE_RESPONSE_DATA_ERROR = NewMapError("300**", "地图服务响应失败")
	QUOTA_PLAN_RUN_OUT         = NewMapError("40000", "地图服务余额耗尽")
	ILLEGAL_CONTENT            = NewMapError("20012", "查询信息存在非法内容")
	SERVER_IS_BUSY             = NewMapError("10016", "服务器负载过高")
	UNKNOWN_ERROR              = NewMapError("00000", "未知的错误")
	// Other errors.
)

func MatchError(code string) *AMapError {
	// 针对匹配
	if code[0:3] == "300" {
		return ENGINE_RESPONSE_DATA_ERROR
	}
	switch code {
	case "40000":
		return QUOTA_PLAN_RUN_OUT
	case "20012":
		return ILLEGAL_CONTENT
	case "10016":
		return SERVER_IS_BUSY
	// case "OTHER...":
	// 	return OTHER...
	default:
		return UNKNOWN_ERROR
	}
}
