Param(
    # IIS ApplicationPool available actions
    [Parameter(Mandatory = $true)]
    [ValidateSet('get-state', 'restart', 'start', 'stop')]
    [string]$Action,
    # Application information name. User can set app-pool
    [Parameter(ParameterSetName = "ApplicationPool", HelpMessage = "Application pool name")]
    [ValidateScript({ $_ -ne "" })]
    [string]$AppPoolName
)

function Get-AppPoolState {
    Param (
        [Parameter(Mandatory = $true)]
        [string]$AppPoolName
    )
    $pool = Get-IISAppPool -Name $AppPoolName
    if ($null -eq $pool) {
        throw "IIS application pool '$AppPoolName' does not exist."
    }
    return $pool.State
}

function Restart-AppPool {
    Param (
        [Parameter(Mandatory = $true)]
        [string]$AppPoolName
    )
    $state = Get-AppPoolState -AppPoolName $AppPoolName
    if ($state -eq "Stopped") {
        throw "IIS application pool '$AppPoolName' is stopped and cannot be restarted."
    }
    Restart-WebAppPool -Name $AppPoolName
    return "IIS application pool '$AppPoolName' restarted."
}

function Start-AppPool {
    Param (
        [Parameter(Mandatory = $true)]
        [string]$AppPoolName
    )
    Start-WebAppPool -Name $AppPoolName
    return "IIS application pool '$AppPoolName' started."
}

function Stop-AppPool {
    Param (
        [Parameter(Mandatory = $true)]
        [string]$AppPoolName
    )
    $state = Get-AppPoolState -AppPoolName $AppPoolName
    if ($state -eq "Stopped") {
        return "IIS application pool '$AppPoolName' already stopped."
    }
    Stop-WebAppPool -Name $AppPoolName
    $sleep = 5
    while ($true) {
        $__pid = Get-WmiObject -Class win32_process -Filter "name='w3wp.exe'" | Where-Object { ($_.CommandLine).Split('\"')[1] -eq $AppPoolName } | ForEach-Object { $_.ProcessId }
        if (-not $__pid) { break }
        if ($sleep -gt 60) {
            Stop-Process $__pid -Force
            if ($__pid) {
                throw "Process for '$AppPoolName' still running after kill."
            }
            return "IIS application pool '$AppPoolName' force killed."
        }
        Start-Sleep -Seconds $sleep
        $sleep += 5
    }
    return "IIS application pool '$AppPoolName' stopped."
}

function Test-AppPoolExists {
    Param (
        [Parameter(Mandatory = $true)]
        [string]$AppPoolName
    )
    $pool = Get-IISAppPool -Name $AppPoolName -WarningAction:SilentlyContinue
    return ($null -ne $pool)
}

function Invoke-AppPoolAction {
    Param(
        [Parameter(Mandatory = $true)]
        [string]$Action,
        [Parameter(Mandatory = $true)]
        [string]$AppPoolName
    )
    if (-not (Test-AppPoolExists -AppPoolName $AppPoolName)) {
        throw "IIS application pool '$AppPoolName' does not exist."
    }
    switch ($Action) {
        "get-state" { return Get-AppPoolState -AppPoolName $AppPoolName }
        "restart"   { return Restart-AppPool -AppPoolName $AppPoolName }
        "start"     { return Start-AppPool -AppPoolName $AppPoolName }
        "stop"      { return Stop-AppPool -AppPoolName $AppPoolName }
        default     { throw "Unknown action: $Action" }
    }
}

try {
    $result = Invoke-AppPoolAction -Action $Action -AppPoolName $AppPoolName
    Write-Host $result
} catch {
    Write-Error $_
    exit 1
}