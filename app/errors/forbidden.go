package errors

import (
	"fmt"
	"net/http"
)

type errorForbidden struct {
	code    ErrorCode
	message string
}

func NewErrorForbidden(code ErrorCode, message string) SystemError {
	return errorForbidden{
		code:    code,
		message: message,
	}
}

func (e errorForbidden) Type() ErrorType {
	return ErrorTypeForbidden
}

func (e errorForbidden) Code() ErrorCode {
	return e.code
}

func (e errorForbidden) Message() string {
	return e.message
}

func (e errorForbidden) HTTPCode() int {
	return http.StatusForbidden
}

func (e errorForbidden) Error() string {
	return fmt.Sprintf("Type: %v, \tCode: %v, \tMessage: %s",
		ErrorTypeForbidden, e.Code(), e.Message())
}
