# DocOntology Documentation

DocOntology is a SaaS API for storing and selecting JSON schemas used in document extraction pipelines.

---

## Overview

DocOntology provides centralized management of JSON schemas for document classification and extraction systems. It enables:

- **Schema Storage**: Version-controlled JSON schemas with lifecycle management
- **Schema Selection**: Find the best matching schema based on document type and country
- **Reference Data**: Manage document types and country classifications

**Who it's for:**
- Document extraction pipelines that need schema-driven field extraction
- Document classifiers that determine document types and select appropriate schemas
- Multi-tenant systems requiring customer-specific schema customization

---

## Getting Started

### Prerequisites

- API key (request at [docdigitizer.com/contact](https://docdigitizer.com/contact))

### Quick Start

Find the best schema for a document type:

```bash
curl -X POST https://api.docdigitizer.com/registry/schemas/find-best \
  -H "Content-Type: application/json" \
  -d '{"docTypeCode": "Invoice", "countryCode": "PT"}'
```

---

## Core Concepts

### Document Types

Document types categorize the kind of document being processed. Common types include:

| Code | Description |
|------|-------------|
| Invoice | Commercial invoices |
| Receipt | Point-of-sale receipts |
| Contract | Legal contracts and agreements |
| CV | Resumes and curriculum vitae |
| BankStatement | Financial statements |

### Countries

Countries use ISO 3166-1 alpha-2 codes for geographic classification:

| Code | Country |
|------|---------|
| PT | Portugal |
| ES | Spain |
| DE | Germany |
| FR | France |
| US | United States |

### Schemas

Schemas are versioned JSON schema definitions that specify what fields to extract from a document type. Each schema has:

- **Name**: Human-readable identifier
- **Content**: JSON Schema definition
- **Document Type**: What kind of document this schema handles
- **Country**: Geographic scope (optional for generic schemas)
- **Visibility**: Access control level
- **Status**: Lifecycle state

---

## Schema Visibility

Schemas have three visibility levels:

| Visibility | Description | Requirements |
|------------|-------------|--------------|
| `public` | Available to all users | Requires `docTypeCode` AND `countryCode` |
| `community` | Available to logged-in users | Future feature |
| `private` | Owner only | Requires `customerId` |

---

## Schema Lifecycle

Schemas progress through three states:

1. **Draft**: Initial state, can be modified freely
2. **Active**: Published and in use, updates create new versions
3. **Deprecated**: No longer recommended for use

```
draft → active → deprecated
```

---

## API Reference

### Base URL

```
https://api.docdigitizer.com/registry
```

### Public Endpoints

| Method | Path | Description |
|--------|------|-------------|
| GET | `/health` | Health check with database status |
| GET | `/doc-types` | List active document types |
| GET | `/countries` | List active countries |
| GET | `/reference-data` | Get all active doc types and countries |
| POST | `/schemas/find-best` | Find best matching schema |

### Admin Endpoints

| Method | Path | Description |
|--------|------|-------------|
| POST | `/admin/schemas` | Create schema |
| GET | `/admin/schemas` | List schemas with filters |
| GET | `/admin/schemas/:id` | Get schema by ID |
| PATCH | `/admin/schemas/:id` | Update schema |
| POST | `/admin/schemas/:id/activate` | Activate draft schema |
| POST | `/admin/schemas/:id/deprecate` | Deprecate active schema |
| DELETE | `/admin/schemas/:id` | Delete draft schema |

---

## Finding Schemas

The `/schemas/find-best` endpoint finds the most appropriate schema based on your criteria.

### Request

```bash
curl -X POST https://api.docdigitizer.com/registry/schemas/find-best \
  -H "Content-Type: application/json" \
  -d '{
    "docTypeCode": "Invoice",
    "countryCode": "PT",
    "customerId": "optional-customer-id"
  }'
```

### Schema Selection Priority

When finding the best schema, the system follows this priority:

1. **Customer's own schemas** (if `customerId` provided)
2. **Public schemas with exact country match**
3. **Public schemas without country** (generic fallback)

### Response

```json
{
  "schema": {
    "publicId": "sch_abc123",
    "publicVersionId": "schv_xyz789",
    "name": "Invoice Portugal",
    "docTypeCode": "Invoice",
    "countryCode": "PT",
    "visibility": "public",
    "status": "active",
    "content": {
      "type": "object",
      "properties": {
        "invoiceNumber": { "type": "string" },
        "nif": { "type": "string" },
        "totalAmount": { "type": "number" }
      }
    }
  },
  "matchType": "exact"
}
```

---

## Managing Schemas

### Create a Schema

```bash
curl -X POST https://api.docdigitizer.com/registry/admin/schemas \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Invoice Portugal",
    "content": {
      "type": "object",
      "properties": {
        "invoiceNumber": { "type": "string" },
        "nif": { "type": "string" },
        "totalAmount": { "type": "number" }
      }
    },
    "docTypeCode": "Invoice",
    "countryCode": "PT",
    "visibility": "public"
  }'
```

### Activate a Schema

```bash
curl -X POST https://api.docdigitizer.com/registry/admin/schemas/sch_abc123/activate
```

### Deprecate a Schema

```bash
curl -X POST https://api.docdigitizer.com/registry/admin/schemas/sch_abc123/deprecate
```

### List Schemas with Filters

```bash
# List active schemas for invoices
curl "https://api.docdigitizer.com/registry/admin/schemas?status=active&docTypeCode=Invoice"

# List all schemas for Portugal
curl "https://api.docdigitizer.com/registry/admin/schemas?countryCode=PT"
```

---

## CLI Tools

For command-line management, we provide two options:

### schemactl (Cross-platform CLI)

A Go-based CLI tool for Linux, macOS, and Windows.

```bash
# Check connection
schemactl health

# List document types
schemactl doc-types list

# Find best schema
schemactl schemas find-best --doc-type Invoice --country PT
```

### PowerShell Module

For Windows users, a PowerShell module is available.

```powershell
# Test connection
Test-SRConnection

# List document types
Get-SRDocType

# Find best schema
Find-SRSchema -DocTypeCode "Invoice" -CountryCode "PT"
```

See the [cmdlets documentation](cmdlets/README.md) for installation and usage details.

---

## Error Handling

### Error Response Format

```json
{
  "error": {
    "code": "NOT_FOUND",
    "message": "Schema not found"
  }
}
```

### Common Error Codes

| Code | Description | Solution |
|------|-------------|----------|
| `NOT_FOUND` | Resource not found | Check the ID or code |
| `VALIDATION_ERROR` | Invalid request data | Check request format |
| `CONFLICT` | Resource already exists | Use a different code |
| `INVALID_STATE` | Invalid state transition | Check schema status |

---

## Troubleshooting

### No schema found

- Verify the document type code exists
- Check if there's a schema for your country
- Try without country code to get a generic schema

### Cannot activate schema

- Schema must be in draft status
- Public schemas require both docTypeCode and countryCode

### Cannot delete schema

- Only draft schemas can be deleted
- Active or deprecated schemas must be deprecated first

---

## Support

For questions, issues, or custom requirements:

- **Email**: support@docdigitizer.com
- **Website**: [docdigitizer.com/contact](https://docdigitizer.com/contact)
