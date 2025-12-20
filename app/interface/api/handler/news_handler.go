// Package handler accepts incoming requests from external traffic.
package handler

import (
	"net/http"

	"github.com/financial_advisor/app/interface/api/presenter"
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
// @Success 200 {object} presenter.New
// @Router /news [get]
func (hdl *NewsHandler) Find(req *http.Request) (any, any, error) {
	return presenter.FormNews(), presenter.Paging{}, nil
}
