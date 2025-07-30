package services

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/1k-off/abcd-lite/internal/config"
	"github.com/1k-off/abcd-lite/internal/server/domain"
	"github.com/gofiber/storage/badger/v2"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

var (
	packageLimits = config.PackageLimits{
		MaxProjects: 10,
		MaxWebsites: 10,
	}

	personalPlanLimits = config.PackageLimits{
		MaxProjects: 1,
		MaxWebsites: 5,
	}

	unlimitedLimits = config.PackageLimits{
		MaxProjects: 0,
		MaxWebsites: 0,
	}
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
	service := NewProjectService(createTestStorage("TestGetProjects_EmptyStorage"), packageLimits)
	projects, err := service.GetProjects()
	assert.NoError(t, err)
	assert.Empty(t, projects)
}
func TestGetProject_NotFound(t *testing.T) {
	service := NewProjectService(createTestStorage("TestGetProject_NotFound"), packageLimits)
	defer cleanupTestStorage(service.storage)
	result, err := service.GetProject("non-existent-id")
	assert.Error(t, err)
	assert.Equal(t, domain.Project{}, result)
}
func TestCreateProject_Success(t *testing.T) {
	service := NewProjectService(createTestStorage("TestCreateProject_Success"), packageLimits)
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
	service := NewProjectService(createTestStorage("TestCreateProject_WithNilSlices"), packageLimits)
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
	service := NewProjectService(createTestStorage("TestCreateProject_MultipleProjects"), packageLimits)
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
	service := NewProjectService(createTestStorage("TestUpdateProject_Success"), packageLimits)
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
	service := NewProjectService(createTestStorage("TestUpdateProject_NotFound"), packageLimits)
	defer cleanupTestStorage(service.storage)
	project := domain.Project{
		ID:   "non-existent-id",
		Name: "Test Project",
	}
	err := service.UpdateProject(project)
	assert.Error(t, err)
}

func TestDeleteProject_Success(t *testing.T) {
	service := NewProjectService(createTestStorage("TestDeleteProject_Success"), packageLimits)
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
	service := NewProjectService(createTestStorage("TestDeleteProject_NotFound"), packageLimits)
	defer cleanupTestStorage(service.storage)
	err := service.DeleteProject("non-existent-id")
	assert.Error(t, err)
}

func TestDeleteProject_KeysNotFound(t *testing.T) {
	service := NewProjectService(createTestStorage("TestDeleteProject_KeysNotFound"), packageLimits)
	defer cleanupTestStorage(service.storage)
	project := domain.Project{}

	err := service.DeleteProject(project.ID)

	assert.Error(t, err)
}

func TestAPIKeyCreationAndUniqueness(t *testing.T) {
	service := NewProjectService(createTestStorage("TestAPIKeyCreationAndUniqueness"), packageLimits)
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
	service := NewProjectService(createTestStorage("TestAPIKeyDeletionByID"), packageLimits)
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

func TestCreateProject_ProjectLimitExceeded(t *testing.T) {
	service := NewProjectService(createTestStorage("TestCreateProject_ProjectLimitExceeded"), personalPlanLimits)
	defer cleanupTestStorage(service.storage)

	// Create first project (should succeed)
	project1 := domain.Project{Name: "Project 1"}
	created1, err := service.CreateProject(project1)
	assert.NoError(t, err)
	assert.NotEmpty(t, created1.ID)

	// Try to create second project (should fail)
	project2 := domain.Project{Name: "Project 2"}
	_, err = service.CreateProject(project2)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "project limit exceeded")
}

func TestCreateProject_WebsiteLimitExceeded(t *testing.T) {
	service := NewProjectService(createTestStorage("TestCreateProject_WebsiteLimitExceeded"), personalPlanLimits)
	defer cleanupTestStorage(service.storage)

	// Create project with 6 websites (should fail - limit is 5)
	project := domain.Project{
		Name:     "Project with too many sites",
		IISSites: []string{"site1.com", "site2.com", "site3.com", "site4.com", "site5.com", "site6.com"},
	}
	_, err := service.CreateProject(project)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "website limit exceeded")
}

