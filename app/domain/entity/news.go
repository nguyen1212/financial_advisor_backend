// Package entity represents core business entity
package entity

import "time"

type NewsStatus int

const (
	NewsStatusUnknown NewsStatus = iota
	NewsStatusAdded
	NewsStatusSynced
)

func (status NewsStatus) String() string {
	switch status {
	case NewsStatusAdded:
		return "added"
	case NewsStatusSynced:
		return "synced"
	default:
		return "unknown"
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

type News struct {
	ID uint64 `gorm:"primaryKey"`

	Title     string
	Thumbnail string
	Link      string
	Status    NewsStatus
	Category  NewsCategory
	Content   string

	Publisher   Publisher `gorm:"foreignKey:PublisherID;->"`
	PublisherID uint64

	PublishedAt time.Time
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
