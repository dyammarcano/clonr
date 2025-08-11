package git

import "github.com/dyammarcano/clonr/internal/db"

func RemoveRepo(url string) error {
	initDB, err := db.InitDB()
	if err != nil {
		return err
	}

	return initDB.RemoveRepoByURL(url)
}
