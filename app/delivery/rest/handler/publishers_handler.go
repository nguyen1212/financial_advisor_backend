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

type PublisherHandler struct{}

func NewPublisherHandler() *PublisherHandler {
	return &PublisherHandler{}
}

// Find function in handler to return list of publishers
// @Summary Find publishers
// @Description Return list of publishers from a specific day range with pagination
// @Tags publishers
// @Produce json
// @Param	queryBody	query	payload.PublishersFindRequest	false	"query params to find publishers"
// @Success 200 {array} presenter.Publisher
// @Router /publishers [get]
func (hdl *PublisherHandler) Find(
	req *http.Request,
) (any, any, error) {
	var (
		uc   = registry.InjectPublishersFindUsecase()
		payl payload.PublishersFindRequest
	)

	if err := payload.ShouldBindWith(req, &payl, binding.Query); err != nil {
		return nil, nil, errors.NewErrorBadRequest(
			errors.ErrorCodeBadRequest,
			"invalid request payload",
		)
	}

	publishers, err := uc.Execute(req.Context(), payl.ToDTO())
	if err != nil {
		return nil, nil, err
	}

	return presenter.FormPublishers(publishers), presenter.Paging{}, nil
}

// Get function in handler to return publisher detail
// @Summary Get publisher by id
// @Description Return publisher detail by id
// @Tags news
// @Produce json
// @Success 200 {object} presenter.Publisher
// @Router /publishers/{id} [get]
func (hdl *PublisherHandler) Get(req *http.Request) (any, error) {
	var (
		uc             = registry.InjectPublisherGetUsecase()
		publisherIDStr = req.Context().Value("id").(string)
	)

	publisherID, err := strconv.Atoi(publisherIDStr)
	if err != nil {
		return nil, errors.NewErrorBadRequest(
			errors.ErrorCodeBadRequest,
			"invalid id param",
		)
	}

	publisher, err := uc.Execute(req.Context(), uint64(publisherID))
	if err != nil {
		return nil, err
	}

	return presenter.FormPublisher(publisher), nil
}

// Create function in handler to create publisher
// @Summary Create publisher
// @Description Create publishers by submitting necessary data
// @Tags publishers
// @Accept json
// @Produce json
// @Param	 payloadBody   body   payload.PublisherCreateRequest true    "Body of request"
// @Success 200 {object} presenter.Publisher
// @Router /publishers [post]
func (hdl *PublisherHandler) Create(
	req *http.Request,
) (any, error) {
	var (
		uc   = registry.InjectPublisherCreateUsecase()
		payl payload.PublisherCreateRequest
	)

	if err := payload.ShouldBindWith(req, &payl, binding.JSON); err != nil {
		return nil, errors.NewErrorBadRequest(
			errors.ErrorCodeBadRequest,
			"invalid request payload",
		)
	}

	publisher, err := uc.Execute(req.Context(), payl.ToDTO())
	if err != nil {
		return nil, err
	}

	return presenter.FormPublisher(publisher), nil
}
