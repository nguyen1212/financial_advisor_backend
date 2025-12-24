package payload

import (
	"github.com/financial_advisor/app/usecases/dto"
)

type PublishersFindRequest struct {
	Paging
}

func (r PublishersFindRequest) ToDTO() dto.PublishersFindRequest {
	d := dto.PublishersFindRequest{
		Paging: dto.PagingRequest{
			Page: r.Page,
			Size: r.Size,
		},
	}

	if d.Paging.Size == 0 {
		d.Paging.Size = 30
	}

	if d.Paging.Page == 0 {
		d.Paging.Page = 1
	}

	return d
}

type PublisherCreateRequest struct {
	Name        string `json:"name" validate:"required"`
	Domain      string `json:"domain" validate:"required,url"`
	Description string `json:"description" validate:"omitempty"`
}

func (r PublisherCreateRequest) Validate() error {
	return validate(r)
}

func (r PublisherCreateRequest) ToDTO() dto.PublisherCreateRequest {
	d := dto.PublisherCreateRequest{}

	d.Name = r.Name
	d.Domain = r.Domain
	d.Description = r.Description

	return d
}
