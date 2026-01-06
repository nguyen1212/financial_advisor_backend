// Package worker represents the logic with combination of queue and consumer
package worker

//go:generate mockgen -destination=./mock/mock_$GOFILE -source=$GOFILE -package=mock

import (
	"context"
	"errors"
	"sync"

	"github.com/financial_advisor/app/services/consumer"
	"github.com/financial_advisor/app/services/queue"
	"github.com/sirupsen/logrus"
)

var (
	once     sync.Once
	instance *worker
)

type I interface {
	Run([]byte) error
	Close() error
}

type worker struct {
	queue    queue.I
	consumer consumer.I
}

// Init initializes the queue and run it in background
func Init(
	ctx context.Context,
	queue queue.I,
	consumer consumer.I,
	workersNo int,
) {
	once.Do(func() {
		instance = &worker{
			queue:    queue,
			consumer: consumer,
		}

		for range workersNo {
			go func() {
				defer func() {
					if r := recover(); r != nil {
						logrus.Errorf("memory queue worker recovered from panic: %v", r)
					}
				}()

				for {
					// check for context cancellation (graceful shutdown)
					select {
					case <-ctx.Done():
						return
					default:
					}

					msg, isClosed := instance.queue.GetMsg()
					if !isClosed {
						break
					}

					instance.consumer.Execute(ctx, msg)
				}
			}()
		}
	})
}

// Get currently only one singleton worker instance
// we can extend to multiple types of worker by removing the sync.Once pattern
func Get() I {
	return instance
}

func (w *worker) Run(job []byte) error {
	return w.queue.Enqueue(job)
}

func (w *worker) Close() error {
	var closeErr error

	if err := w.queue.Close(); err != nil {
		closeErr = err
	}

	if err := w.consumer.Close(); err != nil {
		closeErr = errors.Join(closeErr, err)
	}

	return closeErr
}
