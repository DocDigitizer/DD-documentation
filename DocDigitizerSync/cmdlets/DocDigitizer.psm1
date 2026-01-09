#Requires -Version 5.1

<#
.SYNOPSIS
    DocDigitizer PowerShell Module

.DESCRIPTION
    A PowerShell module for interacting with the DocDigitizer document processing API.
    Provides cmdlets for sending documents, testing connections, and managing configuration.

.NOTES
    Module: DocDigitizer
    Author: DocDigitizer Team
#>

# Get the module path
$ModulePath = $PSScriptRoot

# Import private functions (not exported)
$PrivateFunctions = @(Get-ChildItem -Path "$ModulePath\Private\*.ps1" -ErrorAction SilentlyContinue)
foreach ($Function in $PrivateFunctions) {
    try {
        . $Function.FullName
        Write-Verbose "Imported private function: $($Function.BaseName)"
    }
    catch {
        Write-Error "Failed to import private function $($Function.FullName): $_"
    }
}

# Import public functions (exported)
$PublicFunctions = @(Get-ChildItem -Path "$ModulePath\Public\*.ps1" -ErrorAction SilentlyContinue)
foreach ($Function in $PublicFunctions) {
    try {
        . $Function.FullName
        Write-Verbose "Imported public function: $($Function.BaseName)"
    }
    catch {
        Write-Error "Failed to import public function $($Function.FullName): $_"
    }
}

# Export only public functions
Export-ModuleMember -Function $PublicFunctions.BaseName

# Display module loaded message
Write-Host "[DocDigitizer] Module loaded - $($PublicFunctions.Count) commands available" -ForegroundColor DarkCyan
