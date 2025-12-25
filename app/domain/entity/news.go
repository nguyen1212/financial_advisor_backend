// Package entity represents core business entity
package entity

import (
	"encoding/hex"
	"fmt"
	"time"
)

type NewsStatus int

const (
	NewsStatusUnknown NewsStatus = iota
	NewsStatusAdded
	NewsStatusSynced
	NewsStatusFailed
)

func (status NewsStatus) String() string {
	switch status {
	case NewsStatusAdded:
		return "added"
	case NewsStatusSynced:
		return "synced"
	case NewsStatusFailed:
		return "failed"
	default:
		return "unknown"
	}
}

func ToNewsStatus(status string) NewsStatus {
	switch status {
	case "added":
		return NewsStatusAdded
	case "synced":
		return NewsStatusSynced
	case "failed":
		return NewsStatusFailed
	default:
		return NewsStatusUnknown
	}
}

type NewsCategory int

const (
	NewsCategoryUnknown NewsCategory = iota
	NewsCategoryFinance
	NewsCategoryMilitary
)

func (category NewsCategory) String() string {
	switch category {
	case NewsCategoryFinance:
		return "finance"
	case NewsCategoryMilitary:
		return "military"
	default:
		return "unknown"
	}
}

func ToNewsCategory(category string) NewsCategory {
	switch category {
	case "finance":
		return NewsCategoryFinance
	case "military":
		return NewsCategoryMilitary
	default:
		return NewsCategoryUnknown
	}
}

type News struct {
	ID uint64 `gorm:"primaryKey"`

	Title     string
	Author    string
	Thumbnail string
	URL       string
	HashedURL []byte
	Status    NewsStatus
	Category  NewsCategory

	FilePath string
	FileSize int64
	Content  string `gorm:"-"`

	Publisher   Publisher `gorm:"foreignKey:PublisherID;->"`
	PublisherID uint64

	PublishedAt *time.Time
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func (news *News) StoragePath() string {
	return fmt.Sprintf(
		"data/scraped/publishers/%d/news/%s.txt",
		news.PublisherID,
		hex.EncodeToString(news.HashedURL),
	)
}

func (news *News) StorageDir() string {
	return fmt.Sprintf(
		"data/scraped/publishers/%d/news",
		news.PublisherID,
	)
}

func (news *News) ToMap() map[string]any {
	return map[string]any{
		"author":       news.Author,
		"file_path":    news.FilePath,
		"file_size":    news.FileSize,
		"published_at": news.PublishedAt,
		"status":       news.Status,
		"title":        news.Title,
		"thumbnail":    news.Thumbnail,
	}
}
