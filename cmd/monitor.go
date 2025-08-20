package cmd

import (
	"github.com/dyammarcano/clonr/internal/server"
	"github.com/spf13/cobra"
)

// monitorCmd represents the monitor command
var monitorCmd = &cobra.Command{
	Use:   "monitor",
	Short: "Start the clonr server to monitor and manage repositories via API.",
	Long: `Start the clonr built-in server. This allows you to monitor and manage registered repositories 
via a web API or other integrations. Useful for automation, dashboards, or remote management.`,
	RunE: server.StartServer,
}

func init() {
	rootCmd.AddCommand(monitorCmd)
}
