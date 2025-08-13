package nerds

import (
	"fmt"
	"log"
	"os/exec"
	"sort"
	"strings"
	"time"

	"github.com/AlecAivazis/survey/v2"
	git2 "github.com/dyammarcano/clonr/internal/git"
	"github.com/dyammarcano/clonr/internal/model"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/spf13/cobra"
)

// Tool represents a nerd tool with an ID and a display name
type Tool struct {
	ID   string
	Name string
}

// RegisteredTools is where you add new nerd tools
var RegisteredTools = []Tool{
	{"stats", "Calcular estadísticas"},
	{"lint", "Ejecutar linter"},
	{"update", "Actualizar dependencias"},
	{"contributors", "Encontrar principales contribuidores"},
	{"most-modified", "Listar archivos más modificados"},
}

// StatsData almacena todos los datos de estadísticas del repositorio.
type StatsData struct {
	CommitsByUser     map[string]int
	FileModifications map[string]int
	LinesAdded        int
	LinesDeleted      int
	CommitsByWeekday  map[time.Weekday]int
}

// getRepoStats itera sobre los commits de un repositorio y recopila datos estadísticos.
func getRepoStats(repoPath string) (*StatsData, error) {
	r, err := git.PlainOpen(repoPath)
	if err != nil {
		return nil, fmt.Errorf("error al abrir el repositorio: %w", err)
	}

	commitIterator, err := r.Log(&git.LogOptions{})
	if err != nil {
		return nil, fmt.Errorf("error al obtener el historial de commits: %w", err)
	}

	stats := &StatsData{
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

func Tools(cmd *cobra.Command, args []string) error {
	// Step 1: List repositories
	list, err := git2.ListRepos()
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
		Message: "Selecciona los repositorios para ejecutar herramientas nerds:",
		Options: repoOptions,
	}, &selectedRepos); err != nil {
		return fmt.Errorf("selection prompt failed: %w", err)
	}

	if len(selectedRepos) == 0 {
		log.Println("No se seleccionó ningún repositorio.")
		return nil
	}

	// Step 3: Select nerd tools
	toolOptions := make([]string, len(RegisteredTools))

	for i, t := range RegisteredTools {
		toolOptions[i] = fmt.Sprintf("[%s] %s", t.ID, t.Name)
	}

	var selectedTools []string

	if err := survey.AskOne(&survey.MultiSelect{
		Message: "Selecciona las herramientas nerds a ejecutar:",
		Options: toolOptions,
	}, &selectedTools); err != nil {
		return fmt.Errorf("tool selection prompt failed: %w", err)
	}

	if len(selectedTools) == 0 {
		log.Println("No se seleccionó ninguna herramienta.")
		return nil
	}

	// Step 4: Confirm action
	confirm := false
	if err := survey.AskOne(&survey.Confirm{
		Message: fmt.Sprintf("¿Seguro que quieres ejecutar %d herramientas sobre %d repositorios?", len(selectedTools), len(selectedRepos)),
	}, &confirm); err != nil {
		return fmt.Errorf("confirmation prompt failed: %w", err)
	}

	if !confirm {
		log.Println("Operación cancelada.")
		return nil
	}

	// Step 5: Run tools
	for _, repoSel := range selectedRepos {
		uid := extractBracketValue(repoSel)
		repo := repos[uid]

		for _, toolSel := range selectedTools {
			toolID := extractBracketValue(toolSel)
			log.Printf("Ejecutando herramienta '%s' sobre %s (%s)\n", toolID, repo.UID, repo.Path)
			if err := RunToolByID(toolID, repo); err != nil {
				log.Printf("Error ejecutando herramienta %s: %v", toolID, err)
			}
		}
	}

	return nil
}

// extractBracketValue safely extracts the value between first "[" and "]"
func extractBracketValue(s string) string {
	start := strings.Index(s, "[")
	end := strings.Index(s, "]")
	if start == -1 || end == -1 || start >= end {
		return ""
	}
	return s[start+1 : end]
}

