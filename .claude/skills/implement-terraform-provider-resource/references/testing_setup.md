# Testing setup reference

Utilities used by provider tests in `internal/provider`.

## `isIntegrationTestMode()`

- Defined in `internal/provider/utils.go`.
- Returns true when `TF_ACC` is `1`.
- Skips acceptance-style tests during normal `make test` / `go test` runs.

## `testAccPreCheck(t)`

- Defined in `internal/provider/provider_test.go`.
- Skips unless `TF_ACC=1`. Requires `N8N_ENDPOINT` and `N8N_API_KEY` for live-API acceptance tests.
- Prefer `make testacc` against a **licensed** n8n for project APIs (`feat:projectRole:admin`). Community `make testacc-docker` starts real n8n and skips those tests on 403.

## `testAccProtoV6ProviderFactories`

- Map key must match the provider local name in HCL: `"n8n"`.
- Uses `providerserver.NewProtocol6WithError(New("test")())`.

## Adding acceptance tests

- Build a `provider "n8n" { ... }` block (or rely on `N8N_ENDPOINT` / `N8N_API_KEY`).
- Example resource address: `n8n_project.example` (see `examples/resources/n8n_project/resource.tf`).
