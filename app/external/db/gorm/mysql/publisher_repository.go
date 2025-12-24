// Package mysql defines mysql gorm repository implementations
package mysql

import (
	"context"
	"errors"

	"github.com/financial_advisor/app/domain/entity"
	"github.com/financial_advisor/app/domain/repository"
	"github.com/financial_advisor/app/domain/repository/specifications"
	appError "github.com/financial_advisor/app/errors"
	"gorm.io/gorm"
)

type publisherRepository struct {
	db *gorm.DB
}

func NewPublisherRepository(db *gorm.DB) repository.PublisherRepository {
	return &publisherRepository{db: db}
}

func (r *publisherRepository) Count(
	ctx context.Context,
	spec specifications.I,
) (int64, error) {
	var (
		count int64
		tx    = spec.Query(r.db)
	)

	if err := tx.WithContext(ctx).Count(&count).Error; err != nil {
		return 0, err
	}

	return count, nil
}

func (r *publisherRepository) Find(
	ctx context.Context,
	spec specifications.I,
	paging specifications.PagingI,
) ([]entity.Publisher, error) {
	var (
		publisherList []entity.Publisher
		tx            = spec.Query(r.db).WithContext(ctx).
				Limit(paging.Limit()).
				Offset(paging.Offset())
	)

	if err := tx.Find(&publisherList).Error; err != nil {
		return nil, err
	}

	return publisherList, nil
}

func (r *publisherRepository) Get(
	ctx context.Context,
	spec specifications.I,
) (entity.Publisher, error) {
	var (
		publisher entity.Publisher
		tx        = spec.Query(r.db)
	)

	if err := tx.WithContext(ctx).Take(&publisher).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return entity.Publisher{}, appError.ErrNotFound
		}

		return entity.Publisher{}, err
	}

	return publisher, nil
}

func (r *publisherRepository) Create(
	ctx context.Context,
	publisher *entity.Publisher,
) error {
	if err := r.db.WithContext(ctx).Create(publisher).Error; err != nil {
		return err
	}

	return nil
}
