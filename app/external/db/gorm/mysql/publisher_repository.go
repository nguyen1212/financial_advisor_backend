// Package mysql defines mysql gorm repository implementations
package mysql

import (
	"context"
	"errors"
	"fmt"

	"github.com/financial_advisor/app/domain/entity"
	"github.com/financial_advisor/app/domain/repository"
	"github.com/financial_advisor/app/domain/repository/specifications"
	appError "github.com/financial_advisor/app/errors"
	gormAdapter "github.com/financial_advisor/app/external/db/gorm"
	"gorm.io/gorm"
)

type publisherRepository struct {
	db *gorm.DB
}

func NewPublisherRepository(db *gormAdapter.MySQL) repository.PublisherRepository {
	return &publisherRepository{db: db.DB()}
}

func (r *publisherRepository) Count(
	ctx context.Context,
	spec specifications.I,
) (int64, error) {
	var (
		count int64
		tx    = r.db.WithContext(ctx)
	)

	sql, err := spec.ToCount()
	if err != nil {
		return 0, fmt.Errorf("convert to raw count sql: %w", err)
	}

	if err := tx.Raw(sql).Count(&count).Error; err != nil {
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
		tx            = r.db.WithContext(ctx)
	)

	sql, err := spec.ToFind(paging)
	if err != nil {
		return nil, fmt.Errorf("convert to raw find sql: %w", err)
	}

	if err := tx.Raw(sql).Find(&publisherList).Error; err != nil {
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
		tx        = r.db.WithContext(ctx)
	)

	sql, err := spec.ToGet()
	if err != nil {
		return entity.Publisher{}, fmt.Errorf("convert to raw get sql: %w", err)
	}

	if err := tx.Raw(sql).Take(&publisher).Error; err != nil {
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
