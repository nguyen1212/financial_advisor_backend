// Package consumer holds internal logic for message consumers and processors.
package consumer

//go:generate mockgen -destination=./mock/mock_$GOFILE -source=$GOFILE -package=mock

import "github.com/financial_advisor/app/services/queue"

type Manager interface {
	Execute(queue.Message)
}

type Processor interface {
	Execute([]byte) error
}
