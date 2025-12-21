// Package framework defines gin specific implementations of framework components.
package framework

import (
	"github.com/financial_advisor/app/interface/api/handler"
	"github.com/gin-gonic/gin"
)

type handlerWrapper struct {
	baseHandler handler.Handler
}

func newHandlerWrapper(base handler.Handler) *handlerWrapper {
	return &handlerWrapper{baseHandler: base}
}

func (hdl *handlerWrapper) Find(ctx *gin.Context) {
	data, paging, err := hdl.baseHandler.Find(ctx.Request)
	if err != nil {
		RenderErrors(ctx, err)

		return
	}

	RenderData(ctx, data, paging)
}
