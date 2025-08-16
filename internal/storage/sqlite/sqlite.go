package sqlite

import (
	"os"
	"path/filepath"
	"time"

	"github.com/1k-off/abcd-lite/internal/storage"
	"github.com/gofiber/storage/sqlite3"
)

// SQLiteDriver implements storage.StorageDriver using GoFiber SQLite storage
type SQLiteDriver struct {
	storage *sqlite3.Storage
}

// NewSQLiteDriver creates a new SQLite driver
func NewSQLiteDriver(config storage.StorageConfig) (*SQLiteDriver, error) {
	// Ensure the parent directory exists
	dbDir := filepath.Dir(config.Path)
	if err := os.MkdirAll(dbDir, 0755); err != nil {
		return nil, err
	}

	sqliteStorage := sqlite3.New(sqlite3.Config{
		Database:   config.Path,
		Reset:      false,
		GCInterval: 10 * time.Minute,
	})

	return &SQLiteDriver{storage: sqliteStorage}, nil
}

func (d *SQLiteDriver) Get(key string) ([]byte, error) {
	return d.storage.Get(key)
}

func (d *SQLiteDriver) Set(key string, value []byte, ttl time.Duration) error {
	return d.storage.Set(key, value, ttl)
}

func (d *SQLiteDriver) Delete(key string) error {
	return d.storage.Delete(key)
}

func (d *SQLiteDriver) Close() error {
	return d.storage.Close()
}

// NewSQLiteProjectStorage creates a new SQLite-based project storage using generic implementation
func NewSQLiteProjectStorage(config storage.StorageConfig) (storage.ProjectStorage, error) {
	driver, err := NewSQLiteDriver(config)
	if err != nil {
		return nil, err
	}
	return storage.NewGenericProjectStorage(driver), nil
}

// NewSQLiteSettingsStorage creates a new SQLite-based settings storage using generic implementation
func NewSQLiteSettingsStorage(config storage.StorageConfig) (storage.SettingsStorage, error) {
	driver, err := NewSQLiteDriver(config)
	if err != nil {
		return nil, err
	}
	return storage.NewGenericSettingsStorage(driver), nil
}
