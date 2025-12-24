// Package memoryqueue implements logic of queue interface using memory queue
package memoryqueue

import (
	"errors"
	"sync"

	"github.com/financial_advisor/app/services/consumer"
	"github.com/financial_advisor/app/services/queue"
	"github.com/sirupsen/logrus"
)

type memQueue struct {
	ch       chan queue.Message
	consumer consumer.Manager
}

var (
	once   sync.Once
	mQueue memQueue
)

// Init initializes the queue and run it in background
func Init(globalWorkersNo int, consumer consumer.Manager) {
	once.Do(func() {
		mQueue = memQueue{
			ch:       make(chan queue.Message, globalWorkersNo),
			consumer: consumer,
		}

		for range globalWorkersNo {
			go func() {
				defer func() {
					if r := recover(); r != nil {
						logrus.Errorf("memory queue worker recovered from panic: %v", r)
					}
				}()

				for message := range mQueue.ch {
					mQueue.consumer.Execute(message)
				}
			}()
		}
	})
}

func Get() queue.I {
	return &mQueue
}

func (m *memQueue) Enqueue(msg queue.Message) error {
	if m.ch == nil {
		return errors.New("memory queue is not initialized")
	}

	m.ch <- msg

	return nil
}

func (m *memQueue) Close() error {
	select {
	case _, ok := <-m.ch:
		if !ok {
			close(m.ch)
		}
	default:
		close(m.ch)
	}

	return nil
}
