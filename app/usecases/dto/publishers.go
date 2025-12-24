package dto

import (
	"github.com/financial_advisor/app/domain/entity"
)

type Publisher struct {
	ID          uint64
	Name        string
	Domain      string
	Description string
}

func ToDtoPublisher(publisher entity.Publisher) Publisher {
	return Publisher{
		ID:          publisher.ID,
		Name:        publisher.Name,
		Domain:      publisher.Domain,
		Description: publisher.Description,
	}
}

type PublishersFindRequest struct {
	Paging PagingRequest
}

type PublisherCreateRequest struct {
	Name        string
	Domain      string
	Description string
}
