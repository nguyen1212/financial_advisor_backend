package consumer

import (
	"github.com/financial_advisor/app/services/queue"
	"github.com/sirupsen/logrus"
)

type ephemeralManager struct {
	mProcessors map[queue.MessageType]Processor
}

// NewEphemeralManager creates new ephemeral job manager
// which processes jobs without state tracking
func NewEphemeralManager(processors map[queue.MessageType]Processor) Manager {
	return &ephemeralManager{
		mProcessors: processors,
	}
}

func (m *ephemeralManager) Execute(msg queue.Message) {
	processor, ok := m.mProcessors[msg.Type]
	if !ok {
		logrus.WithField("message_type", msg.Type.String()).
			Errorln("processor for message is not available")

		return
	}

	if err := processor.Execute(msg.Body); err != nil {
		logrus.WithField("message_type", msg.Type.String()).
			WithError(err).
			Errorln("process message")
	}
}
