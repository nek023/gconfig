package cmd

import (
	"fmt"

	"github.com/nek023/gconfig/internal/gcloud"
	"github.com/samber/lo"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(exportCmd)
}

var exportCmd = &cobra.Command{
	Use:   "export [configuration]",
	Short: "Export gcloud configurations to YAML",
	Long: `Export gcloud configurations in YAML format to stdout.
Exports either a specific configuration or all configurations.

Examples:
  gconfig export                    # Export all configurations
  gconfig export my-config          # Export specific configuration
  gconfig export > configs.yaml     # Save to file
  gconfig export | pbcopy           # Copy to clipboard (macOS)

Output format:
  • Single configuration: YAML document
  • Multiple configurations: Multi-document YAML (separated by ---)`,
	Args: cobra.MaximumNArgs(1),
	RunE: runExport,
}

func runExport(cmd *cobra.Command, args []string) error {
	gcloudCmd := gcloud.NewCommand(&gcloud.CommandOption{
		Stdin:  cmd.InOrStdin(),
		Stdout: cmd.OutOrStdout(),
		Stderr: cmd.ErrOrStderr(),
	})

	var configNames []string
	if len(args) > 0 {
		configNames = append(configNames, args[0])
	} else {
		configs, err := gcloudCmd.ListConfigurations()
		if err != nil {
			return fmt.Errorf("failed to list configurations: %w", err)
		}
		if len(configs) == 0 {
			return fmt.Errorf("no configurations found")
		}

		configNames = lo.Map(configs, func(c gcloud.Configuration, index int) string {
			return c.Name
		})
	}

	return exportConfigurations(gcloudCmd, configNames)
}

func exportConfigurations(gcloudCmd gcloud.Command, configNames []string) error {
	for i, configName := range configNames {
		if i > 0 {
			fmt.Println("---")
		}

		output, err := gcloudCmd.DescribeConfiguration(configName)
		if err != nil {
			return fmt.Errorf("failed to export configuration '%s': %w", configName, err)
		}
		fmt.Println(output)
	}

	return nil
}
