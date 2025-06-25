package deployment

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"github.com/1k-off/abcd-lite/internal/oras"
)

type Options struct {
	Destination string
	Clean       bool
	Exclude     []string
}

func DefaultOptions() Options {
	return Options{
		Destination: "",
		Clean:       false,
		Exclude:     []string{},
	}
}

func NewOptions(clean bool, destination string, exclude []string) Options {
	return Options{
		Destination: destination,
		Clean:       clean,
		Exclude:     exclude,
	}
}

func Deploy(o Options, info PackageInfo) error {
	if o.Destination == "" {
		return errors.New("destination is required")
	}

	if o.Clean {
		if err := clean(o.Destination, o.Exclude); err != nil {
			return fmt.Errorf("failed to clean destination: %w", err)
		}
	}
	cmd := oras.OrasCommand{
		Type:      oras.Pull,
		Reference: info.PackageRef,
		Username:  info.Credentials.Username,
		Password:  info.Credentials.Password,
		Output:    o.Destination,
	}

	_, err := oras.Run(cmd)
	if err != nil {
		return fmt.Errorf("oras pull failed: %w", err)
	}

	return nil
}

func clean(destination string, exclude []string) error {
	files, err := os.ReadDir(destination)
	if err != nil {
		return fmt.Errorf("failed to read destination: %w", err)
	}

	// Normalize exclude list: remove trailing slashes
	normalizedExclude := make([]string, len(exclude))
	for i, ex := range exclude {
		normalizedExclude[i] = strings.TrimRight(ex, "/\\")
	}

	for _, file := range files {
		name := file.Name()
		if !slices.Contains(normalizedExclude, name) {
			os.RemoveAll(filepath.Join(destination, name))
		}
	}
	return nil
}
