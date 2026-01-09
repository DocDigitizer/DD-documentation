function Get-DocDigitizerHelp {
    <#
    .SYNOPSIS
        Displays help and documentation for the DocDigitizer PowerShell module.

    .DESCRIPTION
        Shows an overview of all available commands in the DocDigitizer module,
        with quick usage examples and tips for getting more detailed help.

    .PARAMETER Command
        Show detailed help for a specific command.
        Valid values: Send-DocDigitizerDocument, Test-DocDigitizerConnection,
        Get-DocDigitizerConfig, Set-DocDigitizerConfig

    .PARAMETER Examples
        Show usage examples for all commands.

    .PARAMETER Online
        Open the README documentation in the default browser.

    .EXAMPLE
        Get-DocDigitizerHelp

        Shows overview of all available commands.

    .EXAMPLE
        Get-DocDigitizerHelp -Command Send-DocDigitizerDocument

        Shows detailed help for the Send-DocDigitizerDocument command.

    .EXAMPLE
        Get-DocDigitizerHelp -Examples

        Shows usage examples for all commands.

    .OUTPUTS
        Help text displayed to the console.
    #>
    [CmdletBinding(DefaultParameterSetName = 'Overview')]
    param(
        [Parameter(ParameterSetName = 'Command', Position = 0)]
        [ValidateSet('Send-DocDigitizerDocument', 'Test-DocDigitizerConnection', 'Get-DocDigitizerConfig', 'Set-DocDigitizerConfig')]
        [string]$Command,

        [Parameter(ParameterSetName = 'Examples')]
        [switch]$Examples,

        [Parameter(ParameterSetName = 'Online')]
        [switch]$Online
    )

    # Module path for README
    $modulePath = Split-Path $PSScriptRoot -Parent
    $readmePath = Join-Path $modulePath "README.md"

    if ($Online) {
        if (Test-Path $readmePath) {
            Write-Host "Opening README.md..." -ForegroundColor Cyan
            Start-Process $readmePath
        }
        else {
            Write-Warning "README.md not found at: $readmePath"
        }
        return
    }

    if ($Command) {
        Get-Help $Command -Detailed
        return
    }

    if ($Examples) {
        Write-Host ""
        Write-Host "=" * 60 -ForegroundColor DarkCyan
        Write-Host "  DocDigitizer PowerShell Module - Examples" -ForegroundColor Cyan
        Write-Host "=" * 60 -ForegroundColor DarkCyan
        Write-Host ""

        Write-Host "QUICK START" -ForegroundColor Yellow
        Write-Host "-----------"
        Write-Host '  # Test connection to API'
        Write-Host '  Test-DocDigitizerConnection' -ForegroundColor Green
        Write-Host ''
        Write-Host '  # Process a single document'
        Write-Host '  Send-DocDigitizerDocument -FilePath "invoice.pdf"' -ForegroundColor Green
        Write-Host ''
        Write-Host '  # Process and save result'
        Write-Host '  Send-DocDigitizerDocument -FilePath "invoice.pdf" -SaveExtraction' -ForegroundColor Green
        Write-Host ''

        Write-Host "BATCH PROCESSING" -ForegroundColor Yellow
        Write-Host "----------------"
        Write-Host '  # Process all PDFs in current folder'
        Write-Host '  Get-ChildItem *.pdf | Send-DocDigitizerDocument -SaveExtraction' -ForegroundColor Green
        Write-Host ''
        Write-Host '  # Process with shared context'
        Write-Host '  $ctx = [guid]::NewGuid()' -ForegroundColor Green
        Write-Host '  Get-ChildItem *.pdf | Send-DocDigitizerDocument -ContextId $ctx' -ForegroundColor Green
        Write-Host ''

        Write-Host "PIPELINE SELECTION" -ForegroundColor Yellow
        Write-Host "------------------"
        Write-Host '  # Use OCR-based pipeline'
        Write-Host '  Send-DocDigitizerDocument -FilePath "doc.pdf" -Pipeline "MainPipelineWithOCR"' -ForegroundColor Green
        Write-Host ''
        Write-Host '  # Use vision-based pipeline (no OCR)'
        Write-Host '  Send-DocDigitizerDocument -FilePath "doc.pdf" -Pipeline "MainPipelineWithFile"' -ForegroundColor Green
        Write-Host ''

        Write-Host "CONFIGURATION" -ForegroundColor Yellow
        Write-Host "-------------"
        Write-Host '  # Set production URL'
        Write-Host '  Set-DocDigitizerConfig -BaseUrl "https://your-server.run.app"' -ForegroundColor Green
        Write-Host ''
        Write-Host '  # Persist configuration'
        Write-Host '  Set-DocDigitizerConfig -BaseUrl "https://your-server.run.app" -Persist' -ForegroundColor Green
        Write-Host ''
        Write-Host '  # View current config'
        Write-Host '  Get-DocDigitizerConfig' -ForegroundColor Green
        Write-Host ''

        Write-Host "LOG LEVELS" -ForegroundColor Yellow
        Write-Host "----------"
        Write-Host '  # Minimal (default) - essential data only'
        Write-Host '  Send-DocDigitizerDocument -FilePath "doc.pdf" -LogLevel Minimal' -ForegroundColor Green
        Write-Host ''
        Write-Host '  # Full - complete execution details'
        Write-Host '  Send-DocDigitizerDocument -FilePath "doc.pdf" -LogLevel Full' -ForegroundColor Green
        Write-Host ''

        return
    }

    # Default: Show overview
    Write-Host ""
    Write-Host "=" * 60 -ForegroundColor DarkCyan
    Write-Host "  DocDigitizer PowerShell Module" -ForegroundColor Cyan
    Write-Host "=" * 60 -ForegroundColor DarkCyan
    Write-Host ""
    Write-Host "A PowerShell module for interacting with the DocDigitizer API."
    Write-Host ""

    Write-Host "AVAILABLE COMMANDS" -ForegroundColor Yellow
    Write-Host "------------------"
    Write-Host ""

    $commands = @(
        @{ Name = "Send-DocDigitizerDocument"; Desc = "Send a PDF document for processing" }
        @{ Name = "Test-DocDigitizerConnection"; Desc = "Test API connectivity" }
        @{ Name = "Get-DocDigitizerConfig"; Desc = "View current configuration" }
        @{ Name = "Set-DocDigitizerConfig"; Desc = "Set default configuration values" }
        @{ Name = "Get-DocDigitizerHelp"; Desc = "Show this help (you are here)" }
    )

    foreach ($cmd in $commands) {
        Write-Host "  $($cmd.Name)" -ForegroundColor Green -NoNewline
        Write-Host " - $($cmd.Desc)"
    }

    Write-Host ""
    Write-Host "QUICK START" -ForegroundColor Yellow
    Write-Host "-----------"
    Write-Host '  Test-DocDigitizerConnection                    # Check API is running'
    Write-Host '  Send-DocDigitizerDocument -FilePath "doc.pdf"  # Process a document'
    Write-Host ""

    Write-Host "GET MORE HELP" -ForegroundColor Yellow
    Write-Host "-------------"
    Write-Host '  Get-DocDigitizerHelp -Examples                 # Show usage examples'
    Write-Host '  Get-DocDigitizerHelp -Command <CommandName>    # Detailed command help'
    Write-Host '  Get-DocDigitizerHelp -Online                   # Open README.md'
    Write-Host '  Get-Help Send-DocDigitizerDocument -Full       # PowerShell native help'
    Write-Host ""

    Write-Host "ENVIRONMENT VARIABLES" -ForegroundColor Yellow
    Write-Host "---------------------"
    Write-Host "  DOCDIGITIZER_APIKEY   - Your API key (required - get one at docdigitizer.com/contact)"
    Write-Host "  DOCDIGITIZER_URL      - API base URL (default: https://apix.docdigitizer.com/sync)"
    Write-Host "  DOCDIGITIZER_PIPELINE - Default pipeline name"
    Write-Host "  DOCDIGITIZER_LOGLEVEL - Default log level (Minimal/Medium/Full)"
    Write-Host "  DOCDIGITIZER_TIMEOUT  - Request timeout in seconds (default: 300)"
    Write-Host ""
}
