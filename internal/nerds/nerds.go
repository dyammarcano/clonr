package nerds

import (
	"fmt"
	"log"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/AlecAivazis/survey/v2"
	"github.com/dyammarcano/clonr/internal/bolt"
	"github.com/dyammarcano/clonr/internal/model"
	"github.com/dyammarcano/clonr/internal/params"
	"github.com/dyammarcano/clonr/internal/svc"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// Tool represents a nerd tool with an ID and a display name
type Tool struct {
	ID   string
	Name string
}

func Tools(_ *cobra.Command, _ []string) error {
	// Step 1: List repositories
	list, err := svc.ListRepos()
	if err != nil {
		return fmt.Errorf("failed to list repositories: %w", err)
	}

	if len(list) == 0 {
		log.Println("No hay repositorios registrados.")
		return nil
	}

	// Map UID → Repository and build display options
	repos := make(map[string]model.Repository, len(list))
	repoOptions := make([]string, len(list))
	for i, repo := range list {
		repos[repo.UID] = repo
		repoOptions[i] = fmt.Sprintf("[%s] [%s] %s -> %s",
			repo.UID,
			repo.UpdatedAt.Format(time.DateTime),
			repo.URL,
			repo.Path,
		)
	}

	// Step 2: Select repositories
	var selectedRepos []string
	if err := survey.AskOne(&survey.MultiSelect{
		Message: "Selecciona los repositorios para generar estadísticas:",
		Options: repoOptions,
	}, &selectedRepos); err != nil {
		return fmt.Errorf("selection prompt failed: %w", err)
	}

	if len(selectedRepos) == 0 {
		log.Println("No se seleccionó ningún repositorio.")
		return nil
	}

	// Step 3: Confirm action
	confirm := false
	if err := survey.AskOne(&survey.Confirm{
		Message: fmt.Sprintf("¿Seguro que quieres generar estadísticas para %d repositorios?", len(selectedRepos)),
	}, &confirm); err != nil {
		return fmt.Errorf("confirmation prompt failed: %w", err)
	}

	if !confirm {
		log.Println("Operación cancelada.")
		return nil
	}

	// Step 4: Run stats for each selected repo
	for _, repoSel := range selectedRepos {
		uid := extractBracketValue(repoSel)
		repo := repos[uid]

		log.Printf("Generando estadísticas.")
		if err := RunStats(repo); err != nil {
			log.Printf("❌ Error en %s: %v", repo.Path, err)
		}
	}

	return nil
}

// RunStats dispatches tool logic by its ID
func RunStats(repo model.Repository) error {
	metricsDB, err := bolt.InitBolt(filepath.Join(params.AppdataDir, "metrics"), repo.UID)
	if err != nil {
		return err
	}
	defer func(metricsDB *bolt.Bolt) {
		_ = metricsDB.Close()
	}(metricsDB)

	log.Printf("[STATS] Calculando estadísticas para %s", repo.Path)

	stats, err := svc.GetRepoStats(repo.Path)
	if err != nil {
		return err
	}

	if err := metricsDB.SaveStats(stats); err != nil {
		return err
	}

	if viper.GetBool("json") {
		data := stats.Bytes()
		if len(data) == 0 {
			return fmt.Errorf("no se pudieron generar estadísticas para el repositorio %s", repo.Path)
		}
		log.Printf("Estadísticas generadas: %s", string(data))

		return nil
	}

	// --- Commits por usuario ---
	log.Printf("--- Commits por usuario ---\n\n")
	for idx := range stats.CommitsByUser {
		log.Printf("- %s: %d commits\n", stats.CommitsByUser[idx].Item, stats.CommitsByUser[idx].Count)
	}

	// --- Archivos más modificados (Top 100) ---
	log.Printf("--- Archivos más modificados (Top 100) ---\n\n")

	fileStats := make([]model.Content, 0, len(stats.FileModifications))

	for _, item := range stats.FileModifications {
		fileStats = append(fileStats, model.Content{Item: item.Item, Count: item.Count})
	}

	sort.Slice(fileStats, func(i, j int) bool {
		return fileStats[i].Count > fileStats[j].Count
	})

	limit := 100
	if len(fileStats) < limit {
		limit = len(fileStats)
	}

	for i := 0; i < limit; i++ {
		log.Printf("- %s: %d modificaciones\n", fileStats[i].Item, fileStats[i].Count)
	}

	// --- Líneas de código ---
	log.Printf("--- Líneas de código ---\n\n")
	log.Printf("- Líneas añadidas: %d\n", stats.LinesAdded)
	log.Printf("- Líneas eliminadas: %d\n", stats.LinesDeleted)
	log.Printf("- Total de cambios: %d\n", stats.LinesAdded+stats.LinesDeleted)

	// --- Commits por día de la semana ---
	log.Printf("--- Commits por día de la semana ---\n\n")

	days := []time.Weekday{
		time.Sunday, time.Monday, time.Tuesday, time.Wednesday,
		time.Thursday, time.Friday, time.Saturday,
	}

	for _, day := range days {
		count := stats.CommitsByWeekday[day]
		log.Printf("- %s: %d commits\n", day, count)
	}

	return nil
}

// extractBracketValue safely extracts the value between the first "[" and "]"
func extractBracketValue(s string) string {
	start := strings.Index(s, "[")
	end := strings.Index(s, "]")
	if start == -1 || end == -1 || start >= end {
		return ""
	}
	return s[start+1 : end]
}
