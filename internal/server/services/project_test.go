package services

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/1k-off/abcd-lite/internal/server/domain"
	"github.com/gofiber/storage/badger/v2"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

func createTestStorage(testName string) *badger.Storage {
	testDir := filepath.Join(os.TempDir(), "abcd-lite-test-storage", testName)
	os.RemoveAll(testDir)
	storage := badger.New(badger.Config{
		Database: testDir,
		Reset:    true,
	})
	return storage
}

func cleanupTestStorage(storage *badger.Storage) {
	if storage != nil {
		storage.Close()
	}
	testDir := filepath.Join(os.TempDir(), "abcd-lite-test-storage")
	os.RemoveAll(testDir)
}

func makeTestAPIKey(key string) domain.APIKey {
	hash, _ := bcrypt.GenerateFromPassword([]byte(key), bcrypt.DefaultCost)
	return domain.APIKey{
		ID:        uuid.NewString(),
		Hash:      string(hash),
		CreatedAt: "2024-01-01T00:00:00Z",
		Prefix:    key[:4],
		Suffix:    key[len(key)-4:],
	}
}

func TestGetProjects_EmptyStorage(t *testing.T) {
	service := NewProjectService(createTestStorage("TestGetProjects_EmptyStorage"))
	projects, err := service.GetProjects()
	assert.NoError(t, err)
	assert.Empty(t, projects)
}
func TestGetProject_NotFound(t *testing.T) {
	service := NewProjectService(createTestStorage("TestGetProject_NotFound"))
	defer cleanupTestStorage(service.storage)
	result, err := service.GetProject("non-existent-id")
	assert.Error(t, err)
	assert.Equal(t, domain.Project{}, result)
}
func TestCreateProject_Success(t *testing.T) {
	service := NewProjectService(createTestStorage("TestCreateProject_Success"))
	defer cleanupTestStorage(service.storage)
	project := domain.Project{
		Name:     "Test Project",
		IISSites: []string{"site.com"},
		APIKeys:  []domain.APIKey{makeTestAPIKey("key12345")},
	}

	result, err := service.CreateProject(project)

	assert.NoError(t, err)
	assert.NotEmpty(t, result.ID)
	assert.Equal(t, project.Name, result.Name)
	assert.Equal(t, project.IISSites, result.IISSites)
	assert.Equal(t, len(project.APIKeys), len(result.APIKeys))
	assert.NotEmpty(t, result.CreatedAt)
	assert.NotEmpty(t, result.UpdatedAt)
	assert.Equal(t, result.CreatedAt, result.UpdatedAt)
}

func TestCreateProject_WithNilSlices(t *testing.T) {
	service := NewProjectService(createTestStorage("TestCreateProject_WithNilSlices"))
	defer cleanupTestStorage(service.storage)
	project := domain.Project{
		Name: "Test Project",
	}

	result, err := service.CreateProject(project)

	assert.NoError(t, err)
	assert.NotEmpty(t, result.ID)
	assert.Equal(t, project.Name, result.Name)
	assert.NotNil(t, result.IISSites)
	assert.NotNil(t, result.APIKeys)
	assert.Empty(t, result.IISSites)
	assert.Empty(t, result.APIKeys)
}

func TestCreateProject_MultipleProjects(t *testing.T) {
	service := NewProjectService(createTestStorage("TestCreateProject_MultipleProjects"))
	defer cleanupTestStorage(service.storage)

	project1 := domain.Project{Name: "Project 1"}
	project2 := domain.Project{Name: "Project 2"}

	_, err := service.CreateProject(project1)
	assert.NoError(t, err)

	_, err = service.CreateProject(project2)
	assert.NoError(t, err)

	projects, err := service.GetProjects()
	assert.NoError(t, err)
	assert.Len(t, projects, 2)
}

func TestUpdateProject_Success(t *testing.T) {
	service := NewProjectService(createTestStorage("TestUpdateProject_Success"))
	defer cleanupTestStorage(service.storage)

	originalProject := domain.Project{
		Name:     "Original Name",
		IISSites: []string{"oldsite.com"},
		APIKeys:  []domain.APIKey{makeTestAPIKey("oldkey123")},
	}

	created, err := service.CreateProject(originalProject)
	require.NoError(t, err)

	updatedProject := domain.Project{
		ID:       created.ID,
		Name:     "Updated Name",
		IISSites: []string{"newsite.com"},
		APIKeys:  []domain.APIKey{makeTestAPIKey("newkey123")},
	}

	err = service.UpdateProject(updatedProject)
	assert.NoError(t, err)

	result, err := service.GetProject(created.ID)
	assert.NoError(t, err)
	assert.Equal(t, updatedProject.Name, result.Name)
	assert.Equal(t, updatedProject.IISSites, result.IISSites)
	assert.Equal(t, len(updatedProject.APIKeys), len(result.APIKeys))
	assert.Equal(t, created.CreatedAt, result.CreatedAt)
	assert.NotEmpty(t, result.UpdatedAt)
}