// RunToolByID dispatches tool logic by its ID
func RunToolByID(id string, repo model.Repository) error {
	switch id {
	case "stats":
		log.Printf("[STATS] Calculando estadísticas para %s", repo.Path)
		stats, err := getRepoStats(repo.Path)
		if err != nil {
			return err
		}

		log.Println("\n--- Estadísticas del Repositorio ---")
		log.Println("\nCommits por usuario:")
		for user, count := range stats.CommitsByUser {
			log.Printf("- %s: %d\n", user, count)
		}

		log.Println("\nArchivos más modificados (Top 10):")
		type fileStat struct {
			name  string
			count int
		}
		fileStats := make([]fileStat, 0, len(stats.FileModifications))
		for name, count := range stats.FileModifications {
			fileStats = append(fileStats, fileStat{name, count})
		}
		sort.Slice(fileStats, func(i, j int) bool {
			return fileStats[i].count > fileStats[j].count
		})
		limit := 10
		if len(fileStats) < limit {
			limit = len(fileStats)
		}
		for i := 0; i < limit; i++ {
			log.Printf("- %s: %d modificaciones\n", fileStats[i].name, fileStats[i].count)
		}

		log.Println("\nLíneas de código:")
		log.Printf("- Líneas añadidas: %d\n", stats.LinesAdded)
		log.Printf("- Líneas eliminadas: %d\n", stats.LinesDeleted)
		log.Printf("- Total de cambios: %d\n", stats.LinesAdded+stats.LinesDeleted)

		log.Println("\nCommits por día de la semana:")
		days := []time.Weekday{
			time.Sunday, time.Monday, time.Tuesday, time.Wednesday,
			time.Thursday, time.Friday, time.Saturday,
		}
		for _, day := range days {
			count := stats.CommitsByWeekday[day]
			log.Printf("- %s: %d commits\n", day, count)
		}
	case "lint":
		log.Printf("[LINT] Ejecutando linter para %s", repo.Path)
		cmd := exec.Command("golangci-lint", "run", "./...")
		cmd.Dir = repo.Path
		output, err := cmd.CombinedOutput()
		if err != nil {
			log.Printf("Error ejecutando linter: %v", err)
			log.Printf("Salida de linter:\n%s\n", string(output))
			return err
		}
		log.Printf("Linter ejecutado con éxito:\n%s\n", string(output))
	case "update":
		log.Printf("[UPDATE] Actualizando dependencias para %s", repo.Path)
		cmd := exec.Command("go", "mod", "tidy")
		cmd.Dir = repo.Path
		output, err := cmd.CombinedOutput()
		if err != nil {
			log.Printf("Error actualizando dependencias: %v", err)
			log.Printf("Salida de actualización:\n%s\n", string(output))
			return err
		}
		log.Printf("Dependencias actualizadas con éxito:\n%s\n", string(output))
	case "contributors":
		log.Printf("[CONTRIBUTORS] Buscando contribuidores para %s", repo.Path)
		stats, err := getRepoStats(repo.Path)
		if err != nil {
			return err
		}
		log.Println("\n--- Principales Contribuidores ---")
		for user, count := range stats.CommitsByUser {
			log.Printf("- %s: %d commits\n", user, count)
		}
	case "most-modified":
		log.Printf("[MOST-MODIFIED] Buscando archivos más modificados para %s", repo.Path)
		stats, err := getRepoStats(repo.Path)
		if err != nil {
			return err
		}
		log.Println("\n--- Archivos más modificados (Top 10) ---")
		type fileStat struct {
			name  string
			count int
		}
		fileStats := make([]fileStat, 0, len(stats.FileModifications))
		for name, count := range stats.FileModifications {
			fileStats = append(fileStats, fileStat{name, count})
		}
		sort.Slice(fileStats, func(i, j int) bool {
			return fileStats[i].count > fileStats[j].count
		})
		limit := 10
		if len(fileStats) < limit {
			limit = len(fileStats)
		}
		for i := 0; i < limit; i++ {
			log.Printf("- %s: %d modificaciones\n", fileStats[i].name, fileStats[i].count)
		}
	default:
		return fmt.Errorf("herramienta desconocida: %s", id)
	}
	return nil
}
