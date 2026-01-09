function Set-DocDigitizerConfig {
    <#
    .SYNOPSIS
        Sets the DocDigitizer module configuration.

    .DESCRIPTION
        Configures default values for the DocDigitizer module by setting
        environment variables. These persist for the current PowerShell session.

    .PARAMETER ApiKey
        Your DocDigitizer API key. Get one at https://docdigitizer.com/contact

    .PARAMETER BaseUrl
        Default base URL for the DocDigitizer API.

    .PARAMETER Pipeline
        Default pipeline name to use for document processing.

    .PARAMETER LogLevel
        Default log level for responses (Minimal, Medium, Full).

    .PARAMETER Timeout
        Default timeout in seconds for API requests.

    .PARAMETER Persist
        If specified, saves configuration to user profile for future sessions.
        Creates/updates $PROFILE with environment variable settings.

    .EXAMPLE
        Set-DocDigitizerConfig -ApiKey "your-api-key-here"

        Sets your API key for the current session.

    .EXAMPLE
        Set-DocDigitizerConfig -ApiKey "your-api-key-here" -Persist

        Sets your API key and saves it to your PowerShell profile.

    .EXAMPLE
        Set-DocDigitizerConfig -BaseUrl "https://api.example.com"

        Sets the API URL for the current session.

    .EXAMPLE
        Set-DocDigitizerConfig -Pipeline "MainPipelineWithOCR" -LogLevel Full -Persist

        Sets defaults and persists them to the PowerShell profile.

    .EXAMPLE
        Set-DocDigitizerConfig -Timeout 600

        Sets a 10-minute timeout for long-running document processing.
    #>
    [CmdletBinding(SupportsShouldProcess)]
    param(
        [Parameter()]
        [string]$ApiKey,

        [Parameter()]
        [string]$BaseUrl,

        [Parameter()]
        [string]$Pipeline,

        [Parameter()]
        [ValidateSet('Minimal', 'Medium', 'Full')]
        [string]$LogLevel,

        [Parameter()]
        [ValidateRange(10, 3600)]
        [int]$Timeout,

        [Parameter()]
        [switch]$Persist
    )

    Write-CommandLog -CommandName 'Set-DocDigitizerConfig' -Message 'Updating configuration'

    $changes = @()

    if ($ApiKey) {
        $env:DOCDIGITIZER_APIKEY = $ApiKey
        $changes += "DOCDIGITIZER_APIKEY=********"
        Write-Verbose "Set DOCDIGITIZER_APIKEY"
    }

    if ($BaseUrl) {
        $env:DOCDIGITIZER_URL = $BaseUrl
        $changes += "DOCDIGITIZER_URL=$BaseUrl"
        Write-Verbose "Set DOCDIGITIZER_URL to $BaseUrl"
    }

    if ($Pipeline) {
        $env:DOCDIGITIZER_PIPELINE = $Pipeline
        $changes += "DOCDIGITIZER_PIPELINE=$Pipeline"
        Write-Verbose "Set DOCDIGITIZER_PIPELINE to $Pipeline"
    }

    if ($LogLevel) {
        $env:DOCDIGITIZER_LOGLEVEL = $LogLevel
        $changes += "DOCDIGITIZER_LOGLEVEL=$LogLevel"
        Write-Verbose "Set DOCDIGITIZER_LOGLEVEL to $LogLevel"
    }

    if ($Timeout) {
        $env:DOCDIGITIZER_TIMEOUT = $Timeout.ToString()
        $changes += "DOCDIGITIZER_TIMEOUT=$Timeout"
        Write-Verbose "Set DOCDIGITIZER_TIMEOUT to $Timeout"
    }

    if ($changes.Count -eq 0) {
        Write-Warning "No configuration changes specified."
        return
    }

    Write-Host "Configuration updated for current session:" -ForegroundColor Green
    $changes | ForEach-Object { Write-Host "  $_" -ForegroundColor Cyan }

    if ($Persist) {
        if ($PSCmdlet.ShouldProcess($PROFILE, "Add DocDigitizer configuration")) {
            $profileContent = @"

# DocDigitizer PowerShell Module Configuration
`$env:DOCDIGITIZER_APIKEY = '$($env:DOCDIGITIZER_APIKEY)'
`$env:DOCDIGITIZER_URL = '$($env:DOCDIGITIZER_URL)'
`$env:DOCDIGITIZER_PIPELINE = '$($env:DOCDIGITIZER_PIPELINE)'
`$env:DOCDIGITIZER_LOGLEVEL = '$($env:DOCDIGITIZER_LOGLEVEL)'
`$env:DOCDIGITIZER_TIMEOUT = '$($env:DOCDIGITIZER_TIMEOUT)'
"@
            # Ensure profile directory exists
            $profileDir = Split-Path $PROFILE -Parent
            if (-not (Test-Path $profileDir)) {
                New-Item -ItemType Directory -Path $profileDir -Force | Out-Null
            }

            # Check if config already exists in profile
            if (Test-Path $PROFILE) {
                $existingContent = Get-Content $PROFILE -Raw
                if ($existingContent -match '# DocDigitizer PowerShell Module Configuration') {
                    # Replace existing config block
                    $pattern = '# DocDigitizer PowerShell Module Configuration[\s\S]*?(?=\n#|\n\n|\z)'
                    $existingContent = $existingContent -replace $pattern, $profileContent.TrimStart()
                    Set-Content -Path $PROFILE -Value $existingContent
                }
                else {
                    Add-Content -Path $PROFILE -Value $profileContent
                }
            }
            else {
                Set-Content -Path $PROFILE -Value $profileContent.TrimStart()
            }

            Write-Host "Configuration persisted to: $PROFILE" -ForegroundColor Green
        }
    }
}
