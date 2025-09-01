package fzf

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// Command handles fzf-based interactive selection
type Command struct{}

// Options configures fzf behavior
type Options struct {
	Query       string   // Initial query string
	MultiSelect bool     // Allow multiple selections
	Args        []string // Additional fzf arguments
}

// NewCommand creates a new fzf selector
func NewCommand() *Command {
	return &Command{}
}

// IsInstalled checks if fzf is installed on the system
func IsInstalled() bool {
	_, err := exec.LookPath("fzf")
	return err == nil
}

// Run displays items using fzf and returns the selected item(s)
func (s *Command) Run(input string, options *Options) (string, error) {
	if !IsInstalled() {
		return "", fmt.Errorf("fzf is not installed. Please install fzf first")
	}

	if options == nil {
		options = &Options{}
	}

	args := []string{}
	if !options.MultiSelect {
		args = append(args, "+m") // Single selection mode
	}
	if options.Query != "" {
		args = append(args, "-q", options.Query)
	}
	args = append(args, options.Args...)

	fzfCmd := exec.Command("fzf", args...)
	fzfCmd.Stdin = strings.NewReader(input)
	fzfCmd.Stderr = os.Stderr

	output, err := fzfCmd.Output()
	if err != nil {
		var exitErr *exec.ExitError
		if errors.As(err, &exitErr) && exitErr.ExitCode() == 130 {
			// User cancelled (Ctrl+C)
			return "", nil
		}
		return "", fmt.Errorf("failed to execute fzf: %w", err)
	}

	return strings.TrimSpace(string(output)), nil
}
