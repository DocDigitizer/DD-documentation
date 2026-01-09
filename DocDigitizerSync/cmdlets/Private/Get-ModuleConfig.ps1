function Get-ModuleConfig {
    <#
    .SYNOPSIS
        Gets the module configuration from environment or defaults.
    .DESCRIPTION
        Internal function that retrieves configuration values from environment
        variables or returns sensible defaults.
    #>
    [CmdletBinding()]
    param()

    $config = @{
        BaseUrl  = if ($env:DOCDIGITIZER_URL) { $env:DOCDIGITIZER_URL } else { "https://apix.docdigitizer.com/sync" }
        ApiKey   = $env:DOCDIGITIZER_APIKEY
        Pipeline = if ($env:DOCDIGITIZER_PIPELINE) { $env:DOCDIGITIZER_PIPELINE } else { $null }
        LogLevel = if ($env:DOCDIGITIZER_LOGLEVEL) { $env:DOCDIGITIZER_LOGLEVEL } else { $null }
        Timeout  = if ($env:DOCDIGITIZER_TIMEOUT) { [int]$env:DOCDIGITIZER_TIMEOUT } else { 300 }
    }

    return [PSCustomObject]$config
}
