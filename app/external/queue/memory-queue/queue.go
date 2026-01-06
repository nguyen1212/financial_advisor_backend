// Package memoryqueue implements logic of queue interface using memory queue
package memoryqueue

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/financial_advisor/app/services/queue"
)

type memQueue struct {
	ch chan queue.Message
}

func New(workersNo int) queue.I {
	return &memQueue{
		ch: make(chan queue.Message, workersNo),
	}
}

func (m *memQueue) GetMsg() (queue.Message, bool) {
	msg, ok := <-m.ch

	return msg, ok
}

func (m *memQueue) Enqueue(msg []byte) error {
	if m.ch == nil {
		return errors.New("memory queue is not initialized")
	}

	var parsedMsg queue.Message

	if err := json.Unmarshal(msg, &parsedMsg); err != nil {
		return fmt.Errorf("message is incorrect format: %w", err)
	}

	m.ch <- parsedMsg

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
