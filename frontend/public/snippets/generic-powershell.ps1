# Variables
$ABCDLITE_URL = "{{ABCDLITE_URL}}"
$PROJECT_ID = "{{PROJECT_ID}}"
$DEPLOY_KEY = "{{DEPLOY_KEY}}"
$SITE_NAME = "{{SITE_NAME}}"
$PACKAGE_REF = "{{PACKAGE_REF}}"
$EXCLUDE = @("put", "your", "exclusions/here")

# JSON payload
$payload = @{
    project_id = $PROJECT_ID
    deploy_key = $DEPLOY_KEY
    site_name = $SITE_NAME
    stopAppPoolBeforeDeploy = $true
    start_app_pool_after_deploy = $true
    clean_deployment = $true
    exclude = $EXCLUDE
    package_info = @{
        package_ref = $PACKAGE_REF
    }
}

# Send request
Invoke-RestMethod -Uri "$ABCDLITE_URL/deploy/iis" -Method Post -ContentType "application/json" -Body ($payload | ConvertTo-Json -Depth 4) 