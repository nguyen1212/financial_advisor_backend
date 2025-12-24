// Package specifications define implementation leveraging gorm library
package specifications

import (
	"github.com/financial_advisor/app/domain/entity"
	"github.com/financial_advisor/app/domain/repository/specifications"
	"gorm.io/gorm"
)

type publisherByID struct {
	id uint64
}

func NewPublisherByID(id uint64) specifications.I {
	return &publisherByID{id: id}
}

func (q *publisherByID) Query(db *gorm.DB) *gorm.DB {
	return db.Model(&entity.Publisher{}).Where("id = ?", q.id)
}
