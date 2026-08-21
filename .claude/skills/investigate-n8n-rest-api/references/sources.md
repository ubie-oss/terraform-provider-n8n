# Canonical sources for n8n Public API investigation

Use these URLs when following the investigate-n8n-rest-api skill. Prefer **live**
fetches over memory. Do **not** vendor or commit a full OpenAPI snapshot into
this skill directory.

## Source order

1. **Primary A** — resolved instance OpenAPI at `{base}/api/v1/openapi.yml`
2. **Primary B** — GitHub split OpenAPI (fallback when no instance)
3. **Secondary** — official docs under https://docs.n8n.io/connect/n8n-api
4. **Optional verify** — `GET /api/v1/discover`

On conflicts, prefer OpenAPI (A or B) over docs; document drift; never invent.

## Primary A: Instance OpenAPI (preferred)

Live n8n instances (Cloud and self-hosted) often serve a **resolved** OpenAPI
bundle:

```text
{base}/api/v1/openapi.yml
```

Example pattern:

```text
https://internal.users.n8n.cloud/api/v1/openapi.yml
```

### How to fetch

1. Take the instance base from `N8N_ENDPOINT`, provider config, or the user.
2. Ensure the API root ends with `/api/v1` (same normalization idea as
   `internal/n8n.NormalizeEndpoint`).
3. `GET {base}/api/v1/openapi.yml` with `accept: application/yaml` or without
   special Accept. Send `X-N8N-API-KEY` if the instance requires authentication.
4. Record the exact URL and `info.version` in the contract report.
5. Extract the target resource’s paths and `#/components/schemas/...` refs from
   the same file. Do not commit the downloaded YAML into the skill tree.

### Path-overlap note (workflows)

Resolved instance specs have been observed to expose overlapping workflow paths:

| Pattern                  | Examples                                                                                     |
| ------------------------ | -------------------------------------------------------------------------------------------- |
| Classic `{id}`           | `/workflows/{id}`, activate, deactivate, publish, unpublish, archive, unarchive, transfer    |
| Decorator `{workflowId}` | `/workflows/{workflowId}`, `/workflows/{workflowId}/history`, `/workflows/{workflowId}/tags` |

Some decorator-routed operations set `x-decorator-routed: true` and
`x-eov-operation-id: unreachable` while still publishing a real `operationId`.
List all matching ops; prefer documented Public CRUD paths for Terraform; note
duplicates under Gaps/risks.

## Primary B: OpenAPI in n8n-io/n8n (fallback)

When no instance is reachable, use the split OpenAPI under:

```text
packages/cli/src/public-api/v1/
```

### Index file

| Item      | URL pattern                                                                                     |
| --------- | ----------------------------------------------------------------------------------------------- |
| Raw index | `https://raw.githubusercontent.com/n8n-io/n8n/<ref>/packages/cli/src/public-api/v1/openapi.yml` |
| Browse    | `https://github.com/n8n-io/n8n/blob/<ref>/packages/cli/src/public-api/v1/openapi.yml`           |

- `<ref>`: user-pinned **git tag** or release when known; otherwise `master`.
- Record `info.version` and the git ref used.
- `servers`: `/api/v1` and `{url}/api/v1`.

### How to follow `$ref`s

1. Open the index `paths:` map. Each entry points at a relative file, for example:
   - `./handlers/workflows/spec/paths/workflows.yml`
   - `./handlers/workflows/spec/paths/workflows.id.yml`
2. Resolve relative to the directory of the file that contains the `$ref`.
3. Fetch the path YAML. Extract for each HTTP method:
   - `x-eov-operation-id` (and/or `operationId` if present)
   - **`x-required-scope`** (often missing from human docs)
   - `summary` / `description`
   - `parameters`, `requestBody`, `responses`
4. Follow schema `$ref`s into:
   - `handlers/<resource>/spec/schemas/`
   - `shared/spec/schemas/` (via `shared/spec/schemas/_index.yml`)
   - `shared/spec/responses/` (for example `unauthorized.yml`, `forbidden.yml`)
5. Security schemes in the index include:
   - **`ApiKeyAuth`**: header `X-N8N-API-KEY` — **use this for Terraform**
   - `BearerAuth`: HTTP bearer JWT — out of scope for API-key provider auth
   - `CookieAuth`: cookie `n8n-auth` — editor/session; not the Public API key path

### Do not use

