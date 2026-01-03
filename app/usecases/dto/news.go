// Package dto helps to transfer data between usecases and handler
package dto

import (
	"time"

	"github.com/financial_advisor/app/domain/entity"
)

type News struct {
	ID          uint64
	Title       string
	Thumbnail   string
	Status      string
	Category    string
	Publisher   string
	Content     string
	URL         string
	Author      string
	PublishedAt *time.Time
}

func ToDtoNews(news entity.News) News {
	return News{
		ID:          news.ID,
		Title:       news.Title,
		Thumbnail:   news.Thumbnail,
		Status:      news.Status.String(),
		Category:    news.Category.String(),
		Publisher:   news.Publisher.Name,
		Content:     news.Content,
		Author:      news.Author,
		PublishedAt: news.PublishedAt,
		URL:         news.URL,
	}
}

type NewsFindRequest struct {
	From   time.Time
	To     time.Time
	Status *entity.NewsStatus
	Paging PagingRequest
}

type NewsCreateRequest struct {
	URL      string
	Category entity.NewsCategory
}

type NewsSearchSuggestionsRequest struct {
	Keywords []string
}

type NewsSearchRequest struct {
	Keywords []string
	Paging   PagingRequest
}

type NewsGetRequest struct {
	HighlightKeywords []string
}
