package usecases

import (
	"context"
	"errors"
	"fmt"

	"github.com/financial_advisor/app/domain/repository"
	appErrors "github.com/financial_advisor/app/errors"
	"github.com/financial_advisor/app/external/db/gorm/specifications"
	"github.com/financial_advisor/app/usecases/dto"
)

type PublisherGetUsecase interface {
	Execute(
		ctx context.Context,
		publisherID uint64,
	) (dto.Publisher, error)
}

type publisherGetUsecase struct {
	publisherRepo repository.PublisherRepository
}

func NewPublisherGetUsecase(
	publisherRepo repository.PublisherRepository,
) PublisherGetUsecase {
	return &publisherGetUsecase{
		publisherRepo: publisherRepo,
	}
}

func (uc *publisherGetUsecase) Execute(
	ctx context.Context,
	publisherID uint64,
) (dto.Publisher, error) {
	publisher, err := uc.publisherRepo.Get(
		ctx,
		specifications.NewPublisherByID(
			publisherID,
		),
	)
	if err != nil {
		if errors.Is(err, appErrors.ErrNotFound) {
			return dto.Publisher{}, appErrors.NewErrorBadRequest(
				appErrors.ErrorCodePublisherNotFound,
				"publisher not found",
			)
		}

		return dto.Publisher{}, fmt.Errorf("find publisher by date range: %w", err)
	}

	return dto.ToDtoPublisher(publisher), nil
}
