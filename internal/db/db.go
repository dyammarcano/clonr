package db

import (
	"path/filepath"
	"time"

	"github.com/dyammarcano/clonr/internal/model"
	"github.com/dyammarcano/clonr/internal/params"

	"github.com/google/uuid"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Database struct {
	*gorm.DB
}

func InitDB() (*Database, error) {
	path := filepath.Join(params.AppdataDir, "clonr.db")

	db, err := gorm.Open(sqlite.Open(path), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	if err := db.AutoMigrate(&model.Repository{}); err != nil {
		return nil, err
	}

	return &Database{DB: db}, nil
}

func (d *Database) SaveRepo(url, path string) error {
	return d.Create(&model.Repository{
		UID:       uuid.NewString(),
		URL:       url,
		Path:      path,
		ClonedAt:  time.Now(),
		UpdatedAt: time.Now(),
	}).Error
}

func (d *Database) GetAllRepos() ([]model.Repository, error) {
	var repos []model.Repository

	return repos, d.Find(&repos).Error
}

func (d *Database) RemoveRepoByURL(url string) error {
	return d.Where("url = ?", url).Delete(&model.Repository{}).Error
}
