package usecases

import (
	"context"
	"errors"
	"fmt"

	"github.com/financial_advisor/app/domain/repository"
	appErrors "github.com/financial_advisor/app/errors"
	"github.com/financial_advisor/app/external/db/goqu/specifications"
)

type NewsDeleteUsecase interface {
	Execute(
		ctx context.Context,
		newsID uint64,
	) error
}

type newsDeleteUsecase struct {
	newsRepo             repository.NewsRepository
	newsWithFullTextRepo repository.NewsWithFullTextRepository
}

func NewNewsDeleteUsecase(
	newsRepo repository.NewsRepository,
	newsWithFullTextRepo repository.NewsWithFullTextRepository,
) NewsDeleteUsecase {
	return &newsDeleteUsecase{
		newsRepo:             newsRepo,
		newsWithFullTextRepo: newsWithFullTextRepo,
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

	if news.NewsWithFullTextID != 0 {
		if err := uc.newsWithFullTextRepo.Delete(ctx, news.NewsWithFullTextID); err != nil {
			return fmt.Errorf("delete news full text: %w", err)
		}
	}

	if err := uc.newsRepo.Delete(
		ctx,
		newsID,
	); err != nil {
		return fmt.Errorf("find news by date range: %w", err)
	}

	return nil
}
