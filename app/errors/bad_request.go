package errors

import (
	"fmt"
	"net/http"
)

type errorBadRequest struct {
	code    ErrorCode
	message string
}

func NewErrorBadRequest(code ErrorCode, message string) SystemError {
	return errorBadRequest{
		code:    code,
		message: message,
	}
}

func (e errorBadRequest) Type() ErrorType {
	return ErrorTypeBadRequest
}

func (e errorBadRequest) Code() ErrorCode {
	return e.code
}

func (e errorBadRequest) Message() string {
	return e.message
}

func (e errorBadRequest) HTTPCode() int {
	return http.StatusBadRequest
}

func (e errorBadRequest) Error() string {
	return fmt.Sprintf("Type: %v, \tCode: %v, \tMessage: %s",
		ErrorTypeBadRequest, e.Code(), e.Message())
}
