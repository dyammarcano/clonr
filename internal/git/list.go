package git

import (
	"github.com/dyammarcano/clonr/internal/db"
	"github.com/dyammarcano/clonr/internal/model"
)

func ListRepos() ([]model.Repository, error) {
	initDB, err := db.InitDB()
	if err != nil {
		return nil, err
	}

	return initDB.GetAllRepos()
}
