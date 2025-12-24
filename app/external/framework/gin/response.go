package framework

import (
	"errors"
	"net/http"

	appError "github.com/financial_advisor/app/errors"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type responseError struct {
	Type    string `json:"type"`
	Code    string `json:"code"`
	Message string `json:"message"`
	Param   any    `json:"param,omitempty"`
}

type response struct {
	Data   any             `json:"data,omitempty"`
	Paging any             `json:"paging,omitempty"`
	Errors []responseError `json:"errors,omitempty"`
}

// RenderData returns data response
func RenderData(ctx *gin.Context, data, paging any) {
	ctx.JSON(http.StatusOK, response{
		Data:   data,
		Paging: paging,
	})
}

// RenderErrors clasifies errors to render appropriate error response
// - Bad request: error is of type appError.SystemErrors
// - System error: error is of type appError.SystemError
// - Internal server error: all other errors
func RenderErrors(ctx *gin.Context, err error) {
	var errs appError.SystemErrors

	if errors.As(err, &errs) {
		ctx.JSON(http.StatusBadRequest, response{
			Errors: fromSystemErrors(errs),
		})

		return
	}

	var systemErr appError.SystemError
	if errors.As(err, &systemErr) {
		ctx.JSON(systemErr.HTTPCode(), response{
			Errors: fromSystemErrors(appError.SystemErrors{systemErr}),
		})

		return
	}

	logrus.WithError(err).Errorln("internal server error")

	ctx.JSON(http.StatusInternalServerError, response{
		Errors: []responseError{{
			Type:    string(appError.ErrorTypeInternal),
			Code:    string(appError.ErrorCodeInternal),
			Message: "internal server error",
		}},
	})
}

func fromSystemErrors(errs appError.SystemErrors) []responseError {
	if len(errs) == 0 {
		return nil
	}

	respErrors := make([]responseError, len(errs))

	for i, e := range errs {
		respErrors[i] = responseError{
			Type:    string(e.Type()),
			Code:    string(e.Code()),
			Message: e.Message(),
		}
	}

	return respErrors
}
