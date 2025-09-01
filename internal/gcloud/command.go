package gcloud

import (
	"io"
	"os/exec"
	"strings"
)

// Configuration represents a gcloud configuration with its properties
type Configuration struct {
	Name       string
	IsActive   bool
	Properties map[string]any `yaml:"properties"`
}

// Command defines the interface for interacting with gcloud configurations
type Command interface {
	PrintConfigurations() error
	GetConfigurationsTable() (string, error)
	ListConfigurations() ([]Configuration, error)
	DescribeConfiguration(name string) (string, error)
	ActivateConfiguration(name string) error
	CreateConfiguration(name string) error
	SetConfigValue(configuration, path, value string) error
}

type command struct {
	Stdin  io.Reader
	Stdout io.Writer
	Stderr io.Writer
}

// CommandOption holds the I/O configuration for gcloud commands
type CommandOption struct {
	Stdin  io.Reader
	Stdout io.Writer
	Stderr io.Writer
}

// NewCommand creates a new gcloud command executor with the given I/O configuration
func NewCommand(opt *CommandOption) Command {
	return &command{
		Stdin:  opt.Stdin,
		Stdout: opt.Stdout,
		Stderr: opt.Stderr,
	}
}

func (c *command) runCommand(args ...string) error {
	cmd := exec.Command("gcloud", args...)
	cmd.Stdout = c.Stdout
	cmd.Stderr = c.Stderr

	err := cmd.Run()
	if err != nil {
		return err
	}

	return nil
}

func (c *command) getCommandOutput(args ...string) (string, error) {
	cmd := exec.Command("gcloud", args...)

	b, err := cmd.Output()
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(string(b)), nil
}

// PrintConfigurations displays all gcloud configurations in a formatted table
func (c *command) PrintConfigurations() error {
	return c.runCommand("config", "configurations", "list",
		"--format=table(is_active.yesno(yes=\"*\",no=\"-\"), name, properties.core.account, properties.core.project.yesno(no=\"-\"))")
}

// GetConfigurationsTable returns all gcloud configurations as a formatted table string
func (c *command) GetConfigurationsTable() (string, error) {
	return c.getCommandOutput("config", "configurations", "list",
		"--format=table[no-heading](is_active.yesno(yes=\"*\",no=\"-\"), name, properties.core.account, properties.core.project.yesno(no=\"-\"))")
}

// ListConfigurations retrieves all gcloud configurations as structured data
func (c *command) ListConfigurations() ([]Configuration, error) {
	output, err := c.getCommandOutput("config", "configurations", "list", "--format=csv[no-heading](name,is_active.yesno(yes='true',no='false'),properties.core.account,properties.core.project)")
	if err != nil {
		return nil, err
	}

	lines := strings.Split(output, "\n")
	if len(lines) < 1 {
		return []Configuration{}, nil
	}

	var configs []Configuration
	for i := 0; i < len(lines); i++ {
		line := lines[i]
		if line == "" {
			continue
		}

		parts := strings.Split(line, ",")
		if len(parts) != 4 {
			continue
		}

		name := parts[0]
		isActive := parts[1] == "true"
		account := parts[2]
		project := parts[3]

		configs = append(configs, Configuration{
			Name:     name,
			IsActive: isActive,
			Properties: map[string]any{
				"core": map[string]any{
					"account": account,
					"project": project,
				},
			},
		})
	}

	return configs, nil
}

// DescribeConfiguration returns the YAML representation of a specific configuration
func (c *command) DescribeConfiguration(name string) (string, error) {
	return c.getCommandOutput("config", "configurations", "describe", name)
}

// ActivateConfiguration switches to the specified gcloud configuration
func (c *command) ActivateConfiguration(name string) error {
	return c.runCommand("config", "configurations", "activate", name)
}

// CreateConfiguration creates a new gcloud configuration with the specified name
func (c *command) CreateConfiguration(name string) error {
	return c.runCommand("config", "configurations", "create", name)
}

// SetConfigValue sets a property value for the specified configuration
func (c *command) SetConfigValue(configuration, path, value string) error {
	return c.runCommand("config", "set", path, value, "--configuration", configuration)
}
