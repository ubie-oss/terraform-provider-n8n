# n8n HTTP surfaces (what Terraform may use)

Agents confuse several n8n HTTP-ish surfaces. Only the Public REST API is the
contract for a Terraform provider.

## 1. Public REST API — `/api/v1` (in scope)

| Attribute                 | Detail                                                                                                                                                      |
| ------------------------- | ----------------------------------------------------------------------------------------------------------------------------------------------------------- |
| Purpose                   | Programmatic CRUD/list for workflows, credentials, tags, users, projects, variables, data tables, executions, and related resources                         |
| Auth for providers        | Header `X-N8N-API-KEY`                                                                                                                                      |
| Machine-readable contract | Prefer `{instance}/api/v1/openapi.yml` (resolved). Fallback: OpenAPI under `packages/cli/src/public-api/v1/` in [n8n-io/n8n](https://github.com/n8n-io/n8n) |
| Docs                      | https://docs.n8n.io/connect/n8n-api.md                                                                                                                      |
| Terraform                 | **Use this surface only** for provider client design                                                                                                        |

Also listed in OpenAPI (do **not** use for API-key Terraform auth):

- Bearer JWT (`BearerAuth`)
- Cookie `n8n-auth` (`CookieAuth`) — session/editor style

This provider’s HTTP client in `internal/n8n` sends **`X-N8N-API-KEY`**. If
investigation finds different header behavior in code, flag it in the contract
report.

## 2. Editor / internal API — `/rest/...` (out of scope)

| Attribute | Detail                                                                      |
| --------- | --------------------------------------------------------------------------- |
| Purpose   | UI and editor traffic                                                       |
| Auth      | Typically session cookie (`n8n-auth`), not a stable Public API key contract |
| Docs      | Not the Public API reference                                                |
| Terraform | **Do not** implement provider resources against `/rest`                     |

If investigation finds only `/rest` endpoints for a feature, report that the
feature is **not** available on the Public API (or not documented there) and
stop unless the user explicitly asked about the editor API.

## 3. MCP — for example `/mcp-server/http` (out of scope)

| Attribute | Detail                                                               |
| --------- | -------------------------------------------------------------------- |
| Purpose   | Model Context Protocol tools for agents to build/run workflows       |
| Auth      | MCP / OAuth / connector-specific; not the Public REST OpenAPI schema |
| Terraform | **Do not** treat MCP tool lists as REST schemas or CRUD contracts    |

MCP is useful for automation product behavior; it is not a substitute for
OpenAPI when designing Terraform resources.

## 4. n8n CLI (out of scope for schema)

| Attribute | Detail                                                                                            |
| --------- | ------------------------------------------------------------------------------------------------- |
| Purpose   | Command-line wrapper around the Public API (docs recommend it for CI/CD and AI agent integration) |
| Docs      | Linked from https://docs.n8n.io/connect/n8n-api.md                                                |
| Terraform | **Do not** use CLI help output or CLI commands as the OpenAPI contract                            |

CLI may be useful after the contract is known; investigation must still start
from OpenAPI + docs. For Community Edition vs licensed behavior observed via CLI against
`docker-compose.dev.yml`, see
[`dev/docs/n8n-ce-public-api-and-cli-limits.md`](../../../../dev/docs/n8n-ce-public-api-and-cli-limits.md)
(do not treat that matrix as OpenAPI).

## Quick decision

```text
User asks about an n8n HTTP endpoint or "API"
  → Is the path under /api/v1 (or documented as Public API)?
      YES → investigate with instance/GitHub OpenAPI + docs (this skill)
      NO, /rest → out of scope for Terraform; document and stop
      NO, MCP  → out of scope for Terraform; document and stop
      NO, CLI  → out of scope for schema; document and use OpenAPI instead
```
