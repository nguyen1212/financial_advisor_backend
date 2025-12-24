package presenter

import (
	"testing"

	"github.com/financial_advisor/app/usecases/dto"
	"github.com/stretchr/testify/assert"
)

func TestPublisher_FormPublisher(t *testing.T) {
	t.Parallel()
	assert.Equal(t, Publisher{
		ID:          1,
		Name:        "Financial Times",
		Domain:      "https://www.ft.com",
		Description: "Leading financial news publication",
	}, FormPublisher(dto.Publisher{
		ID:          1,
		Name:        "Financial Times",
		Domain:      "https://www.ft.com",
		Description: "Leading financial news publication",
	}))
}

func TestPublisher_FormPublishers(t *testing.T) {
	t.Parallel()
	publisherDtos := []dto.Publisher{
		{
			ID:          1,
			Name:        "Financial Times",
			Domain:      "https://www.ft.com",
			Description: "Leading financial news publication",
		},
		{
			ID:          2,
			Name:        "Reuters",
			Domain:      "https://www.reuters.com",
			Description: "International news agency",
		},
	}
	expected := []Publisher{
		{
			ID:          1,
			Name:        "Financial Times",
			Domain:      "https://www.ft.com",
			Description: "Leading financial news publication",
		},
		{
			ID:          2,
			Name:        "Reuters",
			Domain:      "https://www.reuters.com",
			Description: "International news agency",
		},
	}

	assert.Equal(t, expected, FormPublishers(publisherDtos))
}