package cmd

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"os"
	"strings"

	"github.com/1k-off/abcd-lite/internal/config"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"golang.org/x/crypto/bcrypt"
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Configuration utilities",
}

func init() {
	configCmd.AddCommand(&cobra.Command{
		Use:   "generate",
		Short: "Generate a default config file",
		Run: func(cmd *cobra.Command, args []string) {
			generateConfigFile()
		},
	})
}

func generateConfigFile() {
	adminToken := generateRandomToken(32)
	hash, err := bcrypt.GenerateFromPassword([]byte(adminToken), bcrypt.DefaultCost)
	if err != nil {
		fmt.Println("Failed to hash admin token:", err)
		os.Exit(1)
	}
	jwtSecret := generateRandomToken(32)

	v := viper.New()
	config.SetDefaults(v)
	v.Set("app.admin_token", string(hash))
	v.Set("app.jwt_secret", jwtSecret)
	v.SetConfigType("yaml")

	os.MkdirAll("configs", 0o755)
	configPath := "configs/config.yml"
	if err := v.WriteConfigAs(configPath); err != nil {
		fmt.Println("Failed to write config file:", err)
		os.Exit(1)
	}

	fmt.Println("Config file generated at:", configPath)
	fmt.Println(strings.Repeat("-", 60))
	fmt.Println("Your generated admin token (save this securely, it will not be shown again):")
	fmt.Println(adminToken)
	fmt.Println(strings.Repeat("-", 60))
}

func generateRandomToken(length int) string {
	b := make([]byte, length)
	_, err := rand.Read(b)
	if err != nil {
		panic("failed to generate random token")
	}
	return base64.RawURLEncoding.EncodeToString(b)[:length]
}
