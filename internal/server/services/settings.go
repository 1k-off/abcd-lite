package services

import (
	"github.com/1k-off/abcd-lite/internal/storage"
)

type SettingsService interface {
	DeleteSetting(key string) error
	GetAdminToken() (string, error)
	SetAdminToken(token string) error
	GetJwtSecret() (string, error)
	SetJwtSecret(secret string) error
}

type DefaultSettingsService struct {
	storage storage.SettingsStorage
}

func NewSettingsService(storage storage.SettingsStorage) *DefaultSettingsService {
	return &DefaultSettingsService{storage: storage}
}

func (s *DefaultSettingsService) DeleteSetting(key string) error {
	return s.storage.DeleteSetting(key)
}

func (s *DefaultSettingsService) GetAdminToken() (string, error) {
	return s.storage.GetSetting("admin_token")
}

func (s *DefaultSettingsService) SetAdminToken(token string) error {
	return s.storage.SetSetting("admin_token", token)
}

func (s *DefaultSettingsService) GetJwtSecret() (string, error) {
	return s.storage.GetSetting("jwt_secret")
}

func (s *DefaultSettingsService) SetJwtSecret(secret string) error {
	return s.storage.SetSetting("jwt_secret", secret)
}
