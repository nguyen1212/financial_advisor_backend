package worker

import (
	"context"

	"github.com/financial_advisor/app/config"
	memoryqueue "github.com/financial_advisor/app/external/queue/memory-queue"
	"github.com/financial_advisor/app/registry"
	"github.com/financial_advisor/app/services/queue"
	"github.com/financial_advisor/cmd/internal/shutdown"
	"github.com/sirupsen/logrus"
)

type task struct {
	inMemQueue queue.I
}

func Init(ctx context.Context) error {
	memoryqueue.Init(
		config.Get().WorkerConcurrency,
		registry.InjectEphemeralConsumerService(),
	)

	logrus.WithField("worker_concurrency", config.Get().WorkerConcurrency).
		WithField("queue_type", "in-memory").
		Infoln("worker initialized")

	shutdown.Get().Add(&task{
		inMemQueue: memoryqueue.Get(),
	})

	return nil
}

func (t *task) Name() string {
	return "in-memory-queue-worker"
}

func (t *task) Shutdown() error {
	return t.inMemQueue.Close()
}
