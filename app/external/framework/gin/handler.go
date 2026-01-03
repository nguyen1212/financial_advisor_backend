// Package framework defines gin specific implementations of framework components.
package framework

import (
	"context"
	"net/http"

	"github.com/financial_advisor/app/delivery/rest/handler"
	"github.com/financial_advisor/app/errors"
	"github.com/gin-gonic/gin"
)

func reflectRequest(req *http.Request, ctx *gin.Context) *http.Request {
	reqCtx := ctx.Request.Context()

	// Add key/value pairs from Gin context to the request context
	for key, value := range ctx.Keys {
		reqCtx = context.WithValue(reqCtx, key, value)
	}

	for _, param := range ctx.Params {
		reqCtx = context.WithValue(reqCtx, param.Key, param.Value)
	}

	// Create a new request with the updated context
	newReq := ctx.Request.WithContext(reqCtx)

	return newReq
}

type baseHandlerWrapper struct {
	findHandler   handler.FindHandler
	getHandler    handler.GetHandler
	createHandler handler.CreateHandler
	deleteHandler handler.DeleteHandler
}

type newsHandlerWrapper struct {
	baseHandlerWrapper
	searchHandler handler.SearchHandler
}

func newNewsHandlerWrapper(
	newsHandler *handler.NewsHandler,
) *newsHandlerWrapper {
	h := newsHandlerWrapper{
		searchHandler: newsHandler,
	}

	h.findHandler = newsHandler
	h.getHandler = newsHandler
	h.createHandler = newsHandler
	h.deleteHandler = newsHandler

	return &h
}

func (hdl *newsHandlerWrapper) GetSearchSuggestions(ctx *gin.Context) {
	if hdl.searchHandler == nil {
		RenderErrors(ctx, errors.NewErrorNotImplemented(
			errors.ErrorCodeNotImplemented,
			"method search is not implemented",
		))

		return
	}

	req := reflectRequest(ctx.Request, ctx)

	data, err := hdl.searchHandler.GetSearchSuggestions(req)
	if err != nil {
		RenderErrors(ctx, err)

		return
	}

	RenderData(ctx, data, nil)
}

func (hdl *newsHandlerWrapper) Search(ctx *gin.Context) {
	if hdl.searchHandler == nil {
		RenderErrors(ctx, errors.NewErrorNotImplemented(
			errors.ErrorCodeNotImplemented,
			"method search is not implemented",
		))

		return
	}

	req := reflectRequest(ctx.Request, ctx)

	data, err := hdl.searchHandler.Search(req)
	if err != nil {
		RenderErrors(ctx, err)

		return
	}

	RenderData(ctx, data, nil)
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

	req := reflectRequest(ctx.Request, ctx)

	data, paging, err := hdl.findHandler.Find(req)
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

	req := reflectRequest(ctx.Request, ctx)

	data, err := hdl.getHandler.Get(req)
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

	req := reflectRequest(ctx.Request, ctx)

	data, err := hdl.createHandler.Create(req)
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

	req := reflectRequest(ctx.Request, ctx)

	err := hdl.deleteHandler.Delete(req)
	if err != nil {
		RenderErrors(ctx, err)

		return
	}

	RenderData(ctx, gin.H{"message": "deleted successfully"}, nil)
}
