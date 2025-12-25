package usecases

//go:generate mockgen -destination=./mock/mock_$GOFILE -source=$GOFILE -package=mock

import (
	"context"
	"fmt"

	"github.com/financial_advisor/app/domain/entity"
	"github.com/financial_advisor/app/domain/repository"
	"github.com/financial_advisor/app/external/db/gorm/specifications"
)

type fallbackScrapperUsecase struct {
	newsRepo repository.NewsRepository
}

type FallbackScrapperUsecase interface {
	Execute(ctx context.Context, job entity.WebScrapperJob) error
}

func NewFallbackScrapperUsecase(newsRepo repository.NewsRepository) FallbackScrapperUsecase {
	return fallbackScrapperUsecase{newsRepo}
}

func (uc fallbackScrapperUsecase) Execute(
	ctx context.Context,
	job entity.WebScrapperJob,
) error {
	news, err := uc.newsRepo.Get(ctx, specifications.NewNewsByID(job.NewsID))
	if err != nil {
		return fmt.Errorf("get news by id: %w", err)
	}

	news.Status = entity.NewsStatusFailed

	if err := uc.newsRepo.Update(ctx, &news); err != nil {
		return fmt.Errorf("update news status to failed: %w", err)
	}

	return nil
}
