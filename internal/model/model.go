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

// StatsData almacena todos los datos de estadísticas del repositorio.
type StatsData struct {
	CommitsByUser     map[string]int
	FileModifications map[string]int
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
