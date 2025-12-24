// Package payload presents the request payload
package payload

import (
	"time"

	"github.com/financial_advisor/app/domain/entity"
	"github.com/financial_advisor/app/usecases/dto"
)

type NewsFindRequest struct {
	From time.Time `form:"from"`
	To   time.Time `form:"to"`
	Paging
}

func (r NewsFindRequest) ToDTO() dto.NewsFindRequest {
	return dto.NewsFindRequest{
		From: r.From,
		To:   r.To,
		Paging: dto.PagingRequest{
			Page: r.Page,
			Size: r.Size,
		},
	}
}

type NewsCreateRequest struct {
	URL      string `json:"url" validate:"required,url"`
	Category string `json:"category" validate:"required,oneof=finance military"`
}

func (r NewsCreateRequest) Validate() error {
	return validate(r)
}

func (r NewsCreateRequest) ToDTO() dto.NewsCreateRequest {
	return dto.NewsCreateRequest{
		URL:      r.URL,
		Category: entity.ToNewsCategory(r.Category),
	}
}
