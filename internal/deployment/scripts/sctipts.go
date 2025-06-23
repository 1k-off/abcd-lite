package scripts

import (
	"embed"
	"fmt"
	"os"
)

const (
	ApplicationPoolScript = "ApplicationPool.ps1"
)

//go:embed *.ps1
var psScripts embed.FS

// ExtractScript extracts the embedded PowerShell script to a temporary file
func ExtractScript(scriptName string) (string, error) {
	scriptContent, err := psScripts.ReadFile(scriptName)
	if err != nil {
		return "", fmt.Errorf("failed to read embedded script %s: %w", scriptName, err)
	}

	tmpFile, err := os.CreateTemp("", fmt.Sprintf("abcd-lite-%s-*.ps1", scriptName))
	if err != nil {
		return "", fmt.Errorf("failed to create temp file: %w", err)
	}

	if _, err := tmpFile.Write(scriptContent); err != nil {
		tmpFile.Close()
		os.Remove(tmpFile.Name())
		return "", fmt.Errorf("failed to write script to temp file: %w", err)
	}

	if err := tmpFile.Close(); err != nil {
		os.Remove(tmpFile.Name())
		return "", fmt.Errorf("failed to close temp file: %w", err)
	}

	return tmpFile.Name(), nil
}
