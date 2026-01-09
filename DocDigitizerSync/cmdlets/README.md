# DocDigitizer PowerShell Module

A PowerShell module for processing documents with the DocDigitizer API. Extract data from invoices, receipts, and other documents using OCR and AI.

## Installation

### Step 1: Open PowerShell

Open **PowerShell** or **Windows Terminal** on your computer.

### Step 2: Choose a folder for the module

Navigate to where you want to install the module. For example:

```powershell
cd C:\Tools
```

### Step 3: Clone the repository

```powershell
git clone https://github.com/DocDigitizer/dd-cmdlets.git
```

This creates a `dd-cmdlets` folder with all the module files.

### Step 4: Get your API key

To use the DocDigitizer API, you need an API key. Request one at:

**https://docdigitizer.com/contact**

### Step 5: Import the module

```powershell
Import-Module .\dd-cmdlets\DocDigitizer.psd1
```

You should see:
```
[DocDigitizer] Module loaded - 5 commands available
```

### Step 6: Set your API key

```powershell
Set-DocDigitizerConfig -ApiKey "your-api-key-here"
```

To save it permanently (so you don't have to set it every time):

```powershell
Set-DocDigitizerConfig -ApiKey "your-api-key-here" -Persist
```

### Step 7: Test the connection

```powershell
Test-DocDigitizerConnection
```

You should see:
```
Url       : https://apix.docdigitizer.com/sync
Connected : True
Response  : I am alive
Latency   : ...
```

**Done!** You're ready to process documents.

---

## Quick Start

### Process a document

```powershell
Send-DocDigitizerDocument -FilePath "C:\Documents\invoice.pdf"
```

### Process and save the result

```powershell
Send-DocDigitizerDocument -FilePath "C:\Documents\invoice.pdf" -SaveExtraction
```

This creates `invoice_extraction.json` in the same folder as your PDF.

### Process multiple documents

```powershell
Get-ChildItem "C:\Documents\*.pdf" | Send-DocDigitizerDocument -SaveExtraction
```

---

## Available Commands

| Command | Description |
|---------|-------------|
| `Send-DocDigitizerDocument` | Process a PDF document |
| `Test-DocDigitizerConnection` | Test if the API is reachable |
| `Get-DocDigitizerConfig` | Show current settings |
| `Set-DocDigitizerConfig` | Change settings |
| `Get-DocDigitizerHelp` | Show help |

---

## Examples

### Basic usage

```powershell
# Process a single invoice
Send-DocDigitizerDocument -FilePath "invoice.pdf"

# Process and save result to JSON file
Send-DocDigitizerDocument -FilePath "invoice.pdf" -SaveExtraction

# Process all PDFs in a folder
Get-ChildItem *.pdf | Send-DocDigitizerDocument -SaveExtraction
```

### Working with results

```powershell
# Get result as JSON
$json = Send-DocDigitizerDocument -FilePath "invoice.pdf"

# Convert to PowerShell object
$result = $json | ConvertFrom-Json

# View extractions
$result.output.extractions
```

### Get more details in response

```powershell
# Minimal response (default)
Send-DocDigitizerDocument -FilePath "invoice.pdf"

# Medium - includes metadata
Send-DocDigitizerDocument -FilePath "invoice.pdf" -LogLevel Medium

# Full - complete execution details
Send-DocDigitizerDocument -FilePath "invoice.pdf" -LogLevel Full
```

---

## Troubleshooting

### "git" is not recognized

You need to install Git first:
1. Download from https://git-scm.com/download/win
2. Install with default options
3. Restart PowerShell and try again

### Module not loading

Make sure you're in the correct folder:

```powershell
# Check current folder
Get-Location

# List files to verify dd-cmdlets exists
Get-ChildItem
```

### Connection failed

Check your internet connection and try:

```powershell
Test-DocDigitizerConnection -Verbose
```

### Timeout on large documents

Increase the timeout (default is 5 minutes):

```powershell
Send-DocDigitizerDocument -FilePath "large.pdf" -TimeoutSec 600
```

---

## Requirements

- Windows PowerShell 5.1 or PowerShell 7+
- Internet connection
- Git (for installation)

---

## Documentation

- [Installation Guide](installation.html)
- [Commands Reference](commands.html)
- [Scripting Examples](examples.html)

---

## Support

For issues or questions, contact DocDigitizer support.
