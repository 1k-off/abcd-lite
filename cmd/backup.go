package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/1k-off/abcd-lite/internal/backup"
	"github.com/1k-off/abcd-lite/internal/config"
	"github.com/spf13/cobra"
)

var backupTimestamped bool

var backupCmd = &cobra.Command{
	Use:   "backup [output-directory]",
	Short: "Create a backup of all storage data",
	Long: `Create a backup of all storage data to JSON files.

Use --timestamped flag to create files with timestamps (e.g., projects_2024-01-15_14-30-25.json).

Examples:
  abcd-lite backup ./backups
  abcd-lite backup --timestamped ./backups`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		outputDir := args[0]

		cfg, err := config.Load()
		if err != nil {
			return fmt.Errorf("failed to load config: %w", err)
		}

		storageManager, err := cfg.GetStorage()
		if err != nil {
			return fmt.Errorf("failed to initialize storage: %w", err)
		}
		defer storageManager.Close()

		if outputDir == "" {
			return fmt.Errorf("output directory cannot be empty")
		}

		// Convert to absolute path if relative
		if !filepath.IsAbs(outputDir) {
			absPath, err := filepath.Abs(outputDir)
			if err != nil {
				return fmt.Errorf("failed to resolve output directory path: %w", err)
			}
			outputDir = absPath
		}

		// Check if output directory exists, create if not
		if _, err := os.Stat(outputDir); os.IsNotExist(err) {
			if err := os.MkdirAll(outputDir, 0755); err != nil {
				return fmt.Errorf("failed to create output directory %s: %w", outputDir, err)
			}
			fmt.Printf("Created output directory: %s\n", outputDir)
		}

		fmt.Printf("Starting backup to: %s\n", outputDir)

		if backupTimestamped {
			if err := backup.ExportAllStorageWithTimestamp(storageManager.ProjectStorage, storageManager.SettingsStorage, outputDir); err != nil {
				return fmt.Errorf("backup failed: %w", err)
			}
		} else {
			if err := backup.ExportAllStorage(storageManager.ProjectStorage, storageManager.SettingsStorage, outputDir); err != nil {
				return fmt.Errorf("backup failed: %w", err)
			}
		}

		var projectsPath, settingsPath string
		if backupTimestamped {
			timestamp := time.Now().Format("2006-01-02_15-04-05")
			projectsPath = filepath.Join(outputDir, fmt.Sprintf("projects_%s.json", timestamp))
			settingsPath = filepath.Join(outputDir, fmt.Sprintf("settings_%s.json", timestamp))
		} else {
			projectsPath = filepath.Join(outputDir, "projects.json")
			settingsPath = filepath.Join(outputDir, "settings.json")
		}

		fmt.Println("Backup completed successfully!")
		fmt.Printf("Projects exported to: %s\n", projectsPath)
		fmt.Printf("Settings exported to: %s\n", settingsPath)

		if projectsInfo, err := os.Stat(projectsPath); err == nil {
			fmt.Printf("Projects file size: %d bytes\n", projectsInfo.Size())
		}
		if settingsInfo, err := os.Stat(settingsPath); err == nil {
			fmt.Printf("Settings file size: %d bytes\n", settingsInfo.Size())
		}

		return nil
	},
}

func init() {
	backupCmd.Flags().BoolVar(&backupTimestamped, "timestamped", false, "Create timestamped backup files")
}
