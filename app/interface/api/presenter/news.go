// Package presenter handle logic to transform data returned to client side.
package presenter

import (
	"github.com/financial_advisor/app/interface/api/dto"
)

type New struct {
	Title     string `json:"title"`
	Thumbnail string `json:"thumbnail,omitempty"`
	Status    string `json:"status"`
	Content   string `json:"content,omitempty"`
}

func FormNew(news dto.News) New {
	return New{
		Title:     news.Title,
		Thumbnail: news.Thumbnail,
		Status:    news.Status,
		Content:   news.Content,
	}
}

func FormNews(news []dto.News) []New {
	result := make([]New, len(news))

	for i := range news {
		result[i] = FormNew(news[i])
	}

	return result
}
