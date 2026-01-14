package worker

import "context"

//go:generate mockgen -destination=./mock/mock_$GOFILE -source=$GOFILE -package=mock

type I interface {
	Run(context.Context, []byte) error
	Close() error
}
