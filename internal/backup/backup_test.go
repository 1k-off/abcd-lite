package backup

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/1k-off/abcd-lite/internal/server/domain"
	"github.com/1k-off/abcd-lite/internal/storage"
)

// MockProjectStorage implements ProjectStorage for testing
type MockProjectStorage struct {
	projects []domain.Project
}

func (m *MockProjectStorage) GetProjects() ([]domain.Project, error) {
	return m.projects, nil
}

func (m *MockProjectStorage) GetProject(id string) (domain.Project, error) {
	for _, p := range m.projects {
		if p.ID == id {
			return p, nil
		}
	}
	return domain.Project{}, storage.ErrKeyNotFound
}

func (m *MockProjectStorage) CreateProject(project domain.Project) (domain.Project, error) {
	m.projects = append(m.projects, project)
	return project, nil
}

func (m *MockProjectStorage) UpdateProject(project domain.Project) error {
	for i, p := range m.projects {
		if p.ID == project.ID {
			m.projects[i] = project
			return nil
		}
	}
	return storage.ErrKeyNotFound
}

func (m *MockProjectStorage) DeleteProject(id string) error {
	for i, p := range m.projects {
		if p.ID == id {
			m.projects = append(m.projects[:i], m.projects[i+1:]...)
			return nil
		}
	}
	return storage.ErrKeyNotFound
}

func (m *MockProjectStorage) Delete(key string) error {
	// For testing, treat this as DeleteProject
	if strings.HasPrefix(key, "project:") {
		projectID := strings.TrimPrefix(key, "project:")
		return m.DeleteProject(projectID)
	}
	return nil
}

func (m *MockProjectStorage) Exists(key string) (bool, error) {
	// For testing, check if project exists
	if strings.HasPrefix(key, "project:") {
		projectID := strings.TrimPrefix(key, "project:")
		for _, p := range m.projects {
			if p.ID == projectID {
				return true, nil
			}
		}
		return false, nil
	}
	return false, nil
}

func (m *MockProjectStorage) Get(key string) ([]byte, error) {
	// For testing, return project keys or project data
	if key == "project:keys" {
		keys := make([]string, 0, len(m.projects))
		for _, p := range m.projects {
			keys = append(keys, p.ID)
		}
		data, _ := json.Marshal(keys)
		return data, nil
	}
	if strings.HasPrefix(key, "project:") {
		projectID := strings.TrimPrefix(key, "project:")
		project, err := m.GetProject(projectID)
		if err != nil {
			return nil, err
		}
		data, _ := json.Marshal(project)
		return data, nil
	}
	return nil, storage.ErrKeyNotFound
}

func (m *MockProjectStorage) Set(key string, value []byte, ttl time.Duration) error {
	// For testing, handle project keys
	if key == "project:keys" {
		var keys []string
		if err := json.Unmarshal(value, &keys); err != nil {
			return err
		}
		// Update projects list based on keys
		m.projects = make([]domain.Project, 0)
		for _, key := range keys {
			if project, err := m.GetProject(key); err == nil {
				m.projects = append(m.projects, project)
			}
		}
		return nil
	}
	return nil
}

func (m *MockProjectStorage) Ping() error {
	return nil
}

func (m *MockProjectStorage) DeleteAllProjects() error {
	m.projects = []domain.Project{}
	return nil
}

func (m *MockProjectStorage) GetProjectKeys() ([]byte, error) {
	return []byte("[]"), nil
}

func (m *MockProjectStorage) SetProjectKeys(keys []byte) error {
	return nil
}

func (m *MockProjectStorage) AddAPIKey(projectID string, apiKey domain.APIKey) error {
	return nil
}

func (m *MockProjectStorage) RemoveAPIKey(projectID, keyID string) error {
	return nil
}

func (m *MockProjectStorage) GetAPIKeys(projectID string) ([]domain.APIKey, error) {
	return nil, nil
}

// MockSettingsStorage implements SettingsStorage for testing
type MockSettingsStorage struct {
	settings map[string]string
}

