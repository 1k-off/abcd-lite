package oras

import (
	"fmt"
	"strings"

	"oras.land/oras/cmd/oras/root"
)

// CommandType represents the supported ORAS commands.
type CommandType string

const (
	Pull CommandType = "pull"
)

// OrasCommand holds the command type and its arguments.
type OrasCommand struct {
	Type        CommandType
	Registry    string
	Reference   string
	Username    string
	Password    string
	Output      string
	Concurrency int
}

// Args builds the CLI arguments for the given command.
func (c OrasCommand) Args() ([]string, error) {
	switch c.Type {
	case Pull:
		return c.pull()
	default:
		return nil, fmt.Errorf("unsupported command type: %s", c.Type)
	}
}

// Run executes the given OrasCommand and returns the output and any error encountered.
func Run(cmd OrasCommand) (string, error) {
	args, err := cmd.Args()
	if err != nil {
		return "", err
	}
	return runOrasCommand(args)
}

func (c OrasCommand) pull() ([]string, error) {
	if c.Reference == "" {
		return nil, fmt.Errorf("reference must be specified for Pull command")
	}
	args := []string{"pull"}
	if c.Output != "" {
		args = append(args, "--output", c.Output)
	}
	if c.Concurrency > 0 {
		args = append(args, "--concurrency", fmt.Sprintf("%d", c.Concurrency))
	}
	if c.Username != "" && c.Password != "" {
		args = append(args, "--username", c.Username, "--password", c.Password)
	}
	args = append(args, c.Reference)
	return args, nil
}

// runOrasCommand executes an oras CLI command with the given arguments.
// It returns the output and any error encountered.
func runOrasCommand(args []string) (string, error) {
	cmd := root.New()
	cmd.SetArgs(args)
	builder := &strings.Builder{}
	cmd.SetOut(builder)
	cmd.SetErr(builder)
	err := cmd.Execute()
	output := builder.String()
	if err != nil {
		return output, fmt.Errorf("failed to execute oras command: %w\nOutput: %s", err, output)
	}
	return output, nil
}
