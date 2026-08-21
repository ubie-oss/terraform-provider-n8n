---
name: investigate-n8n-rest-api
description: >
  Investigate n8n Public REST API (/api/v1) from instance OpenAPI, GitHub
  OpenAPI, and official docs. Use when mapping n8n HTTP endpoints, schemas,
  auth, pagination, scopes, or errors for Terraform; when checking public API
  vs /rest vs MCP vs CLI; when looking up X-N8N-API-KEY, workflows, credentials,
  tags, executions, users, projects, variables, or data tables.
compatibility: Prefer a reachable n8n instance OpenAPI at {base}/api/v1/openapi.yml; else GitHub n8n-io/n8n OpenAPI. Also needs docs.n8n.io. Do not use n8n MCP or CLI as the REST contract.
---

# SOP: Investigate n8n Public REST API

## Purpose

Extract an accurate, versioned **n8n Public API** (`/api/v1`) contract for one
resource before writing Go client methods or Terraform schemas. Prefer a live
**resolved** OpenAPI from the target instance, then official docs. Emit a
Terraform-oriented contract report; do not invent paths, fields, or auth schemes.

This skill does **not** implement Terraform resources or the Go client. It does
**not** vendor a full OpenAPI dump into the skill directory.

## Workflow Checklist

- [ ] **1. Collect target** — resource name, Terraform intent, optional n8n version / instance URL
- [ ] **2. Confirm Public API** — only `/api/v1`; reject `/rest`, MCP, and CLI as provider contracts
- [ ] **3. Fetch OpenAPI** — instance `{base}/api/v1/openapi.yml` if reachable; else GitHub split YAML
- [ ] **4. Extract operations** — method, path, operationId / `x-eov-operation-id`, `x-required-scope`, schemas, errors
- [ ] **5. Fetch docs markdown** — overview + auth + pagination + resource page; note drift
- [ ] **6. Merge contract** — prefer OpenAPI on conflicts; note path overlaps / decorator-routed ops
- [ ] **7. Optional live verify** — `GET /api/v1/discover` (and other GETs only)
- [ ] **8. Write report** — fill [assets/contract-report.template.md](assets/contract-report.template.md)

## Hard Rules

1. **Invent nothing** — every path, field, and error code must appear in OpenAPI or official docs.
2. **Public API only** — Terraform maps to `/api/v1`. Treat `/rest`, MCP, and **n8n CLI** as out of scope for schema (see [references/api-surfaces.md](references/api-surfaces.md)).
3. **Auth for the provider** — use header `X-N8N-API-KEY`. Spec may also list Bearer JWT and cookie `n8n-auth`; those are not the Terraform API-key path.
4. **Verify client auth** — this provider’s client in `internal/n8n` must send `X-N8N-API-KEY`. Flag in the report only if code drifts away from that.
5. **Do not use n8n CLI as schema** — CLI wraps the Public API for automation; the REST OpenAPI is the contract.
6. **Optional live calls** — prefer `GET /api/v1/discover` (and other GETs). No POST/PUT/PATCH/DELETE unless the user explicitly asks.
7. **Secrets** — never log or commit API keys, credential secret bodies, or a downloaded OpenAPI dump into the skill tree.
8. **No third-party playground** — do not send keys through Scalar’s docs proxy.

## Detailed Instructions

### 1. Collect the target

Ask or infer:

1. **Resource** (for example `workflow`, `tag`, `credential`, `variable`). See [references/resource-index.md](references/resource-index.md).
2. **Terraform intent** — resource vs data source vs both.
3. **Instance base URL** — if available (`N8N_ENDPOINT` or user-provided host).
4. **n8n version** — git tag or release if known and no instance OpenAPI is available.

### 2. Stay on the Public API

If the user asked about `/rest/...`, MCP tools, or n8n CLI, document them as **not** a Terraform Public API contract and stop (or continue only for the Public API portion). Details: [references/api-surfaces.md](references/api-surfaces.md).

### 3. Fetch OpenAPI (instance first)

Canonical fetch rules: [references/sources.md](references/sources.md).

**Preferred — resolved instance OpenAPI** when a base URL is available:

```text
{base}/api/v1/openapi.yml
```

Example pattern: `https://internal.users.n8n.cloud/api/v1/openapi.yml`.

