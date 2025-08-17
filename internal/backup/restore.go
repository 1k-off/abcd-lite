package backup

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/1k-off/abcd-lite/internal/server/domain"
	"github.com/1k-off/abcd-lite/internal/server/services"
	"github.com/1k-off/abcd-lite/internal/storage"
)

// RestoreAllStorage restores all storage data from backup files
func RestoreAllStorage(
	projectService services.ProjectService,
	settingsStorage storage.SettingsStorage,
	projectsPath, settingsPath string,
	options RestoreOptions,
) (*RestoreResult, error) {
	result := &RestoreResult{}

	// Restore projects
	if err := restoreProjects(projectService, projectsPath, options, result); err != nil {
		if !options.SkipErrors {
			return result, fmt.Errorf("failed to restore projects: %w", err)
		}
		result.Errors = append(result.Errors, fmt.Sprintf("Projects restore failed: %v", err))
	}

	// Restore settings
	if err := restoreSettings(settingsStorage, settingsPath, options, result); err != nil {
		if !options.SkipErrors {
			return result, fmt.Errorf("failed to restore settings: %w", err)
		}
		result.Errors = append(result.Errors, fmt.Sprintf("Settings restore failed: %v", err))
	}

	return result, nil
}

// RestoreAllStorageWithLimits restores all storage data with license limit enforcement
func RestoreAllStorageWithLimits(
	projectService services.ProjectService,
	settingsStorage storage.SettingsStorage,
	projectsPath, settingsPath string,
	options RestoreOptions,
	maxProjects, maxWebsites int,
) (*RestoreResult, error) {
	result := &RestoreResult{}

	// Restore projects with license limit enforcement
	if err := restoreProjectsWithLimits(projectService, projectsPath, options, result, maxProjects, maxWebsites); err != nil {
		if !options.SkipErrors {
			return result, fmt.Errorf("failed to restore projects: %w", err)
		}
		result.Errors = append(result.Errors, fmt.Sprintf("Projects restore failed: %v", err))
	}

	// Restore settings
	if err := restoreSettings(settingsStorage, settingsPath, options, result); err != nil {
		if !options.SkipErrors {
			return result, fmt.Errorf("failed to restore settings: %w", err)
		}
		result.Errors = append(result.Errors, fmt.Sprintf("Settings restore failed: %v", err))
	}

	return result, nil
}

// restoreProjects restores projects from backup file
func restoreProjects(projectService services.ProjectService, projectsPath string, options RestoreOptions, result *RestoreResult) error {
	// Read and parse projects backup file
	data, err := os.ReadFile(projectsPath)
	if err != nil {
		return fmt.Errorf("failed to read projects backup file: %w", err)
	}

	var backupData struct {
		BackupInfo BackupInfo       `json:"backupInfo"`
		Projects   []domain.Project `json:"projects"`
	}

	if err := json.Unmarshal(data, &backupData); err != nil {
		return fmt.Errorf("failed to parse projects backup file: %w", err)
	}

	result.TotalItems += len(backupData.Projects)

	// Handle overwrite option
	if options.OverwriteExisting {
		if err := projectService.DeleteAllProjects(); err != nil {
			// Log the error but continue - the database might be in an inconsistent state
			result.Warnings = append(result.Warnings, fmt.Sprintf("Warning: failed to delete existing projects: %v", err))
			result.Warnings = append(result.Warnings, "Continuing with restore - existing data may remain")
		} else {
			result.Warnings = append(result.Warnings, "All existing projects were deleted before restore")
		}
	}

	// Restore projects
	for _, project := range backupData.Projects {
		if options.ValidateOnly {
			result.RestoredItems++
			continue
		}

		// Restore project with original ID
		if _, err := projectService.RestoreProject(project); err != nil {
			if options.SkipErrors {
				result.FailedItems++
				result.Errors = append(result.Errors, fmt.Sprintf("Failed to restore project %s: %v", project.ID, err))
				continue
			}
			return fmt.Errorf("failed to restore project %s: %w", project.ID, err)
		}

		result.RestoredItems++
	}

	return nil
}