- A historical n8n-docs path such as `n8n-io/n8n-docs/.../docs/api/v1/openapi.yml` — that bundled file has **404**'d; do not cite it as source of truth.
- Invented paths or fields not present in the fetched YAML.
- Vendoring the full instance OpenAPI into `assets/` or `references/`.

## Secondary: Official docs (markdown)

Docs overview (start here): https://docs.n8n.io/connect/n8n-api.md

That page introduces REST authentication, pagination, the playground, and the
endpoint reference. It also mentions **n8n CLI** as a developer-facing wrapper
around the Public API — useful for agents/CI, **not** a Terraform schema source.

Docs pages are available as Markdown by appending `.md` to the page URL.

| Topic                   | URL                                                          |
| ----------------------- | ------------------------------------------------------------ |
| API overview            | https://docs.n8n.io/connect/n8n-api.md                       |
| Authentication & scopes | https://docs.n8n.io/connect/n8n-api/authentication.md        |
| Pagination              | https://docs.n8n.io/connect/n8n-api/pagination.md            |
| Endpoint reference hub  | https://docs.n8n.io/connect/n8n-api/api-reference.md         |
| API playground notes    | https://docs.n8n.io/connect/n8n-api/use-an-api-playground.md |
| Docs index              | https://docs.n8n.io/llms.txt                                 |

### Per-resource pages (singular path)

Examples (append `.md` when fetching):

| Resource   | Docs path                                         |
| ---------- | ------------------------------------------------- |
| Workflow   | https://docs.n8n.io/connect/n8n-api/workflow.md   |
| Credential | https://docs.n8n.io/connect/n8n-api/credential.md |
| Execution  | https://docs.n8n.io/connect/n8n-api/execution.md  |
| User       | https://docs.n8n.io/connect/n8n-api/user.md       |
| Tags       | https://docs.n8n.io/connect/n8n-api/tags.md       |
| Variables  | https://docs.n8n.io/connect/n8n-api/variables.md  |
| Data table | https://docs.n8n.io/connect/n8n-api/data-table.md |
| Projects   | https://docs.n8n.io/connect/n8n-api/projects.md   |
| Folders    | https://docs.n8n.io/connect/n8n-api/folders.md    |

If unsure of the slug, search [llms.txt](https://docs.n8n.io/llms.txt) or
[sitemap.md](https://docs.n8n.io/sitemap.md). Wrong paths (for example
`api-reference/workflows`) may 404. For OpenAPI tags without a docs page, rely
on OpenAPI and note missing docs.

### Auth and base URL (verified)

- Header: **`X-N8N-API-KEY: <key>`**
- Cloud: `https://<subdomain>.app.n8n.cloud/api/v1`
- Self-hosted: `https://<host>/api/v1` (include path prefix if the instance uses one)
- API is not available on free trial (per docs)
- Enterprise API keys may be limited by scopes listed on the authentication page

### Pagination (verified)

- Default `limit`: 100
- Maximum `limit`: 250
- Cursor: request `cursor`; response `nextCursor`

## Optional live verify: discover

`GET /api/v1/discover` returns a capability map filtered by the caller's API key
scopes. Spec marks `x-required-scope: none`. Useful query params:

- `resource` — for example `workflow`, `tags`, `credential`
- `operation` — for example `read`, `create`, `list`
- `include=schemas` — inline request body schemas

Example:

```bash
curl -sS -H "X-N8N-API-KEY: $N8N_API_KEY" -H "accept: application/json" \
  "$N8N_BASE/api/v1/discover?resource=workflow&include=schemas"
```

Self-hosted built-in Swagger UI (not Cloud): `{host}/{path}/api/v1/docs`.

Do **not** send API keys through the documentation site's Scalar proxy.

## Conflict resolution

| Situation                                          | Action                                                                                     |
| -------------------------------------------------- | ------------------------------------------------------------------------------------------ |
| Docs omit `x-required-scope`                       | Keep scope from OpenAPI                                                                    |
| Path param name differs (`{id}` vs `{workflowId}`) | Prefer OpenAPI; list both; note drift / overlap                                            |
| Discover omits an operation OpenAPI lists          | Prefer OpenAPI; note key scope or edition may hide it                                      |
| Live GET shape differs from schema                 | Prefer OpenAPI for design; record live difference                                          |
| Instance OpenAPI vs GitHub OpenAPI differ          | Prefer **instance** when targeting that deployment; else prefer GitHub ref the user pinned |

Never invent a middle ground that appears in neither source.
