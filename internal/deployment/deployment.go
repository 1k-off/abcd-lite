package deployment

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"slices"
	"strings"

	"github.com/1k-off/abcd-lite/internal/deps"
	"github.com/gofiber/fiber/v3/log"
)

type Options struct {
	Destination string
	Concurrency int
	Clean       bool
	Exclude     []string
}

func DefaultOptions() Options {
	return Options{
		Destination: "",
		Concurrency: 3,
		Clean:       false,
		Exclude:     []string{},
	}
}

func NewOptions(clean bool, destination string, concurrency int, exclude []string) Options {
	return Options{
		Destination: destination,
		Concurrency: concurrency,
		Clean:       clean,
		Exclude:     exclude,
	}
}

func Deploy(o Options, info PackageInfo) error {
	if o.Destination == "" {
		return errors.New("destination is required")
	}
	if o.Concurrency <= 0 {
		return errors.New("concurrency must be greater than 0")
	}

	args := []string{"pull"}
	args = append(args, "--output", o.Destination, "--concurrency", fmt.Sprintf("%d", o.Concurrency))
	if info.Credentials.Username != "" && info.Credentials.Password != "" {
		args = append(args, "--username", info.Credentials.Username, "--password", info.Credentials.Password)
	}
	if info.Credentials.LoginServer == "" {
		info.Credentials.LoginServer = defaultLoginServer
	}
	args = append(args, fmt.Sprintf("%s/%s:%s", info.Credentials.LoginServer, info.Name, info.Version))
	orasPath := deps.OrasBinPath()

	if o.Clean {
		if err := clean(o.Destination, o.Exclude); err != nil {
			return fmt.Errorf("failed to clean destination: %w", err)
		}
	}

	cmd := exec.Command(orasPath, args...)

	var stdout, stderr strings.Builder
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		log.Debug(cmd.String())
		log.Debug(stdout.String())
		log.Debug(stderr.String())

		return fmt.Errorf("failed to deploy package: %w", err)
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
