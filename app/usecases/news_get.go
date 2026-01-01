package usecases

import (
	"context"
	"errors"
	"fmt"

	"github.com/financial_advisor/app/domain/repository"
	appErrors "github.com/financial_advisor/app/errors"
	goquSpec "github.com/financial_advisor/app/external/db/goqu/specifications"
	"github.com/financial_advisor/app/usecases/dto"
)

type NewsGetUsecase interface {
	Execute(
		ctx context.Context,
		newsID uint64,
		req dto.NewsGetRequest,
	) (dto.News, error)
}

type newsGetUsecase struct {
	newsRepo             repository.NewsRepository
	newsWithFullTextRepo repository.NewsWithFullTextRepository
}

func NewNewsGetUsecase(
	newsRepo repository.NewsRepository,
	newsWithFullTextRepo repository.NewsWithFullTextRepository,
) NewsGetUsecase {
	return &newsGetUsecase{
		newsRepo:             newsRepo,
		newsWithFullTextRepo: newsWithFullTextRepo,
	}
}

func (uc *newsGetUsecase) Execute(
	ctx context.Context,
	newsID uint64,
	req dto.NewsGetRequest,
) (dto.News, error) {
	news, err := uc.newsRepo.Get(
		ctx,
		goquSpec.NewNewsByID(
			newsID,
		),
	)
	if err != nil {
		if errors.Is(err, appErrors.ErrNotFound) {
			return dto.News{}, appErrors.NewErrorNotFound(
				appErrors.ErrorCodeNewsNotFound,
				"news not found",
			)
		}

		return dto.News{}, fmt.Errorf("find news by date range: %w", err)
	}

	if news.NewsWithFullTextID == 0 {
		return dto.ToDtoNews(news), nil
	}

	newsWithFullText, err := uc.newsWithFullTextRepo.Get(
		ctx,
		goquSpec.NewNewsWithFullTextByFileID(
			news.NewsWithFullTextID,
			news.FileSize,
			req.HighlightKeywords,
		),
	)
	if err != nil {
		if errors.Is(err, appErrors.ErrNotFound) {
			return dto.ToDtoNews(news), nil
		}

		return dto.News{}, fmt.Errorf("find news with content by document id: %w", err)
	}

	news.Content = newsWithFullText.Content

	return dto.ToDtoNews(news), nil
}
