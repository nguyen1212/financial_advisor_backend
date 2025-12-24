// Package mysql defines mysql gorm repository implementations
package mysql

import (
	"context"
	"errors"

	"github.com/financial_advisor/app/domain/entity"
	"github.com/financial_advisor/app/domain/repository"
	"github.com/financial_advisor/app/domain/repository/specifications"
	appErrors "github.com/financial_advisor/app/errors"
	"gorm.io/gorm"
)

type newsRepository struct {
	db *gorm.DB
}

func NewNewsRepository(db *gorm.DB) repository.NewsRepository {
	return &newsRepository{db: db}
}

func (r *newsRepository) Count(
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

func (r *newsRepository) Find(
	ctx context.Context,
	spec specifications.I,
) ([]entity.News, error) {
	var (
		newsList []entity.News
		tx       = spec.Query(r.db)
	)

	if err := tx.WithContext(ctx).Find(&newsList).Error; err != nil {
		return nil, err
	}

	return newsList, nil
}

func (r *newsRepository) Get(
	ctx context.Context,
	spec specifications.I,
) (entity.News, error) {
	var (
		news entity.News
		tx   = spec.Query(r.db)
	)

	if err := tx.WithContext(ctx).Take(&news).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return entity.News{}, appErrors.ErrNotFound
		}
		return entity.News{}, err
	}

	return news, nil
}

func (r *newsRepository) Create(
	ctx context.Context,
	news *entity.News,
) error {
	if err := r.db.WithContext(ctx).Create(news).Error; err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return appErrors.ErrConflicted
		}

		return err
	}

	return nil
}

func (r *newsRepository) Update(
	ctx context.Context,
	news *entity.News,
) error {
	if err := r.db.WithContext(ctx).Model(news).Updates(news.ToMap()).Error; err != nil {
		return err
	}

	return nil
}

func (r *newsRepository) Delete(
	ctx context.Context,
	newsID uint64,
) error {
	if err := r.db.WithContext(ctx).Delete(&entity.News{ID: newsID}).Error; err != nil {
		return err
	}

	return nil
}
