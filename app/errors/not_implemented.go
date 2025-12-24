package errors

import (
	"fmt"
	"net/http"
)

type errorNotImplemented struct {
	code    ErrorCode
	message string
}

func NewErrorNotImplemented(
	code ErrorCode,
	message string,
) SystemError {
	return errorNotImplemented{
		code:    code,
		message: message,
	}
}

func (e errorNotImplemented) Type() ErrorType {
	return ErrorTypeNotImplemented
}

func (e errorNotImplemented) Code() ErrorCode {
	return e.code
}

func (e errorNotImplemented) Message() string {
	return e.message
}

func (e errorNotImplemented) HTTPCode() int {
	return http.StatusNotImplemented
}

func (e errorNotImplemented) Error() string {
	return fmt.Sprintf("Type: %v, \tCode: %v, \tMessage: %s",
		ErrorTypeNotImplemented, e.Code(), e.Message())
}
