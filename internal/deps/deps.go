package deps

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
)

var binDir = filepath.Join("bin")

// GetDependencies downloads and installs all required dependencies
func GetDependencies() error {
	if err := os.MkdirAll(binDir, 0755); err != nil {
		return fmt.Errorf("failed to create bin directory: %w", err)
	}

	if !checkOrasExists() {
		if err := downloadOras(); err != nil {
			return fmt.Errorf("failed to download ORAS: %w", err)
		}
	}

	return nil
}

func getPlatformInfo() (string, string, error) {
	platform := ""
	arch := ""

	switch runtime.GOOS {
	case "windows":
		platform = "windows"
	case "darwin":
		platform = "darwin"
	case "linux":
		platform = "linux"
	default:
		return platform, arch, fmt.Errorf("unsupported platform: %s", platform)
	}

	switch runtime.GOARCH {
	case "amd64":
		arch = "amd64"
	case "arm64":
		arch = "arm64"
	default:
		return platform, arch, fmt.Errorf("unsupported architecture: %s", arch)
	}

	if platform == "windows" && arch == "arm64" {
		return platform, arch, fmt.Errorf("windows arm64 is not supported")
	}

	return platform, arch, nil
}
