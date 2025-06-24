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

	if !slices.Contains(project.APIKeys, iis.ApiKey) {
		return fmt.Errorf("unauthorized deployment to project %s", project.Name)
	}

	if !slices.Contains(project.IISSites, iis.SiteName) {
		return fmt.Errorf("site name %s not found in project %s", iis.SiteName, iis.ProjectId)
	}

	packageCredentials := deployment.NewCredentials(iis.PackageInfo.Credentials.Username, iis.PackageInfo.Credentials.Password, iis.PackageInfo.Credentials.LoginServer)
	packageInfo := deployment.NewPackageInfo(iis.PackageInfo.Name, iis.PackageInfo.Version, packageCredentials)
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
