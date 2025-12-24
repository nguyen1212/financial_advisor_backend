// Package framework defines gin specific implementations of framework components.
package framework

import (
	"github.com/financial_advisor/app/delivery/rest/handler"
	"github.com/financial_advisor/app/errors"
	"github.com/gin-gonic/gin"
)

type baseHandlerWrapper struct {
	findHandler   handler.FindHandler
	getHandler    handler.GetHandler
	createHandler handler.CreateHandler
	deleteHandler handler.DeleteHandler
}

func newNewsHandlerWrapper(
	newsHandler *handler.NewsHandler,
) *baseHandlerWrapper {
	return &baseHandlerWrapper{
		findHandler:   newsHandler,
		createHandler: newsHandler,
		getHandler:    newsHandler,
		deleteHandler: newsHandler,
	}
}

func newPublisherHandlerWrapper(
	publisherHandler *handler.PublisherHandler,
) *baseHandlerWrapper {
	return &baseHandlerWrapper{
		findHandler:   publisherHandler,
		createHandler: publisherHandler,
		getHandler:    publisherHandler,
	}
}

func (hdl *baseHandlerWrapper) Find(ctx *gin.Context) {
	if hdl.findHandler == nil {
		RenderErrors(ctx, errors.NewErrorNotImplemented(
			errors.ErrorCodeNotImplemented,
			"method find is not implemented",
		))

		return
	}

	data, paging, err := hdl.findHandler.Find(ctx.Request)
	if err != nil {
		RenderErrors(ctx, err)

		return
	}

	RenderData(ctx, data, paging)
}

func (hdl *baseHandlerWrapper) Get(ctx *gin.Context) {
	if hdl.getHandler == nil {
		RenderErrors(ctx, errors.NewErrorNotImplemented(
			errors.ErrorCodeNotImplemented,
			"method get is not implemented",
		))

		return
	}

	data, err := hdl.getHandler.Get(ctx.Request)
	if err != nil {
		RenderErrors(ctx, err)

		return
	}

	RenderData(ctx, data, nil)
}

func (hdl *baseHandlerWrapper) Create(ctx *gin.Context) {
	if hdl.createHandler == nil {
		RenderErrors(ctx, errors.NewErrorNotImplemented(
			errors.ErrorCodeNotImplemented,
			"method create is not implemented",
		))

		return
	}

	data, err := hdl.createHandler.Create(ctx.Request)
	if err != nil {
		RenderErrors(ctx, err)

		return
	}

	RenderData(ctx, data, nil)
}

func (hdl *baseHandlerWrapper) Delete(ctx *gin.Context) {
	if hdl.deleteHandler == nil {
		RenderErrors(ctx, errors.NewErrorNotImplemented(
			errors.ErrorCodeNotImplemented,
			"method delete is not implemented",
		))
		return
	}

	err := hdl.deleteHandler.Delete(ctx.Request)
	if err != nil {
		RenderErrors(ctx, err)

		return
	}

	RenderData(ctx, gin.H{"message": "deleted successfully"}, nil)
}
