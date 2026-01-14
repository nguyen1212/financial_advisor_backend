package worker

import (
	"context"
	"encoding/json"
	"errors"
	"testing"

	"github.com/financial_advisor/app/domain/entity"
	mockRepo "github.com/financial_advisor/app/domain/repository/mock"
	"github.com/financial_advisor/app/domain/repository/specifications"
	goquSpec "github.com/financial_advisor/app/external/db/goqu/specifications"
	mockConsumer "github.com/financial_advisor/app/services/consumer/mock"
	queuePkg "github.com/financial_advisor/app/services/queue"
	mockQueue "github.com/financial_advisor/app/services/queue/mock"
	"github.com/financial_advisor/app/services/uuid/mock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func Test_jobWorker_start(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		var (
			ctx  = context.Background()
			msg1 = queuePkg.Message{
				Type: queuePkg.MessageTypeWebScrapper,
				UUID: "uuid",
			}
			job = entity.Job{
				UUID:   msg1.UUID,
				Status: entity.JobStatusFailed,
			}
			msg2 = msg1

			queue         = mockQueue.NewMockI(mockCtrl)
			consumer      = mockConsumer.NewMockI(mockCtrl)
			jobRepo       = mockRepo.NewMockJobsRepository(mockCtrl)
			uuidGenerator = mock.NewMockUUIDGenerator(mockCtrl)

			worker = &jobWorker{
				queue:         queue,
				consumer:      consumer,
				jobRepo:       jobRepo,
				uuidGenerator: uuidGenerator,
			}
		)

		queue.EXPECT().Dequeue().Return(msg1, true)

		consumer.EXPECT().Execute(ctx, msg1).Return(nil)

		jobRepo.EXPECT().Get(
			ctx,
			specifications.CustomMatcher(specifications.SpecMatcher(goquSpec.JobByUUID(msg1.UUID))),
		).Return(job, nil)

		job.Status = entity.JobStatusCompleted
		jobRepo.EXPECT().Update(ctx, &job).Return(nil)

		queue.EXPECT().Dequeue().Return(msg2, false)

		worker.start(ctx)
	})

	t.Run("consumer execute fails", func(t *testing.T) {
		t.Parallel()

		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		var (
			ctx = context.Background()
			msg = queuePkg.Message{
				Type: queuePkg.MessageTypeWebScrapper,
				UUID: "uuid-failed",
			}
			job = entity.Job{
				UUID:   msg.UUID,
				Status: entity.JobStatusProcessing,
			}
			consumerErr = errors.New("consumer execution failed")

			queue         = mockQueue.NewMockI(mockCtrl)
			consumer      = mockConsumer.NewMockI(mockCtrl)
			jobRepo       = mockRepo.NewMockJobsRepository(mockCtrl)
			uuidGenerator = mock.NewMockUUIDGenerator(mockCtrl)

			worker = &jobWorker{
				queue:         queue,
				consumer:      consumer,
				jobRepo:       jobRepo,
				uuidGenerator: uuidGenerator,
			}
		)

		queue.EXPECT().Dequeue().Return(msg, true)

		consumer.EXPECT().Execute(ctx, msg).Return(consumerErr)

		jobRepo.EXPECT().Get(
			ctx,
			specifications.CustomMatcher(specifications.SpecMatcher(goquSpec.JobByUUID(msg.UUID))),
		).Return(job, nil)

		job.Status = entity.JobStatusFailed
		job.Result.Error = consumerErr.Error()
		resultEnc, _ := json.Marshal(job.Result)
		job.ResultEnc = resultEnc

		jobRepo.EXPECT().Update(ctx, &job).Return(nil)

		var emptyMsg queuePkg.Message
		queue.EXPECT().Dequeue().Return(emptyMsg, false)

		worker.start(ctx)
	})

	t.Run("job repository get fails", func(t *testing.T) {
		t.Parallel()

		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		var (
			ctx = context.Background()
			msg = queuePkg.Message{
				Type: queuePkg.MessageTypeWebScrapper,
				UUID: "uuid-get-error",
			}
			getErr = errors.New("failed to get job")

			queue         = mockQueue.NewMockI(mockCtrl)
			consumer      = mockConsumer.NewMockI(mockCtrl)
			jobRepo       = mockRepo.NewMockJobsRepository(mockCtrl)
			uuidGenerator = mock.NewMockUUIDGenerator(mockCtrl)

			worker = &jobWorker{
				queue:         queue,
				consumer:      consumer,
				jobRepo:       jobRepo,
				uuidGenerator: uuidGenerator,
			}
		)

		queue.EXPECT().Dequeue().Return(msg, true)

		consumer.EXPECT().Execute(ctx, msg).Return(nil)

		jobRepo.EXPECT().Get(
			ctx,
			specifications.CustomMatcher(specifications.SpecMatcher(goquSpec.JobByUUID(msg.UUID))),
		).Return(entity.Job{}, getErr)

		var emptyMsg queuePkg.Message
		queue.EXPECT().Dequeue().Return(emptyMsg, false)

		worker.start(ctx)
	})

	t.Run("job repository update fails", func(t *testing.T) {
		t.Parallel()

		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		var (
			ctx = context.Background()
			msg = queuePkg.Message{
				Type: queuePkg.MessageTypeWebScrapper,
				UUID: "uuid-update-error",
			}
			job = entity.Job{
				UUID:   msg.UUID,
				Status: entity.JobStatusProcessing,
			}
			updateErr = errors.New("failed to update job")

			queue         = mockQueue.NewMockI(mockCtrl)
			consumer      = mockConsumer.NewMockI(mockCtrl)
			jobRepo       = mockRepo.NewMockJobsRepository(mockCtrl)
			uuidGenerator = mock.NewMockUUIDGenerator(mockCtrl)

			worker = &jobWorker{
				queue:         queue,
				consumer:      consumer,
				jobRepo:       jobRepo,
				uuidGenerator: uuidGenerator,
			}
		)

		queue.EXPECT().Dequeue().Return(msg, true)

		consumer.EXPECT().Execute(ctx, msg).Return(nil)

		jobRepo.EXPECT().Get(
			ctx,
			specifications.CustomMatcher(specifications.SpecMatcher(goquSpec.JobByUUID(msg.UUID))),
		).Return(job, nil)

		job.Status = entity.JobStatusCompleted
		jobRepo.EXPECT().Update(ctx, &job).Return(updateErr)

		var emptyMsg queuePkg.Message
		queue.EXPECT().Dequeue().Return(emptyMsg, false)

		worker.start(ctx)
	})

	t.Run("context cancelled", func(t *testing.T) {
		t.Parallel()

		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		ctx, cancel := context.WithCancel(context.Background())
		cancel() // Cancel context immediately

		var (
			queue         = mockQueue.NewMockI(mockCtrl)
			consumer      = mockConsumer.NewMockI(mockCtrl)
			jobRepo       = mockRepo.NewMockJobsRepository(mockCtrl)
			uuidGenerator = mock.NewMockUUIDGenerator(mockCtrl)

			worker = &jobWorker{
				queue:         queue,
				consumer:      consumer,
				jobRepo:       jobRepo,
				uuidGenerator: uuidGenerator,
			}
		)

		worker.start(ctx)
	})
}

