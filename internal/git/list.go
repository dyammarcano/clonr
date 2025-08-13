package git

import (
	"fmt"
	"os"
	"time"

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

func PrettiListRepos(print bool) ([]model.Repository, error) {
	list, err := ListRepos()
	if err != nil {
		return nil, err
	}

	if print {
		for _, repo := range list {
			_, _ = fmt.Fprintf(os.Stdout, "[%s] [%s] %s -> %s\n", repo.UID, repo.UpdatedAt.Format(time.DateTime), repo.URL, repo.Path)
		}
	}

	return list, nil
}
