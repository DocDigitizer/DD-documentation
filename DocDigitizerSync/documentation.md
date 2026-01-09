# DocDigitizer Sync Documentation

DocDigitizer Sync is an intelligent document processing API that transforms unstructured documents into structured, actionable data through OCR, classification, and AI-powered extraction.

---

## Overview

DocDigitizer Sync provides:

- **OCR Extraction**: Extract text from scanned documents using Google Cloud Vision
- **Document Classification**: Automatically identify document types (invoice, receipt, contract, etc.)
- **Data Extraction**: Extract structured fields based on document type using AI
- **Multi-Document Processing**: Automatically split and process multi-document PDFs

---

## Getting Started

### Prerequisites

- API key (request at [docdigitizer.com/contact](https://docdigitizer.com/contact))
- PDF documents to process

### Quick Start

Send a document for processing:

```bash
curl -X POST https://apix.docdigitizer.com/sync \
  -H "X-API-Key: YOUR_API_KEY" \
  -F "file=@invoice.pdf"
```

---

## API Reference

### Base URL

```
https://apix.docdigitizer.com/sync
```

### Authentication

Include your API key in the request header:

```
X-API-Key: your-api-key-here
```

### Endpoints

| Method | Path | Description |
|--------|------|-------------|
| POST | `/` | Process document(s) |
| GET | `/` | Health check |

---

## Processing Documents

### Request

Send a PDF file using multipart form data:

```bash
curl -X POST https://apix.docdigitizer.com/sync \
  -H "X-API-Key: YOUR_API_KEY" \
  -F "file=@document.pdf"
```

### Request Parameters

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `file` | File | Yes | PDF document to process |
| `pipeline` | String | No | Specific pipeline to use (default: auto-select) |
| `documentId` | UUID | No | Custom document identifier |
| `contextId` | UUID | No | Custom context identifier for grouping |

### Response

```json
{
  "requestToken": "ABC1234",
  "state": "Completed",
  "pipeline": "MainPipelineWithOCR",
  "pageCount": 3,
  "output": {
    "extractions": [
      {
        "documentType": "Invoice",
        "fields": {
          "invoiceNumber": "INV-2024-001",
          "date": "2024-01-15",
          "totalAmount": 1250.00,
          "vendorName": "Acme Corp"
        }
      }
    ]
  },
  "timers": {
    "total": 2345.67,
    "ocr": 1200.00,
    "classification": 450.00,
    "extraction": 695.67
  }
}
```

### Response Fields

| Field | Description |
|-------|-------------|
| `requestToken` | Unique trace ID for this request (use for support inquiries) |
| `state` | Processing status: `Completed`, `Failed`, `Partial` |
| `pipeline` | Pipeline used for processing |
| `pageCount` | Number of pages in the document |
| `output.extractions` | Array of extracted documents with fields |
| `timers` | Processing time breakdown (milliseconds) |

---

## Response Headers

Each response includes timing and tracing headers:

| Header | Description |
|--------|-------------|
| `X-DD-TraceId` | Request trace ID (same as `requestToken`) |
| `X-DD-Timer-Total` | Total processing time (ms) |
| `X-DD-Timer-OCR` | OCR processing time (ms) |
| `X-DD-Timer-Classification` | Classification time (ms) |

---

## Supported Document Types

DocDigitizer Sync can classify and extract data from various document types:

| Document Type | Description |
|---------------|-------------|
| Invoice | Commercial invoices with line items |
| Receipt | Point-of-sale receipts |
| Contract | Legal contracts and agreements |
| CV | Resumes and curriculum vitae |
| ID Document | Identity documents |
| Bank Statement | Financial statements |

The system automatically detects the document type and applies the appropriate extraction schema.

---

## Multi-Document Processing

When you submit a PDF containing multiple documents (e.g., several invoices in one file), DocDigitizer Sync automatically:

1. Detects document boundaries
2. Splits the PDF into individual documents
3. Classifies each document separately
4. Extracts data from each document
5. Returns results grouped by document

---

## Error Handling

### Error Response Format

```json
{
  "requestToken": "XYZ7890",
  "state": "Failed",
  "error": {
    "code": "INVALID_FILE",
    "message": "The uploaded file is not a valid PDF"
  }
}
```

### Common Error Codes

| Code | Description | Solution |
|------|-------------|----------|
| `INVALID_FILE` | File is not a valid PDF | Ensure file is a valid PDF document |
| `FILE_TOO_LARGE` | File exceeds size limit | Split into smaller files |
| `OCR_FAILED` | OCR processing failed | Check document quality |
| `TIMEOUT` | Processing timed out | Retry or contact support |
| `UNAUTHORIZED` | Invalid API key | Check your API key |

---

## PowerShell Tools

For easy integration with Windows environments, use our PowerShell module.

See the [cmdlets documentation](cmdlets/README.md) for:

- Installation instructions
- Available commands
- Usage examples

### Quick Example

```powershell
Import-Module .\DocDigitizer.psd1

# Test connection
Test-DocDigitizerConnection

# Process a document
Send-DocDigitizerDocument -FilePath "invoice.pdf"

# Process and save result
Send-DocDigitizerDocument -FilePath "invoice.pdf" -SaveExtraction
```

---

## Troubleshooting

### Document not recognized

- Ensure the document is clear and readable
- Check that text is not too small or blurry
- Verify the PDF is not corrupted

### Slow processing

- Large documents take longer to process
- Complex documents with many pages may timeout
- Consider splitting very large documents

### Missing fields in extraction

- Some fields may not be present in all documents
- Field names may vary by document type
- Contact support for custom extraction requirements

### Connection issues

- Verify your API key is correct
- Check your internet connection
- Ensure the API endpoint is accessible from your network

---

## Support

For questions, issues, or custom requirements:

- **Email**: support@docdigitizer.com
- **Website**: [docdigitizer.com/contact](https://docdigitizer.com/contact)

When contacting support, include your `requestToken` for faster resolution.
