function Write-CommandLog {
    <#
    .SYNOPSIS
        Writes a formatted command execution log message.

    .DESCRIPTION
        Internal helper function to display consistent loading/execution
        indicators when DocDigitizer commands are invoked.

    .PARAMETER CommandName
        The name of the command being executed.

    .PARAMETER Message
        Optional additional message to display.

    .PARAMETER Status
        The type of log: Start, Complete, or Info.
    #>
    [CmdletBinding()]
    param(
        [Parameter(Mandatory)]
        [string]$CommandName,

        [Parameter()]
        [string]$Message,

        [Parameter()]
        [ValidateSet('Start', 'Complete', 'Info')]
        [string]$Status = 'Start'
    )

    $timestamp = Get-Date -Format 'HH:mm:ss'

    switch ($Status) {
        'Start' {
            $icon = '>'
            $color = 'Cyan'
            $statusText = 'Executing'
        }
        'Complete' {
            $icon = '+'
            $color = 'Green'
            $statusText = 'Completed'
        }
        'Info' {
            $icon = '*'
            $color = 'Yellow'
            $statusText = 'Info'
        }
    }

    $logLine = "[$timestamp] [$icon] $CommandName"
    if ($Message) {
        $logLine += " - $Message"
    }

    Write-Host $logLine -ForegroundColor $color
}
