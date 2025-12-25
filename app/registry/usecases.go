// Package registry defines dependency injection
package registry

import (
	"github.com/financial_advisor/app/external/db/gorm"
	"github.com/financial_advisor/app/external/db/gorm/mysql"
	"github.com/financial_advisor/app/external/hasher"
	memoryqueue "github.com/financial_advisor/app/external/queue/memory-queue"
	"github.com/financial_advisor/app/usecases"
)

func InjectNewsFindUsecase() usecases.NewsFindUsecase {
	db := gorm.Get()

	return usecases.NewNewsFindUsecase(
		mysql.NewNewsRepository(db),
	)
}

func InjectPublishersFindUsecase() usecases.PublishersFindUsecase {
	db := gorm.Get()

	return usecases.NewPublishersFindUsecase(
		mysql.NewPublisherRepository(db),
	)
}

func InjectPublisherCreateUsecase() usecases.PublisherCreateUsecase {
	db := gorm.Get()

	return usecases.NewPublisherCreateUsecase(
		mysql.NewPublisherRepository(db),
	)
}

func InjectVnExpressScrapperUsecase() usecases.WebScrapperUsecase {
	db := gorm.Get()

	return usecases.NewVnExpressScrapperUsecase(mysql.NewNewsRepository(db))
}

func InjectNewsCreateUsecase() usecases.NewsCreateUsecase {
	db := gorm.Get()
	memqueue := memoryqueue.Get()
	hasher := hasher.NewMD5()

	return usecases.NewNewsCreateUsecase(
		mysql.NewNewsRepository(db),
		mysql.NewPublisherRepository(db),
		memqueue,
		hasher,
	)
}

func InjectNewsGetUsecase() usecases.NewsGetUsecase {
	db := gorm.Get()

	return usecases.NewNewsGetUsecase(
		mysql.NewNewsRepository(db),
	)
}

func InjectNewsDeleteUsecase() usecases.NewsDeleteUsecase {
	db := gorm.Get()

	return usecases.NewNewsDeleteUsecase(
		mysql.NewNewsRepository(db),
	)
}

func InjectPublisherGetUsecase() usecases.PublisherGetUsecase {
	db := gorm.Get()

	return usecases.NewPublisherGetUsecase(
		mysql.NewPublisherRepository(db),
	)
}

func InjectFallbackScrapperUsecase() usecases.FallbackScrapperUsecase {
	db := gorm.Get()

	return usecases.NewFallbackScrapperUsecase(
		mysql.NewNewsRepository(db),
	)
}
