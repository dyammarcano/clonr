package server

import (
	"fmt"
	"net/http"

	"github.com/dyammarcano/clonr/internal/db"
	"github.com/dyammarcano/clonr/internal/git"

	"github.com/gin-gonic/gin"
)

func StartServer() error {
	initDB, err := db.InitDB()
	if err != nil {
		return fmt.Errorf("starting server: %w", err)
	}

	r := gin.Default()

	r.GET("/repos", func(c *gin.Context) {
		repos, err := initDB.GetAllRepos()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, repos)
	})

	r.POST("/repos/update-all", func(c *gin.Context) {
		repos, err := initDB.GetAllRepos()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		var results = make(map[string]string)
		for _, repo := range repos {
			err := git.PullRepo(repo.Path)
			if err != nil {
				results[repo.URL] = "error: " + err.Error()
			} else {
				results[repo.URL] = "updated"
			}
		}
		c.JSON(http.StatusOK, results)
	})

	return r.Run(":4000")
}
