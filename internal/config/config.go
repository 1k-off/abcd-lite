package config

import (
	"github.com/gofiber/fiber/v3/log"
	"github.com/spf13/viper"
)

type Config struct {
	App      App      `mapstructure:"app"`
	Database Database `mapstructure:"database"`
	Log      Log      `mapstructure:"log"`
}

type App struct {
	Port  string `mapstructure:"port"`
	Env   string `mapstructure:"env"`
	Debug bool   `mapstructure:"debug"`
}

type Database struct {
	Path string `mapstructure:"path"`
}

type Log struct {
	Level string `mapstructure:"level"`
}

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

	return config, nil
}

func setDefaults(v *viper.Viper) {
	v.SetDefault("app.port", "8900")
	v.SetDefault("app.env", "development")
	v.SetDefault("app.debug", false)
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
