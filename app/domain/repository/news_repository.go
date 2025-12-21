// Package repository presents database interaction
package repository

import (
	"context"

	"github.com/financial_advisor/app/domain/entity"
	"github.com/financial_advisor/app/domain/repository/specifications"
	"gorm.io/gorm"
)

type NewsRepository[T gorm.DB] interface {
	Find(ctx context.Context, spec specifications.I[T]) ([]entity.News, error)
}
