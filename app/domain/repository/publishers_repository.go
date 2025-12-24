// Package repository presents database interaction
package repository

//go:generate mockgen -destination=./mock/mock_$GOFILE -source=$GOFILE -package=mock

import (
	"context"

	"github.com/financial_advisor/app/domain/entity"
	"github.com/financial_advisor/app/domain/repository/specifications"
)

type PublisherRepository interface {
	Count(
		ctx context.Context,
		spec specifications.I,
	) (int64, error)
	Find(
		ctx context.Context,
		spec specifications.I,
		paging specifications.PagingI,
	) ([]entity.Publisher, error)
	Get(
		ctx context.Context,
		spec specifications.I,
	) (entity.Publisher, error)
	Create(
		ctx context.Context,
		obj *entity.Publisher,
	) error
}
