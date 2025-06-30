package config

import (
	"fmt"
	"io/fs"

	"github.com/gofiber/fiber/v3/log"
	"github.com/spf13/viper"
)

type Config struct {
	App      App `mapstructure:"app"`
	Database Database
	Log      Log `mapstructure:"log"`
}

type App struct {
	AdminToken     string   `mapstructure:"admin_token"`
	JwtSecret      string   `mapstructure:"jwt_secret"`
	Port           string   `mapstructure:"port"`
	Env            string   `mapstructure:"env"`
	Debug          bool     `mapstructure:"debug"`
	AllowedOrigins []string `mapstructure:"allowed_origins"`
}

type Database struct {
	Path    string
	GeoIPFS fs.FS
}

type Log struct {
	Level string `mapstructure:"level"`
}

const (
	AppEnvProduction  = "production"
	AppEnvDevelopment = "development"
)

// Load loads the configuration from file and environment variables
func Load() (*Config, error) {
	v := viper.New()

	// Set default values
	setDefaults(v)

	// Read config file
	v.SetConfigName("config")
	v.SetConfigType("yml")
	v.AddConfigPath("./configs")
	v.AddConfigPath(".")

	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, err
		}
	}

	// Read environment variables
	v.AutomaticEnv()

	config := &Config{}
	if err := v.Unmarshal(config); err != nil {
		return nil, err
	}

	log.SetLevel(getLogLevel(config.Log.Level))

	if config.App.Debug {
		log.SetLevel(log.LevelDebug)
	}

	if err := validateEnv(config.App.Env); err != nil {
		return nil, err
	}

	return config, nil
}

func setDefaults(v *viper.Viper) {
	v.SetDefault("app.port", "8900")
	v.SetDefault("app.env", "production")
	v.SetDefault("app.debug", false)
	v.SetDefault("app.allowed_origins", []string{"http://localhost:5173", "http://localhost:8900"})
	v.SetDefault("database.path", "./data")
	v.SetDefault("log.level", "info")
}

func getLogLevel(level string) log.Level {
	switch level {
	case "debug", "DEBUG":
		return log.LevelDebug
	case "info", "INFO":
		return log.LevelInfo
	case "warn", "WARN":
		return log.LevelWarn
	case "error", "ERROR":
		return log.LevelError
	case "fatal", "FATAL":
		return log.LevelFatal
	}
	return log.LevelInfo
}

func validateEnv(env string) error {
	if env != AppEnvProduction && env != AppEnvDevelopment {
		return fmt.Errorf("invalid environment: %s", env)
	}
	return nil
}
