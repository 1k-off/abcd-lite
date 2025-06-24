package services

import (
	"errors"
	"testing"

	"github.com/1k-off/abcd-lite/internal/server/domain"
)

type mockProjectService struct {
	project domain.Project
	getErr  error
}

func (m *mockProjectService) GetProjects() ([]domain.Project, error) { return nil, nil }
func (m *mockProjectService) GetProject(id string) (domain.Project, error) {
	return m.project, m.getErr
}
func (m *mockProjectService) CreateProject(project domain.Project) (domain.Project, error) {
	return domain.Project{}, nil
}
func (m *mockProjectService) UpdateProject(project domain.Project) error { return nil }
func (m *mockProjectService) DeleteProject(id string) error              { return nil }

func TestIISDeploymentService_ProjectNotFound(t *testing.T) {
	mockSvc := &mockProjectService{getErr: errors.New("not found")}
	service := &DefaultIISDeploymentService{projectService: mockSvc}

	iis := domain.IIS{ProjectId: "p1", ApiKey: "key1", SiteName: "site1"}
	err := service.Deploy(iis)
	if err == nil || err.Error() != "not found" {
		t.Errorf("expected not found error, got %v", err)
	}
}

func TestIISDeploymentService_UnauthorizedAPIKey(t *testing.T) {
	mockProj := domain.Project{ID: "p1", Name: "p1", APIKeys: []string{"key1"}, IISSites: []string{"site1"}}
	mockSvc := &mockProjectService{project: mockProj}
	service := &DefaultIISDeploymentService{projectService: mockSvc}

	iis := domain.IIS{ProjectId: "p1", ApiKey: "badkey", SiteName: "site1"}
	err := service.Deploy(iis)
	if err == nil || err.Error() != "unauthorized deployment to project p1" {
		t.Errorf("expected unauthorized error, got %v", err)
	}
}

func TestIISDeploymentService_SiteNotFound(t *testing.T) {
	mockProj := domain.Project{ID: "p1", APIKeys: []string{"key1"}, IISSites: []string{"site1"}}
	mockSvc := &mockProjectService{project: mockProj}
	service := &DefaultIISDeploymentService{projectService: mockSvc}

	iis := domain.IIS{ProjectId: "p1", ApiKey: "key1", SiteName: "notfound"}
	err := service.Deploy(iis)
	if err == nil || err.Error() != "site name notfound not found in project p1" {
		t.Errorf("expected site not found error, got %v", err)
	}
}

func TestIISDeploymentService_EmptyAPIKey(t *testing.T) {
	mockProj := domain.Project{ID: "p1", Name: "p1", APIKeys: []string{"key1"}, IISSites: []string{"site1"}}
	mockSvc := &mockProjectService{project: mockProj}
	service := &DefaultIISDeploymentService{projectService: mockSvc}

	iis := domain.IIS{ProjectId: "p1", ApiKey: "", SiteName: "site1"}
	err := service.Deploy(iis)
	if err == nil || err.Error() != "unauthorized deployment to project p1" {
		t.Errorf("expected unauthorized error for empty API key, got %v", err)
	}
}
