package zError

import (
	"fmt"
)

type Error interface {
	error
	GetCode() int
	GetMessage() string
}

type BaseError struct {
	Code    int
	Message string
}

func (e *BaseError) Error() string {
	return fmt.Sprintf("error code: %d, message: %s", e.Code, e.Message)
}

func (e *BaseError) GetCode() int {
	return e.Code
}

func (e *BaseError) GetMessage() string {
	return e.Message
}

func New(msg string) Error {
	return &BaseError{
		Code:    0,
		Message: msg,
	}
}

func NewWithCode(code int, msg string) Error {
	return &BaseError{
		Code:    code,
		Message: msg,
	}
}

func Errorf(format string, args ...interface{}) Error {
	return &BaseError{
		Code:    0,
		Message: fmt.Sprintf(format, args...),
	}
}
