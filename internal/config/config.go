package config

import "github.com/spf13/viper"

type Config struct {
	App      App      `mapstructure:"app"`
	Database Database `mapstructure:"database"`
}

type App struct {
	Port string `mapstructure:"port"`
	Env  string `mapstructure:"env"`
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

	return config, nil
}

func setDefaults(v *viper.Viper) {
	v.SetDefault("app.port", "8900")
	v.SetDefault("app.env", "development")
	v.SetDefault("database.path", "./data")
	v.SetDefault("log.level", "info")
}
