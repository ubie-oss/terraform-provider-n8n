# n8n Public API contract report

Fill every section from OpenAPI and official docs. Prefer OpenAPI on conflicts
(instance OpenAPI when targeting that deployment; else GitHub). Cite URLs and
`info.version`. Do not invent fields.

## Target

- **Resource name:**
- **Terraform intent:** (resource / data source / both)
- **Instance base URL (if any):**
- **n8n version / git ref investigated:**
- **OpenAPI `info.version`:**

## Sources

| Role                     | URL or path | Notes                  |
| ------------------------ | ----------- | ---------------------- |
| Instance OpenAPI URL     |             | used / skipped         |
| GitHub OpenAPI ref       |             | used / skipped         |
| Path / schema extracts   |             |                        |
| Docs pages               |             | start: connect/n8n-api |
| Optional discover / live |             | skipped / used         |

## Auth

- **Provider auth scheme:** `X-N8N-API-KEY` (yes/no — must be yes for Public API keys)
- **Other schemes seen in spec:** (Bearer / Cookie — note out of scope for API-key provider)
- **Scopes required (Enterprise):** list `x-required-scope` values
- **Client auth header verified:** (`internal/n8n` uses `X-N8N-API-KEY` — yes / drift)
- **Edition / trial limits:**

## Base URL

- **Cloud pattern:**
- **Self-hosted pattern:**
- **Instance under test (if any):**

## Operations

| Method | Path | operationId / x-eov-operation-id | x-required-scope | Request schema | Response schema | Errors |
| ------ | ---- | -------------------------------- | ---------------- | -------------- | --------------- | ------ |
|        |      |                                  |                  |                |                 |        |

## Schemas (summary)

### Create / update request

- Required properties:
- Optional properties:
- Write-only / secret properties:
- Read-only / computed properties returned on write:

### Read response

- Identity fields (for Terraform `id` / import):
- Nested objects notes:

## Pagination

- Applies to list? (yes/no)
- Query params (`limit`, `cursor`, others):
- Response cursor field (`nextCursor`):
- Defaults / max:

## Terraform CRUD coverage

| Lifecycle          | Covered by | Notes |
| ------------------ | ---------- | ----- |
| Create             |            |       |
| Read               |            |       |
| Update             |            |       |
| Delete             |            |       |
| Import             |            |       |
| List (data source) |            |       |

Extra operations (activate, publish, archive, transfer, retry, etc.):

## Gaps and risks

- Missing Update/Delete:
- Docs vs OpenAPI drift:
- Duplicate / decorator-routed paths:
- CLI / MCP not used as schema: (confirmed)
- Client auth header verified:
- Secrets in state:
- Destructive or async operations:

## Recommended client methods

List intended Go client methods (names only), mapped to operations above. Do not
implement them in this skill.

1.
2.
3.

## Conclusion

One short paragraph: is this resource suitable for a Terraform resource/data
source with the Public API as investigated? What must be decided next?
