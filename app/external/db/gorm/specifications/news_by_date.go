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
	status    *entity.NewsStatus
}

func NewNewsByDate(
	startDate,
	endDate time.Time,
	status *entity.NewsStatus,
) specifications.I {
	return &newsByDate{
		startDate: startDate,
		endDate:   endDate,
		status:    status,
	}
}

func (q *newsByDate) Query(db *gorm.DB) *gorm.DB {
	tx := db.Model(&entity.News{})
	if !q.startDate.IsZero() {
		tx = tx.Where("created_at >= ?", q.startDate)
	}

	if !q.endDate.IsZero() {
		tx = tx.Where("created_at <= ?", q.endDate)
	}

	if q.status != nil {
		tx = tx.Where("status = ?", *q.status)
	}

	return tx
}
