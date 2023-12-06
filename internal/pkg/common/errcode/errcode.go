package errcode

type ErrCode struct {
	Code    int
	ErrCode string
	Msg     string
}

func New(code int, errCode, msg string) *ErrCode {
	return &ErrCode{code, errCode, msg}
}

func (e *ErrCode) Error() string {
	return e.Msg
}

func (e *ErrCode) SetMsg(msg string) *ErrCode {
	e.Msg = msg
	return e
}
