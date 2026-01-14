// Package consumer holds implementation logic of different consumer types
package consumer

import (
	"context"
	"errors"

	"github.com/financial_advisor/app/services/queue"
	"github.com/sirupsen/logrus"
)

// manager processes jobs without state tracking
type manager struct {
	mProcessors map[queue.MessageType]Processor
}

// NewManager creates new ephemeral job manager
func NewManager(processors map[queue.MessageType]Processor) I {
	return &manager{
		mProcessors: processors,
	}
}

func (m *manager) Execute(ctx context.Context, msg queue.Message) error {
	processor, ok := m.mProcessors[msg.Type]
	if !ok {
		logrus.WithField("message_type", msg.Type.String()).
			Errorln("processor for message is not available")

		return errors.New("processor for message is not available: " + msg.Type.String())
	}

	if err := processor.Execute(ctx, msg.Body); err != nil {
		logrus.WithField("message_type", msg.Type.String()).
			WithError(err).
			Errorln("process message")
	}

	return nil
}

func (m *manager) Close() error {
	return nil
}
