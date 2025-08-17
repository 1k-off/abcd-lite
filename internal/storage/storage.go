package storage

import (
	"errors"
	"time"
)

// StorageDriver defines the interface for different storage backends
type StorageDriver interface {
	Get(key string) ([]byte, error)
	Set(key string, value []byte, ttl time.Duration) error
	Delete(key string) error
	Close() error
}

// Storage defines the interface for all storage operations
type Storage interface {
	// Basic key-value operations
	Get(key string) ([]byte, error)
	Set(key string, value []byte, ttl time.Duration) error
	Delete(key string) error
	Exists(key string) (bool, error)

	// Health check
	Ping() error
}

// StorageType represents the type of storage backend
type StorageType string

const (
	StorageTypeBadger StorageType = "badger"
	StorageTypeSQLite StorageType = "sqlite"
)

// StorageConfig holds configuration for storage backends
type StorageConfig struct {
	Type StorageType `json:"type"`
	Path string      `json:"path"`
}

// Storage errors
var (
	ErrUnsupportedStorageType = errors.New("unsupported storage type")
	ErrKeyNotFound            = errors.New("key not found")
)
