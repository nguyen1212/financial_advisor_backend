// Package handler accepts incoming requests from external traffic.
package handler

//go:generate mockgen -destination=./mock/mock_$GOFILE -source=$GOFILE -package=mock

import (
	"net/http"
)

type FindHandler interface {
	Find(*http.Request) (any, any, error)
}

type GetHandler interface {
	Get(*http.Request) (any, error)
}

type CreateHandler interface {
	Create(*http.Request) (any, error)
}

type DeleteHandler interface {
	Delete(*http.Request) error
}

type SearchHandler interface {
	GetSearchSuggestions(*http.Request) (any, error)
	Search(*http.Request) (any, error)
}
