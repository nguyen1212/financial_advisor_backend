// Package consumer holds implementation logic of different consumer types
package consumer

import (
	"context"

	"github.com/financial_advisor/app/services/consumer"
	"github.com/financial_advisor/app/services/queue"
	"github.com/sirupsen/logrus"
)

// ephemeralManager processes jobs without state tracking
type ephemeralManager struct {
	mProcessors map[queue.MessageType]consumer.Processor
}

// NewEphemeralManager creates new ephemeral job manager
// which processes jobs without state tracking
func NewEphemeralManager(processors map[queue.MessageType]consumer.Processor) consumer.I {
	return &ephemeralManager{
		mProcessors: processors,
	}
}

func (m *ephemeralManager) Execute(ctx context.Context, msg queue.Message) {
	processor, ok := m.mProcessors[msg.Type]
	if !ok {
		logrus.WithField("message_type", msg.Type.String()).
			Errorln("processor for message is not available")

		return
	}

	if err := processor.Execute(ctx, msg.Body); err != nil {
		logrus.WithField("message_type", msg.Type.String()).
			WithError(err).
			Errorln("process message")
	}
}

func (m *ephemeralManager) Close() error {
	return nil
}
