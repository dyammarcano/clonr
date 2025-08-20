package core

import (
	"fmt"
	"log"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/dyammarcano/clonr/internal/database"
	"github.com/spf13/cobra"
)

// if dest dir is a dot clones into the current dir, if not,
// then clone into specified dir when dest dir not exists use default dir, saved in db

func CloneRepo(_ *cobra.Command, args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("repository URL is required")
	}

	uri, err := url.ParseRequestURI(args[0])
	if err != nil {
		return fmt.Errorf("invalid repository URL: %w", err)
	}

	if uri.Scheme != "http" && uri.Scheme != "https" {
		return fmt.Errorf("invalid repository URL: %s", uri.String())
	}

	if uri.Host == "" {
		return fmt.Errorf("invalid repository URL: %s", uri.String())
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

	initDB, err := database.InitDB()
	if err != nil {
		return fmt.Errorf("starting server: %w", err)
	}

	absPath, err := filepath.Abs(pathStr)
	if err != nil {
		return fmt.Errorf("error determining absolute path: %w", err)
	}

	savePath := filepath.Join(absPath, extractRepoName(uri.String()))

	runCmd := exec.Command("git", "clone", uri.String(), savePath)

	output, err := runCmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("git clone error: %v - %s", err, string(output))
	}

	if err := initDB.SaveRepo(uri.String(), savePath); err != nil {
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