func TestUpdateProject_NotFound(t *testing.T) {
	service := NewProjectService(createTestStorage("TestUpdateProject_NotFound"))
	defer cleanupTestStorage(service.storage)
	project := domain.Project{
		ID:   "non-existent-id",
		Name: "Test Project",
	}
	err := service.UpdateProject(project)
	assert.Error(t, err)
}

func TestDeleteProject_Success(t *testing.T) {
	service := NewProjectService(createTestStorage("TestDeleteProject_Success"))
	defer cleanupTestStorage(service.storage)
	project1 := domain.Project{Name: "Project 1"}
	project2 := domain.Project{Name: "Project 2"}
	project3 := domain.Project{Name: "Project 3"}

	created1, err := service.CreateProject(project1)
	require.NoError(t, err)
	created2, err := service.CreateProject(project2)
	require.NoError(t, err)
	created3, err := service.CreateProject(project3)
	require.NoError(t, err)

	err = service.DeleteProject(created2.ID)
	assert.NoError(t, err)

	_, err = service.GetProject(created2.ID)
	assert.Error(t, err)

	result1, err := service.GetProject(created1.ID)
	assert.NoError(t, err)
	assert.Equal(t, created1.ID, result1.ID)

	result3, err := service.GetProject(created3.ID)
	assert.NoError(t, err)
	assert.Equal(t, created3.ID, result3.ID)

	projects, err := service.GetProjects()
	assert.NoError(t, err)
	assert.Len(t, projects, 2)
}

func TestDeleteProject_NotFound(t *testing.T) {
	service := NewProjectService(createTestStorage("TestDeleteProject_NotFound"))
	defer cleanupTestStorage(service.storage)
	err := service.DeleteProject("non-existent-id")
	assert.Error(t, err)
}

func TestDeleteProject_KeysNotFound(t *testing.T) {
	service := NewProjectService(createTestStorage("TestDeleteProject_KeysNotFound"))
	defer cleanupTestStorage(service.storage)
	project := domain.Project{}

	err := service.DeleteProject(project.ID)

	assert.Error(t, err)
}

func TestAPIKeyCreationAndUniqueness(t *testing.T) {
	service := NewProjectService(createTestStorage("TestAPIKeyCreationAndUniqueness"))
	defer cleanupTestStorage(service.storage)
	project := domain.Project{
		Name:    "API Key Project",
		APIKeys: []domain.APIKey{},
	}
	created, err := service.CreateProject(project)
	require.NoError(t, err)
	// Add two API keys
	key1, meta1, err := service.AddAPIKey(created.ID)
	require.NoError(t, err)
	key2, meta2, err := service.AddAPIKey(created.ID)
	require.NoError(t, err)
	assert.NotEmpty(t, key1)
	assert.NotEmpty(t, key2)
	assert.NotEqual(t, key1, key2)
	assert.NotEmpty(t, meta1.ID)
	assert.NotEmpty(t, meta2.ID)
	assert.NotEqual(t, meta1.ID, meta2.ID)
	// Fetch project and check keys
	updated, err := service.GetProject(created.ID)
	require.NoError(t, err)
	assert.Len(t, updated.APIKeys, 2)
	ids := map[string]bool{}
	for _, k := range updated.APIKeys {
		assert.NotEmpty(t, k.ID)
		ids[k.ID] = true
	}
	assert.Len(t, ids, 2)
}

func TestAPIKeyDeletionByID(t *testing.T) {
	service := NewProjectService(createTestStorage("TestAPIKeyDeletionByID"))
	defer cleanupTestStorage(service.storage)
	project := domain.Project{
		Name:    "API Key Delete Project",
		APIKeys: []domain.APIKey{},
	}
	created, err := service.CreateProject(project)
	require.NoError(t, err)
	_, meta1, err := service.AddAPIKey(created.ID)
	require.NoError(t, err)
	_, meta2, err := service.AddAPIKey(created.ID)
	require.NoError(t, err)
	// Delete first key
	err = service.RemoveAPIKey(created.ID, meta1.ID)
	assert.NoError(t, err)
	updated, err := service.GetProject(created.ID)
	require.NoError(t, err)
	assert.Len(t, updated.APIKeys, 1)
	assert.Equal(t, meta2.ID, updated.APIKeys[0].ID)
}
