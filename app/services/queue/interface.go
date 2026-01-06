// Package queue defines the interface for a queue service.
package queue

//go:generate mockgen -destination=./mock/mock_$GOFILE -source=$GOFILE -package=mock

type (
	Message struct {
		Body []byte
		Type MessageType
	}

	I interface {
		GetMsg() (Message, bool)
		Enqueue([]byte) error
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
