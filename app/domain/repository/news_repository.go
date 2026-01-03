// Package repository presents database interaction
package repository

//go:generate mockgen -destination=./mock/mock_$GOFILE -source=$GOFILE -package=mock

import (
	"context"

	"github.com/financial_advisor/app/domain/entity"
	"github.com/financial_advisor/app/domain/repository/specifications"
)

type NewsRepository interface {
	Count(ctx context.Context, spec specifications.I) (int64, error)
	Find(
		ctx context.Context,
		spec specifications.I,
		paging specifications.PagingI,
	) ([]entity.News, error)
	Get(ctx context.Context, spec specifications.I) (entity.News, error)
	Create(ctx context.Context, news *entity.News) error
	Update(ctx context.Context, news *entity.News) error
	Delete(ctx context.Context, newsID uint64) error
}

type NewsWithFullTextRepository interface {
	Find(
		ctx context.Context,
		spec specifications.I,
		paging specifications.PagingI,
	) ([]entity.NewsWithFullText, error)
	FindSearchSuggestions(
		ctx context.Context,
		specs specifications.I,
		paging specifications.PagingI,
	) ([]string, error)
	Get(ctx context.Context, spec specifications.I) (entity.NewsWithFullText, error)
	Create(ctx context.Context, news *entity.NewsWithFullText) error
	Delete(ctx context.Context, id uint64) error
}
