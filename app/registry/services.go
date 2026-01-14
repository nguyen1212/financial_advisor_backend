package registry

import (
	"github.com/financial_advisor/app/domain/entity"
	"github.com/financial_advisor/app/services/consumer"
	"github.com/financial_advisor/app/services/consumer/processor"
	"github.com/financial_advisor/app/services/queue"
	"github.com/financial_advisor/app/usecases"
	"github.com/google/wire"
)

var ConsumerServiceSet = wire.NewSet(
	singletonSet,
	repositorySet,
	serviceSet,
	ProvideProcessors,
)

func ProvideProcessors() map[queue.MessageType]consumer.Processor {
	var (
		mUcs = map[entity.WebDomain]usecases.WebScrapperUsecase{
			entity.WebDomainVnExpress: InjectVnExpressScrapperUsecase(),
		}
		fallbackUc = InjectFallbackScrapperUsecase()

		mProcessors = map[queue.MessageType]consumer.Processor{
			queue.MessageTypeWebScrapper: processor.NewWebScrapperProcessor(mUcs, fallbackUc),
		}
	)

	return mProcessors
}
