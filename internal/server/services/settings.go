package services

import (
	"github.com/gofiber/storage/badger/v2"
)

type SettingsService interface {
	GetSetting(key string) (string, error)
	SetSetting(key, value string) error
	DeleteSetting(key string) error
	GetAdminToken() (string, error)
	SetAdminToken(token string) error
	GetJwtSecret() (string, error)
	SetJwtSecret(secret string) error
}

type DefaultSettingsService struct {
	storage *badger.Storage
}

func NewSettingsService(storage *badger.Storage) *DefaultSettingsService {
	return &DefaultSettingsService{storage: storage}
}

func (s *DefaultSettingsService) GetSetting(key string) (string, error) {
	data, err := s.storage.Get("settings:" + key)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func (s *DefaultSettingsService) SetSetting(key, value string) error {
	return s.storage.Set("settings:"+key, []byte(value), 0)
}

func (s *DefaultSettingsService) DeleteSetting(key string) error {
	return s.storage.Delete("settings:" + key)
}

func (s *DefaultSettingsService) GetAdminToken() (string, error) {
	return s.GetSetting("admin_token")
}

func (s *DefaultSettingsService) SetAdminToken(token string) error {
	return s.SetSetting("admin_token", token)
}

func (s *DefaultSettingsService) GetJwtSecret() (string, error) {
	return s.GetSetting("jwt_secret")
}

func (s *DefaultSettingsService) SetJwtSecret(secret string) error {
	return s.SetSetting("jwt_secret", secret)
}
