package svc

import (
	"fmt"
	"os"
	"time"

	"github.com/dyammarcano/clonr/internal/db"
	"github.com/dyammarcano/clonr/internal/model"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
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

// GetRepoStats itera sobre los commits de un repositorio y recopila datos estadísticos.
func GetRepoStats(repoPath string) (*model.StatsData, error) {
	r, err := git.PlainOpen(repoPath)
	if err != nil {
		return nil, fmt.Errorf("error al abrir el repositorio: %w", err)
	}

	commitIterator, err := r.Log(&git.LogOptions{})
	if err != nil {
		return nil, fmt.Errorf("error al obtener el historial de commits: %w", err)
	}

	stats := &model.StatsData{
		CommitsByUser:     make(map[string]int),
		FileModifications: make(map[string]int),
		CommitsByWeekday:  make(map[time.Weekday]int),
	}

	err = commitIterator.ForEach(func(commit *object.Commit) error {
		stats.CommitsByUser[commit.Author.Email]++
		stats.CommitsByWeekday[commit.Author.When.Weekday()]++

		if commit.NumParents() == 0 {
			return nil
		}

		parent, err := commit.Parent(0)
		if err != nil {
			return err
		}

		patch, err := parent.Patch(commit)
		if err != nil {
			return err
		}

		for _, fs := range patch.Stats() {
			stats.FileModifications[fs.Name]++
			stats.LinesAdded += fs.Addition
			stats.LinesDeleted += fs.Deletion
		}

		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("error al iterar sobre los commits: %w", err)
	}

	return stats, nil
}
