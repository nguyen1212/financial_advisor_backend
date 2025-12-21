// Package payload presents the request payload
package payload

import (
	"time"

	"github.com/financial_advisor/app/interface/api/dto"
)

type NewsFindRequest struct {
	From time.Time `form:"from"`
	To   time.Time `form:"to"`
}

func (r NewsFindRequest) ToDTO() dto.NewsFindRequest {
	return dto.NewsFindRequest{
		From: r.From,
		To:   r.To,
	}
}
