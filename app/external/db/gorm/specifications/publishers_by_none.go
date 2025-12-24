// Package specifications define implementation leveraging gorm library
package specifications

import (
	"github.com/financial_advisor/app/domain/entity"
	"github.com/financial_advisor/app/domain/repository/specifications"
	"gorm.io/gorm"
)

type publishersByNone struct{}

func NewPublishersByNone() specifications.I {
	return &publishersByNone{}
}

func (q *publishersByNone) Query(db *gorm.DB) *gorm.DB {
	return db.Model(&entity.Publisher{})
}
