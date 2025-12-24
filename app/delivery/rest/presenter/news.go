// Package presenter handle logic to transform data returned to client side.
package presenter

import (
	"time"

	"github.com/financial_advisor/app/usecases/dto"
)

type New struct {
	ID          uint64 `json:"id,omitempty"`
	Title       string `json:"title"`
	Thumbnail   string `json:"thumbnail,omitempty"`
	Status      string `json:"status"`
	Content     string `json:"content,omitempty"`
	Author      string `json:"author,omitempty"`
	PublishedAt string `json:"published_at,omitempty"`
	URL         string `json:"url,omitempty"`
}

func FormNew(news dto.News) New {
	return New{
		ID:        news.ID,
		Title:     news.Title,
		Thumbnail: news.Thumbnail,
		Status:    news.Status,
		Content:   news.Content,
		URL:       news.URL,
		Author:    news.Author,
		PublishedAt: func() string {
			if news.PublishedAt != nil {
				return news.PublishedAt.Format(time.DateTime)
			}

			return ""
		}(),
	}
}

func FormNews(news []dto.News) []New {
	result := make([]New, len(news))

	for i := range news {
		result[i] = FormNew(news[i])
	}

	return result
}
