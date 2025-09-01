package cmd

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "gconfig",
	Short: "Enhanced gcloud configuration management tool",
	Long: `gconfig is a CLI tool that enhances Google Cloud SDK configuration management
with interactive selection, bulk import/export, and streamlined workflows.

Features:
  • Interactive configuration switching with fzf
  • Export configurations to YAML format
  • Import configurations from YAML files
  • List all available configurations

Requirements:
  • gcloud CLI installed and configured
  • fzf (for interactive switching)`,
}

// Execute runs the root command and handles any errors
func Execute() error {
	return rootCmd.Execute()
}
