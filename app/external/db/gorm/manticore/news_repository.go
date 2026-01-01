// Package manticore holds implementation of manticore database
package manticore

import (
	"context"
	"errors"
	"fmt"

	"github.com/doug-martin/goqu/v9"
	_ "github.com/doug-martin/goqu/v9/dialect/mysql"
	"github.com/financial_advisor/app/domain/entity"
	"github.com/financial_advisor/app/domain/repository"
	"github.com/financial_advisor/app/domain/repository/specifications"
	appErrors "github.com/financial_advisor/app/errors"
	gormAdapter "github.com/financial_advisor/app/external/db/gorm"
	"gorm.io/gorm"
)

type newsRepository struct {
	db *gorm.DB
}

func NewNewsRepository(db *gormAdapter.Manticore) repository.NewsWithFullTextRepository {
	return &newsRepository{db: db.DB()}
}

func (r *newsRepository) Find(
	ctx context.Context,
	spec specifications.I,
	paging specifications.PagingI,
) ([]entity.NewsWithFullText, error) {
	var (
		newsList []entity.NewsWithFullText
		tx       = r.db.WithContext(ctx)
	)

	sql, err := spec.ToFind(paging)
	if err != nil {
		return nil, fmt.Errorf("convert to raw find sql: %w", err)
	}

	if err := tx.WithContext(ctx).Raw(sql).Find(&newsList).Error; err != nil {
		return nil, err
	}

	return newsList, nil
}

func (r *newsRepository) FindSearchSuggestions(
	ctx context.Context,
	spec specifications.I,
	paging specifications.PagingI,
) ([]string, error) {
	var (
		tx          = r.db.WithContext(ctx)
		suggestions []string
	)

	sql, err := spec.ToFind(paging)
	if err != nil {
		return nil, err
	}

	if err = tx.Raw(sql).Find(&suggestions).Error; err != nil {
		return nil, err
	}

	return suggestions, nil
}

func (r *newsRepository) Get(
	ctx context.Context,
	spec specifications.I,
) (entity.NewsWithFullText, error) {
	var (
		tx   = r.db.WithContext(ctx)
		news entity.NewsWithFullText
	)

	sql, err := spec.ToGet()
	if err != nil {
		return entity.NewsWithFullText{}, err
	}

	if err = tx.Raw(sql).Take(&news).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return entity.NewsWithFullText{}, appErrors.ErrNotFound
		}
		return entity.NewsWithFullText{}, err
	}

	return news, nil
}

func (r *newsRepository) Create(ctx context.Context, news *entity.NewsWithFullText) error {
	tx := r.db.WithContext(ctx).Begin()
	defer tx.Rollback()

	sql, _, err := goqu.Dialect("mysql").
		Insert("news").
		Rows(goqu.Record{
			"news_id": news.NewsID,
			"title":   news.Title,
			"content": news.Content,
		}).
		ToSQL()
	if err != nil {
		return fmt.Errorf("failed to build insert sql: %w", err)
	}

	var id uint64

	if err := tx.Exec(sql).Error; err != nil {
		return err
	}

	if err := tx.Raw(`SELECT LAST_INSERT_ID()`).Scan(&id).Error; err != nil {
		return err
	}

	if err := tx.Commit().Error; err != nil {
		return err
	}

	news.ID = id

	return nil
}

func (r *newsRepository) Delete(
	ctx context.Context,
	id uint64,
) error {
	return r.db.WithContext(ctx).Delete(&entity.NewsWithFullText{ID: id}).Error
}