func Test_jobWorker_Run(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		var (
			ctx = context.Background()
			msg = queuePkg.Message{
				Type: queuePkg.MessageTypeWebScrapper,
				Body: []byte(`{"test": "data"}`),
			}
			msgBytes, _ = json.Marshal(msg)
			uuid        = "generated-uuid"

			queue         = mockQueue.NewMockI(mockCtrl)
			consumer      = mockConsumer.NewMockI(mockCtrl)
			jobRepo       = mockRepo.NewMockJobsRepository(mockCtrl)
			uuidGenerator = mock.NewMockUUIDGenerator(mockCtrl)

			worker = &jobWorker{
				queue:         queue,
				consumer:      consumer,
				jobRepo:       jobRepo,
				uuidGenerator: uuidGenerator,
			}
		)

		var uuidBytes [16]byte
		uuidGenerator.EXPECT().GetUUID().Return(uuid, uuidBytes, nil)

		expectedJob := &entity.Job{
			UUID:    uuid,
			Status:  entity.JobStatusNew,
			Type:    entity.ToJobType(msg.Type.String()),
			Payload: msg.Body,
		}
		jobRepo.EXPECT().Create(ctx, expectedJob).Return(nil)

		expectedMsg := queuePkg.Message{
			UUID: uuid,
			Type: msg.Type,
			Body: msg.Body,
		}
		queue.EXPECT().Enqueue(expectedMsg).Return(nil)

		err := worker.Run(ctx, msgBytes)
		require.NoError(t, err)
	})

	t.Run("invalid json", func(t *testing.T) {
		t.Parallel()

		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		var (
			ctx           = context.Background()
			invalidJSON   = []byte(`{"invalid json`)
			queue         = mockQueue.NewMockI(mockCtrl)
			consumer      = mockConsumer.NewMockI(mockCtrl)
			jobRepo       = mockRepo.NewMockJobsRepository(mockCtrl)
			uuidGenerator = mock.NewMockUUIDGenerator(mockCtrl)

			worker = &jobWorker{
				queue:         queue,
				consumer:      consumer,
				jobRepo:       jobRepo,
				uuidGenerator: uuidGenerator,
			}
		)

		err := worker.Run(ctx, invalidJSON)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "message is incorrect format")
	})

	t.Run("uuid generation fails", func(t *testing.T) {
		t.Parallel()

		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		var (
			ctx = context.Background()
			msg = queuePkg.Message{
				Type: queuePkg.MessageTypeWebScrapper,
				Body: []byte(`{"test": "data"}`),
			}
			msgBytes, _ = json.Marshal(msg)
			uuidErr     = errors.New("uuid generation failed")

			queue         = mockQueue.NewMockI(mockCtrl)
			consumer      = mockConsumer.NewMockI(mockCtrl)
			jobRepo       = mockRepo.NewMockJobsRepository(mockCtrl)
			uuidGenerator = mock.NewMockUUIDGenerator(mockCtrl)

			worker = &jobWorker{
				queue:         queue,
				consumer:      consumer,
				jobRepo:       jobRepo,
				uuidGenerator: uuidGenerator,
			}
		)

		var uuidBytes [16]byte
		uuidGenerator.EXPECT().GetUUID().Return("", uuidBytes, uuidErr)

		err := worker.Run(ctx, msgBytes)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "generate uuid")
	})

	t.Run("job creation fails", func(t *testing.T) {
		t.Parallel()

		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		var (
			ctx = context.Background()
			msg = queuePkg.Message{
				Type: queuePkg.MessageTypeWebScrapper,
				Body: []byte(`{"test": "data"}`),
			}
			msgBytes, _ = json.Marshal(msg)
			uuid        = "generated-uuid"
			createErr   = errors.New("failed to create job")

			queue         = mockQueue.NewMockI(mockCtrl)
			consumer      = mockConsumer.NewMockI(mockCtrl)
			jobRepo       = mockRepo.NewMockJobsRepository(mockCtrl)
			uuidGenerator = mock.NewMockUUIDGenerator(mockCtrl)

			worker = &jobWorker{
				queue:         queue,
				consumer:      consumer,
				jobRepo:       jobRepo,
				uuidGenerator: uuidGenerator,
			}
		)

		var uuidBytes [16]byte
		uuidGenerator.EXPECT().GetUUID().Return(uuid, uuidBytes, nil)

		expectedJob := &entity.Job{
			UUID:    uuid,
			Status:  entity.JobStatusNew,
			Type:    entity.ToJobType(msg.Type.String()),
			Payload: msg.Body,
		}
		jobRepo.EXPECT().Create(ctx, expectedJob).Return(createErr)

		err := worker.Run(ctx, msgBytes)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "create job")
	})

	t.Run("enqueue fails", func(t *testing.T) {
		t.Parallel()

		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		var (
			ctx = context.Background()
			msg = queuePkg.Message{
				Type: queuePkg.MessageTypeWebScrapper,
				Body: []byte(`{"test": "data"}`),
			}
			msgBytes, _  = json.Marshal(msg)
			uuid         = "generated-uuid"
			enqueueErr   = errors.New("failed to enqueue")

			queue         = mockQueue.NewMockI(mockCtrl)
			consumer      = mockConsumer.NewMockI(mockCtrl)
			jobRepo       = mockRepo.NewMockJobsRepository(mockCtrl)
			uuidGenerator = mock.NewMockUUIDGenerator(mockCtrl)

			worker = &jobWorker{
				queue:         queue,
				consumer:      consumer,
				jobRepo:       jobRepo,
				uuidGenerator: uuidGenerator,
			}
		)

		var uuidBytes [16]byte
		uuidGenerator.EXPECT().GetUUID().Return(uuid, uuidBytes, nil)

		expectedJob := &entity.Job{
			UUID:    uuid,
			Status:  entity.JobStatusNew,
			Type:    entity.ToJobType(msg.Type.String()),
			Payload: msg.Body,
		}
		jobRepo.EXPECT().Create(ctx, expectedJob).Return(nil)

		expectedMsg := queuePkg.Message{
			UUID: uuid,
			Type: msg.Type,
			Body: msg.Body,
		}
		queue.EXPECT().Enqueue(expectedMsg).Return(enqueueErr)

		err := worker.Run(ctx, msgBytes)
		require.Error(t, err)
		assert.Equal(t, enqueueErr, err)
	})
}

