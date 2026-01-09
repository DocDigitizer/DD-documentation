@{
    # Script module file associated with this manifest
    RootModule        = 'DocDigitizer.psm1'

    # Version number of this module
    ModuleVersion     = '1.0.0'

    # ID used to uniquely identify this module
    GUID              = 'a1b2c3d4-e5f6-7890-abcd-ef1234567890'

    # Author of this module
    Author            = 'DocDigitizer Team'

    # Company or vendor of this module
    CompanyName       = 'DocDigitizer'

    # Copyright statement for this module
    Copyright         = '(c) 2025 DocDigitizer. All rights reserved.'

    # Description of the functionality provided by this module
    Description       = 'PowerShell module for interacting with the DocDigitizer document processing API. Send PDF documents for OCR, classification, and data extraction.'

    # Minimum version of the PowerShell engine required by this module
    PowerShellVersion = '5.1'

    # Functions to export from this module
    FunctionsToExport = @(
        'Send-DocDigitizerDocument'
        'Test-DocDigitizerConnection'
        'Get-DocDigitizerConfig'
        'Set-DocDigitizerConfig'
        'Get-DocDigitizerHelp'
    )

    # Cmdlets to export from this module
    CmdletsToExport   = @()

    # Variables to export from this module
    VariablesToExport = @()

    # Aliases to export from this module
    AliasesToExport   = @()

    # Private data to pass to the module specified in RootModule
    PrivateData       = @{
        PSData = @{
            # Tags applied to this module for discoverability
            Tags         = @('DocDigitizer', 'OCR', 'PDF', 'DocumentProcessing', 'API')

            # License URI for this module
            LicenseUri   = ''

            # Project site URI for this module
            ProjectUri   = ''

            # Release notes for this module
            ReleaseNotes = @'
## Version 1.0.0
- Initial release
- Send-DocDigitizerDocument: Send PDF documents for processing
- Test-DocDigitizerConnection: Test API connectivity
- Get-DocDigitizerConfig: View current configuration
- Set-DocDigitizerConfig: Configure module defaults
- Get-DocDigitizerHelp: Show help and usage examples
'@
        }
    }
}
