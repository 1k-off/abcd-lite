package config

import (
	"github.com/1k-off/abcd-lite/internal/storage"
	"github.com/1k-off/abcd-lite/internal/storage/badger"
	"github.com/1k-off/abcd-lite/internal/storage/sqlite"
)

// DefaultStorageConfig returns the default storage configuration
func DefaultStorageConfig() storage.StorageConfig {
	return storage.StorageConfig{
		Type: storage.StorageTypeSQLite,
		Path: "data/abcd.db",
	}
}

// StorageManager manages shared storage instances
type StorageManager struct {
	ProjectStorage  storage.ProjectStorage
	SettingsStorage storage.SettingsStorage
	Driver          storage.StorageDriver
}

// Close closes the shared storage driver
func (sm *StorageManager) Close() error {
	return sm.Driver.Close()
}

// GetStorage returns a storage manager with shared storage instances
func (c *Config) GetStorage() (*StorageManager, error) {
	storageConfig := storage.StorageConfig{
		Type: storage.StorageType(c.Storage.Type),
		Path: c.Storage.Path,
	}

	var driver storage.StorageDriver
	var err error

	switch storageConfig.Type {
	case storage.StorageTypeBadger:
		driver, err = badger.NewBadgerDriver(storageConfig)
	case storage.StorageTypeSQLite:
		driver, err = sqlite.NewSQLiteDriver(storageConfig)
	default:
		return nil, storage.ErrUnsupportedStorageType
	}

	if err != nil {
		return nil, err
	}

	// Create both storages using the same driver instance
	projectStorage := storage.NewGenericProjectStorage(driver)
	settingsStorage := storage.NewGenericSettingsStorage(driver)

	return &StorageManager{
		ProjectStorage:  projectStorage,
		SettingsStorage: settingsStorage,
		Driver:          driver,
	}, nil
}
