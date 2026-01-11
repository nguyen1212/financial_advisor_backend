// Package consumer holds internal logic for message consumers and processors.
package consumer

//go:generate mockgen -destination=./mock/mock_$GOFILE -source=$GOFILE -package=mock

import (
	"context"

	"github.com/financial_advisor/app/services/queue"
)

type I interface {
	Execute(context.Context, queue.Message)
	Close() error
}

type Processor interface {
	Execute(context.Context, []byte) error
}