func Test_jobWorker_Close(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		var (
			queue         = mockQueue.NewMockI(mockCtrl)
			consumer      = mockConsumer.NewMockI(mockCtrl)
			jobRepo       = mockRepo.NewMockJobsRepository(mockCtrl)
			uuidGenerator = mock.NewMockUUIDGenerator(mockCtrl)

			worker = &jobWorker{
				queue:         queue,
				consumer:      consumer,
				jobRepo:       jobRepo,
				uuidGenerator: uuidGenerator,
			}
		)

		queue.EXPECT().Close().Return(nil)
		consumer.EXPECT().Close().Return(nil)

		err := worker.Close()
		require.NoError(t, err)
	})

	t.Run("queue close fails", func(t *testing.T) {
		t.Parallel()

		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		var (
			queueErr      = errors.New("queue close failed")
			queue         = mockQueue.NewMockI(mockCtrl)
			consumer      = mockConsumer.NewMockI(mockCtrl)
			jobRepo       = mockRepo.NewMockJobsRepository(mockCtrl)
			uuidGenerator = mock.NewMockUUIDGenerator(mockCtrl)

			worker = &jobWorker{
				queue:         queue,
				consumer:      consumer,
				jobRepo:       jobRepo,
				uuidGenerator: uuidGenerator,
			}
		)

		queue.EXPECT().Close().Return(queueErr)
		consumer.EXPECT().Close().Return(nil)

		err := worker.Close()
		require.Error(t, err)
		assert.Contains(t, err.Error(), "queue close failed")
	})

	t.Run("consumer close fails", func(t *testing.T) {
		t.Parallel()

		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		var (
			consumerErr   = errors.New("consumer close failed")
			queue         = mockQueue.NewMockI(mockCtrl)
			consumer      = mockConsumer.NewMockI(mockCtrl)
			jobRepo       = mockRepo.NewMockJobsRepository(mockCtrl)
			uuidGenerator = mock.NewMockUUIDGenerator(mockCtrl)

			worker = &jobWorker{
				queue:         queue,
				consumer:      consumer,
				jobRepo:       jobRepo,
				uuidGenerator: uuidGenerator,
			}
		)

		queue.EXPECT().Close().Return(nil)
		consumer.EXPECT().Close().Return(consumerErr)

		err := worker.Close()
		require.Error(t, err)
		assert.Contains(t, err.Error(), "consumer close failed")
	})

	t.Run("both close fail", func(t *testing.T) {
		t.Parallel()

		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		var (
			queueErr      = errors.New("queue close failed")
			consumerErr   = errors.New("consumer close failed")
			queue         = mockQueue.NewMockI(mockCtrl)
			consumer      = mockConsumer.NewMockI(mockCtrl)
			jobRepo       = mockRepo.NewMockJobsRepository(mockCtrl)
			uuidGenerator = mock.NewMockUUIDGenerator(mockCtrl)

			worker = &jobWorker{
				queue:         queue,
				consumer:      consumer,
				jobRepo:       jobRepo,
				uuidGenerator: uuidGenerator,
			}
		)

		queue.EXPECT().Close().Return(queueErr)
		consumer.EXPECT().Close().Return(consumerErr)

		err := worker.Close()
		require.Error(t, err)
		assert.Contains(t, err.Error(), "queue close failed")
		assert.Contains(t, err.Error(), "consumer close failed")
	})
}
