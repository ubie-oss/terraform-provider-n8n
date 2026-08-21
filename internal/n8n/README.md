# `internal/n8n`

HTTP client for the [n8n Public API](https://docs.n8n.io/connect/n8n-api/).

## Behavior

- **`Client.HTTP`**: standard `*http.Client` whose transport applies, in order: **max concurrent in-flight requests** (`max_concurrent_requests`), **token-bucket rate limit** (`requests_per_second`), then **`X-N8N-API-KEY`** using the provider `api_key`.
- **Endpoint**: instance URL is normalized to end with `/api/v1`.
- **Defaults** (when provider omits optional attributes): `10` concurrent requests, `10` requests per second.

## Layout

- `client.go` / `request.go` / `errors.go`: shared HTTP, JSON, 404, and license-403 helpers
- `models/`: request and response structs
- `api/v1/<resource>/`: one file per Public API versioned operation (for example `api/v1/projects/list_projects_v1.go`)
- `services/`: orchestration for awkward APIs (project Read is list-then-find; there is no `GET /projects/{id}`; folder list uses skip/take, not cursor; credential list uses cursor like projects)
- `controllers/`: Terraform-facing orchestrators used by resources and data sources

Terraform resources and data sources call [`controllers/`](controllers), which wrap these services. Do not call `services` from `internal/provider`.

Provider attributes are defined in [`internal/provider/provider.go`](../provider/provider.go).
