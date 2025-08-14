package model

import (
	"encoding/json"
	"time"
)

type Repository struct {
	ID          uint `gorm:"primaryKey"`
	UID         string
	URL         string
	Path        string
	ClonedAt    time.Time
	UpdatedAt   time.Time
	LastChecked time.Time
}

type Content struct {
	Item  string
	Count int
}

type StatsData struct {
	CommitsByUser     []Content
	FileModifications []Content
	LinesAdded        int
	LinesDeleted      int
	CommitsByWeekday  map[time.Weekday]int
}

func (s *StatsData) Bytes() []byte {
	d, err := json.Marshal(s)
	if err != nil {
		return nil
	}
	return d
}