func TestCreateProject_WebsiteLimitWithinBounds(t *testing.T) {
	service := NewProjectService(createTestStorage("TestCreateProject_WebsiteLimitWithinBounds"), personalPlanLimits)
	defer cleanupTestStorage(service.storage)

	// Create project with 5 websites (should succeed - exactly at limit)
	project := domain.Project{
		Name:     "Project with max sites",
		IISSites: []string{"site1.com", "site2.com", "site3.com", "site4.com", "site5.com"},
	}
	created, err := service.CreateProject(project)
	assert.NoError(t, err)
	assert.NotEmpty(t, created.ID)
	assert.Len(t, created.IISSites, 5)
}

func TestUpdateProject_WebsiteLimitExceeded(t *testing.T) {
	service := NewProjectService(createTestStorage("TestUpdateProject_WebsiteLimitExceeded"), personalPlanLimits)
	defer cleanupTestStorage(service.storage)

	// Create project with 3 websites
	project := domain.Project{
		Name:     "Test Project",
		IISSites: []string{"site1.com", "site2.com", "site3.com"},
	}
	created, err := service.CreateProject(project)
	require.NoError(t, err)

	// Try to update with 6 websites (should fail - would exceed limit of 5)
	updatedProject := domain.Project{
		ID:       created.ID,
		Name:     "Updated Project",
		IISSites: []string{"site1.com", "site2.com", "site3.com", "site4.com", "site5.com", "site6.com"},
	}
	err = service.UpdateProject(updatedProject)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "website limit exceeded")
}

func TestUpdateProject_WebsiteLimitWithinBounds(t *testing.T) {
	service := NewProjectService(createTestStorage("TestUpdateProject_WebsiteLimitWithinBounds"), personalPlanLimits)
	defer cleanupTestStorage(service.storage)

	// Create project with 2 websites
	project := domain.Project{
		Name:     "Test Project",
		IISSites: []string{"site1.com", "site2.com"},
	}
	created, err := service.CreateProject(project)
	require.NoError(t, err)

	// Update with 5 websites (should succeed - exactly at limit)
	updatedProject := domain.Project{
		ID:       created.ID,
		Name:     "Updated Project",
		IISSites: []string{"site1.com", "site2.com", "site3.com", "site4.com", "site5.com"},
	}
	err = service.UpdateProject(updatedProject)
	assert.NoError(t, err)

	// Verify the update
	result, err := service.GetProject(created.ID)
	assert.NoError(t, err)
	assert.Len(t, result.IISSites, 5)
}

func TestUpdateProject_WebsiteCountUnchanged(t *testing.T) {
	service := NewProjectService(createTestStorage("TestUpdateProject_WebsiteCountUnchanged"), personalPlanLimits)
	defer cleanupTestStorage(service.storage)

	// Create project with 4 websites
	project := domain.Project{
		Name:     "Test Project",
		IISSites: []string{"site1.com", "site2.com", "site3.com", "site4.com"},
	}
	created, err := service.CreateProject(project)
	require.NoError(t, err)

	// Update with same number of websites but different names (should succeed)
	updatedProject := domain.Project{
		ID:       created.ID,
		Name:     "Updated Project",
		IISSites: []string{"new-site1.com", "new-site2.com", "new-site3.com", "new-site4.com"},
	}
	err = service.UpdateProject(updatedProject)
	assert.NoError(t, err)

	// Verify the update
	result, err := service.GetProject(created.ID)
	assert.NoError(t, err)
	assert.Len(t, result.IISSites, 4)
	assert.Equal(t, []string{"new-site1.com", "new-site2.com", "new-site3.com", "new-site4.com"}, result.IISSites)
}

func TestUpdateProject_WebsiteCountDecreased(t *testing.T) {
	service := NewProjectService(createTestStorage("TestUpdateProject_WebsiteCountDecreased"), personalPlanLimits)
	defer cleanupTestStorage(service.storage)

	// Create project with 4 websites
	project := domain.Project{
		Name:     "Test Project",
		IISSites: []string{"site1.com", "site2.com", "site3.com", "site4.com"},
	}
	created, err := service.CreateProject(project)
	require.NoError(t, err)

	// Update with fewer websites (should succeed)
	updatedProject := domain.Project{
		ID:       created.ID,
		Name:     "Updated Project",
		IISSites: []string{"site1.com", "site2.com"},
	}
	err = service.UpdateProject(updatedProject)
	assert.NoError(t, err)

	// Verify the update
	result, err := service.GetProject(created.ID)
	assert.NoError(t, err)
	assert.Len(t, result.IISSites, 2)
}