- Normalize `{base}` so the path ends with `/api/v1` (same idea as `internal/n8n.NormalizeEndpoint`).
- Send `X-N8N-API-KEY` if the instance requires auth to serve the spec.
- This file is typically a **resolved** bundle (inline schemas). Extract the target resource’s paths directly; no `$ref` chase required unless a `$ref` remains.

**Fallback — GitHub split OpenAPI** when no instance is reachable:

```text
https://raw.githubusercontent.com/n8n-io/n8n/<ref>/packages/cli/src/public-api/v1/openapi.yml
```

Use `<ref>` = user-pinned tag, or `master`. Follow path/schema `$ref`s as described in [references/sources.md](references/sources.md).

Record which source was used, `info.version`, and the instance URL or git ref.

### 4. Extract operations and schemas

1. Select paths for the target resource (see [references/resource-index.md](references/resource-index.md)).
2. For each HTTP method, record: path, `operationId` and/or `x-eov-operation-id`, **`x-required-scope`**, request/response schemas, error responses.
3. **Path overlap / decorator-routed ops:** Instance OpenAPI may list both `/workflows/{id}` and `/workflows/{workflowId}` (and tags/history under `{workflowId}`), and may set `x-decorator-routed: true` with `x-eov-operation-id: unreachable` while still exposing a real `operationId`. List **all** matching operations, prefer documented Public CRUD paths for Terraform mapping, and note duplicates under Gaps/risks.
4. Note extra lifecycle ops (for example workflow `activate` / `deactivate` / `publish` / `unpublish` / `archive` / `unarchive` / `transfer`) separately from core CRUD.

### 5. Fetch official docs markdown

Start from the overview: https://docs.n8n.io/connect/n8n-api.md

Also fetch authentication, pagination, and the matching resource page listed in [references/sources.md](references/sources.md). Append `.md` to docs.n8n.io URLs.

If docs and OpenAPI disagree on path params, field names, or verbs, **prefer OpenAPI** and note the drift in the report.

### 6. Merge into a Terraform-oriented contract

Map operations to Create / Read / Update / Delete / Import / List. Call out:

- Required vs computed attributes
- Secrets or write-only fields
- Pagination (`limit` / `cursor` / `nextCursor`)
- Edition or scope limits (Enterprise scopes, free-trial API unavailable)
- Duplicate or decorator-routed paths

### 7. Optional live verify

When an instance URL and API key are available:

```bash
curl -sS -H "X-N8N-API-KEY: $N8N_API_KEY" -H "accept: application/json" \
  "$N8N_BASE/api/v1/discover?resource=<resource>&include=schemas"
```

Use results only to confirm availability for that key. Do not replace OpenAPI with discover output when they conflict without documenting it.

### 8. Write the contract report

Copy [assets/contract-report.template.md](assets/contract-report.template.md) and fill every section. Cite the OpenAPI source URL (instance or GitHub) and `info.version`.

## Example

**Input:** "Investigate n8n workflows REST API for a Terraform resource."

**Actions:**

1. Fetch `{instance}/api/v1/openapi.yml` if available; else GitHub OpenAPI.
2. Extract `/workflows`, `/workflows/{id}`, activate/deactivate/publish/archive/transfer, and any `{workflowId}` duplicates.
3. Fetch https://docs.n8n.io/connect/n8n-api.md, authentication, pagination, and https://docs.n8n.io/connect/n8n-api/workflow.md.
4. Fill the contract report with CRUD coverage and gaps (activate/publish as extra ops; decorator-routed duplicates called out).

**Output:** A completed contract report, not Go or Terraform code.

## Success Criteria

- [ ] Report cites instance OpenAPI URL **or** GitHub git ref, plus `info.version`.
- [ ] Every listed operation includes method, path, operation id(s), and `x-required-scope` when present in the spec.
- [ ] Auth section specifies `X-N8N-API-KEY` for provider use; client auth verified or drift flagged.
- [ ] Pagination and error codes are documented from sources, not guessed.
- [ ] Terraform CRUD coverage and gaps (including path overlaps) are explicit.
- [ ] No invented endpoints; `/rest`, MCP, and CLI are not treated as the provider contract.

## Additional Resources

- [Canonical sources and fetch order](references/sources.md)
- [API surfaces: public vs /rest vs MCP vs CLI](references/api-surfaces.md)
- [Resource tag to path index](references/resource-index.md)
- [Contract report template](assets/contract-report.template.md)

## Optional / Related

After the report is complete, use it as input when implementing a Terraform resource or data source (for example with `implement-terraform-provider-resource`). That skill is **not** required to finish this investigation.
