package db

import (
	"time"

	"clonr/internal/model"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Database struct {
	*gorm.DB
}

func InitDB() (*Database, error) {
	db, err := gorm.Open(sqlite.Open("clonr.db"), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	if err := db.AutoMigrate(&model.Repository{}); err != nil {
		return nil, err
	}

	return &Database{
		DB: db,
	}, nil
}

func (d *Database) SaveRepo(url, path string) error {
	repo := model.Repository{
		URL:      url,
		Path:     path,
		ClonedAt: time.Now(),
	}
	return d.DB.Create(&repo).Error
}

func (d *Database) GetAllRepos() ([]model.Repository, error) {
	var repos []model.Repository
	err := d.DB.Find(&repos).Error
	return repos, err
}
