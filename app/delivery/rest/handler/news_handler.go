package handler

import (
	"net/http"
	"strconv"

	"github.com/financial_advisor/app/delivery/rest/payload"
	"github.com/financial_advisor/app/delivery/rest/presenter"
	"github.com/financial_advisor/app/errors"
	"github.com/financial_advisor/app/registry"
	"github.com/gin-gonic/gin/binding"
)

type NewsHandler struct{}

func NewNewsHandler() *NewsHandler {
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

	if err := payl.Validate(); err != nil {
		return nil, nil, err
	}

	news, paging, err := uc.Execute(req.Context(), payl.ToDTO())
	if err != nil {
		return nil, nil, err
	}

	return presenter.FormNews(news), presenter.Paging{Total: paging.Total}, nil
}

// Get function in handler to return news detail
// @Summary Get news by id
// @Description Return news detail by id
// @Tags news
// @Produce json
// @Success 200 {object} presenter.New
// @Router /news/{id} [get]
func (hdl *NewsHandler) Get(req *http.Request) (any, error) {
	var (
		uc        = registry.InjectNewsGetUsecase()
		newsIDStr = req.Context().Value("id").(string)
		query     payload.NewsGetRequest
	)

	newsID, err := strconv.Atoi(newsIDStr)
	if err != nil {
		return nil, errors.NewErrorBadRequest(
			errors.ErrorCodeBadRequest,
			"invalid id param",
		)
	}

	if err = payload.ShouldBindWith(req, &query, binding.Query); err != nil {
		return nil, errors.NewErrorBadRequest(
			errors.ErrorCodeBadRequest,
			"invalid query payload",
		)
	}

	if err := query.Validate(); err != nil {
		return nil, err
	}

	news, err := uc.Execute(req.Context(), uint64(newsID), query.ToDTO())
	if err != nil {
		return nil, err
	}

	return presenter.FormNew(news), nil
}

// Create function in handler to create a news article
// @Summary Create news
// @Description Create news article by submitting necessary data
// @Tags news
// @Accept json
// @Produce json
// @Param	 payloadBody   body   payload.NewsCreateRequest true    "Body of request"
// @Success 200 {object} presenter.New
// @Router /news [post]
func (hdl *NewsHandler) Create(req *http.Request) (any, error) {
	var (
		uc   = registry.InjectNewsCreateUsecase()
		payl payload.NewsCreateRequest
	)

	if err := payload.ShouldBindWith(req, &payl, binding.JSON); err != nil {
		return nil, errors.NewErrorBadRequest(
			errors.ErrorCodeBadRequest,
			"invalid request payload",
		)
	}

	if err := payl.Validate(); err != nil {
		return nil, err
	}

	news, err := uc.Execute(req.Context(), payl.ToDTO())
	if err != nil {
		return nil, err
	}

	return presenter.FormNew(news), nil
}

// Delete function in handler to delete an article
// @Summary Delete news
// @Description Delete news article by id
// @Tags news
// @Accept json
// @Produce json
// @Param	queryBody	query	payload.NewsSearchRequest	false	"query params to provide keywords"
// @Success 200 {object} any
// @Router /news/{id} [delete]
func (hdl *NewsHandler) Delete(req *http.Request) error {
	var (
		uc        = registry.InjectNewsDeleteUsecase()
		newsIDStr = req.Context().Value("id").(string)
	)

	newsID, err := strconv.Atoi(newsIDStr)
	if err != nil {
		return errors.NewErrorBadRequest(
			errors.ErrorCodeBadRequest,
			"invalid id param",
		)
	}

	if err := uc.Execute(req.Context(), uint64(newsID)); err != nil {
		return err
	}

	return nil
}

// GetSearchSuggestions function in handler to return news search suggestions
// @Summary Get search suggestions by keywords
// @Description Return news search suggestions by keywords
// @Tags news
// @Produce json
// @Param	queryBody	query	payload.NewsSearchSuggestionRequest	false	"query params to provide keywords"
// @Success 200 {array} string
// @Router /news/search/suggestions [get]
func (hdl *NewsHandler) GetSearchSuggestions(req *http.Request) (any, error) {
	var (
		uc    = registry.InjectNewsSearchSuggestionGetUsecase()
		query payload.NewsSearchSuggestionRequest
	)

	if err := payload.ShouldBindWith(req, &query, binding.Query); err != nil {
		return nil, errors.NewErrorBadRequest(
			errors.ErrorCodeBadRequest,
			"invalid id param",
		)
	}

	suggestions, err := uc.Execute(req.Context(), query.ToDTO())
	if err != nil {
		return nil, err
	}

	return suggestions, nil
}

// Search function in handler to return news by keywords
// @Summary Find news by keywords
// @Description Return news search suggestions by keywords
// @Tags news
// @Produce json
// @Param	queryBody	query	payload.NewsSearchRequest	false	"query params to provide keywords"
// @Success 200 {array} presenter.New
// @Router /news/search [get]
func (hdl *NewsHandler) Search(req *http.Request) (any, error) {
	var (
		uc    = registry.InjectNewsSearchUsecase()
		query payload.NewsSearchRequest
	)

	if err := payload.ShouldBindWith(req, &query, binding.Query); err != nil {
		return nil, errors.NewErrorBadRequest(
			errors.ErrorCodeBadRequest,
			"invalid id param",
		)
	}

	news, err := uc.Execute(req.Context(), query.ToDTO())
	if err != nil {
		return nil, err
	}

	return presenter.FormNews(news), nil
}
