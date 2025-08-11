package model

import "time"

type Repository struct {
	ID          uint `gorm:"primaryKey"`
	UID         string
	URL         string
	Path        string
	ClonedAt    time.Time
	UpdatedAt   time.Time
	LastChecked time.Time
}
