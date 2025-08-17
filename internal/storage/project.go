package storage

import (
	"encoding/json"
	"time"

	"github.com/1k-off/abcd-lite/internal/server/domain"
)

// ProjectStorage defines project-specific storage operations
type ProjectStorage interface {
	Storage

	// Project operations
	GetProjects() ([]domain.Project, error)
	GetProject(id string) (domain.Project, error)
	CreateProject(project domain.Project) (domain.Project, error)
	UpdateProject(project domain.Project) error
	DeleteProject(id string) error
	GetProjectKeys() ([]byte, error)
	SetProjectKeys(keys []byte) error

	// API Key operations
	AddAPIKey(projectID string, apiKey domain.APIKey) error
	RemoveAPIKey(projectID, keyID string) error
	GetAPIKeys(projectID string) ([]domain.APIKey, error)
}

// GenericProjectStorage provides a single implementation that works with any StorageDriver
type GenericProjectStorage struct {
	driver StorageDriver
}

// NewGenericProjectStorage creates a generic project storage that works with any driver
func NewGenericProjectStorage(driver StorageDriver) *GenericProjectStorage {
	return &GenericProjectStorage{driver: driver}
}

func (s *GenericProjectStorage) Get(key string) ([]byte, error) {
	return s.driver.Get(key)
}

func (s *GenericProjectStorage) Set(key string, value []byte, ttl time.Duration) error {
	return s.driver.Set(key, value, ttl)
}

func (s *GenericProjectStorage) Delete(key string) error {
	return s.driver.Delete(key)
}

func (s *GenericProjectStorage) Exists(key string) (bool, error) {
	_, err := s.driver.Get(key)
	if err != nil {
		if err.Error() == "key not found" {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func (s *GenericProjectStorage) Ping() error {
	// Most storage backends don't have a ping method, so we'll just return nil
	return nil
}

func (s *GenericProjectStorage) GetProjects() ([]domain.Project, error) {
	projects := make([]domain.Project, 0)

	keys, err := s.driver.Get("project:keys")
	if err != nil {
		return projects, err
	}

	var projectKeys []string
	if err := json.Unmarshal(keys, &projectKeys); err != nil {
		if len(keys) == 0 {
			emptyKeys, _ := json.Marshal([]string{})
			s.driver.Set("project:keys", emptyKeys, 0)
			return projects, nil
		}
		return projects, err
	}

	for _, key := range projectKeys {
		data, err := s.driver.Get("project:" + key)
		if err != nil {
			continue
		}

		var project domain.Project
		if err := json.Unmarshal(data, &project); err != nil {
			continue
		}
		projects = append(projects, project)
	}

	return projects, nil
}

func (s *GenericProjectStorage) GetProject(id string) (domain.Project, error) {
	data, err := s.driver.Get("project:" + id)
	if err != nil {
		return domain.Project{}, err
	}
	var project domain.Project
	if err := json.Unmarshal(data, &project); err != nil {
		return domain.Project{}, err
	}
	return project, nil
}

func (s *GenericProjectStorage) CreateProject(project domain.Project) (domain.Project, error) {
	data, err := json.Marshal(project)
	if err != nil {
		return domain.Project{}, err
	}

	if err := s.driver.Set("project:"+project.ID, data, 0); err != nil {
		return domain.Project{}, err
	}

	// Update project keys list
	keys, err := s.driver.Get("project:keys")
	var projectKeys []string

	if err != nil {
		if err.Error() == "key not found" {
			projectKeys = []string{}
		} else {
			return domain.Project{}, err
		}
	} else {
		if err := json.Unmarshal(keys, &projectKeys); err != nil {
			projectKeys = []string{}
		}
	}

	projectKeys = append(projectKeys, project.ID)
	keysData, err := json.Marshal(projectKeys)
	if err != nil {
		return domain.Project{}, err
	}

	if err := s.driver.Set("project:keys", keysData, 0); err != nil {
		return domain.Project{}, err
	}

	return project, nil
}

func (s *GenericProjectStorage) UpdateProject(project domain.Project) error {
	data, err := json.Marshal(project)
	if err != nil {
		return err
	}

	return s.driver.Set("project:"+project.ID, data, 0)
}

func (s *GenericProjectStorage) DeleteProject(id string) error {
	if err := s.driver.Delete("project:" + id); err != nil {
		return err
	}

	keys, err := s.driver.Get("project:keys")
	if err != nil {
		return err
	}

	var projectKeys []string
	if err := json.Unmarshal(keys, &projectKeys); err != nil {
		return err
	}

	newKeys := make([]string, 0, len(projectKeys))
	for _, key := range projectKeys {
		if key != id {
			newKeys = append(newKeys, key)
		}
	}

	keysData, err := json.Marshal(newKeys)
	if err != nil {
		return err
	}

	return s.driver.Set("project:keys", keysData, 0)
}

func (s *GenericProjectStorage) GetProjectKeys() ([]byte, error) {
	return s.driver.Get("project:keys")
}

func (s *GenericProjectStorage) SetProjectKeys(keys []byte) error {
	return s.driver.Set("project:keys", keys, 0)
}

func (s *GenericProjectStorage) AddAPIKey(projectID string, apiKey domain.APIKey) error {
	project, err := s.GetProject(projectID)
	if err != nil {
		return err
	}

	project.APIKeys = append(project.APIKeys, apiKey)
	return s.UpdateProject(project)
}

func (s *GenericProjectStorage) RemoveAPIKey(projectID, keyID string) error {
	project, err := s.GetProject(projectID)
	if err != nil {
		return err
	}

	newKeys := make([]domain.APIKey, 0, len(project.APIKeys))
	for _, k := range project.APIKeys {
		if k.ID != keyID {
			newKeys = append(newKeys, k)
		}
	}
	project.APIKeys = newKeys
	return s.UpdateProject(project)
}

func (s *GenericProjectStorage) GetAPIKeys(projectID string) ([]domain.APIKey, error) {
	project, err := s.GetProject(projectID)
	if err != nil {
		return nil, err
	}
	return project.APIKeys, nil
}
