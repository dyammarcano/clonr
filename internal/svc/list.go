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

	// Use maps internally for counting
	commitsByUser := make(map[string]int)
	fileModifications := make(map[string]int)
	commitsByWeekday := make(map[time.Weekday]int)

	var linesAdded, linesDeleted int

	err = commitIterator.ForEach(func(commit *object.Commit) error {
		commitsByUser[commit.Author.Email]++
		commitsByWeekday[commit.Author.When.Weekday()]++

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
			fileModifications[fs.Name]++
			linesAdded += fs.Addition
			linesDeleted += fs.Deletion
		}

		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("error al iterar sobre los commits: %w", err)
	}

	// Convert maps to []Content for the final struct
	stats := &model.StatsData{
		CommitsByUser:     mapToContentSlice(commitsByUser),
		FileModifications: mapToContentSlice(fileModifications),
		LinesAdded:        linesAdded,
		LinesDeleted:      linesDeleted,
		CommitsByWeekday:  commitsByWeekday,
	}

	return stats, nil
}

// Helper to convert map[string]int → []Content
func mapToContentSlice(m map[string]int) []model.Content {
	slice := make([]model.Content, 0, len(m))
	for name, count := range m {
		slice = append(slice, model.Content{Item: name, Count: count})
	}
	return slice
}
