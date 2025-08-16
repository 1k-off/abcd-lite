package badger

import (
	"os"
	"path/filepath"
	"time"

	"github.com/1k-off/abcd-lite/internal/storage"
	"github.com/gofiber/storage/badger/v2"
)

// BadgerDriver implements storage.StorageDriver using GoFiber BadgerDB storage
type BadgerDriver struct {
	storage *badger.Storage
}

// NewBadgerDriver creates a new BadgerDB driver
func NewBadgerDriver(config storage.StorageConfig) (*BadgerDriver, error) {
	// Ensure the parent directory exists
	dbDir := filepath.Dir(config.Path)
	if err := os.MkdirAll(dbDir, 0755); err != nil {
		return nil, err
	}

	badgerStorage := badger.New(badger.Config{
		Database:   config.Path,
		Reset:      false,
		GCInterval: 10 * time.Minute,
	})

	return &BadgerDriver{storage: badgerStorage}, nil
}

func (d *BadgerDriver) Get(key string) ([]byte, error) {
	return d.storage.Get(key)
}

func (d *BadgerDriver) Set(key string, value []byte, ttl time.Duration) error {
	return d.storage.Set(key, value, ttl)
}

func (d *BadgerDriver) Delete(key string) error {
	return d.storage.Delete(key)
}

func (d *BadgerDriver) Close() error {
	return d.storage.Close()
}

// NewBadgerProjectStorage creates a new BadgerDB-based project storage using generic implementation
func NewBadgerProjectStorage(config storage.StorageConfig) (storage.ProjectStorage, error) {
	driver, err := NewBadgerDriver(config)
	if err != nil {
		return nil, err
	}
	return storage.NewGenericProjectStorage(driver), nil
}

// NewBadgerSettingsStorage creates a new BadgerDB-based settings storage using generic implementation
func NewBadgerSettingsStorage(config storage.StorageConfig) (storage.SettingsStorage, error) {
	driver, err := NewBadgerDriver(config)
	if err != nil {
		return nil, err
	}
	return storage.NewGenericSettingsStorage(driver), nil
}
