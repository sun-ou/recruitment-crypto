package pkg

import (
	"fmt"
	"net/http"
)

var (
	Success       = NewError(0, "success")
	ServerError   = NewError(100001, "internal error")
	InvaildParams = NewError(100002, "invaild params")
)

var errCode = map[int]string{}

type Error struct {
	code    int      `json:"code"`
	msg     string   `json:"msg"`
	details []string `json:"details"`
}

func NewError(code int, msg string) *Error {
	if _, ok := errCode[code]; ok {
		panic(fmt.Sprintf("error code %d is already defined", code))
	}
	errCode[code] = msg
	return &Error{code: code, msg: msg}
}

func (e *Error) Error() string {
	return fmt.Sprintf("code: %d, msg: %s", e.code, errCode[e.code])
}

func (e *Error) Code() int { return e.code }

func (e *Error) Msg() string { return e.msg }

func (e *Error) Details() []string { return e.details }

func (e *Error) WitchDetails(details ...string) *Error {
	newError := *e
	newError.details = []string{}
	newError.details = append(newError.details, details...)
	return &newError
}

func (e *Error) StatusCode() int {
	switch e.Code() {
	case Success.Code():
		return http.StatusOK
	case InvaildParams.Code():
		return http.StatusBadRequest
	default:
		return http.StatusInternalServerError
	}
}
