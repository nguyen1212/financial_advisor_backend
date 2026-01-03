// Package usecases define implementation of business logic
package usecases

import (
	"context"
	"fmt"

	"github.com/financial_advisor/app/domain/entity"
	"github.com/financial_advisor/app/domain/repository"
	"github.com/financial_advisor/app/external/db/goqu/specifications"
	"github.com/financial_advisor/app/usecases/dto"
	"github.com/samber/lo"
)

type NewsFindUsecase interface {
	Execute(
		ctx context.Context,
		req dto.NewsFindRequest,
	) ([]dto.News, dto.PagingResponse, error)
}

type newsFindUsecase struct {
	newsRepo             repository.NewsRepository
	newsWithFullTextRepo repository.NewsWithFullTextRepository
}

func NewNewsFindUsecase(
	newsRepo repository.NewsRepository,
	newsWithFullTextRepo repository.NewsWithFullTextRepository,
) NewsFindUsecase {
	return &newsFindUsecase{
		newsRepo:             newsRepo,
		newsWithFullTextRepo: newsWithFullTextRepo,
	}
}

func (uc *newsFindUsecase) Execute(
	ctx context.Context,
	req dto.NewsFindRequest,
) ([]dto.News, dto.PagingResponse, error) {
	total, err := uc.newsRepo.Count(
		ctx,
		specifications.NewsByDate(
			req.From,
			req.To,
			req.Status,
		),
	)
	if err != nil {
		return nil, dto.PagingResponse{}, fmt.Errorf("count news by date range: %w", err)
	}

	if total == 0 {
		return []dto.News{}, dto.PagingResponse{Total: 0}, nil
	}

	news, err := uc.newsRepo.Find(
		ctx,
		specifications.NewsByDate(
			req.From,
			req.To,
			req.Status,
		),
		specifications.ToPaging(req.Paging.Size, req.Paging.Page),
	)
	if err != nil {
		return nil, dto.PagingResponse{}, fmt.Errorf("find news by date range: %w", err)
	}

	if len(news) == 0 {
		return []dto.News{}, dto.PagingResponse{}, nil
	}

	var (
		fileIDs = make([]uint64, len(news))
		dtoNews = make([]dto.News, len(news))
	)

	for i := range news {
		fileIDs[i] = news[i].NewsWithFullTextID
	}

	newsWithFullText, err := uc.newsWithFullTextRepo.Find(
		ctx,
		specifications.NewNewsWithFullTextByFileIDs(fileIDs, 256),
		specifications.ToPaging(len(fileIDs), 1),
	)
	if err != nil {
		return nil, dto.PagingResponse{}, fmt.Errorf("find news with full text by file IDs: %w", err)
	}

	mNewsWithFullTextByOriginalID := lo.SliceToMap(newsWithFullText, func(n entity.NewsWithFullText) (uint64, string) {
		return n.ID, n.Content
	})

	for i := range news {
		news[i].Content = mNewsWithFullTextByOriginalID[news[i].NewsWithFullTextID]

		dtoNews[i] = dto.ToDtoNews(news[i])
	}

	return dtoNews, dto.PagingResponse{Total: int(total)}, nil
}
