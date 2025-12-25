// Package payload presents the request payload
package payload

import (
	"time"

	"github.com/financial_advisor/app/domain/entity"
	"github.com/financial_advisor/app/usecases/dto"
)

type NewsFindRequest struct {
	From   time.Time `form:"from"`
	To     time.Time `form:"to"`
	Status string    `form:"status" validate:"omitempty,oneof=added synced failed"`
	Paging
}

func (r NewsFindRequest) Validate() error {
	return validate(r)
}

func (r NewsFindRequest) ToDTO() dto.NewsFindRequest {
	d := dto.NewsFindRequest{
		From: r.From,
		To:   r.To,
		Paging: dto.PagingRequest{
			Page: r.Page,
			Size: r.Size,
		},
	}

	if r.Status != "" {
		status := entity.ToNewsStatus(r.Status)

		d.Status = &status
	}

	return d
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