func (m *MockSettingsStorage) GetSetting(key string) (string, error) {
	if value, exists := m.settings[key]; exists {
		return value, nil
	}
	return "", storage.ErrKeyNotFound
}

func (m *MockSettingsStorage) SetSetting(key, value string) error {
	if m.settings == nil {
		m.settings = make(map[string]string)
	}
	m.settings[key] = value
	return nil
}

func (m *MockSettingsStorage) DeleteSetting(key string) error {
	delete(m.settings, key)
	return nil
}

func (m *MockSettingsStorage) GetAllSettings() (map[string]string, error) {
	if m.settings == nil {
		return make(map[string]string), nil
	}
	return m.settings, nil
}

func (m *MockSettingsStorage) GetSettingKeys() ([]string, error) {
	if m.settings == nil {
		return []string{}, nil
	}
	keys := make([]string, 0, len(m.settings))
	for k := range m.settings {
		keys = append(keys, k)
	}
	return keys, nil
}

func (m *MockSettingsStorage) HasSetting(key string) (bool, error) {
	_, exists := m.settings[key]
	return exists, nil
}

func (m *MockSettingsStorage) Get(key string) ([]byte, error) {
	if value, exists := m.settings[key]; exists {
		return []byte(value), nil
	}
	return nil, storage.ErrKeyNotFound
}

func (m *MockSettingsStorage) Set(key string, value []byte, ttl time.Duration) error {
	m.settings[key] = string(value)
	return nil
}

func (m *MockSettingsStorage) Delete(key string) error {
	delete(m.settings, key)
	return nil
}

func (m *MockSettingsStorage) Exists(key string) (bool, error) {
	_, exists := m.settings[key]
	return exists, nil
}

func (m *MockSettingsStorage) Ping() error {
	return nil
}

func TestExportAllStorage(t *testing.T) {
	// Create test data
	projects := []domain.Project{
		{
			ID:        "test1",
			Name:      "Test Project 1",
			IISSites:  []string{"site1", "site2"},
			APIKeys:   []domain.APIKey{},
			CreatedAt: time.Now().Format(time.RFC3339),
			UpdatedAt: time.Now().Format(time.RFC3339),
		},
		{
			ID:        "test2",
			Name:      "Test Project 2",
			IISSites:  []string{"site3"},
			APIKeys:   []domain.APIKey{},
			CreatedAt: time.Now().Format(time.RFC3339),
			UpdatedAt: time.Now().Format(time.RFC3339),
		},
	}

	settings := map[string]string{
		"setting1": "value1",
		"setting2": "value2",
	}

	// Create mock storages
	projectStorage := &MockProjectStorage{projects: projects}
	settingsStorage := &MockSettingsStorage{settings: settings}

	// Create temporary directory
	tempDir, err := os.MkdirTemp("", "backup_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Test export
	if err := ExportAllStorage(projectStorage, settingsStorage, tempDir); err != nil {
		t.Fatalf("Export failed: %v", err)
	}

	// Verify files were created
	projectsFile := filepath.Join(tempDir, "projects.json")
	settingsFile := filepath.Join(tempDir, "settings.json")

	if _, err := os.Stat(projectsFile); os.IsNotExist(err) {
		t.Error("Projects file was not created")
	}

	if _, err := os.Stat(settingsFile); os.IsNotExist(err) {
		t.Error("Settings file was not created")
	}

	// Verify backup info
	projectsInfo, err := GetBackupInfo(projectsFile)
	if err != nil {
		t.Fatalf("Failed to get projects backup info: %v", err)
	}

	if projectsInfo.ItemCount != 2 {
		t.Errorf("Expected 2 projects, got %d", projectsInfo.ItemCount)
	}

	settingsInfo, err := GetBackupInfo(settingsFile)
	if err != nil {
		t.Fatalf("Failed to get settings backup info: %v", err)
	}

	if settingsInfo.ItemCount != 2 {
		t.Errorf("Expected 2 settings, got %d", settingsInfo.ItemCount)
	}
}
