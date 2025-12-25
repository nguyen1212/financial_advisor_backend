package usecases

import (
	"context"
	"fmt"

	"github.com/dchest/validator"
	"github.com/financial_advisor/app/domain/entity"
	"github.com/financial_advisor/app/domain/repository"
	"github.com/financial_advisor/app/errors"
	appErrors "github.com/financial_advisor/app/errors"
	"github.com/financial_advisor/app/external/db/gorm/specifications"
	"github.com/financial_advisor/app/usecases/dto"
)

type PublisherCreateUsecase interface {
	Execute(
		ctx context.Context,
		req dto.PublisherCreateRequest,
	) (dto.Publisher, error)
}

func NewPublisherCreateUsecase(
	publisherRepo repository.PublisherRepository,
) PublisherCreateUsecase {
	return &publisherCreateUsecase{
		publisherRepo: publisherRepo,
	}
}

type publisherCreateUsecase struct {
	publisherRepo repository.PublisherRepository
}

func (uc *publisherCreateUsecase) Execute(
	ctx context.Context,
	req dto.PublisherCreateRequest,
) (dto.Publisher, error) {
	if !validator.IsValidDomain(req.Domain) {
		return dto.Publisher{}, appErrors.NewErrorBadRequest(
			appErrors.ErrorCodeBadRequest,
			"invalid domain format",
		)
	}

	domain := validator.NormalizeDomain(req.Domain)

	countPublisher, err := uc.publisherRepo.Count(
		ctx,
		specifications.NewPublisherByDomain(domain),
	)
	if err != nil {
		return dto.Publisher{}, fmt.Errorf("count existing publisher by domain: %w", err)
	}

	if countPublisher > 0 {
		return dto.Publisher{}, errors.NewErrorConflicted(
			errors.ErrorCodeConflicted,
			"publisher with the given domain already exists",
		)
	}

	newPublisher := entity.Publisher{
		Name:        req.Name,
		Domain:      domain,
		Description: req.Description,
	}

	if err := uc.publisherRepo.Create(
		ctx,
		&newPublisher,
	); err != nil {
		return dto.Publisher{}, fmt.Errorf("create new publisher: %w", err)
	}

	return dto.ToDtoPublisher(newPublisher), nil
}
