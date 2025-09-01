package cmd

import (
	"fmt"
	"strings"

	"github.com/nek023/gconfig/internal/fzf"
	"github.com/nek023/gconfig/internal/gcloud"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(switchCmd)
}

var switchCmd = &cobra.Command{
	Use:   "switch [query]",
	Short: "Switch gcloud configuration interactively",
	Long: `Switch between gcloud configurations using an interactive fzf selector.
The currently active configuration is highlighted in red.

Examples:
  gconfig switch              # Show all configurations
  gconfig switch prod         # Filter configurations containing "prod"
  gconfig switch my-project   # Filter configurations containing "my-project"

Requirements:
  • fzf must be installed`,
	Args: cobra.MaximumNArgs(1),
	RunE: runSwitch,
}

func runSwitch(cmd *cobra.Command, args []string) error {
	gcloudCmd := gcloud.NewCommand(&gcloud.CommandOption{
		Stdin:  cmd.InOrStdin(),
		Stdout: cmd.OutOrStdout(),
		Stderr: cmd.ErrOrStderr(),
	})
	fzfCmd := fzf.NewCommand()

	if !fzf.IsInstalled() {
		return fmt.Errorf("fzf is not installed")
	}

	configTable, err := gcloudCmd.GetConfigurationsTable()
	if err != nil {
		return fmt.Errorf("failed to get configurations table: %w", err)
	}
	configTable = highlightActiveConfiguration(configTable)

	var query string
	if len(args) > 0 {
		query = args[0]
	}

	selected, err := fzfCmd.Run(configTable, &fzf.Options{Query: query})
	if err != nil {
		return err
	}

	configName, err := extractConfigurationName(selected)
	if err != nil {
		return err
	}

	if err := gcloudCmd.ActivateConfiguration(configName); err != nil {
		return fmt.Errorf("failed to switch configuration: %w", err)
	}

	return nil
}

func highlightActiveConfiguration(input string) string {
	const (
		prefix = "*"
		red    = "\033[01;31m" // Bold red
		reset  = "\033[0m"     // Reset color
	)

	lines := strings.Split(input, "\n")
	var result []string

	for _, line := range lines {
		if strings.HasPrefix(line, prefix) {
			// Highlight the prefix in red
			highlighted := red + prefix + reset + line[1:]
			result = append(result, highlighted)
		} else {
			result = append(result, line)
		}
	}

	return strings.Join(result, "\n")
}

func extractConfigurationName(selected string) (string, error) {
	fields := strings.Fields(selected)
	if len(fields) < 2 {
		return "", fmt.Errorf("invalid selection result")
	}
	return fields[1], nil
}
