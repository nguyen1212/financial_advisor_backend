//go:build wireinject
// +build wireinject

// The build tag makes sure the stub is not built in the final build.

package registry

import (
	"github.com/financial_advisor/app/usecases"
	"github.com/google/wire"
)

func InjectNewsFindUsecase() usecases.NewsFindUsecase {
	wire.Build(NewsFindUsecaseSet)

	return nil
}

func InjectPublishersFindUsecase() usecases.PublishersFindUsecase {
	wire.Build(PublishersFindUsecaseSet)

	return nil
}

func InjectPublisherCreateUsecase() usecases.PublisherCreateUsecase {
	wire.Build(PublisherCreateUsecaseSet)

	return nil
}

func InjectVnExpressScrapperUsecase() usecases.WebScrapperUsecase {
	wire.Build(VnExpressScrapperUsecaseSet)

	return nil
}

func InjectNewsCreateUsecase() usecases.NewsCreateUsecase {
	wire.Build(NewsCreateUsecaseSet)

	return nil
}

func InjectNewsGetUsecase() usecases.NewsGetUsecase {
	wire.Build(NewsGetUsecaseSet)

	return nil
}

func InjectNewsDeleteUsecase() usecases.NewsDeleteUsecase {
	wire.Build(NewsDeleteUsecaseSet)

	return nil
}

func InjectPublisherGetUsecase() usecases.PublisherGetUsecase {
	wire.Build(PublisherGetUsecaseSet)

	return nil
}

func InjectFallbackScrapperUsecase() usecases.FallbackScrapperUsecase {
	wire.Build(FallbackScrapperUsecaseSet)

	return nil
}

func InjectNewsSearchSuggestionGetUsecase() usecases.NewsSearchSuggestionsGetUsecase {
	wire.Build(NewsSearchSuggestionsGetUsecaseSet)

	return nil
}

func InjectNewsSearchUsecase() usecases.NewsSearchUsecase {
	wire.Build(NewsSearchUsecaseSet)

	return nil
}
