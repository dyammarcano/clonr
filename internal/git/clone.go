package git

import (
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"
)

func CloneRepo(url string) (string, error) {
	// Define path destino, ej: ./repos/<nombre-repo> + timestamp opcional
	repoName := extractRepoName(url) // función que extrae "repo" de la URL
	path := filepath.Join("repos", repoName)

	cmd := exec.Command("git", "clone", url, path)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("git clone error: %v - %s", err, string(output))
	}

	return path, nil
}

func extractRepoName(url string) string {
	parts := strings.Split(url, "/")
	last := parts[len(parts)-1]
	return strings.TrimSuffix(last, ".git")
}
