package deployment

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
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
	// Normalize the destination path
	destPath := filepath.Clean(destination)
	// Check if destination exists
	info, err := os.Stat(destPath)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("failed to read destination: %w", err)
		}
		return err
	}
	if !info.IsDir() {
		return fmt.Errorf("destination is not a directory: %s", destPath)
	}
	// Normalize exclude list: remove trailing slashes and use filepath separators
	normalizedExclude := make(map[string]struct{})
	parentDirs := make(map[string]struct{})
	for _, ex := range exclude {
		exPath := filepath.Clean(strings.TrimRight(ex, "/\\"))
		normalizedExclude[exPath] = struct{}{}
		// Add all parent directories of the excluded path
		for p := filepath.Dir(exPath); p != "." && p != string(os.PathSeparator); p = filepath.Dir(p) {
			parentDirs[p] = struct{}{}
			if p == filepath.Dir(p) {
				break // reached root
			}
		}
	}

	var dirsToDelete []string

	err = filepath.WalkDir(destination, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if path == destination {
			return nil // don't delete the root
		}
		relPath, err := filepath.Rel(destination, path)
		if err != nil {
			return err
		}
		// If this path is in the exclude set, skip deleting it and, if dir, skip its children
		if _, found := normalizedExclude[relPath]; found {
			if d.IsDir() {
				return filepath.SkipDir // skip this dir and its children
			}
			return nil // skip this file
		}
		// If any ancestor of this path is in the exclude set, skip deleting, but keep traversing
		for p := filepath.Dir(relPath); p != "." && p != string(os.PathSeparator); p = filepath.Dir(p) {
			if _, found := normalizedExclude[p]; found {
				return nil // skip deleting, but keep traversing
			}
			if p == filepath.Dir(p) {
				break // reached root
			}
		}
		if d.IsDir() {
			// Only delete if not a parent of any excluded path
			if _, isParent := parentDirs[relPath]; !isParent {
				dirsToDelete = append(dirsToDelete, path)
			}
			return nil // don't delete now, delete after walk
		}
		// Remove file
		return os.RemoveAll(path)
	})
	if err != nil {
		return err
	}
	// Delete directories in reverse order (deepest first)
	for i := len(dirsToDelete) - 1; i >= 0; i-- {
		err := os.RemoveAll(dirsToDelete[i])
		if err != nil {
			return err
		}
	}
	return nil
}
