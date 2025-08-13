package git

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/dyammarcano/clonr/internal/db"
	"github.com/dyammarcano/clonr/internal/model"
)

// UpdateAllRepos pulls the latest changes for all repositories in the clonr database.
func UpdateAllRepos() {
	dbConn, err := db.InitDB()
	if err != nil {
		log.Printf("Failed to open database: %v\n", err)
		return
	}

	repos, err := dbConn.GetAllRepos()
	if err != nil {
		log.Printf("Failed to get repositories: %v\n", err)
		return
	}

	for _, repo := range repos {
		_ = UpdateRepo(repo.Path)
	}
}

func MoveRepo(uid string, repos []model.Repository, destFolder string) error {
	var repo *model.Repository
	for i := range repos {
		if repos[i].UID == uid {
			repo = &repos[i]
			break
		}
	}

	if repo == nil {
		return fmt.Errorf("repository with UID %s not found", uid)
	}

	// Ensure destination exists
	if err := os.MkdirAll(destFolder, 0o755); err != nil {
		return fmt.Errorf("failed to create destination folder: %w", err)
	}

	absDestFolder, err := filepath.Abs(destFolder)
	if err != nil {
		return fmt.Errorf("failed to get absolute path of destination folder: %w", err)
	}

	destPath := filepath.Join(absDestFolder, filepath.Base(repo.Path))

	dbConn, err := db.InitDB()
	if err != nil {
		return fmt.Errorf("failed to init database: %w", err)
	}

	// Move folder
	if err := os.Rename(repo.Path, destPath); err != nil {
		return fmt.Errorf("failed to move repository: %w", err)
	}

	fmt.Printf("Repository %s moved to %s\n", repo.UID, destPath)

	repo.Path = destPath
	repo.UpdatedAt = time.Now()
	repo.LastChecked = time.Now()

	if err := dbConn.UpdateRepoFields(repo.UID, repo.Path, repo.UpdatedAt, repo.LastChecked); err != nil {
		return fmt.Errorf("failed to update repository in database: %w", err)
	}

	return nil
}

func UpdateRepo(path string) error {
	log.Printf("Updating %s...", path)

	cmd := exec.Command("git", "pull", "origin")
	cmd.Dir = path

	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Printf("[pull error] %v: %s\n", err, string(output))
		return err
	}

	log.Printf("[updated] %s\n", output)

	return nil
}