func TestUpdateProject_NoWebsites(t *testing.T) {
	service := NewProjectService(createTestStorage("TestUpdateProject_NoWebsites"), personalPlanLimits)
	defer cleanupTestStorage(service.storage)

	// Create project with 3 websites
	project := domain.Project{
		Name:     "Test Project",
		IISSites: []string{"site1.com", "site2.com", "site3.com"},
	}
	created, err := service.CreateProject(project)
	require.NoError(t, err)

	// Update with no websites (should succeed)
	updatedProject := domain.Project{
		ID:       created.ID,
		Name:     "Updated Project",
		IISSites: []string{},
	}
	err = service.UpdateProject(updatedProject)
	assert.NoError(t, err)

	// Verify the update
	result, err := service.GetProject(created.ID)
	assert.NoError(t, err)
	assert.Len(t, result.IISSites, 0)
}

func TestCreateProject_NoWebsites(t *testing.T) {
	service := NewProjectService(createTestStorage("TestCreateProject_NoWebsites"), personalPlanLimits)
	defer cleanupTestStorage(service.storage)

	// Create project with no websites (should succeed)
	project := domain.Project{
		Name:     "Test Project",
		IISSites: []string{},
	}
	created, err := service.CreateProject(project)
	assert.NoError(t, err)
	assert.NotEmpty(t, created.ID)
	assert.Len(t, created.IISSites, 0)
}

func TestCreateProject_AddWebsitesLater(t *testing.T) {
	service := NewProjectService(createTestStorage("TestCreateProject_AddWebsitesLater"), personalPlanLimits)
	defer cleanupTestStorage(service.storage)

	// Create project with no websites
	project := domain.Project{
		Name:     "Test Project",
		IISSites: []string{},
	}
	created, err := service.CreateProject(project)
	require.NoError(t, err)

	// Update with 5 websites (should succeed)
	updatedProject := domain.Project{
		ID:       created.ID,
		Name:     "Updated Project",
		IISSites: []string{"site1.com", "site2.com", "site3.com", "site4.com", "site5.com"},
	}
	err = service.UpdateProject(updatedProject)
	assert.NoError(t, err)

	// Verify the update
	result, err := service.GetProject(created.ID)
	assert.NoError(t, err)
	assert.Len(t, result.IISSites, 5)
}

func TestCreateProject_AddWebsitesLaterExceedsLimit(t *testing.T) {
	service := NewProjectService(createTestStorage("TestCreateProject_AddWebsitesLaterExceedsLimit"), personalPlanLimits)
	defer cleanupTestStorage(service.storage)

	// Create project with no websites
	project := domain.Project{
		Name:     "Test Project",
		IISSites: []string{},
	}
	created, err := service.CreateProject(project)
	require.NoError(t, err)

	// Try to update with 6 websites (should fail)
	updatedProject := domain.Project{
		ID:       created.ID,
		Name:     "Updated Project",
		IISSites: []string{"site1.com", "site2.com", "site3.com", "site4.com", "site5.com", "site6.com"},
	}
	err = service.UpdateProject(updatedProject)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "website limit exceeded")
}

func TestCreateProject_UnlimitedProjects(t *testing.T) {
	service := NewProjectService(createTestStorage("TestCreateProject_UnlimitedProjects"), unlimitedLimits)
	defer cleanupTestStorage(service.storage)

	// Create multiple projects (should all succeed)
	for i := 1; i <= 10; i++ {
		project := domain.Project{
			Name:     fmt.Sprintf("Project %d", i),
			IISSites: []string{fmt.Sprintf("site%d.com", i)},
		}
		created, err := service.CreateProject(project)
		assert.NoError(t, err)
		assert.NotEmpty(t, created.ID)
	}

	// Verify all projects were created
	projects, err := service.GetProjects()
	assert.NoError(t, err)
	assert.Len(t, projects, 10)
}

