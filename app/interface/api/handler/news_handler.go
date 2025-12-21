// Package handler accepts incoming requests from external traffic.
package handler

import (
	"net/http"

	"github.com/financial_advisor/app/errors"
	"github.com/financial_advisor/app/interface/api/payload"
	"github.com/financial_advisor/app/interface/api/presenter"
	"github.com/financial_advisor/app/registry"
	"github.com/gin-gonic/gin/binding"
)

type NewsHandler struct{}

func NewNewsHandler() Handler {
	return &NewsHandler{}
}

// Find function in handler to return list of news
// @Summary Find news
// @Description Return list of news from a specific day range with pagination
// @Tags news
// @Produce json
// @Param	queryBody	query	payload.NewsFindRequest	false	"query params to find news"
// @Success 200 {array} presenter.New
// @Router /news [get]
func (hdl *NewsHandler) Find(
	req *http.Request,
) (any, any, error) {
	var (
		uc   = registry.InjectNewsFindUsecase()
		payl payload.NewsFindRequest
	)

	if err := payload.ShouldBindWith(req, &payl, binding.Query); err != nil {
		return nil, nil, errors.NewErrorBadRequest(
			errors.ErrorCodeBadRequest,
			"invalid request payload",
		)
	}

	news, err := uc.Execute(req.Context(), payl.ToDTO())
	if err != nil {
		return nil, nil, err
	}

	return presenter.FormNews(news), presenter.Paging{}, nil
}
