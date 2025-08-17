package cmd

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/1k-off/abcd-lite/internal/backup"
	"github.com/1k-off/abcd-lite/internal/config"
	"github.com/1k-off/abcd-lite/internal/server/services"
	"github.com/gofiber/fiber/v3/log"
	"github.com/spf13/cobra"
)

var (
	restoreOverwrite  bool
	restoreValidate   bool
	restoreSkipErrors bool
)

var restoreCmd = &cobra.Command{
	Use:   "restore [backup-folder]",
	Short: "Restore data from backup files",
	Long: `Restore data from backup files to the storage.

The command automatically detects both standard backup files (projects.json, settings.json) and 
timestamped backup files (projects_YYYY-MM-DD_HH-MM-SS.json, settings_YYYY-MM-DD_HH-MM-SS.json).

Use --overwrite to replace existing data, --validate-only to check backup files without restoring,
or --skip-errors to continue restore even if some items fail.

Examples:
  abcd-lite restore ./backups                           # Use backup folder
  abcd-lite restore --overwrite ./backups               # Overwrite existing data
  abcd-lite restore --validate-only ./backups           # Validate without restoring`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		// Single argument: backup folder
		backupFolder := args[0]

		// Find backup files in the folder
		projectsBackupPath, settingsBackupPath, err := findBackupFilesInFolder(backupFolder)
		if err != nil {
			return fmt.Errorf("failed to find backup files in folder: %w", err)
		}

		fmt.Printf("Using backup folder: %s\n", backupFolder)
		fmt.Printf("Found projects backup: %s\n", filepath.Base(projectsBackupPath))
		fmt.Printf("Found settings backup: %s\n", filepath.Base(settingsBackupPath))

		cfg, err := config.Load()
		if err != nil {
			return fmt.Errorf("failed to load config: %w", err)
		}

		storageManager, err := cfg.GetStorage()
		if err != nil {
			return fmt.Errorf("failed to initialize storage: %w", err)
		}
		defer storageManager.Close()

		// Validate backup files exist
		if err := validateBackupFiles(projectsBackupPath, settingsBackupPath); err != nil {
			return fmt.Errorf("backup file validation failed: %w", err)
		}

		// Show backup information
		if err := showBackupInfo(projectsBackupPath, settingsBackupPath); err != nil {
			return fmt.Errorf("failed to get backup info: %w", err)
		}

		// Create restore options
		options := backup.RestoreOptions{
			OverwriteExisting: restoreOverwrite,
			ValidateOnly:      restoreValidate,
			SkipErrors:        restoreSkipErrors,
		}

		if restoreValidate {
			fmt.Println("Validation mode - no data will be restored")
		}

		// Perform the restore with license limit enforcement
		fmt.Println("Starting restore...")
		fmt.Println("Security: License limits are enforced during restore")

		limits := getCurrentLimits(cfg)
		projectService := services.NewProjectService(storageManager.ProjectStorage, limits)

		result, err := backup.RestoreAllStorageWithLimits(
			projectService,
			storageManager.SettingsStorage,
			projectsBackupPath,
			settingsBackupPath,
			options,
			limits.MaxProjects,
			limits.MaxWebsites,
		)

		if err != nil {
			return fmt.Errorf("restore failed: %w", err)
		}

		displayRestoreResults(result)

		return nil
	},
}

// getCurrentLimits returns the current license limits
func getCurrentLimits(cfg *config.Config) config.PackageLimits {
	limits, _, err := config.Limits(context.Background(), cfg.App.LicenseKey, cfg.App.LicenseFile)
	if err != nil {
		// If license validation fails, fall back to personal plan limits
		log.Warnf("Failed to get license limits: %v, using personal plan defaults", err)
		return config.PackageLimits{
			MaxProjects: config.PersonalPlanProjects(),
			MaxWebsites: config.PersonalPlanWebsites(),
		}
	}
	return limits
}

func init() {
	restoreCmd.Flags().BoolVar(&restoreOverwrite, "overwrite", false, "Overwrite existing data")
	restoreCmd.Flags().BoolVar(&restoreValidate, "validate-only", false, "Validate backup files without restoring")
	restoreCmd.Flags().BoolVar(&restoreSkipErrors, "skip-errors", false, "Continue restore even if some items fail")
}

