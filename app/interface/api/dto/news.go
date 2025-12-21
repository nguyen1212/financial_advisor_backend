// Package dto helps to transfer data between usecases and handler
package dto

import (
	"time"

	"github.com/financial_advisor/app/domain/entity"
)

type News struct {
	Title     string
	Thumbnail string
	Status    string
	Category  string
	Publisher string
	Content   string
}

func ToDtoNews(news entity.News) News {
	return News{
		Title:     news.Title,
		Thumbnail: news.Thumbnail,
		Status:    news.Status.String(),
		Category:  news.Category.String(),
		Publisher: news.Publisher.Name,
		Content:   news.Content,
	}
}

type NewsFindRequest struct {
	From time.Time
	To   time.Time
}
