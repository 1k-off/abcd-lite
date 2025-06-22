package services

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/1k-off/abcd-lite/internal/server/domain"
	"github.com/gofiber/storage/badger/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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
		APIKeys:  []string{"key"},
	}

	result, err := service.CreateProject(project)

	assert.NoError(t, err)
	assert.NotEmpty(t, result.ID)
	assert.Equal(t, project.Name, result.Name)
	assert.Equal(t, project.IISSites, result.IISSites)
	assert.Equal(t, project.APIKeys, result.APIKeys)
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
		APIKeys:  []string{"oldkey"},
	}

	created, err := service.CreateProject(originalProject)
	require.NoError(t, err)

	updatedProject := domain.Project{
		ID:       created.ID,
		Name:     "Updated Name",
		IISSites: []string{"newsite.com"},
		APIKeys:  []string{"newkey"},
	}

	err = service.UpdateProject(updatedProject)
	assert.NoError(t, err)

	result, err := service.GetProject(created.ID)
	assert.NoError(t, err)
	assert.Equal(t, updatedProject.Name, result.Name)
	assert.Equal(t, updatedProject.IISSites, result.IISSites)
	assert.Equal(t, updatedProject.APIKeys, result.APIKeys)
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
	project := domain.Project{Name: "Test Project"}

	err := service.DeleteProject(project.ID)

	assert.Error(t, err)
}