// validateBackupFiles checks that the backup files exist and are valid
func validateBackupFiles(projectsPath, settingsPath string) error {
	if _, err := os.Stat(projectsPath); os.IsNotExist(err) {
		return fmt.Errorf("projects backup file not found: %s", projectsPath)
	}

	if _, err := os.Stat(settingsPath); os.IsNotExist(err) {
		return fmt.Errorf("settings backup file not found: %s", settingsPath)
	}

	if err := backup.ValidateBackupFile(projectsPath); err != nil {
		return fmt.Errorf("invalid projects backup file: %w", err)
	}

	if err := backup.ValidateBackupFile(settingsPath); err != nil {
		return fmt.Errorf("invalid settings backup file: %w", err)
	}

	return nil
}

// findBackupFilesInFolder automatically detects backup files in a folder
func findBackupFilesInFolder(folderPath string) (string, string, error) {
	// Check if folder exists
	if _, err := os.Stat(folderPath); os.IsNotExist(err) {
		return "", "", fmt.Errorf("backup folder not found: %s", folderPath)
	}

	// First try to find standard backup files
	projectsPath := filepath.Join(folderPath, "projects.json")
	settingsPath := filepath.Join(folderPath, "settings.json")

	// If standard files don't exist, look for timestamped files
	if _, err := os.Stat(projectsPath); os.IsNotExist(err) {
		projectsPath, err = findLatestBackupFile(folderPath, "projects")
		if err != nil {
			return "", "", fmt.Errorf("projects backup file not found in folder: %s", folderPath)
		}
	}

	if _, err := os.Stat(settingsPath); os.IsNotExist(err) {
		settingsPath, err = findLatestBackupFile(folderPath, "settings")
		if err != nil {
			return "", "", fmt.Errorf("settings backup file not found in folder: %s", folderPath)
		}
	}

	return projectsPath, settingsPath, nil
}

// findLatestBackupFile finds the most recent timestamped backup file
func findLatestBackupFile(folderPath, prefix string) (string, error) {
	entries, err := os.ReadDir(folderPath)
	if err != nil {
		return "", fmt.Errorf("failed to read folder: %w", err)
	}

	var latestFile string
	var latestTime time.Time

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		name := entry.Name()
		if !strings.HasPrefix(name, prefix) || !strings.HasSuffix(name, ".json") {
			continue
		}

		// Extract timestamp from filename (e.g., projects_2024-01-15_14-30-25.json)
		if strings.Contains(name, "_") {
			filePath := filepath.Join(folderPath, name)
			info, err := entry.Info()
			if err != nil {
				continue
			}

			if info.ModTime().After(latestTime) {
				latestTime = info.ModTime()
				latestFile = filePath
			}
		}
	}

	if latestFile == "" {
		return "", fmt.Errorf("no timestamped %s backup files found", prefix)
	}

	return latestFile, nil
}

// showBackupInfo displays information about the backup files
func showBackupInfo(projectsPath, settingsPath string) error {
	fmt.Println("Backup files:")

	// Projects backup info
	projectsInfo, err := backup.GetBackupInfo(projectsPath)
	if err != nil {
		return fmt.Errorf("failed to get projects backup info: %w", err)
	}
	fmt.Printf("  Projects: %s (%s, %d items)\n",
		filepath.Base(projectsPath),
		projectsInfo.Timestamp.Format("2006-01-02 15:04:05"),
		projectsInfo.ItemCount)

	// Settings backup info
	settingsInfo, err := backup.GetBackupInfo(settingsPath)
	if err != nil {
		return fmt.Errorf("failed to get settings backup info: %w", err)
	}
	fmt.Printf("  Settings: %s (%s, %d items)\n",
		filepath.Base(settingsPath),
		settingsInfo.Timestamp.Format("2006-01-02 15:04:05"),
		settingsInfo.ItemCount)

	return nil
}

// displayRestoreResults shows the results of the restore operation
func displayRestoreResults(result *backup.RestoreResult) {
	fmt.Println("\nRestore completed!")
	fmt.Printf("Total items processed: %d\n", result.TotalItems)
	fmt.Printf("Items restored: %d\n", result.RestoredItems)
	fmt.Printf("Items skipped: %d\n", result.SkippedItems)
	fmt.Printf("Items failed: %d\n", result.FailedItems)

	if len(result.Warnings) > 0 {
		fmt.Println("\nWarnings:")
		for _, warning := range result.Warnings {
			fmt.Printf("  - %s\n", warning)
		}
	}

	if len(result.Errors) > 0 {
		fmt.Println("\nErrors:")
		for _, err := range result.Errors {
			fmt.Printf("  - %s\n", err)
		}
	}

	if result.FailedItems == 0 && len(result.Errors) == 0 {
		fmt.Println("\n✅ Restore completed successfully!")
	} else if result.RestoredItems > 0 {
		fmt.Println("\n⚠️  Restore completed with some issues")
	} else {
		fmt.Println("\n❌ Restore failed")
	}
}
