package cmd

import (
	"fmt"

	"github.com/nek023/gconfig/internal/gcloud"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(listCmd)
}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all gcloud configurations",
	Long: `Display all gcloud configurations in a formatted table.
Shows configuration name, active status, account, and project.

Output format:
  • Active configuration marked with '*'
  • Inactive configurations marked with '-'
  • Shows associated GCP account and project

Examples:
  gconfig list                     # Display all configurations
  gconfig list | grep prod         # Filter output
  
Equivalent to:
  gcloud config configurations list`,
	RunE: runList,
}

// runList handles the list command execution
func runList(cmd *cobra.Command, args []string) error {
	gcloudCmd := gcloud.NewCommand(&gcloud.CommandOption{
		Stdin:  cmd.InOrStdin(),
		Stdout: cmd.OutOrStdout(),
		Stderr: cmd.ErrOrStderr(),
	})
	err := gcloudCmd.PrintConfigurations()
	if err != nil {
		return fmt.Errorf("failed to list configurations: %w", err)
	}
	return nil
}
