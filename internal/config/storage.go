package config

import "github.com/gofiber/storage/badger/v2"

func (cfg *Config) Storage() *badger.Storage {
	return badger.New(badger.Config{
		Database: cfg.Database.Path + "/abcd.db",
		Reset:    false,
	})
}
