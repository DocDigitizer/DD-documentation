function Get-DocDigitizerConfig {
    <#
    .SYNOPSIS
        Gets the current DocDigitizer module configuration.

    .DESCRIPTION
        Displays the current configuration values for the DocDigitizer module,
        including environment variables and their sources.

    .EXAMPLE
        Get-DocDigitizerConfig

        Displays all current configuration values.

    .OUTPUTS
        PSCustomObject with current configuration settings.
    #>
    [CmdletBinding()]
    param()

    Write-CommandLog -CommandName 'Get-DocDigitizerConfig' -Message 'Reading configuration'

    $config = Get-ModuleConfig

    # Add source information
    $result = [PSCustomObject]@{
        BaseUrl         = $config.BaseUrl
        BaseUrlSource   = if ($env:DOCDIGITIZER_URL) { 'Environment' } else { 'Default' }
        Pipeline        = if ($config.Pipeline) { $config.Pipeline } else { '(not set - will use server default)' }
        PipelineSource  = if ($env:DOCDIGITIZER_PIPELINE) { 'Environment' } else { 'Default' }
        LogLevel        = if ($config.LogLevel) { $config.LogLevel } else { '(not set - will use Minimal)' }
        LogLevelSource  = if ($env:DOCDIGITIZER_LOGLEVEL) { 'Environment' } else { 'Default' }
        Timeout         = $config.Timeout
        TimeoutSource   = if ($env:DOCDIGITIZER_TIMEOUT) { 'Environment' } else { 'Default' }
    }

    return $result
}
