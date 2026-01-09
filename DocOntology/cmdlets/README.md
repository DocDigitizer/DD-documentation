# Schema Registry

PowerShell module for managing JSON schemas, document types, and countries through the Schema Registry API.

## Installation

```powershell
# Clone the repository
git clone https://github.com/DocDigitizer/schemactl.git

# Install the module
cd schemactl/powershell
.\Install.ps1 -AddToProfile
```

Restart PowerShell or run `Import-Module SchemaRegistry`.

## Quick Start

```powershell
# Test connection
Test-SRConnection

# List document types
Get-SRDocType

# List countries
Get-SRCountry

# List schemas
Get-SRSchema

# Get help
Get-SRHelp
```

## Configuration

```powershell
# Optional - set API URL and key
$env:SCHEMACTL_API_URL = "https://your-api-url"
$env:SCHEMACTL_API_KEY = "your-api-key"
```

| Variable | Description | Default |
|----------|-------------|---------|
| `SCHEMACTL_API_URL` | API base URL | Built-in default |
| `SCHEMACTL_API_KEY` | API key for authentication | (none) |
| `SCHEMACTL_TIMEOUT` | Request timeout in seconds | 30 |

## Commands

### Document Types

```powershell
Get-SRDocType                                    # List all
Get-SRDocType -Code "invoice"                    # Get one
New-SRDocType -Code "invoice" -Name "Invoice"    # Create
Set-SRDocType -Code "invoice" -Name "Tax Invoice"# Update
Remove-SRDocType -Code "invoice"                 # Delete
```

### Countries

```powershell
Get-SRCountry                                    # List all
Get-SRCountry -Code "PT"                         # Get one
New-SRCountry -Code "PT" -Name "Portugal"        # Create
Set-SRCountry -Code "PT" -Name "Portuguese Republic"  # Update
Remove-SRCountry -Code "PT"                      # Delete
```

### Schemas

```powershell
# List and filter
Get-SRSchema
Get-SRSchema -Status active
Get-SRSchema -DocTypeCode "invoice" -CountryCode "PT"

# Get specific schema
Get-SRSchema -Id "sch_abc123"

# Create schema
$content = @{
    type = "object"
    properties = @{
        invoiceNumber = @{ type = "string" }
        amount = @{ type = "number" }
    }
}
New-SRSchema -Name "Invoice Schema" -DocTypeCode "invoice" -Content $content

# Update
Set-SRSchema -Id "sch_abc123" -Name "Updated Name"

# Lifecycle: draft -> active -> deprecated
Enable-SRSchema -Id "sch_abc123"   # Activate
Disable-SRSchema -Id "sch_abc123"  # Deprecate

# Delete (draft only)
Remove-SRSchema -Id "sch_abc123"

# Find best match
Find-SRSchema -DocTypeCode "invoice" -CountryCode "PT"

# Get all versions
Get-SRSchemaVersion -Id "sch_abc123"
```

### Utilities

```powershell
Test-SRConnection    # Health check
Get-SRReferenceData  # All doc types and countries
Get-SRHelp           # Full help guide
Get-SRHelp -Examples # All examples
```

## Getting Help

```powershell
Get-SRHelp                        # Overview
Get-SRHelp -Examples              # All examples
Get-Help Get-SRSchema -Full       # Detailed help for a command
Get-Command -Module SchemaRegistry # List all commands
```

---

## Command Reference

| Command | Description |
|---------|-------------|
| `Get-SRDocType` | List or get document types |
| `New-SRDocType` | Create a document type |
| `Set-SRDocType` | Update a document type |
| `Remove-SRDocType` | Delete a document type |
| `Get-SRCountry` | List or get countries |
| `New-SRCountry` | Create a country |
| `Set-SRCountry` | Update a country |
| `Remove-SRCountry` | Delete a country |
| `Get-SRSchema` | List or get schemas |
| `Get-SRSchemaVersion` | Get all versions of a schema |
| `New-SRSchema` | Create a schema (draft) |
| `Set-SRSchema` | Update a schema |
| `Enable-SRSchema` | Activate a schema |
| `Disable-SRSchema` | Deprecate a schema |
| `Remove-SRSchema` | Delete a draft schema |
| `Find-SRSchema` | Find best matching schema |
| `Invoke-SRSchemaMatch` | Upload file for classification |
| `Test-SRConnection` | Check API health |
| `Get-SRReferenceData` | Get all reference data |
| `Get-SRHelp` | Show help guide |

---

## Cross-Platform CLI

A standalone CLI (`schemactl`) is also available for Linux, macOS, and Windows. See the [releases page](https://github.com/DocDigitizer/schemactl/releases) for binaries.
