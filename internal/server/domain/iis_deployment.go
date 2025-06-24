package domain

type IIS struct {
	ProjectId               string      `json:"project_id"`
	ApiKey                  string      `json:"api_key"`
	SiteName                string      `json:"site_name"`
	AppPoolName             string      `json:"app_pool_name"`
	StopAppPoolBeforeDeploy bool        `json:"stop_app_pool_before_deploy"`
	StartAppPoolAfterDeploy bool        `json:"start_app_pool_after_deploy"`
	CleanDeployment         bool        `json:"clean_deployment"`
	Exclude                 []string    `json:"exclude"`
	PackageInfo             PackageInfo `json:"package_info"`
}
