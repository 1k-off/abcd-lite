package deps

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
)

const (
	orasVersion = "1.2.3"
	orasOwner   = "oras-project"
	orasRepo    = "oras"
)

// OrasBinPath returns the path to the ORAS binary
func OrasBinPath() string {
	orasPath := filepath.Join(binDir, "oras")
	if runtime.GOOS == "windows" {
		orasPath += ".exe"
	}
	return orasPath
}

// downloadOras downloads and extracts ORAS binary for the current platform
func downloadOras() error {
	platform, arch, err := getPlatformInfo()
	if err != nil {
		return err
	}

	// Map platform and set extension
	var ext string
	switch platform {
	case "windows":
		platform = "windows"
		ext = "zip"
	case "darwin":
		platform = "darwin"
		ext = "tar.gz"
	case "linux":
		platform = "linux"
		ext = "tar.gz"
	default:
		return fmt.Errorf("oras not supported for %s %s", platform, arch)
	}

	filename := fmt.Sprintf("oras_%s_%s_%s.%s", orasVersion, platform, arch, ext)
	url := fmt.Sprintf("https://github.com/%s/%s/releases/download/v%s/%s",
		orasOwner, orasRepo, orasVersion, filename)

	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("failed to download ORAS: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to download ORAS: status code %d", resp.StatusCode)
	}

	tmpFile, err := os.CreateTemp("", "oras-*."+ext)
	if err != nil {
		return fmt.Errorf("failed to create temp file: %w", err)
	}
	defer os.Remove(tmpFile.Name())
	defer tmpFile.Close()

	if _, err := io.Copy(tmpFile, resp.Body); err != nil {
		return fmt.Errorf("failed to save downloaded file: %w", err)
	}

	if _, err := tmpFile.Seek(0, 0); err != nil {
		return fmt.Errorf("failed to reset file pointer: %w", err)
	}
	// Extract the archive
	if ext == "zip" {
		if err := extractZip(tmpFile); err != nil {
			return fmt.Errorf("failed to extract zip: %w", err)
		}
	} else {
		if err := extractTarGz(tmpFile); err != nil {
			return fmt.Errorf("failed to extract tar.gz: %w", err)
		}
	}

	return nil
}

// checkOrasExists checks if ORAS exists and is executable by running 'oras version'
func checkOrasExists() bool {
	if _, err := os.Stat(OrasBinPath()); err != nil {
		return false
	}
	cmd := exec.Command(OrasBinPath(), "version")
	if err := cmd.Run(); err != nil {
		return false
	}

	return true
}
