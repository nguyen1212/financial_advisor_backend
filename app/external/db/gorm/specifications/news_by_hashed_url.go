package specifications

import (
	"github.com/financial_advisor/app/domain/entity"
	"github.com/financial_advisor/app/domain/repository/specifications"
	"gorm.io/gorm"
)

type newsByHashedURL struct {
	hashedURL []byte
}

func NewNewsByHashedURL(hashedURL []byte) specifications.I {
	return newsByHashedURL{hashedURL: hashedURL}
}

func (q newsByHashedURL) Query(db *gorm.DB) *gorm.DB {
	tx := db.Model(&entity.News{}).Where("hashed_url = ?", q.hashedURL)

	return tx
}
