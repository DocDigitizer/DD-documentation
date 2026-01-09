function Format-Json {
    <#
    .SYNOPSIS
        Formats JSON with cleaner indentation.
    .DESCRIPTION
        Internal function that formats JSON using jq if available,
        otherwise falls back to PowerShell's default formatting.
    #>
    [CmdletBinding()]
    param(
        [Parameter(Mandatory, ValueFromPipeline)]
        [string]$Json
    )

    # Check if jq is available
    $jqAvailable = $null -ne (Get-Command jq -ErrorAction SilentlyContinue)

    if ($jqAvailable) {
        try {
            $formatted = $Json | jq .
            if ($LASTEXITCODE -eq 0 -and $formatted) {
                return $formatted -join "`n"
            }
        }
        catch {
            # Fall through to default
        }
    }

    # Fallback: return as-is (PowerShell's default formatting)
    return $Json
}
