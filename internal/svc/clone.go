package svc

import (
	"fmt"
	"io"
	"log"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/dyammarcano/clonr/internal/db"
	"github.com/spf13/cobra"
)

// cloneOutputMsg es un mensaje para actualizar la TUI con la salida del comando git.
type cloneOutputMsg string

// cloneFinishedMsg es un mensaje enviado cuando el comando git ha terminado.
type cloneFinishedMsg struct{ err error }

// TUIModel es el modelo de datos para nuestra interfaz de usuario en terminal.
type TUIModel struct {
	repoURL  string
	repoPath string
	spinner  spinner.Model
	output   strings.Builder
	status   string
	err      error
	quitting bool
}

// Initializa un nuevo TUIModel.
func newTUIModel(repoURL, repoPath string) TUIModel {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = s.Style.Foreground(lipgloss.Color("205"))
	return TUIModel{
		repoURL:  repoURL,
		repoPath: repoPath,
		spinner:  s,
		status:   fmt.Sprintf("Clonando %s a %s...", repoURL, repoPath),
	}
}

// Init inicializa el modelo y ejecuta el comando git clone.
func (m TUIModel) Init() tea.Cmd {
	return tea.Batch(m.startClone(), m.spinner.Tick)
}

// Update maneja los mensajes entrantes para actualizar el modelo.
func (m TUIModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			m.quitting = true
			return m, tea.Quit
		}
	case cloneOutputMsg:
		m.output.WriteString(string(msg))
		m.status = "Clonando..."
	case cloneFinishedMsg:
		m.err = msg.err
		if m.err == nil {
			m.status = "Clonación completada con éxito."
		} else {
			m.status = "Clonación fallida."
		}
		return m, tea.Quit
	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	}
	return m, nil
}

// View renderiza la interfaz de usuario.
func (m TUIModel) View() string {
	if m.quitting {
		return "Operación cancelada.\n"
	}

	var s string
	if m.err != nil {
		s = fmt.Sprintf("\n%s %s\n%s\n", m.spinner.View(), m.status, m.err)
	} else {
		s = fmt.Sprintf("\n%s %s\n\n%s\n", m.spinner.View(), m.status, m.output.String())
	}
	return s
}

// startClone ejecuta el comando git clone y envía la salida al modelo.
func (m TUIModel) startClone() tea.Cmd {
	return func() tea.Msg {
		runCmd := exec.Command("git", "clone", m.repoURL, m.repoPath)
		stdout, err := runCmd.StdoutPipe()
		if err != nil {
			return cloneFinishedMsg{fmt.Errorf("error creando pipe de stdout: %w", err)}
		}
		stderr, err := runCmd.StderrPipe()
		if err != nil {
			return cloneFinishedMsg{fmt.Errorf("error creando pipe de stderr: %w", err)}
		}

		if err := runCmd.Start(); err != nil {
			return cloneFinishedMsg{fmt.Errorf("error al iniciar el comando git clone: %w", err)}
		}

		// Leer la salida de stdout y stderr en tiempo real.
		reader := io.MultiReader(stdout, stderr)
		buf := make([]byte, 1024)
		for {
			n, err := reader.Read(buf)
			if n > 0 {
				if msg := string(buf[:n]); msg != "" {
					tea.Println(cloneOutputMsg(msg))
				}
			}
			if err == io.EOF {
				break
			}
			if err != nil {
				return cloneFinishedMsg{fmt.Errorf("error leyendo la salida del comando: %w", err)}
			}
		}

		err = runCmd.Wait()
		return cloneFinishedMsg{err}
	}
}

func CloneRepo(_ *cobra.Command, args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("repository URL is required")
	}

	u, err := url.Parse(strings.TrimSpace(args[0]))
	if err != nil {
		return fmt.Errorf("error parsing repository URL: %w", err)
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

	savePath := filepath.Join(absPath, extractRepoName(u.String()))

	// Inicia el programa TUI para el git clone.
	model := newTUIModel(u.String(), savePath)
	p := tea.NewProgram(model)
	finalModel, err := p.Run()
	if err != nil {
		return fmt.Errorf("error running TUI program: %w", err)
	}

	// Comprueba si la operación fue exitosa.
	if finalModel.(TUIModel).err != nil {
		return finalModel.(TUIModel).err
	}

	if err := initDB.SaveRepo(u.String(), savePath); err != nil {
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
