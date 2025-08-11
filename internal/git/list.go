package git

import (
	"clonr/internal/db"
	"clonr/internal/model"
)

func ListRepos() ([]*model.Repository, error) {
	initDB, err := db.InitDB()
	if err != nil {
		return nil, err
	}

	repos, err := initDB.GetAllRepos()
	if err != nil {
		return nil, err
	}

	var repoPointers []*model.Repository
	for i := range repos {
		repoPointers = append(repoPointers, &repos[i])
	}

	return repoPointers, nil
}
