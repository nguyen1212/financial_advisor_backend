package errors

import (
	"fmt"
	"net/http"
)

type errorNotFound struct {
	code    ErrorCode
	message string
}

func NewErrorNotFound(code ErrorCode, message string) SystemError {
	return errorNotFound{
		code:    code,
		message: message,
	}
}

func (e errorNotFound) Type() ErrorType {
	return ErrorTypeNotFound
}

func (e errorNotFound) Code() ErrorCode {
	return e.code
}

func (e errorNotFound) Message() string {
	return e.message
}

func (e errorNotFound) HTTPCode() int {
	return http.StatusNotFound
}

func (e errorNotFound) Error() string {
	return fmt.Sprintf("Type: %v, \tCode: %v, \tMessage: %s",
		ErrorTypeNotFound, e.Code(), e.Message())
}
