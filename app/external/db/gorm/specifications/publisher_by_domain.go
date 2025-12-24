package specifications

import (
	"github.com/financial_advisor/app/domain/entity"
	"github.com/financial_advisor/app/domain/repository/specifications"
	"gorm.io/gorm"
)

type publisherByDomain struct {
	domain string
}

func NewPublisherByDomain(domain string) specifications.I {
	return &publisherByDomain{
		domain: domain,
	}
}

func (q *publisherByDomain) Query(db *gorm.DB) *gorm.DB {
	return db.Model(&entity.Publisher{}).Where("domain = ?", q.domain)
}
