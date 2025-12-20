// Package entity represents core business entity
package entity

import "time"

type NewsStatus int

const (
	NewsStatusUnknown NewsStatus = iota
	NewsStatusAdded
	NewsStatusSynced
)

type NewsDomain int

const (
	NewsDomainFinance NewsDomain = iota
	NewsDomainMilitary
)

type News struct {
	ID uint64 `gorm:"primaryKey"`

	Title     string
	Thumbnail string
	Status    NewsStatus
	Domain    NewsDomain
	Publisher string

	CreatedAt time.Time
	UpdatedAt time.Time
}
