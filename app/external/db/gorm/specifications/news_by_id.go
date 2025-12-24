package specifications

import (
	"github.com/financial_advisor/app/domain/entity"
	"github.com/financial_advisor/app/domain/repository/specifications"
	"gorm.io/gorm"
)

type newsByID struct {
	id         uint64
	preloaders []string
}

func NewNewsByID(id uint64, preloaders ...string) specifications.I {
	return newsByID{id: id, preloaders: preloaders}
}

func (q newsByID) Query(db *gorm.DB) *gorm.DB {
	tx := db.Model(&entity.News{}).Where("id = ?", q.id)

	for i := range q.preloaders {
		tx = tx.Preload(q.preloaders[i])
	}

	return tx
}
