package model

import "time"

type Repository struct {
	ID          uint   `gorm:"primaryKey"`
	URL         string `gorm:"uniqueIndex"`
	Path        string
	ClonedAt    time.Time
	LastChecked time.Time
	Updated     bool
}
