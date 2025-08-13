package db

import (
	"crypto/sha256"
	"fmt"
	"path/filepath"
	"time"

	"github.com/dyammarcano/clonr/internal/model"
	"github.com/dyammarcano/clonr/internal/params"

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
		UID:       fmt.Sprintf("%x", sha256.Sum256([]byte(url)))[:10],
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

func (d *Database) UpdateRepoFields(uid, path string, updatedAt, lastChecked time.Time) error {
	return d.Model(&model.Repository{}).
		Where("uid = ?", uid).
		Updates(map[string]any{
			"path":         path,
			"updated_at":   updatedAt,
			"last_checked": lastChecked,
		}).Error
}
