package handler

import (
	"context"
	"net/http"
)

type findHandler interface {
	Find(*http.Request) (any, any, error)
}

type getHandler interface {
	Get(ctx context.Context)
}

type Handler interface {
	findHandler
}
