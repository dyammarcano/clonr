package git

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/dyammarcano/clonr/internal/db"
	"github.com/go-git/go-git/v5"
	"github.com/spf13/cobra"
)

func MapRepos(cmd *cobra.Command, args []string) error {
	rootDir := "."

	if len(args) > 0 {
		rootDir = args[0]
	}

	dbConn, err := db.InitDB()
	if err != nil {
		return fmt.Errorf("failed to open database: %v", err)
	}

	found := 0
	err = filepath.Walk(rootDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() && info.Name() == ".git" {
			repoPath := filepath.Dir(path)
			repoUrl := extractRepoURL(repoPath)
			if repoUrl == "unknown" {
				repoUrl = extractRepoURL2(repoPath)
			}

			dbErr := dbConn.SaveRepo(repoUrl, repoPath)
			if dbErr == nil {
				cmd.Println("Added:", repoPath)
				found++
			}
			// Don't recurse into .git
			return filepath.SkipDir
		}
		return nil
	})
	if err != nil {
		cmd.Println("Error searching for git repos:", err)
	}
	cmd.Printf("%d repositories mapped.\n", found)

	return nil
}

func extractRepoURL(repoPath string) string {
	repo, err := git.PlainOpen(repoPath)
	if err != nil {
		return "unknown"
	}

	remotes, err := repo.Remotes()
	if err != nil || len(remotes) == 0 {
		return "unknown"
	}

	for _, remote := range remotes {
		if remote.Config().Name == "origin" {
			if len(remote.Config().URLs) > 0 {
				return remote.Config().URLs[0]
			}
		}
	}

	// fallback: return first remote URL if origin not found
	if len(remotes) > 0 && len(remotes[0].Config().URLs) > 0 {
		return remotes[0].Config().URLs[0]
	}

	return "unknown"
}

func extractRepoURL2(repoPath string) string {
	// Abrir el repo existente
	r, err := git.PlainOpen(repoPath)
	if err != nil {
		return "unknown"
	}

	// Obtener remote "origin"
	rem, err := r.Remote("origin")
	if err != nil {
		return "unknown"
	}

	urls := rem.Config().URLs
	if len(urls) > 0 {
		return urls[0]
	}

	return "unknown"
}
