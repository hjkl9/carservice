package errcode

type ErrCode struct {
	Code    int
	ErrCode string
	Msg     string
}

func (e *ErrCode) Error() string {
	return e.Msg
}

func New(code int, errCode, msg string) error {
	return &ErrCode{code, errCode, msg}
}
