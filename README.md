# DocDigitizer Documentation

Welcome to the official documentation repository for DocDigitizer products.

---

## Products

### DocDigitizer Sync

Intelligent document processing API that transforms unstructured documents into structured, actionable data through OCR, classification, and AI-powered extraction.

**Key Features:**
- OCR extraction via Google Cloud Vision
- AI-powered document classification
- LLM-based field extraction
- Multi-document PDF splitting

| Resource | Description |
|----------|-------------|
| [Usage Guide](./DocDigitizerSync/index.html) | User-friendly HTML guide |
| [Documentation](./DocDigitizerSync/documentation.md) | Technical documentation |
| [API Specification](./DocDigitizerSync/openapi.yaml) | OpenAPI 3.0 spec |
| [PowerShell Cmdlets](./DocDigitizerSync/cmdlets/) | Command-line tools |

---

### DocOntology

SaaS API for storing and selecting JSON schemas used in document extraction pipelines.

**Key Features:**
- Version-controlled JSON schemas
- Lifecycle management (draft, active, deprecated)
- Smart schema selection by document type and country
- Multi-tenant support with visibility controls

| Resource | Description |
|----------|-------------|
| [Usage Guide](./DocOntology/index.html) | User-friendly HTML guide |
| [Documentation](./DocOntology/documentation.md) | Technical documentation |
| [API Specification](./DocOntology/openapi.yaml) | OpenAPI 3.0 spec |
| [PowerShell Cmdlets](./DocOntology/cmdlets/) | Command-line tools |

---

## Command Line Tools

### DocDigitizer Sync - PowerShell Module

```powershell
# Import the module
Import-Module .\DocDigitizerSync\cmdlets\DocDigitizer.psd1

# Test connection
Test-DocDigitizerConnection

# Process a document
Send-DocDigitizerDocument -FilePath "invoice.pdf"
```

**Available Commands:**

| Command | Description |
|---------|-------------|
| `Send-DocDigitizerDocument` | Process a PDF document |
| `Test-DocDigitizerConnection` | Test API connectivity |
| `Get-DocDigitizerConfig` | View current settings |
| `Set-DocDigitizerConfig` | Update settings |
| `Get-DocDigitizerHelp` | Show help |

### DocOntology - PowerShell Module

```powershell
# Import the module
Import-Module .\DocOntology\cmdlets\SchemaRegistry.psd1

# Test connection
Test-SRConnection

# List document types
Get-SRDocType

# Find best schema
Find-SRSchema -DocTypeCode "Invoice" -CountryCode "PT"
```

**Available Commands:**

| Command | Description |
|---------|-------------|
| `Get-SRDocType` | List document types |
| `Get-SRCountry` | List countries |
| `Get-SRSchema` | List or get schemas |
| `Find-SRSchema` | Find best matching schema |
| `New-SRSchema` | Create a schema |
| `Enable-SRSchema` | Activate a schema |
| `Test-SRConnection` | Check API health |

---

## Getting Started

1. **For Document Processing**: Start with the [DocDigitizer Sync Usage Guide](./DocDigitizerSync/index.html)
2. **For Schema Management**: Start with the [DocOntology Usage Guide](./DocOntology/index.html)

## API Access

To use the DocDigitizer APIs, you need an API key. Request one at:

**[docdigitizer.com/contact](https://docdigitizer.com/contact)**

---

## Repository Structure

```
DD-documentation/
├── README.md                          # This file
├── DocDigitizerSync/
│   ├── documentation.md               # Technical documentation
│   ├── index.html                     # Usage guide
│   ├── openapi.yaml                   # API specification
│   └── cmdlets/                       # PowerShell module
│       ├── README.md
│       ├── DocDigitizer.psd1
│       ├── DocDigitizer.psm1
│       ├── Public/                    # Public cmdlets
│       └── Private/                   # Helper functions
└── DocOntology/
    ├── documentation.md               # Technical documentation
    ├── index.html                     # Usage guide
    ├── openapi.yaml                   # API specification
    └── cmdlets/                       # PowerShell module
        └── README.md
```

---

## Support

For questions, issues, or custom requirements:

- **Email**: support@docdigitizer.com
- **Website**: [docdigitizer.com/contact](https://docdigitizer.com/contact)

---

Copyright (c) 2025 DocDigitizer. All rights reserved.
