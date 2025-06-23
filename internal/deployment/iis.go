package deployment

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/1k-off/abcd-lite/internal/deployment/scripts"
)

type IIS struct {
	SiteName                string
	AppPoolName             string
	StopAppPoolBeforeDeploy bool
	StartAppPoolAfterDeploy bool
	CleanDeployment         bool
	Exclude                 []string
}

func IISDefault() IIS {
	return IIS{
		SiteName:                "",
		AppPoolName:             "",
		StopAppPoolBeforeDeploy: false,
		StartAppPoolAfterDeploy: true,
	}
}

const (
	IISAppPoolActionGetState = "get-state"
	IISAppPoolActionRestart  = "restart"
	IISAppPoolActionStart    = "start"
	IISAppPoolActionStop     = "stop"

	IISAppPoolStateRunning = "Running"
	IISAppPoolStateStopped = "Stopped"
	IISAppPoolStateUnknown = "Unknown"
)

func (i *IIS) Deploy(p PackageInfo) error {
	if i.SiteName == "" && i.AppPoolName == "" {
		return fmt.Errorf("site name and app pool name are not set")
	}

	if i.SiteName != "" && i.AppPoolName == "" {
		appPoolName, err := i.getAppPoolNameFromSite()
		if err != nil {
			return fmt.Errorf("failed to get app pool name: %w", err)
		}
		i.AppPoolName = appPoolName
	}

	siteRoot, err := i.getIISWebsiteRoot(i.SiteName)
	if err != nil {
		return fmt.Errorf("failed to get IIS website root: %w", err)
	}

	if i.StopAppPoolBeforeDeploy {
		if err := i.executeAppPoolAction(IISAppPoolActionStop); err != nil {
			return fmt.Errorf("failed to stop app pool: %w", err)
		}
	}

	deploymentOptions := NewOptions(i.CleanDeployment, siteRoot, 3, i.Exclude)
	if err := Deploy(deploymentOptions, p); err != nil {
		return fmt.Errorf("failed to deploy: %w", err)
	}

	if i.StartAppPoolAfterDeploy {
		if err := i.executeAppPoolAction(IISAppPoolActionStart); err != nil {
			return fmt.Errorf("failed to start app pool: %w", err)
		}
	}

	return nil
}

// getAppPoolNameFromSite retrieves the application pool name from the site name using PowerShell
func (i *IIS) getAppPoolNameFromSite() (string, error) {
	script := `
		$s = Get-IISSite -Name "%s" -WarningAction:SilentlyContinue
		if ($null -eq "$s") {
			Write-Host "IIS site %s does not exist"
			exit 1
		}
		$p = ((Get-IISSite -Name "%s").Applications | Where-Object {$_.Path -eq "/" }).ApplicationPoolName
		Write-Output $p
	`

	scriptContent := fmt.Sprintf(script, i.SiteName, i.SiteName, i.SiteName)

	cmd := exec.Command("powershell", "-Command", scriptContent)
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to get app pool name: %w", err)
	}

	appPoolName := strings.TrimSpace(string(output))
	if appPoolName == "" {
		return "", fmt.Errorf("no application pool found for site: %s", i.SiteName)
	}

	return appPoolName, nil
}

// executeAppPoolAction executes a PowerShell action on the application pool
func (i *IIS) executeAppPoolAction(action string) error {
	if err := validateIISAppPoolAction(action); err != nil {
		return fmt.Errorf("failed to validate app pool action: %w", err)
	}

	scriptPath, err := scripts.ExtractScript(scripts.ApplicationPoolScript)
	if err != nil {
		return fmt.Errorf("failed to extract script: %w", err)
	}
	defer os.Remove(scriptPath)

	cmd := exec.Command("powershell", "-ExecutionPolicy", "Bypass", "-File", scriptPath, "-Action", action, "-AppPoolName", i.AppPoolName)

	var stdout, stderr strings.Builder
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		errorMsg := fmt.Sprintf("application pool %s failed: %v", action, err)
		if stderr.Len() > 0 {
			errorMsg += fmt.Sprintf(" - stderr: %s", stderr.String())
		}
		if stdout.Len() > 0 {
			errorMsg += fmt.Sprintf(" - stdout: %s", stdout.String())
		}
		return errors.New(errorMsg)
	}

	return nil
}

// getAppPoolState gets the current state of the application pool
func (i *IIS) getAppPoolState() (string, error) {
	scriptPath, err := scripts.ExtractScript(scripts.ApplicationPoolScript)
	if err != nil {
		return "", fmt.Errorf("failed to extract script: %w", err)
	}
	defer os.Remove(scriptPath)

	cmd := exec.Command("powershell", "-ExecutionPolicy", "Bypass", "-File", scriptPath, "-Action", IISAppPoolActionGetState, "-AppPoolName", i.AppPoolName)

	var stdout, stderr strings.Builder
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return IISAppPoolStateUnknown, fmt.Errorf("failed to get app pool state: %w", err)
	}

	output := stdout.String()
	lines := strings.Split(output, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line != "" {
			if line == IISAppPoolStateRunning || line == IISAppPoolStateStopped {
				return line, nil
			}
			return IISAppPoolStateUnknown, fmt.Errorf("failed to parse app pool state from output: %s", output)
		}
	}

	return IISAppPoolStateUnknown, fmt.Errorf("failed to parse app pool state from output: %s", output)
}

func validateIISAppPoolAction(action string) error {
	switch action {
	case IISAppPoolActionRestart, IISAppPoolActionStart, IISAppPoolActionStop:
		return nil
	default:
		return fmt.Errorf("invalid action: %s", action)
	}
}

func (i *IIS) getIISWebsiteRoot(siteName string) (string, error) {
	script := `
		$s = Get-IISSite -Name "%s" -WarningAction:SilentlyContinue
		if ($null -eq $s) {
			Write-Host "IIS site %s does not exist"
			exit 1
		}
		$p = ((Get-IISSite -Name "%s").Applications | Where-Object {$_.Path -eq "/" }).VirtualDirectories["/"].PhysicalPath
		Write-Output $p
	`

	scriptContent := fmt.Sprintf(script, i.SiteName, i.SiteName, i.SiteName)

	cmd := exec.Command("powershell", "-Command", scriptContent)
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to get IIS website root: %w", err)
	}

	return strings.TrimSpace(string(output)), nil
}
