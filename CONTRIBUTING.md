# Contributing

This repository is the Terraform provider for n8n (`ubie-oss/n8n`), built with the Terraform Plugin Framework. It manages n8n **team projects**, **folders**, and **credentials**, not workflows.

## Prerequisites

- Go, using the version declared in `go.mod` (install Go yourself; it is not provided via mise here)
- Terraform, using the version declared in `.terraform-version`
- GNU Make

Optional: [mise](https://mise.jdx.dev/) — run `mise trust` in the repo root if needed, then `mise install` to install the **Trunk CLI** only from [`mise.toml`](mise.toml). Go stays outside mise. When upgrading Trunk, update both `.trunk/trunk.yaml` `cli.version` and `mise.toml` `trunk = "..."`.

Quality checks use **Trunk** (`trunk check` / `trunk fmt`) plus **`gosec`** and **`deadcode`** via `GNUmakefile` targets (`make lint`, `make build`). Git hooks: run `make setup-dev` to clear conflicting `core.hooksPath` and run `trunk git-hooks sync` (no pre-commit framework).

## Local Development

Run the common checks:

```shell
make test
go build -v ./
go generate ./...
```

Run formatting before sending changes:

```shell
make format
```

## Repository Structure

- `main.go`: provider server entry point
- `internal/provider`: provider implementation, resources, data sources, docs embedded by tests, and unit tests
- `internal/n8n`: n8n Public API HTTP client (`Client.HTTP`, `X-N8N-API-KEY`, optional rate or concurrency limits) plus `models/`, versioned `api/v1/<resource>/`, `services/`, and `controllers/` (Terraform-facing orchestrators used by resources and data sources)
- `examples`: Terraform examples used by documentation generation
- `docs`: generated or hand-maintained provider documentation
- `tools`: Go tool dependency tracking for code generation and analysis
- `mise.toml`: optional Trunk CLI pin for mise users

## Acceptance tests

Credential CRUD is available on Community Edition. Write-only secret attributes (generic `data` and typed resource secrets) require Terraform **1.11 or later**. Acceptance tests set `TF_ACC_TERRAFORM_VERSION` from `.terraform-version` (see `GNUmakefile`) so a tfenv global version older than 1.11 is not used when tests run in a temp directory.

Community Edition Docker can exercise the Public API. Credential CRUD runs on Community. Team-project CRUD still needs `feat:projectRole:admin`; those tests **skip** on Community with HTTP 403. Folder APIs need `feat:folders` (Registered Community or paid). Custom-role writes need Enterprise. For a fuller CE vs licensed matrix (Public API and n8n CLI probe notes), see [`dev/docs/n8n-ce-public-api-and-cli-limits.md`](dev/docs/n8n-ce-public-api-and-cli-limits.md).

### Docker Compose (Community n8n)

`make testacc-docker` starts pinned n8n from [`docker-compose.dev.yml`](docker-compose.dev.yml), waits for `/healthz/readiness`, then runs [`scripts/bootstrap-n8n.sh`](scripts/bootstrap-n8n.sh). That script is **harness-only**: it calls internal `/rest/owner/setup` (or `/rest/login`) and `/rest/api-keys` to mint a key. The provider never uses `/rest`.

```shell
make testacc-docker
make testacc-docker TESTARGS='-run ^TestAccN8n_'
```

The compose file uses ephemeral SQLite (no volume, no extra Postgres) so each `up` is a clean instance. `N8N_PUBLIC_API_DISABLED` is `false` and `N8N_SECURE_COOKIE` is `false` for HTTP localhost.

### Licensed instance (project and folder CRUD)

Copy `.env.template` to `.env`, set `N8N_ENDPOINT` and `N8N_API_KEY` for a licensed n8n, then:

```shell
make testacc TESTARGS='-run ^TestAccN8nProject_'
make testacc TESTARGS='-run ^TestAccN8nFolder_'
make testacc TESTARGS='-run ^TestAccN8nCredential_'
```

`make test` stays unit-only (`TF_ACC` unset). Do not point acceptance tests at production n8n.
