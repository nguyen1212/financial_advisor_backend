package presenter

import (
	"github.com/financial_advisor/app/usecases/dto"
)

type Publisher struct {
	ID          uint64 `json:"id"`
	Name        string `json:"name"`
	Domain      string `json:"domain"`
	Description string `json:"Description,omitempty"`
}

func FormPublisher(publisher dto.Publisher) Publisher {
	return Publisher{
		ID:          publisher.ID,
		Name:        publisher.Name,
		Domain:      publisher.Domain,
		Description: publisher.Description,
	}
}

func FormPublishers(publishers []dto.Publisher) []Publisher {
	result := make([]Publisher, len(publishers))

	for i := range publishers {
		result[i] = FormPublisher(publishers[i])
	}

	return result
}
