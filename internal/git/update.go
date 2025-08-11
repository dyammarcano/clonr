package git

import (
	"fmt"
	"os/exec"

	"github.com/dyammarcano/clonr/internal/db"
)

// UpdateAllRepos pulls the latest changes for all repositories in the clonr database.
func UpdateAllRepos() {
	dbConn, err := db.InitDB()
	if err != nil {
		fmt.Println("Failed to open database:", err)
		return
	}

	repos, err := dbConn.GetAllRepos()
	if err != nil {
		fmt.Println("Failed to get repositories:", err)
		return
	}

	for _, repo := range repos {
		_ = UpdateRepo(repo.Path)
	}
}

func UpdateRepo(path string) error {
	fmt.Printf("Updating %s... ", path)
	cmd := exec.Command("git", "pull", "origin")
	cmd.Dir = path
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("[pull error] %v: %s\n", err, string(output))
		return err
	}
	fmt.Println("[updated]", string(output))
	return nil
}
