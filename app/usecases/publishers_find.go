package usecases

import (
	"context"
	"fmt"

	"github.com/financial_advisor/app/usecases/dto"
	"github.com/financial_advisor/app/domain/repository"
	"github.com/financial_advisor/app/external/db/gorm/specifications"
)

type PublishersFindUsecase interface {
	Execute(
		ctx context.Context,
		req dto.PublishersFindRequest,
	) ([]dto.Publisher, error)
}

type publishersFindUsecase struct {
	publisherRepo repository.PublisherRepository
}

func NewPublishersFindUsecase(
	publisherRepo repository.PublisherRepository,
) PublishersFindUsecase {
	return &publishersFindUsecase{
		publisherRepo: publisherRepo,
	}
}

func (uc *publishersFindUsecase) Execute(
	ctx context.Context,
	req dto.PublishersFindRequest,
) ([]dto.Publisher, error) {
	publishers, err := uc.publisherRepo.Find(
		ctx,
		specifications.NewPublishersByNone(),
		specifications.ToPaging(req.Paging.Size, req.Paging.Page),
	)
	if err != nil {
		return nil, fmt.Errorf("find news by date range: %w", err)
	}

	dtoPublishers := make([]dto.Publisher, len(publishers))

	for i := range publishers {
		dtoPublishers[i] = dto.ToDtoPublisher(publishers[i])
	}

	return dtoPublishers, nil
}
