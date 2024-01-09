package api

import (
	"fmt"
)

type ApiCoder interface {
	HttpCode() int16
	Code() string
	Message() string
	WithDetails(...string)
	Details() []string
	Error() string
}

// apiCode is used to API error code.
type apiCode struct {
	httpCode int16
	code     string
	msg      string
	details  []string
}

var e error = NewApiCode(400, "no", "")

// NewApiCode is create a error code entity.
func NewApiCode(httpCode int16, code, msg string) ApiCoder {
	details := []string{}
	return &apiCode{httpCode, code, msg, details}
}

// HttpCode returns http status code.
func (ac *apiCode) HttpCode() int16 {
	return ac.httpCode
}

// Code returns error code.
func (ac *apiCode) Code() string {
	return "E" + ac.code
}

// Message returns error message.
func (ac *apiCode) Message() string {
	return ac.msg
}

// WithDetails could put more error details of type of string in response data.
func (ac *apiCode) WithDetails(dts ...string) {
	ac.details = append(ac.details, dts...)
}

// Details returns error details of type of string.
func (ac *apiCode) Details() []string {
	return ac.details
}

func (ac *apiCode) Error() string {
	return fmt.Sprintf("[%s] - %s", ac.Code(), ac.Message())
}

func ParseApiError(err error) ApiCoder {
	return err.(*apiCode)
}
