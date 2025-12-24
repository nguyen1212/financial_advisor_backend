package usecases

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/financial_advisor/app/domain/repository"
	appErrors "github.com/financial_advisor/app/errors"
	"github.com/financial_advisor/app/external/db/gorm/specifications"
)

type NewsDeleteUsecase interface {
	Execute(
		ctx context.Context,
		newsID uint64,
	) error
}

type newsDeleteUsecase struct {
	newsRepo repository.NewsRepository
}

func NewNewsDeleteUsecase(
	newsRepo repository.NewsRepository,
) NewsDeleteUsecase {
	return &newsDeleteUsecase{
		newsRepo: newsRepo,
	}
}

func (uc *newsDeleteUsecase) Execute(
	ctx context.Context,
	newsID uint64,
) error {
	news, err := uc.newsRepo.Get(
		ctx,
		specifications.NewNewsByID(newsID),
	)
	if err != nil {
		if errors.Is(err, appErrors.ErrNotFound) {
			return nil
		}

		return fmt.Errorf("get news by id: %w", err)
	}

	if err := os.Remove(news.StoragePath()); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("remove news file: %w", err)
	}

	if err := uc.newsRepo.Delete(
		ctx,
		newsID,
	); err != nil {
		return fmt.Errorf("find news by date range: %w", err)
	}

	return nil
}
