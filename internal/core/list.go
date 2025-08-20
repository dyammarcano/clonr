package core

import (
	"github.com/dyammarcano/clonr/internal/database"
	"github.com/dyammarcano/clonr/internal/model"
)

func ListRepos() ([]model.Repository, error) {
	initDB, err := database.InitDB()
	if err != nil {
		return nil, err
	}

	return initDB.GetAllRepos()
}
