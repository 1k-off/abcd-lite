package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/1k-off/abcd-lite/internal/config"
	"github.com/1k-off/abcd-lite/internal/server/services"
	"github.com/1k-off/abcd-lite/pkg/util"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"golang.org/x/crypto/bcrypt"
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Configuration utilities",
}

const (
	appModeDefault  = "default"
	appModeParanoid = "paranoid"
)

func init() {
	var paranoid bool
	generateCmd := &cobra.Command{
		Use:   "generate",
		Short: "Generate a default config file",
		Run: func(cmd *cobra.Command, args []string) {
			generateConfigFile(paranoid)
		},
	}
	generateCmd.Flags().BoolVar(&paranoid, "paranoid", false, "Use paranoid IIS web.config")
	configCmd.AddCommand(generateCmd)
}

func generateConfigFile(paranoid bool) {
	configPath := "configs/config.yml"
	if _, err := os.Stat(configPath); err == nil {
		fmt.Printf("Config file already exists at: %s\n", configPath)
		fmt.Println("If you need to update it, use the 'update' command instead.")
		return
	}
	adminToken := util.GenerateRandomToken(32)
	hash, err := bcrypt.GenerateFromPassword([]byte(adminToken), bcrypt.DefaultCost)
	if err != nil {
		fmt.Println("Failed to hash admin token:", err)
		os.Exit(1)
	}
	jwtSecret := util.GenerateRandomToken(32)

	v := viper.New()
	config.SetDefaults(v)
	v.SetConfigType("yaml")

	os.MkdirAll("configs", 0o755)
	configPath = "configs/config.yml"
	if err := v.WriteConfigAs(configPath); err != nil {
		fmt.Println("Failed to write config file:", err)
		os.Exit(1)
	}

	// Store admin token in datastore
	cfg, err := config.Load()
	if err != nil {
		fmt.Println("Failed to load config for datastore initialization:", err)
		os.Exit(1)
	}

	storage := cfg.Storage()
	defer storage.Close()

	settingsService := services.NewSettingsService(storage)
	if err := settingsService.SetAdminToken(string(hash)); err != nil {
		fmt.Println("Failed to store admin token in datastore:", err)
		os.Exit(1)
	}

	if err := settingsService.SetJwtSecret(jwtSecret); err != nil {
		fmt.Println("Failed to store JWT secret in datastore:", err)
		os.Exit(1)
	}

	fmt.Println("Config file generated at:", configPath)
	fmt.Println("Admin token and JWT secret stored in datastore")
	fmt.Println(strings.Repeat("-", 60))
	fmt.Println("Your generated admin token (save this securely, it will not be shown again):")
	fmt.Println(adminToken)
	fmt.Println(strings.Repeat("-", 60))

	mode := appModeDefault

	os.MkdirAll("web", 0o755)
	var srcConfig string
	if paranoid {
		srcConfig = "internal/config/iis/iis_paranoid.config"
		mode = appModeParanoid
	} else {
		srcConfig = "internal/config/iis/iis_default.config"
		mode = appModeDefault
	}
	webConfigPath := "web/web.config"
	port := v.GetString("app.port")
	if port == "" {
		port = "8900"
	}
	if err := RenderIISConfigTemplate(srcConfig, webConfigPath, port); err != nil {
		fmt.Println("Failed to render IIS config:", err)
		os.Exit(1)
	}
	fmt.Printf("web.config generated at: %s with mode: %s\n", webConfigPath, mode)
}
