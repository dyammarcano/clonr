package cmd

import (
	"fmt"
	"log"
	"strings"

	"github.com/dyammarcano/clonr/internal/git"

	"github.com/AlecAivazis/survey/v2"
	"github.com/spf13/cobra"
)

var removeCmd = &cobra.Command{
	Use:   "remove",
	Short: "Remove repositories from the registry (does not delete files)",
	Long: `Interactively select one or more repositories to remove from the clonr registry. This does not delete any 
files from disk, only removes the selected repositories from clonr's management.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		list, err := git.ListRepos()
		if err != nil {
			return err
		}

		if len(list) == 0 {
			log.Printf("No hay repositorios registrados.\n")
			return nil
		}

		options := make([]string, len(list))
		for i, r := range list {
			options[i] = fmt.Sprintf("%s -> %s", r.URL, r.Path)
		}

		var selected []string

		prompt := &survey.MultiSelect{
			Message: "Selecciona los repositorios a eliminar del registro:",
			Options: options,
		}

		if err = survey.AskOne(prompt, &selected); err != nil {
			return err
		}

		if len(selected) == 0 {
			log.Printf("No se seleccionó ningún repositorio.")
			return nil
		}

		confirm := false
		promptConfirm := &survey.Confirm{
			Message: fmt.Sprintf("Seguro quieres eliminar %d repositorios del registro?", len(selected)),
		}

		if err = survey.AskOne(promptConfirm, &confirm); err != nil {
			return err
		}

		if !confirm {
			log.Printf("Operación cancelada.\n")
			return nil
		}

		for _, sel := range selected {
			parts := strings.SplitN(sel, " -> ", 2)
			if len(parts) == 0 {
				continue
			}

			url := parts[0]
			if err := git.RemoveRepo(url); err != nil {
				log.Printf("Error removiendo %s: %v\n", url, err)
			} else {
				log.Printf("Repositorio removido: %s\n", url)
			}
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(removeCmd)
}
