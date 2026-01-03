// Package registry defines dependency injection
package registry

import (
	"github.com/financial_advisor/app/usecases"
	"github.com/google/wire"
)

var (
	NewsFindUsecaseSet = wire.NewSet(
		singletonSet,
		repositorySet,
		usecases.NewNewsFindUsecase,
	)
	PublishersFindUsecaseSet = wire.NewSet(
		singletonSet,
		repositorySet,
		usecases.NewPublishersFindUsecase,
	)
	PublisherCreateUsecaseSet = wire.NewSet(
		singletonSet,
		repositorySet,
		usecases.NewPublisherCreateUsecase,
	)
	VnExpressScrapperUsecaseSet = wire.NewSet(
		singletonSet,
		repositorySet,
		usecases.NewVnExpressScrapperUsecase,
	)
	NewsCreateUsecaseSet = wire.NewSet(
		singletonSet,
		repositorySet,
		serviceSet,
		usecases.NewNewsCreateUsecase,
	)
	NewsGetUsecaseSet = wire.NewSet(
		singletonSet,
		repositorySet,
		usecases.NewNewsGetUsecase,
	)
	NewsDeleteUsecaseSet = wire.NewSet(
		singletonSet,
		repositorySet,
		usecases.NewNewsDeleteUsecase,
	)
	PublisherGetUsecaseSet = wire.NewSet(
		singletonSet,
		repositorySet,
		usecases.NewPublisherGetUsecase,
	)
	FallbackScrapperUsecaseSet = wire.NewSet(
		singletonSet,
		repositorySet,
		usecases.NewFallbackScrapperUsecase,
	)
	NewsSearchSuggestionsGetUsecaseSet = wire.NewSet(
		singletonSet,
		repositorySet,
		usecases.NewNewsSearchSuggestionsGetUsecase,
	)
	NewsSearchUsecaseSet = wire.NewSet(
		singletonSet,
		repositorySet,
		usecases.NewNewsSearchUsecase,
	)
)
