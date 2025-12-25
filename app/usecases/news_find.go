// Package usecases define implementation of business logic
package usecases

import (
	"context"
	"fmt"

	"github.com/financial_advisor/app/domain/repository"
	"github.com/financial_advisor/app/external/db/gorm/specifications"
	"github.com/financial_advisor/app/usecases/dto"
)

type NewsFindUsecase interface {
	Execute(
		ctx context.Context,
		req dto.NewsFindRequest,
	) ([]dto.News, error)
}

type newsFindUsecase struct {
	newsRepo repository.NewsRepository
}

func NewNewsFindUsecase(
	newsRepo repository.NewsRepository,
) NewsFindUsecase {
	return &newsFindUsecase{
		newsRepo: newsRepo,
	}
}

func (uc *newsFindUsecase) Execute(
	ctx context.Context,
	req dto.NewsFindRequest,
) ([]dto.News, error) {
	news, err := uc.newsRepo.Find(
		ctx,
		specifications.NewNewsByDate(
			req.From,
			req.To,
			req.Status,
		),
	)
	if err != nil {
		return nil, fmt.Errorf("find news by date range: %w", err)
	}

	dtoNews := make([]dto.News, len(news))

	for i := range news {
		dtoNews[i] = dto.ToDtoNews(news[i])
	}

	return dtoNews, nil
}
