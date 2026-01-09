function Test-DocDigitizerConnection {
    <#
    .SYNOPSIS
        Tests the connection to the DocDigitizer API.

    .DESCRIPTION
        Performs a health check against the DocIngester API endpoint to verify
        connectivity and service availability.

    .PARAMETER BaseUrl
        Base URL of the DocDigitizer API. Defaults to environment variable
        DOCDIGITIZER_URL or the built-in default (https://apix.docdigitizer.com/sync).

    .PARAMETER Quiet
        Returns only $true or $false instead of detailed status.

    .EXAMPLE
        Test-DocDigitizerConnection

        Tests connection to the default or configured endpoint.

    .EXAMPLE
        Test-DocDigitizerConnection -BaseUrl "https://api.example.com"

        Tests connection to a specific endpoint.

    .EXAMPLE
        if (Test-DocDigitizerConnection -Quiet) { Send-DocDigitizerDocument ... }

        Use in conditional logic.

    .OUTPUTS
        PSCustomObject with connection status, or Boolean if -Quiet is specified.
    #>
    [CmdletBinding()]
    param(
        [Parameter()]
        [string]$BaseUrl,

        [Parameter()]
        [switch]$Quiet
    )

    $config = Get-ModuleConfig
    if (-not $BaseUrl) { $BaseUrl = $config.BaseUrl }

    Write-CommandLog -CommandName 'Test-DocDigitizerConnection' -Message "Testing $BaseUrl"

    $result = [PSCustomObject]@{
        Url       = $BaseUrl
        Connected = $false
        Response  = $null
        Error     = $null
        Latency   = $null
    }

    # Build headers with API key
    $headers = @{}
    if ($config.ApiKey) {
        $headers['x-api-key'] = $config.ApiKey
    }

    try {
        $stopwatch = [System.Diagnostics.Stopwatch]::StartNew()
        $response = Invoke-RestMethod -Uri $BaseUrl -Method Get -Headers $headers -TimeoutSec 10
        $stopwatch.Stop()

        $result.Connected = $true
        $result.Response = $response
        $result.Latency = $stopwatch.ElapsedMilliseconds

        Write-CommandLog -CommandName 'Test-DocDigitizerConnection' -Message "Connected in $($result.Latency)ms" -Status Complete
        Write-Verbose "Connected to $BaseUrl in $($result.Latency)ms"
    }
    catch {
        $result.Error = $_.Exception.Message
        Write-CommandLog -CommandName 'Test-DocDigitizerConnection' -Message "Failed: $($result.Error)" -Status Info
        Write-Verbose "Failed to connect to $BaseUrl : $($result.Error)"
    }

    if ($Quiet) {
        return $result.Connected
    }

    return $result
}
