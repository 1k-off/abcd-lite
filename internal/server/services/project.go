package services

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"time"

	"github.com/1k-off/abcd-lite/internal/config"
	"github.com/1k-off/abcd-lite/internal/server/domain"
	"github.com/1k-off/abcd-lite/internal/server/messages"
	"github.com/gofiber/storage/badger/v2"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type ProjectService interface {
	GetProjects() ([]domain.Project, error)
	GetProject(id string) (domain.Project, error)
	CreateProject(project domain.Project) (domain.Project, error)
	UpdateProject(project domain.Project) error
	DeleteProject(id string) error
	AddAPIKey(projectID string) (string, domain.APIKey, error)
	CheckAPIKey(apiKey, hash string) bool
	RemoveAPIKey(projectID, keyID string) error
}

type DefaultProjectService struct {
	storage *badger.Storage
	limits  config.PackageLimits
}

func NewProjectService(storage *badger.Storage, limits config.PackageLimits) *DefaultProjectService {
	return &DefaultProjectService{storage: storage, limits: limits}
}

func (s *DefaultProjectService) GetProjects() ([]domain.Project, error) {
	projects := make([]domain.Project, 0)

	keys, err := s.storage.Get("project:keys")
	if err != nil {
		return projects, err
	}

	var projectKeys []string
	if err := json.Unmarshal(keys, &projectKeys); err != nil {
		if len(keys) == 0 {
			emptyKeys, _ := json.Marshal([]string{})
			s.storage.Set("project:keys", emptyKeys, 0)
			return projects, nil
		}
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

func (s *DefaultProjectService) GetProject(id string) (domain.Project, error) {
	data, err := s.storage.Get("project:" + id)
	if err != nil {
		return domain.Project{}, err
	}
	var project domain.Project
	if err := json.Unmarshal(data, &project); err != nil {
		return domain.Project{}, err
	}
	return project, nil
}

func (s *DefaultProjectService) CreateProject(project domain.Project) (domain.Project, error) {
	isProjectLimitExceeded, err := s.isProjectLimitExceeded()
	if err != nil {
		return domain.Project{}, err
	}
	if isProjectLimitExceeded {
		return domain.Project{}, errors.New(messages.ErrProjectLimitExceeded)
	}

	// Initialize empty slices if they are nil
	if project.IISSites == nil {
		project.IISSites = make([]string, 0)
	}
	if project.APIKeys == nil {
		project.APIKeys = make([]domain.APIKey, 0)
	}

	if len(project.IISSites) > 0 {
		isWebsiteLimitExceeded, err := s.isWebsiteLimitExceeded(len(project.IISSites))
		if err != nil {
			return domain.Project{}, err
		}
		if isWebsiteLimitExceeded {
			return domain.Project{}, errors.New(messages.ErrWebsiteLimitExceeded)
		}
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

	if len(project.IISSites) > 0 {
		existingSites := len(existingProject.IISSites)
		newSites := len(project.IISSites)

		if newSites > existingSites {
			isWebsiteLimitExceeded, err := s.isWebsiteLimitExceeded(newSites - existingSites)
			if err != nil {
				return err
			}
			if isWebsiteLimitExceeded {
				return errors.New("website limit exceeded")
			}
		}
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

	keys, err := s.storage.Get("project:keys")
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

	if err := s.storage.Set("project:keys", keysData, 0); err != nil {
		return err
	}

	return s.storage.Set("project:keys", keysData, 0)
}

func (s *DefaultProjectService) AddAPIKey(projectID string) (string, domain.APIKey, error) {
	project, err := s.GetProject(projectID)
	if err != nil {
		return "", domain.APIKey{}, err
	}
	apiKey, err := generateAPIKey(20)
	if err != nil {
		return "", domain.APIKey{}, err
	}
	hash, err := hashAPIKey(apiKey)
	if err != nil {
		return "", domain.APIKey{}, err
	}
	apiKeyMeta := domain.APIKey{
		ID:        uuid.NewString(),
		Hash:      hash,
		CreatedAt: time.Now().UTC().Format(time.RFC3339),
		Prefix:    apiKey[:4],
		Suffix:    apiKey[len(apiKey)-4:],
	}
	project.APIKeys = append(project.APIKeys, apiKeyMeta)
	if err := s.UpdateProject(project); err != nil {
		return "", domain.APIKey{}, err
	}
	return apiKey, apiKeyMeta, nil
}

// RemoveAPIKey removes an API key from a project by id
func (s *DefaultProjectService) RemoveAPIKey(projectID, keyID string) error {
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
	project.UpdatedAt = time.Now().UTC().Format(time.RFC3339)
	return s.UpdateProject(project)
}

func (s *DefaultProjectService) CheckAPIKey(apiKey, hash string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(apiKey)) == nil
}

func generateAPIKey(length int) (string, error) {
	b := make([]byte, length)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

func hashAPIKey(apiKey string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(apiKey), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

func (s *DefaultProjectService) isProjectLimitExceeded() (bool, error) {
	projects, err := s.GetProjects()
	if err != nil {
		return false, err
	}
	if len(projects) >= s.limits.MaxProjects && s.limits.MaxProjects != 0 {
		return true, nil
	}
	return false, nil
}

func (s *DefaultProjectService) isWebsiteLimitExceeded(newSites int) (bool, error) {
	projects, err := s.GetProjects()
	if err != nil {
		return false, err
	}
	websites := 0
	for _, project := range projects {
		websites += len(project.IISSites)
	}
	websites += newSites
	if websites > s.limits.MaxWebsites && s.limits.MaxWebsites != 0 {
		return true, nil
	}
	return false, nil
}
