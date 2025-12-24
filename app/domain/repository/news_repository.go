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
	Find(ctx context.Context, spec specifications.I) ([]entity.News, error)
	Get(ctx context.Context, spec specifications.I) (entity.News, error)
	Create(ctx context.Context, news *entity.News) error
	Update(ctx context.Context, news *entity.News) error
	Delete(ctx context.Context, newsID uint64) error
}
