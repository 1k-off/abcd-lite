package storage

import (
	"encoding/json"
	"time"
)

// SettingsStorage defines settings-specific storage operations
type SettingsStorage interface {
	Storage

	// Settings operations
	GetSetting(key string) (string, error)
	SetSetting(key, value string) error
	DeleteSetting(key string) error
	GetAllSettings() (map[string]string, error)
	GetSettingKeys() ([]string, error)
	HasSetting(key string) (bool, error)
}

// GenericSettingsStorage provides a single implementation that works with any StorageDriver
type GenericSettingsStorage struct {
	driver StorageDriver
}

// NewGenericSettingsStorage creates a generic settings storage that works with any driver
func NewGenericSettingsStorage(driver StorageDriver) *GenericSettingsStorage {
	return &GenericSettingsStorage{driver: driver}
}

// GetDriver returns the underlying storage driver for sharing connections
func (s *GenericSettingsStorage) GetDriver() StorageDriver {
	return s.driver
}

func (s *GenericSettingsStorage) Get(key string) ([]byte, error) {
	return s.driver.Get(key)
}

func (s *GenericSettingsStorage) Set(key string, value []byte, ttl time.Duration) error {
	return s.driver.Set(key, value, ttl)
}

func (s *GenericSettingsStorage) Delete(key string) error {
	return s.driver.Delete(key)
}

func (s *GenericSettingsStorage) Exists(key string) (bool, error) {
	_, err := s.driver.Get(key)
	if err != nil {
		if err.Error() == "key not found" {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func (s *GenericSettingsStorage) Ping() error {
	// Most storage backends don't have a ping method, so we'll just return nil
	return nil
}

func (s *GenericSettingsStorage) GetSetting(key string) (string, error) {
	data, err := s.driver.Get("settings:" + key)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func (s *GenericSettingsStorage) SetSetting(key, value string) error {
	// Store the setting
	if err := s.driver.Set("settings:"+key, []byte(value), 0); err != nil {
		return err
	}

	// Update the list of setting keys
	keys, err := s.driver.Get("settings:keys")
	var settingKeys []string

	if err != nil {
		if err.Error() == "key not found" {
			settingKeys = []string{}
		} else {
			return err
		}
	} else {
		if err := json.Unmarshal(keys, &settingKeys); err != nil {
			settingKeys = []string{}
		}
	}

	// Add the key if it doesn't exist
	keyExists := false
	for _, existingKey := range settingKeys {
		if existingKey == key {
			keyExists = true
			break
		}
	}

	if !keyExists {
		settingKeys = append(settingKeys, key)
		keysData, err := json.Marshal(settingKeys)
		if err != nil {
			return err
		}
		if err := s.driver.Set("settings:keys", keysData, 0); err != nil {
			return err
		}
	}

	return nil
}

func (s *GenericSettingsStorage) DeleteSetting(key string) error {
	// Delete the setting
	if err := s.driver.Delete("settings:" + key); err != nil {
		return err
	}

	// Update the list of setting keys
	keys, err := s.driver.Get("settings:keys")
	if err != nil {
		// If no keys list exists, nothing to update
		return nil
	}

	var settingKeys []string
	if err := json.Unmarshal(keys, &settingKeys); err != nil {
		return err
	}

	// Remove the key from the list
	newKeys := make([]string, 0, len(settingKeys))
	for _, existingKey := range settingKeys {
		if existingKey != key {
			newKeys = append(newKeys, existingKey)
		}
	}

	// Update the keys list
	keysData, err := json.Marshal(newKeys)
	if err != nil {
		return err
	}

	return s.driver.Set("settings:keys", keysData, 0)
}

// GetSettingKeys returns all setting keys
func (s *GenericSettingsStorage) GetSettingKeys() ([]string, error) {
	keys, err := s.driver.Get("settings:keys")
	if err != nil {
		if err.Error() == "key not found" {
			return []string{}, nil
		}
		return nil, err
	}

	var settingKeys []string
	if err := json.Unmarshal(keys, &settingKeys); err != nil {
		return []string{}, nil
	}

	return settingKeys, nil
}

// HasSetting checks if a setting exists
func (s *GenericSettingsStorage) HasSetting(key string) (bool, error) {
	_, err := s.driver.Get("settings:" + key)
	if err != nil {
		if err.Error() == "key not found" {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func (s *GenericSettingsStorage) GetAllSettings() (map[string]string, error) {
	settings := make(map[string]string)

	// Get the list of setting keys
	keys, err := s.driver.Get("settings:keys")
	if err != nil {
		// If no keys list exists, return empty map
		if err.Error() == "key not found" {
			return settings, nil
		}
		return nil, err
	}

	var settingKeys []string
	if err := json.Unmarshal(keys, &settingKeys); err != nil {
		// If keys can't be unmarshaled, return empty map
		return settings, nil
	}

	// Retrieve each setting
	for _, key := range settingKeys {
		value, err := s.driver.Get("settings:" + key)
		if err != nil {
			// Skip settings that can't be retrieved
			continue
		}
		settings[key] = string(value)
	}

	return settings, nil
}