// restoreProjectsWithLimits restores projects with license limit enforcement
func restoreProjectsWithLimits(
	projectService services.ProjectService,
	projectsPath string,
	options RestoreOptions,
	result *RestoreResult,
	maxProjects, maxWebsites int,
) error {
	// Read and parse projects backup file
	data, err := os.ReadFile(projectsPath)
	if err != nil {
		return fmt.Errorf("failed to read projects backup file: %w", err)
	}

	var backupData struct {
		BackupInfo BackupInfo       `json:"backupInfo"`
		Projects   []domain.Project `json:"projects"`
	}

	if err := json.Unmarshal(data, &backupData); err != nil {
		return fmt.Errorf("failed to parse projects backup file: %w", err)
	}

	result.TotalItems += len(backupData.Projects)

	// Check license limits before restore
	if !options.ValidateOnly {
		// Check project limit (0 means unlimited)
		if maxProjects > 0 && len(backupData.Projects) > maxProjects {
			return fmt.Errorf("backup contains %d projects, but license only allows %d", len(backupData.Projects), maxProjects)
		}

		// Count total websites across all projects
		totalWebsites := 0
		for _, project := range backupData.Projects {
			totalWebsites += len(project.IISSites)
		}

		// Check website limit (0 means unlimited)
		if maxWebsites > 0 && totalWebsites > maxWebsites {
			return fmt.Errorf("backup contains %d total websites, but license only allows %d", totalWebsites, maxWebsites)
		}
	}

	// Handle overwrite option
	if options.OverwriteExisting {
		if err := projectService.DeleteAllProjects(); err != nil {
			// Log the error but continue - the database might be in an inconsistent state
			result.Warnings = append(result.Warnings, fmt.Sprintf("Warning: failed to delete existing projects: %v", err))
			result.Warnings = append(result.Warnings, "Continuing with restore - existing data may remain")
		} else {
			result.Warnings = append(result.Warnings, "All existing projects were deleted before restore")
		}
	}

	// Restore projects
	for _, project := range backupData.Projects {
		if options.ValidateOnly {
			result.RestoredItems++
			continue
		}

		// Restore project with original ID
		if _, err := projectService.RestoreProject(project); err != nil {
			if options.SkipErrors {
				result.FailedItems++
				result.Errors = append(result.Errors, fmt.Sprintf("Failed to restore project %s: %v", project.ID, err))
				continue
			}
			return fmt.Errorf("failed to restore project %s: %w", project.ID, err)
		}

		result.RestoredItems++
	}

	return nil
}

// restoreSettings restores settings from backup file
func restoreSettings(settingsStorage storage.SettingsStorage, settingsPath string, options RestoreOptions, result *RestoreResult) error {
	// Read and parse settings backup file
	data, err := os.ReadFile(settingsPath)
	if err != nil {
		return fmt.Errorf("failed to read settings backup file: %w", err)
	}

	var backupData struct {
		BackupInfo BackupInfo        `json:"backupInfo"`
		Settings   map[string]string `json:"settings"`
	}

	if err := json.Unmarshal(data, &backupData); err != nil {
		return fmt.Errorf("failed to parse settings backup file: %w", err)
	}

	result.TotalItems += len(backupData.Settings)

	// Handle overwrite option for settings
	if options.OverwriteExisting {
		// Get existing settings keys
		existingKeys, err := settingsStorage.GetSettingKeys()
		if err != nil {
			result.Warnings = append(result.Warnings, fmt.Sprintf("Could not get existing settings keys: %v", err))
		} else {
			// Delete existing settings
			for _, key := range existingKeys {
				if err := settingsStorage.DeleteSetting(key); err != nil {
					result.Warnings = append(result.Warnings, fmt.Sprintf("Could not delete existing setting %s: %v", key, err))
				}
			}
		}
		result.Warnings = append(result.Warnings, "All existing settings were deleted before restore")
	}

	// Restore settings
	for key, value := range backupData.Settings {
		if options.ValidateOnly {
			result.RestoredItems++
			continue
		}

		if err := settingsStorage.SetSetting(key, value); err != nil {
			if options.SkipErrors {
				result.FailedItems++
				result.Errors = append(result.Errors, fmt.Sprintf("Failed to restore setting %s: %v", key, err))
				continue
			}
			return fmt.Errorf("failed to restore setting %s: %w", key, err)
		}

		result.RestoredItems++
	}

	return nil
}
