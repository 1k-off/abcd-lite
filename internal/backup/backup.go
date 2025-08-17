package backup

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/1k-off/abcd-lite/internal/server/domain"
	"github.com/1k-off/abcd-lite/internal/storage"
)

// BackupInfo contains metadata about a backup file
type BackupInfo struct {
	Timestamp time.Time `json:"timestamp"`
	ItemCount int       `json:"itemCount"`
	Version   string    `json:"version"`
}

// RestoreOptions defines options for restore operations
type RestoreOptions struct {
	OverwriteExisting bool
	ValidateOnly      bool
	SkipErrors        bool
}

// RestoreResult contains the results of a restore operation
type RestoreResult struct {
	TotalItems    int      `json:"totalItems"`
	RestoredItems int      `json:"restoredItems"`
	SkippedItems  int      `json:"skippedItems"`
	FailedItems   int      `json:"failedItems"`
	Warnings      []string `json:"warnings"`
	Errors        []string `json:"errors"`
}

// ExportAllStorage exports all storage data to JSON files
func ExportAllStorage(projectStorage storage.ProjectStorage, settingsStorage storage.SettingsStorage, outputDir string) error {
	// Export projects
	if err := exportProjects(projectStorage, outputDir); err != nil {
		return fmt.Errorf("failed to export projects: %w", err)
	}

	// Export settings
	if err := exportSettings(settingsStorage, outputDir); err != nil {
		return fmt.Errorf("failed to export settings: %w", err)
	}

	return nil
}

// ExportAllStorageWithTimestamp exports all storage data to timestamped JSON files
func ExportAllStorageWithTimestamp(projectStorage storage.ProjectStorage, settingsStorage storage.SettingsStorage, outputDir string) error {
	// Generate a single timestamp for both files
	timestamp := time.Now().Format("2006-01-02_15-04-05")

	// Export projects with timestamp
	if err := exportProjectsWithTimestamp(projectStorage, outputDir, timestamp); err != nil {
		return fmt.Errorf("failed to export projects: %w", err)
	}

	// Export settings with timestamp
	if err := exportSettingsWithTimestamp(settingsStorage, outputDir, timestamp); err != nil {
		return fmt.Errorf("failed to export settings: %w", err)
	}

	return nil
}

// exportProjects exports all projects to a JSON file
func exportProjects(projectStorage storage.ProjectStorage, outputDir string) error {
	projects, err := projectStorage.GetProjects()
	if err != nil {
		return fmt.Errorf("failed to get projects: %w", err)
	}

	// Create backup info
	backupInfo := BackupInfo{
		Timestamp: time.Now(),
		ItemCount: len(projects),
		Version:   "1.0",
	}

	// Create export data structure
	exportData := struct {
		BackupInfo BackupInfo       `json:"backupInfo"`
		Projects   []domain.Project `json:"projects"`
	}{
		BackupInfo: backupInfo,
		Projects:   projects,
	}

	// Marshal to JSON
	data, err := json.MarshalIndent(exportData, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal projects: %w", err)
	}

	// Write to file
	outputPath := filepath.Join(outputDir, "projects.json")
	if err := os.WriteFile(outputPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write projects file: %w", err)
	}

	return nil
}

// exportProjectsWithTimestamp exports all projects to a timestamped JSON file
func exportProjectsWithTimestamp(projectStorage storage.ProjectStorage, outputDir, timestamp string) error {
	projects, err := projectStorage.GetProjects()
	if err != nil {
		return fmt.Errorf("failed to get projects: %w", err)
	}

	// Create backup info
	backupInfo := BackupInfo{
		Timestamp: time.Now(),
		ItemCount: len(projects),
		Version:   "1.0",
	}

	// Create export data structure
	exportData := struct {
		BackupInfo BackupInfo       `json:"backupInfo"`
		Projects   []domain.Project `json:"projects"`
	}{
		BackupInfo: backupInfo,
		Projects:   projects,
	}

	// Marshal to JSON
	data, err := json.MarshalIndent(exportData, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal projects: %w", err)
	}

	// Use provided timestamp for filename
	filename := fmt.Sprintf("projects_%s.json", timestamp)
	outputPath := filepath.Join(outputDir, filename)

	// Write to file
	if err := os.WriteFile(outputPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write projects file: %w", err)
	}

	return nil
}

// exportSettings exports all settings to a JSON file
func exportSettings(settingsStorage storage.SettingsStorage, outputDir string) error {
	allSettings, err := settingsStorage.GetAllSettings()
	if err != nil {
		return fmt.Errorf("failed to get settings: %w", err)
	}

	// Create backup info
	backupInfo := BackupInfo{
		Timestamp: time.Now(),
		ItemCount: len(allSettings),
		Version:   "1.0",
	}

	// Create export data structure
	exportData := struct {
		BackupInfo BackupInfo        `json:"backupInfo"`
		Settings   map[string]string `json:"settings"`
	}{
		BackupInfo: backupInfo,
		Settings:   allSettings,
	}

	// Marshal to JSON
	data, err := json.MarshalIndent(exportData, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal settings: %w", err)
	}

	// Write to file
	outputPath := filepath.Join(outputDir, "settings.json")
	if err := os.WriteFile(outputPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write settings file: %w", err)
	}

	return nil
}

// exportSettingsWithTimestamp exports all settings to a timestamped JSON file
func exportSettingsWithTimestamp(settingsStorage storage.SettingsStorage, outputDir, timestamp string) error {
	allSettings, err := settingsStorage.GetAllSettings()
	if err != nil {
		return fmt.Errorf("failed to get settings: %w", err)
	}

	// Create backup info
	backupInfo := BackupInfo{
		Timestamp: time.Now(),
		ItemCount: len(allSettings),
		Version:   "1.0",
	}

	// Create export data structure
	exportData := struct {
		BackupInfo BackupInfo        `json:"backupInfo"`
		Settings   map[string]string `json:"settings"`
	}{
		BackupInfo: backupInfo,
		Settings:   allSettings,
	}

	// Marshal to JSON
	data, err := json.MarshalIndent(exportData, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal settings: %w", err)
	}

	// Use provided timestamp for filename
	filename := fmt.Sprintf("settings_%s.json", timestamp)
	outputPath := filepath.Join(outputDir, filename)

	// Write to file
	if err := os.WriteFile(outputPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write settings file: %w", err)
	}

	return nil
}

// GetBackupFileName generates a timestamped backup filename
func GetBackupFileName(prefix, extension string) string {
	timestamp := time.Now().Format("2006-01-02_15-04-05")
	return fmt.Sprintf("%s_%s.%s", prefix, timestamp, extension)
}

// GetBackupInfo extracts backup information from a backup file
func GetBackupInfo(filePath string) (*BackupInfo, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read backup file: %w", err)
	}

	// Try to parse as projects backup first
	var projectsData struct {
		BackupInfo BackupInfo `json:"backupInfo"`
	}
	if err := json.Unmarshal(data, &projectsData); err == nil && projectsData.BackupInfo.Timestamp != (time.Time{}) {
		return &projectsData.BackupInfo, nil
	}

	// Try to parse as settings backup
	var settingsData struct {
		BackupInfo BackupInfo `json:"backupInfo"`
	}
	if err := json.Unmarshal(data, &settingsData); err == nil && settingsData.BackupInfo.Timestamp != (time.Time{}) {
		return &settingsData.BackupInfo, nil
	}

	return nil, fmt.Errorf("invalid backup file format")
}

// ValidateBackupFile validates that a backup file has the correct structure
func ValidateBackupFile(filePath string) error {
	_, err := GetBackupInfo(filePath)
	return err
}
