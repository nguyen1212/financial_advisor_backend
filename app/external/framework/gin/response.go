package framework

import (
	"net/http"

	"github.com/gin-gonic/gin"
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
