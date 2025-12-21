// Package mysql defines mysql gorm repository implementations
package mysql

import (
	"context"

	"github.com/financial_advisor/app/domain/entity"
	"github.com/financial_advisor/app/domain/repository"
	"github.com/financial_advisor/app/domain/repository/specifications"
	"gorm.io/gorm"
)

type newsRepository struct {
	db *gorm.DB
}

func NewNewsRepository(db *gorm.DB) repository.NewsRepository[gorm.DB] {
	return &newsRepository{db: db}
}

func (r *newsRepository) Find(
	ctx context.Context,
	spec specifications.I[gorm.DB],
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
