package services

import (
	"fmt"
	"slices"

	"github.com/1k-off/abcd-lite/internal/deployment"
	"github.com/1k-off/abcd-lite/internal/server/domain"
)

type IISDeploymentService interface {
	Deploy(iis domain.IIS) error
}

type DefaultIISDeploymentService struct {
	projectService ProjectService
}

func NewIISDeploymentService(p ProjectService) *DefaultIISDeploymentService {
	return &DefaultIISDeploymentService{projectService: p}
}

func (s *DefaultIISDeploymentService) Deploy(iis domain.IIS) error {
	project, err := s.projectService.GetProject(iis.ProjectId)
	if err != nil {
		return err
	}

	// Check API key using bcrypt
	authorized := false
	for _, key := range project.APIKeys {
		if s.projectService.CheckAPIKey(iis.DeployKey, key.Hash) {
			authorized = true
			break
		}
	}
	if !authorized {
		return fmt.Errorf("unauthorized deployment to project %s", project.Name)
	}

	if !slices.Contains(project.IISSites, iis.SiteName) {
		return fmt.Errorf("site name %s not found in project %s", iis.SiteName, iis.ProjectId)
	}

	packageCredentials := deployment.NewCredentials(iis.PackageInfo.Credentials.Username, iis.PackageInfo.Credentials.Password)
	packageInfo := deployment.NewPackageInfo(iis.PackageInfo.PackageRef, packageCredentials)
	iisDeployment := deployment.IIS{
		SiteName:                iis.SiteName,
		AppPoolName:             iis.AppPoolName,
		StopAppPoolBeforeDeploy: iis.StopAppPoolBeforeDeploy,
		StartAppPoolAfterDeploy: iis.StartAppPoolAfterDeploy,
		CleanDeployment:         iis.CleanDeployment,
		Exclude:                 iis.Exclude,
	}

	if err := iisDeployment.Deploy(packageInfo); err != nil {
		return err
	}

	return nil
}
