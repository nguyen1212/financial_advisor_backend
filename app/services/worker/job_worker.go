// Package worker represents the logic with combination of queue and consumer
package worker

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"sync"

	"github.com/financial_advisor/app/domain/entity"
	"github.com/financial_advisor/app/domain/repository"
	"github.com/financial_advisor/app/external/db/goqu/specifications"
	"github.com/financial_advisor/app/services/consumer"
	"github.com/financial_advisor/app/services/queue"
	"github.com/financial_advisor/app/services/uuid"
	"github.com/sirupsen/logrus"
)

var (
	once     sync.Once
	instance *jobWorker
)

type jobWorker struct {
	queue         queue.I
	consumer      consumer.I
	uuidGenerator uuid.UUIDGenerator
	jobRepo       repository.JobsRepository
}

// InitJobWorker initializes the queue and run it in background
func InitJobWorker(
	ctx context.Context,
	queue queue.I,
	consumer consumer.I,
	jobRepo repository.JobsRepository,
	uuidGenerator uuid.UUIDGenerator,
	workersNo int,
) {
	once.Do(func() {
		instance = &jobWorker{
			queue:         queue,
			consumer:      consumer,
			jobRepo:       jobRepo,
			uuidGenerator: uuidGenerator,
		}

		for range workersNo {
			// this worker should run in the background with same lifetime as the application
			go instance.start(ctx)
		}
	})
}

func (w *jobWorker) start(ctx context.Context) {
	defer func() {
		if r := recover(); r != nil {
			logrus.Errorf("worker recovered from panic: %v", r)
		}
	}()

	for {
		// check for context cancellation (graceful shutdown)
		select {
		case <-ctx.Done():
			logrus.Info("worker context done, exiting worker")

			return
		default:
		}

		msg, isClosed := w.queue.Dequeue()
		if !isClosed {
			break
		}

		jobErr := w.consumer.Execute(ctx, msg)

		job, err := w.jobRepo.Get(
			ctx,
			specifications.JobByUUID(msg.UUID),
		)
		if err != nil {
			logrus.Errorf("get job by uuid %s: %v", msg.UUID, err)

			continue
		}

		job.Status = entity.JobStatusCompleted
		if jobErr != nil {
			job.Status = entity.JobStatusFailed
			job.Result.Error = jobErr.Error()

			resultEnc, err := json.Marshal(job.Result)
			if err != nil {
				logrus.Errorf("marshal job result: %v", err)
			} else {
				job.ResultEnc = resultEnc
			}
		}

		if err := w.jobRepo.Update(ctx, &job); err != nil {
			logrus.Errorf("update job status for uuid %s: %v", msg.UUID, err)
		}
	}
}

// Get currently only one singleton worker instance
// we can extend to multiple types of worker by removing the sync.Once pattern
func Get() I {
	return instance
}

func (w *jobWorker) Run(ctx context.Context, job []byte) error {
	var parsedMsg queue.Message

	if err := json.Unmarshal(job, &parsedMsg); err != nil {
		return fmt.Errorf("message is incorrect format: %w", err)
	}

	uuid, _, err := w.uuidGenerator.GetUUID()
	if err != nil {
		return fmt.Errorf("generate uuid: %w", err)
	}

	jobEnt := entity.Job{
		Type:    entity.ToJobType(parsedMsg.Type.String()),
		Status:  entity.JobStatusNew,
		Payload: parsedMsg.Body,
		UUID:    uuid,
	}

	if err := w.jobRepo.Create(ctx, &jobEnt); err != nil {
		return fmt.Errorf("create job: %w", err)
	}

	parsedMsg.UUID = jobEnt.UUID

	return w.queue.Enqueue(parsedMsg)
}

func (w *jobWorker) Close() error {
	var closeErr error

	if err := w.queue.Close(); err != nil {
		closeErr = err
	}

	if err := w.consumer.Close(); err != nil {
		closeErr = errors.Join(closeErr, err)
	}

	return closeErr
}
