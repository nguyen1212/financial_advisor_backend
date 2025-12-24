package consumer

import (
	"errors"
	"testing"

	"github.com/financial_advisor/app/services/consumer/mock"
	"github.com/financial_advisor/app/services/consumer/processor"
	"github.com/financial_advisor/app/services/queue"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func Test_NewEphemeralManager(t *testing.T) {
	t.Parallel()

	mProcessors := map[queue.MessageType]Processor{
		queue.MessageTypeWebScrapper: processor.NewWebScrapperProcessor(nil),
	}

	assert.Equal(
		t,
		&ephemeralManager{
			mProcessors: mProcessors,
		},
		NewEphemeralManager(mProcessors),
	)
}

func Test_ephemeralManager_Execute(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		msg        queue.Message
		setupMocks func(*mock.MockProcessor, []byte)
	}{
		{
			name: "processor is not available",
			msg: queue.Message{
				Type: queue.MessageType(999),
			},
			setupMocks: func(*mock.MockProcessor, []byte) {},
		},
		{
			name: "processor is available and failed to execute message",
			msg: queue.Message{
				Type: queue.MessageTypeWebScrapper,
				Body: []byte("test body"),
			},
			setupMocks: func(processor *mock.MockProcessor, body []byte) {
				processor.EXPECT().Execute(body).Return(errors.New("some error"))
			},
		},
	}

	for i := range tests {
		t.Run(tests[i].name, func(t *testing.T) {
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()

			var (
				webScrapperProcessorMock = mock.NewMockProcessor(mockCtrl)
				mProcessors              = map[queue.MessageType]Processor{
					queue.MessageTypeWebScrapper: webScrapperProcessorMock,
				}
				manager = &ephemeralManager{
					mProcessors: mProcessors,
				}
				msg = tests[i].msg
			)

			tests[i].setupMocks(webScrapperProcessorMock, msg.Body)

			manager.Execute(msg)
		})
	}
}
