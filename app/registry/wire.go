// Package registry defines dependency injection
package registry

import (
	"github.com/financial_advisor/app/external/db/gorm"
	"github.com/financial_advisor/app/external/db/gorm/mysql"
	"github.com/financial_advisor/app/usecases"
)

func InjectNewsFindUsecase() usecases.NewsFindUsecase {
	db := gorm.Get()

	return usecases.NewNewsFindUsecase(
		mysql.NewNewsRepository(db),
	)
}
