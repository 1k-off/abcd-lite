package services

import (
	"encoding/json"
	"time"

	"github.com/1k-off/abcd-lite/internal/server/domain"
	"github.com/gofiber/storage/badger/v2"
	"github.com/google/uuid"
)

type ProjectService interface {
	GetProjects() ([]domain.Project, error)
	CreateProject(project domain.Project) (domain.Project, error)
	UpdateProject(project domain.Project) error
	DeleteProject(id string) error
}

type DefaultProjectService struct {
	storage *badger.Storage
}

func NewProjectService(storage *badger.Storage) *DefaultProjectService {
	return &DefaultProjectService{storage: storage}
}

func (s *DefaultProjectService) GetProjects() ([]domain.Project, error) {
	projects := make([]domain.Project, 0)

	keys, err := s.storage.Get("project:keys")
	if err != nil {
		if err.Error() == "key not found" {
			emptyKeys, _ := json.Marshal([]string{})
			s.storage.Set("project:keys", emptyKeys, 0)
			return projects, nil
		}
		return projects, err
	}

	var projectKeys []string
	if err := json.Unmarshal(keys, &projectKeys); err != nil {
		return projects, err
	}

	for _, key := range projectKeys {
		data, err := s.storage.Get("project:" + key)
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

func (s *DefaultProjectService) CreateProject(project domain.Project) (domain.Project, error) {
	// Initialize empty slices if they are nil
	if project.IISSites == nil {
		project.IISSites = make([]string, 0)
	}
	if project.APIKeys == nil {
		project.APIKeys = make([]string, 0)
	}

	// Generate ID and timestamps
	project.ID = uuid.New().String()
	project.CreatedAt = time.Now().Format(time.RFC3339)
	project.UpdatedAt = project.CreatedAt

	data, err := json.Marshal(project)
	if err != nil {
		return domain.Project{}, err
	}

	if err := s.storage.Set("project:"+project.ID, data, 0); err != nil {
		return domain.Project{}, err
	}

	// Update project keys list
	keys, err := s.storage.Get("project:keys")
	if err != nil {
		if err.Error() != "key not found" {
			return domain.Project{}, err
		}
		keys = []byte("[]")
	}

	var projectKeys []string
	if err := json.Unmarshal(keys, &projectKeys); err != nil {
		return domain.Project{}, err
	}

	projectKeys = append(projectKeys, project.ID)
	keysData, err := json.Marshal(projectKeys)
	if err != nil {
		return domain.Project{}, err
	}

	if err := s.storage.Set("project:keys", keysData, 0); err != nil {
		return domain.Project{}, err
	}

	return project, nil
}

func (s *DefaultProjectService) UpdateProject(project domain.Project) error {
	existingData, err := s.storage.Get("project:" + project.ID)
	if err != nil {
		return err
	}

	var existingProject domain.Project
	if err := json.Unmarshal(existingData, &existingProject); err != nil {
		return err
	}

	project.CreatedAt = existingProject.CreatedAt
	project.UpdatedAt = time.Now().Format(time.RFC3339)

	data, err := json.Marshal(project)
	if err != nil {
		return err
	}

	if err := s.storage.Set("project:"+project.ID, data, 0); err != nil {
		return err
	}

	return nil
}

func (s *DefaultProjectService) DeleteProject(id string) error {
	if err := s.storage.Delete("project:" + id); err != nil {
		return err
	}

	// Update project keys list
	keys, err := s.storage.Get("project:keys")
	if err != nil {
		return err
	}

	var projectKeys []string
	if err := json.Unmarshal(keys, &projectKeys); err != nil {
		return err
	}

	// Remove the deleted project ID from the list
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

	return s.storage.Set("project:keys", keysData, 0)
}
