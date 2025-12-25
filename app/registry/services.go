package registry

import (
	"github.com/financial_advisor/app/domain/entity"
	"github.com/financial_advisor/app/services/consumer"
	"github.com/financial_advisor/app/services/consumer/processor"
	"github.com/financial_advisor/app/services/queue"
	"github.com/financial_advisor/app/usecases"
)

func InjectEphemeralConsumerService() consumer.Manager {
	var (
		mUcs = map[entity.WebDomain]usecases.WebScrapperUsecase{
			entity.WebDomainVnExpress: InjectVnExpressScrapperUsecase(),
		}
		fallbackUc = InjectFallbackScrapperUsecase()

		mProcessors = map[queue.MessageType]consumer.Processor{
			queue.MessageTypeWebScrapper: processor.NewWebScrapperProcessor(mUcs, fallbackUc),
		}

		consumerManager = consumer.NewEphemeralManager(
			mProcessors,
		)
	)

	return consumerManager
}
