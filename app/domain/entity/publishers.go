package entity

import "time"

type Publisher struct {
	ID uint64 `gorm:"primaryKey"`

	Name        string
	Description string
	Domain      string

	CreatedAt time.Time
	UpdatedAt time.Time
}
