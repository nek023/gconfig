package cmd

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/goccy/go-yaml"
	"github.com/nek023/gconfig/internal/gcloud"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(importCmd)
}

var importCmd = &cobra.Command{
	Use:   "import [file]",
	Short: "Import gcloud configurations from YAML",
	Long: `Import gcloud configurations from YAML format.
Reads from a file or stdin and applies settings to gcloud.

Behavior:
  • Existing configurations: Activates and updates properties
  • New configurations: Creates and sets properties
  • Supports single and multi-document YAML files

Examples:
  gconfig import configs.yaml              # Import from file
  cat configs.yaml | gconfig import        # Import from stdin
  gconfig export | gconfig import          # Copy all configurations
  curl -s $URL | gconfig import            # Import from URL

Input format:
  • Accepts output from 'gconfig export'
  • Supports YAML with nested properties`,
	Args: cobra.MaximumNArgs(1),
	RunE: runImport,
}

// runImport handles the import command execution
func runImport(cmd *cobra.Command, args []string) error {
	var (
		data []byte
		err  error
	)
	if len(args) > 0 {
		// ファイルから読み込み
		data, err = os.ReadFile(args[0])
		if err != nil {
			return fmt.Errorf("failed to read file: %w", err)
		}
	} else {
		// 標準入力から読み込み
		data, err = io.ReadAll(os.Stdin)
		if err != nil {
			return fmt.Errorf("failed to read from stdin: %w", err)
		}
	}

	var configs []gcloud.Configuration
	dec := yaml.NewDecoder(bytes.NewReader(data))
	for {
		var config gcloud.Configuration
		err := dec.Decode(&config)
		if err == io.EOF {
			break
		} else if err != nil {
			return fmt.Errorf("failed to decode YAML: %w", err)
		}
		if config.Name == "" {
			continue
		}
		configs = append(configs, config)
	}

	gcloudCmd := gcloud.NewCommand(&gcloud.CommandOption{
		Stdin:  cmd.InOrStdin(),
		Stdout: cmd.OutOrStdout(),
		Stderr: cmd.ErrOrStderr(),
	})
	for _, config := range configs {
		err := importConfiguration(gcloudCmd, config)
		if err != nil {
			return fmt.Errorf("failed to import configuration '%s': %w", config.Name, err)
		}
	}

	return nil
}

// importConfiguration imports a single configuration to gcloud
func importConfiguration(gcloudCmd gcloud.Command, config gcloud.Configuration) error {
	existingConfig, err := gcloudCmd.DescribeConfiguration(config.Name)
	exists := existingConfig != "" && err == nil

	if exists {
		err := gcloudCmd.ActivateConfiguration(config.Name)
		if err != nil {
			return fmt.Errorf("failed to activate configuration: %w", err)
		}
	} else {
		if err := gcloudCmd.CreateConfiguration(config.Name); err != nil {
			return fmt.Errorf("failed to create configuration: %w", err)
		}
	}

	if err := applyProperties(gcloudCmd, config.Name, config.Properties); err != nil {
		return fmt.Errorf("failed to apply properties: %w", err)
	}

	return nil
}

// applyProperties applies all properties to a gcloud configuration
func applyProperties(gcloudCmd gcloud.Command, configName string, properties map[string]any) error {
	return applyPropertiesRecursive(gcloudCmd, configName, "", properties)
}

// applyPropertiesRecursive recursively applies nested properties to a gcloud configuration
func applyPropertiesRecursive(gcloudCmd gcloud.Command, configName, path string, value any) error {
	switch v := value.(type) {
	case map[string]interface{}:
		for key, subValue := range v {
			newPath := key
			if path != "" {
				newPath = path + "/" + key
			}
			err := applyPropertiesRecursive(gcloudCmd, configName, newPath, subValue)
			if err != nil {
				return err
			}
		}
	case string, int, bool, float64:
		if path != "" && v != nil {
			path = strings.TrimPrefix(path, "core/") // core/は省略する
			valueStr := fmt.Sprintf("%v", v)
			err := gcloudCmd.SetConfigValue(configName, path, valueStr)
			if err != nil {
				return fmt.Errorf("failed to set %s: %w", path, err)
			}
		}
	}
	return nil
}
