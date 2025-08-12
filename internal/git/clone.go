package git

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/dyammarcano/clonr/internal/db"

	"github.com/spf13/cobra"
)

func CloneRepo(_ *cobra.Command, args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("repository URL is required")
	}

	url := strings.TrimSpace(args[0])
	if url == "" {
		return fmt.Errorf("repository URL cannot be empty")
	}

	pathStr := "."

	if len(args) > 1 {
		pathStr = args[1]
	}

	if pathStr == "." || pathStr == "./" {
		wd, err := os.Getwd()
		if err != nil {
			return fmt.Errorf("error getting current working directory: %w", err)
		}

		pathStr = wd
	}

	if _, err := os.Stat(pathStr); os.IsNotExist(err) {
		if err := os.MkdirAll(pathStr, os.ModePerm); err != nil {
			return fmt.Errorf("error creating directory %s: %w", pathStr, err)
		}
	}

	initDB, err := db.InitDB()
	if err != nil {
		return fmt.Errorf("starting server: %w", err)
	}

	absPath, err := filepath.Abs(pathStr)
	if err != nil {
		return fmt.Errorf("error determining absolute path: %w", err)
	}

	savePath := filepath.Join(absPath, extractRepoName(url))

	runCmd := exec.Command("git", "clone", url, savePath)

	output, err := runCmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("git clone error: %v - %s", err, string(output))
	}

	if err := initDB.SaveRepo(url, savePath); err != nil {
		return fmt.Errorf("error saving repo to database: %w", err)
	}

	log.Printf("Cloned repo at %s\n", savePath)

	return nil
}

func PullRepo(path string) error {
	cmd := exec.Command("git", "-C", path, "pull")

	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("error git pull: %v, output: %s", err, string(output))
	}

	return nil
}

func extractRepoName(url string) string {
	parts := strings.Split(url, "/")
	last := parts[len(parts)-1]

	return strings.TrimSuffix(last, ".git")
}
