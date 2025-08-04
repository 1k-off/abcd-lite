package cmd

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/1k-off/abcd-lite/internal/config"
	"github.com/1k-off/abcd-lite/internal/server/services"
	"github.com/1k-off/abcd-lite/pkg/util"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"golang.org/x/crypto/bcrypt"
)

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update configuration",
}

var updateWebCmd = &cobra.Command{
	Use:   "web",
	Short: "Update web configuration",
}

var updateWebModeCmd = &cobra.Command{
	Use:   "mode [default|paranoid]",
	Short: "Update web mode (default or paranoid)",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		mode := args[0]
		if mode != appModeDefault && mode != appModeParanoid {
			return errors.New("mode must be 'default' or 'paranoid'")
		}
		var srcConfig string
		if mode == appModeParanoid {
			srcConfig = "internal/config/iis/iis_paranoid.config"
		} else {
			srcConfig = "internal/config/iis/iis_default.config"
		}
		webConfigPath := "web/web.config"
		// Read port from config.yml
		v := viper.New()
		v.SetConfigFile("configs/config.yml")
		if err := v.ReadInConfig(); err != nil {
			return fmt.Errorf("failed to read config: %w", err)
		}
		port := v.GetString("app.port")
		if port == "" {
			port = "8900"
		}
		if err := RenderIISConfigTemplate(srcConfig, webConfigPath, port); err != nil {
			return fmt.Errorf("failed to render IIS config: %w", err)
		}
		fmt.Printf("web.config updated to %s mode.\n", mode)
		fmt.Println("Please restart the service for changes to take effect.")
		return nil
	},
}

var updateWebPortCmd = &cobra.Command{
	Use:   "port [port]",
	Short: "Update web port in both IIS and app config",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		port := args[0]
		if !regexp.MustCompile(`^\d{1,5}$`).MatchString(port) {
			return errors.New("invalid port number format")
		}
		// Validate port range
		var portNum int
		fmt.Sscanf(port, "%d", &portNum)
		if portNum < 1 || portNum > 65535 {
			return errors.New("port must be between 1 and 65535")
		}
		// Update app config
		v := viper.New()
		v.SetConfigFile("configs/config.yml")
		if err := v.ReadInConfig(); err != nil {
			return fmt.Errorf("failed to read config: %w", err)
		}
		v.Set("app.port", port)
		if err := v.WriteConfig(); err != nil {
			return fmt.Errorf("failed to write config: %w", err)
		}
		// Re-render IIS config (web.config) with new port, using current mode
		webConfigPath := "web/web.config"
		// Try to detect mode by checking which template matches web.config
		mode := appModeDefault
		// If web.config matches paranoid template structure, use paranoid
		// (Simple heuristic: look for AllowDeployIIS rule)
		data, err := os.ReadFile(webConfigPath)
		if err == nil && strings.Contains(string(data), "AllowDeployIIS") {
			mode = appModeParanoid
		}
		var srcConfig string
		if mode == appModeParanoid {
			srcConfig = "internal/config/iis/iis_paranoid.config"
		} else {
			srcConfig = "internal/config/iis/iis_default.config"
		}
		if err := RenderIISConfigTemplate(srcConfig, webConfigPath, port); err != nil {
			return fmt.Errorf("failed to render IIS config: %w", err)
		}
		fmt.Printf("Port updated to %s in both app config and web.config.\n", port)
		fmt.Println("Please restart the service for changes to take effect.")
		return nil
	},
}

var updateWebJwtSecretCmd = &cobra.Command{
	Use:   "jwt-secret",
	Short: "Update JWT secret in datastore",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		jwtSecret := util.GenerateRandomToken(32)

		// Initialize storage and settings service
		cfg, err := config.Load()
		if err != nil {
			return fmt.Errorf("failed to load config: %w", err)
		}

		storage := cfg.Storage()
		defer storage.Close()

		settingsService := services.NewSettingsService(storage)
		if err := settingsService.SetJwtSecret(jwtSecret); err != nil {
			return fmt.Errorf("failed to save JWT secret to datastore: %w", err)
		}

		fmt.Println(strings.Repeat("-", 60))
		fmt.Println("Your new JWT secret (save this securely, it will not be shown again):")
		fmt.Println(jwtSecret)
		fmt.Println(strings.Repeat("-", 60))
		fmt.Println("JWT secret updated successfully in datastore.")
		fmt.Println("Please restart the service for changes to take effect.")
		return nil
	},
}

var updatePasswordCmd = &cobra.Command{
	Use:   "password",
	Short: "Update admin password",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Print("Enter new password: ")
		reader := bufio.NewReader(os.Stdin)
		pw, _ := reader.ReadString('\n')
		pw = strings.TrimSpace(pw)
		hash, err := bcrypt.GenerateFromPassword([]byte(pw), bcrypt.DefaultCost)
		if err != nil {
			return fmt.Errorf("failed to hash password: %w", err)
		}

		// Initialize storage and settings service
		cfg, err := config.Load()
		if err != nil {
			return fmt.Errorf("failed to load config: %w", err)
		}

		storage := cfg.Storage()
		defer storage.Close()

		settingsService := services.NewSettingsService(storage)
		if err := settingsService.SetAdminToken(string(hash)); err != nil {
			return fmt.Errorf("failed to save admin token to datastore: %w", err)
		}

		fmt.Println(strings.Repeat("-", 60))
		fmt.Println("Your new admin password (save this securely, it will not be shown again):")
		fmt.Println(pw)
		fmt.Println(strings.Repeat("-", 60))
		fmt.Println("Admin password updated successfully in datastore.")
		return nil
	},
}

func init() {
	updateWebCmd.AddCommand(updateWebModeCmd)
	updateWebCmd.AddCommand(updateWebPortCmd)
	updateWebCmd.AddCommand(updateWebJwtSecretCmd)
	updateCmd.AddCommand(updateWebCmd)
	updateCmd.AddCommand(updatePasswordCmd)
	configCmd.AddCommand(updateCmd)
}
