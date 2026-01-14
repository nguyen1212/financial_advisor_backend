// Package queue defines the interface for a queue service.
package queue

//go:generate mockgen -destination=./mock/mock_$GOFILE -source=$GOFILE -package=mock

type (
	Message struct {
		Body []byte
		Type MessageType
		UUID string
	}

	I interface {
		Dequeue() (Message, bool)
		Enqueue(Message) error
		Close() error
	}

	MessageType int
)

const (
	MessageTypeWebScrapper MessageType = iota
)

func (t MessageType) String() string {
	switch t {
	case MessageTypeWebScrapper:
		return "web_scrapper"
	default:
		return "unknown"
	}
}