func TestCreateProject_UnlimitedWebsites(t *testing.T) {
	service := NewProjectService(createTestStorage("TestCreateProject_UnlimitedWebsites"), unlimitedLimits)
	defer cleanupTestStorage(service.storage)

	// Create project with many websites (should succeed)
	project := domain.Project{
		Name:     "Project with many sites",
		IISSites: []string{"site1.com", "site2.com", "site3.com", "site4.com", "site5.com", "site6.com", "site7.com", "site8.com", "site9.com", "site10.com"},
	}
	created, err := service.CreateProject(project)
	assert.NoError(t, err)
	assert.NotEmpty(t, created.ID)
	assert.Len(t, created.IISSites, 10)
}

func TestUpdateProject_UnlimitedWebsites(t *testing.T) {
	service := NewProjectService(createTestStorage("TestUpdateProject_UnlimitedWebsites"), unlimitedLimits)
	defer cleanupTestStorage(service.storage)

	// Create project with 3 websites
	project := domain.Project{
		Name:     "Test Project",
		IISSites: []string{"site1.com", "site2.com", "site3.com"},
	}
	created, err := service.CreateProject(project)
	require.NoError(t, err)

	// Update with many more websites (should succeed)
	updatedProject := domain.Project{
		ID:       created.ID,
		Name:     "Updated Project",
		IISSites: []string{"site1.com", "site2.com", "site3.com", "site4.com", "site5.com", "site6.com", "site7.com", "site8.com", "site9.com", "site10.com", "site11.com", "site12.com"},
	}
	err = service.UpdateProject(updatedProject)
	assert.NoError(t, err)

	// Verify the update
	result, err := service.GetProject(created.ID)
	assert.NoError(t, err)
	assert.Len(t, result.IISSites, 12)
}

func TestCreateProject_UnlimitedProjectsAndWebsites(t *testing.T) {
	service := NewProjectService(createTestStorage("TestCreateProject_UnlimitedProjectsAndWebsites"), unlimitedLimits)
	defer cleanupTestStorage(service.storage)

	// Create multiple projects with many websites each
	for i := 1; i <= 5; i++ {
		project := domain.Project{
			Name:     fmt.Sprintf("Project %d", i),
			IISSites: []string{fmt.Sprintf("site%d-1.com", i), fmt.Sprintf("site%d-2.com", i), fmt.Sprintf("site%d-3.com", i)},
		}
		created, err := service.CreateProject(project)
		assert.NoError(t, err)
		assert.NotEmpty(t, created.ID)
		assert.Len(t, created.IISSites, 3)
	}

	// Verify all projects were created
	projects, err := service.GetProjects()
	assert.NoError(t, err)
	assert.Len(t, projects, 5)

	// Verify total websites
	totalWebsites := 0
	for _, project := range projects {
		totalWebsites += len(project.IISSites)
	}
	assert.Equal(t, 15, totalWebsites) // 5 projects * 3 websites each
}

func TestUpdateProject_UnlimitedWebsitesSameCount(t *testing.T) {
	service := NewProjectService(createTestStorage("TestUpdateProject_UnlimitedWebsitesSameCount"), unlimitedLimits)
	defer cleanupTestStorage(service.storage)

	// Create project with 5 websites
	project := domain.Project{
		Name:     "Test Project",
		IISSites: []string{"site1.com", "site2.com", "site3.com", "site4.com", "site5.com"},
	}
	created, err := service.CreateProject(project)
	require.NoError(t, err)

	// Update with same number of websites but different names (should succeed)
	updatedProject := domain.Project{
		ID:       created.ID,
		Name:     "Updated Project",
		IISSites: []string{"new-site1.com", "new-site2.com", "new-site3.com", "new-site4.com", "new-site5.com"},
	}
	err = service.UpdateProject(updatedProject)
	assert.NoError(t, err)

	// Verify the update
	result, err := service.GetProject(created.ID)
	assert.NoError(t, err)
	assert.Len(t, result.IISSites, 5)
	assert.Equal(t, []string{"new-site1.com", "new-site2.com", "new-site3.com", "new-site4.com", "new-site5.com"}, result.IISSites)
}

