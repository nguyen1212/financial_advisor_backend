// Package worker initializes and manages the worker service for processing tasks.
package worker

import (
	"context"

	"github.com/financial_advisor/app/config"
	"github.com/financial_advisor/app/external/db/gorm"
	"github.com/financial_advisor/app/external/db/gorm/mysql"
	memoryqueue "github.com/financial_advisor/app/external/queue/memory-queue"
	"github.com/financial_advisor/app/external/uuid"
	"github.com/financial_advisor/app/registry"
	workerService "github.com/financial_advisor/app/services/worker"
	"github.com/financial_advisor/cmd/internal/shutdown"
	"github.com/sirupsen/logrus"
)

type task struct {
	worker workerService.I
}

func Init(ctx context.Context) error {
	workerService.InitJobWorker(
		ctx,
		memoryqueue.New(config.Get().WorkerConcurrency),
		registry.InjectConsumerService(),
		mysql.NewJobsRepository(gorm.GetMySQLIns()),
		uuid.NewUUIDv7(),
		config.Get().WorkerConcurrency,
	)

	logrus.WithField("worker_concurrency", config.Get().WorkerConcurrency).
		WithField("queue_type", "in-memory").
		WithField("worker_type", "stateful").
		Infoln("worker initialized")

	shutdown.Get().Add(&task{
		worker: workerService.Get(),
	})

	return nil
}

func (t *task) Name() string {
	return "ephemeral_worker"
}

func (t *task) Shutdown() error {
	return t.worker.Close()
}
