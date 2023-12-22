package errcode

type ErrCode struct {
	Code    int
	ErrCode string
	Msg     string
	Details []string
}

func New(code int, errCode, msg string) *ErrCode {
	return &ErrCode{code, errCode, msg, nil}
}

func (e *ErrCode) Lazy(msg string, details ...string) *ErrCode {
	e.Msg = msg
	if len(details) > 0 {
		e.Details = details
	}
	return e
}

func (e *ErrCode) Error() string {
	return e.Msg
}

func (e *ErrCode) SetMsg(msg string) *ErrCode {
	e.Msg = msg
	return e
}

func (e *ErrCode) SetDetails(details ...string) *ErrCode {
	e.Details = details
	return e
}
