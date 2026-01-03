// Package usecases define implementation of business logic
package usecases

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/financial_advisor/app/domain/entity"
	"github.com/financial_advisor/app/domain/repository"
	"github.com/financial_advisor/app/external/db/goqu/specifications"
	"github.com/financial_advisor/app/usecases/dto"
)

type NewsSearchUsecase interface {
	Execute(
		ctx context.Context,
		req dto.NewsSearchRequest,
	) ([]dto.News, error)
}

type newsSearchUsecase struct {
	newsRepo             repository.NewsRepository
	newsWithFullTextRepo repository.NewsWithFullTextRepository
}

func NewNewsSearchUsecase(
	newsRepo repository.NewsRepository,
	newsWithFullTextRepo repository.NewsWithFullTextRepository,
) NewsSearchUsecase {
	return &newsSearchUsecase{
		newsRepo:             newsRepo,
		newsWithFullTextRepo: newsWithFullTextRepo,
	}
}

func (uc *newsSearchUsecase) Execute(
	ctx context.Context,
	req dto.NewsSearchRequest,
) ([]dto.News, error) {
	if len(req.Keywords) == 0 || strings.TrimSpace(strings.Join(req.Keywords, "")) == "" {
		return []dto.News{}, nil
	}

	ftsOp := specifications.FullTextSearchOpProximity
	if len(req.Keywords) > 10 {
		ftsOp = specifications.FullTextSearchOpQuorum
	}

	newsWithFullText, err := uc.newsWithFullTextRepo.Find(
		ctx,
		specifications.NewNewsWithFullTextByKeywords(
			req.Keywords,
			256,
			ftsOp,
		),
		specifications.ToPaging(req.Paging.Size, req.Paging.Page),
	)
	if err != nil {
		return nil, fmt.Errorf("find news with full text by file IDs: %w", err)
	}

	// safety check only
	if len(newsWithFullText) == 0 {
		return []dto.News{}, nil
	}

	var (
		ids                   = make([]uint64, len(newsWithFullText))
		mNewsWithFullTextByID = make(map[uint64]*entity.NewsWithFullText)
		dtoNews               = make([]dto.News, len(newsWithFullText))
	)

	for i := range newsWithFullText {
		parsedID, err := strconv.ParseUint(newsWithFullText[i].NewsID, 10, 64)
		if err != nil {
			continue
		}

		ids[i] = parsedID
		mNewsWithFullTextByID[newsWithFullText[i].ID] = &newsWithFullText[i]
	}

	news, err := uc.newsRepo.Find(
		ctx,
		specifications.NewNewsByIDs(ids),
		specifications.ToPaging(len(ids), 1),
	)
	if err != nil {
		return nil, fmt.Errorf("find news by date range: %w", err)
	}

	for i := range news {
		newsPtr := mNewsWithFullTextByID[news[i].NewsWithFullTextID]

		if newsPtr != nil {
			news[i].Content = newsPtr.Content
		}

		dtoNews[i] = dto.ToDtoNews(news[i])
	}

	return dtoNews, nil
}
