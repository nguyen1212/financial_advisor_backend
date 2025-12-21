// Package specifications define implementation leveraging gorm library
package specifications

import (
	"time"

	"github.com/financial_advisor/app/domain/entity"
	"github.com/financial_advisor/app/domain/repository/specifications"
	"gorm.io/gorm"
)

type newsByDate struct {
	startDate time.Time
	endDate   time.Time
}

func NewNewsByDate(startDate, endDate time.Time) specifications.I[gorm.DB] {
	return &newsByDate{
		startDate: startDate,
		endDate:   endDate,
	}
}

func (q *newsByDate) Query(db *gorm.DB) *gorm.DB {
	query := db.Model(&entity.News{})
	if !q.startDate.IsZero() {
		query = query.Where("created_at >= ?", q.startDate)
	}

	if !q.endDate.IsZero() {
		query = query.Where("created_at <= ?", q.endDate)
	}

	return query
}
