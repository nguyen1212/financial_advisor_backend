package registry

import (
	"github.com/financial_advisor/app/domain/entity"
	ephemeralConsumer "github.com/financial_advisor/app/external/consumer"
	"github.com/financial_advisor/app/external/consumer/processor"
	"github.com/financial_advisor/app/services/consumer"
	"github.com/financial_advisor/app/services/queue"
	"github.com/financial_advisor/app/usecases"
)

// func InjectEphemeralConsumer() consumer.I {
// 	var (
// 		mUcs = map[entity.WebDomain]usecases.WebScrapperUsecase{
// 			entity.WebDomainVnExpress: InjectVnExpressScrapperUsecase(),
// 		}
// 		fallbackUc = InjectFallbackScrapperUsecase()
//
// 		mProcessors = map[queue.MessageType]consumer.Processor{
// 			queue.MessageTypeWebScrapper: processor.NewWebScrapperProcessor(mUcs, fallbackUc),
// 		}
//
// 		consumerManager = consumer.NewEphemeralManager(
// 			mProcessors,
// 		)
// 	)
//
// 	return consumerManager
// }

func InjectEphemeralConsumerService() consumer.I {
	var (
		mUcs = map[entity.WebDomain]usecases.WebScrapperUsecase{
			entity.WebDomainVnExpress: InjectVnExpressScrapperUsecase(),
		}
		fallbackUc = InjectFallbackScrapperUsecase()

		mProcessors = map[queue.MessageType]consumer.Processor{
			queue.MessageTypeWebScrapper: processor.NewWebScrapperProcessor(mUcs, fallbackUc),
		}

		consumer = ephemeralConsumer.NewEphemeralManager(
			mProcessors,
		)
	)

	return consumer
}