func TestCreateProject_UnlimitedNoWebsites(t *testing.T) {
	service := NewProjectService(createTestStorage("TestCreateProject_UnlimitedNoWebsites"), unlimitedLimits)
	defer cleanupTestStorage(service.storage)

	// Create project with no websites (should succeed)
	project := domain.Project{
		Name:     "Test Project",
		IISSites: []string{},
	}
	created, err := service.CreateProject(project)
	assert.NoError(t, err)
	assert.NotEmpty(t, created.ID)
	assert.Len(t, created.IISSites, 0)
}

func TestUpdateProject_UnlimitedAddWebsitesLater(t *testing.T) {
	service := NewProjectService(createTestStorage("TestUpdateProject_UnlimitedAddWebsitesLater"), unlimitedLimits)
	defer cleanupTestStorage(service.storage)

	// Create project with no websites
	project := domain.Project{
		Name:     "Test Project",
		IISSites: []string{},
	}
	created, err := service.CreateProject(project)
	require.NoError(t, err)

	// Update with many websites (should succeed)
	updatedProject := domain.Project{
		ID:       created.ID,
		Name:     "Updated Project",
		IISSites: []string{"site1.com", "site2.com", "site3.com", "site4.com", "site5.com", "site6.com", "site7.com", "site8.com", "site9.com", "site10.com"},
	}
	err = service.UpdateProject(updatedProject)
	assert.NoError(t, err)

	// Verify the update
	result, err := service.GetProject(created.ID)
	assert.NoError(t, err)
	assert.Len(t, result.IISSites, 10)
}

func TestCreateProject_UnlimitedMixedLimits(t *testing.T) {
	// Test with unlimited projects but limited websites
	mixedLimits := config.PackageLimits{
		MaxProjects: 0, // Unlimited projects
		MaxWebsites: 5, // Limited websites
	}

	service := NewProjectService(createTestStorage("TestCreateProject_UnlimitedMixedLimits"), mixedLimits)
	defer cleanupTestStorage(service.storage)

	// Create multiple projects (should succeed)
	for i := 1; i <= 3; i++ {
		project := domain.Project{
			Name:     fmt.Sprintf("Project %d", i),
			IISSites: []string{fmt.Sprintf("site%d.com", i)},
		}
		created, err := service.CreateProject(project)
		assert.NoError(t, err)
		assert.NotEmpty(t, created.ID)
	}

	// Verify all projects were created
	projects, err := service.GetProjects()
	assert.NoError(t, err)
	assert.Len(t, projects, 3)

	// Try to create a project with too many websites (should fail)
	projectWithManySites := domain.Project{
		Name:     "Project with too many sites",
		IISSites: []string{"site1.com", "site2.com", "site3.com", "site4.com", "site5.com", "site6.com"},
	}
	_, err = service.CreateProject(projectWithManySites)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "website limit exceeded")
}

func TestCreateProject_LimitedProjectsUnlimitedWebsites(t *testing.T) {
	// Test with limited projects but unlimited websites
	mixedLimits := config.PackageLimits{
		MaxProjects: 1, // Limited projects
		MaxWebsites: 0, // Unlimited websites
	}

	service := NewProjectService(createTestStorage("TestCreateProject_LimitedProjectsUnlimitedWebsites"), mixedLimits)
	defer cleanupTestStorage(service.storage)

	// Create first project with many websites (should succeed)
	project1 := domain.Project{
		Name:     "Project 1",
		IISSites: []string{"site1.com", "site2.com", "site3.com", "site4.com", "site5.com", "site6.com", "site7.com", "site8.com", "site9.com", "site10.com"},
	}
	created1, err := service.CreateProject(project1)
	assert.NoError(t, err)
	assert.NotEmpty(t, created1.ID)
	assert.Len(t, created1.IISSites, 10)

	// Try to create second project (should fail)
	project2 := domain.Project{
		Name:     "Project 2",
		IISSites: []string{"site1.com"},
	}
	_, err = service.CreateProject(project2)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "project limit exceeded")
}
