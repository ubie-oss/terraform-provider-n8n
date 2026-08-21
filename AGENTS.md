# AGENTS.md

Shared guidance for agents working in this repository.

## Overview

Terraform provider for [n8n](https://n8n.io/), written in Go with the HashiCorp Terraform Plugin Framework. It talks to the n8n Public API and currently manages **team projects** (`n8n_project`), **project folders** (`n8n_folder` resource, `n8n_folder` / `n8n_folders` data sources), and **credentials** (generic `n8n_credential` plus typed resources such as `n8n_credential_http_header_auth`, `n8n_credential_slack_api`, `n8n_credential_aws`, `n8n_credential_github_api`, `n8n_credential_open_ai_api`, and `n8n_credential_anthropic_api`; data sources `n8n_credential` / `n8n_credentials` / `n8n_credential_schema`). It does not manage workflows. MCP Authentication is a node option, not a credential type. Nearby type gotchas: `githubApi` vs `githubOAuth2Api`, `googleSheetsOAuth2Api` vs `googleSheetsTriggerOAuth2Api`, `aws` vs `awsAssumeRole`, `oAuth2Api` (lowercase leading `o`), `openAiApi` vs `azureOpenAiApi`, and Azure Entra Cognitive Services OAuth2 as a separate type from `azureOpenAiApi`. Ollama credentials are primarily identified by `base_url` (optional proxy `api_key`).

## mise (optional)

If you use [mise](https://mise.jdx.dev/), run `mise trust` in the repo root on first use if mise refuses to load the config, then `mise install`. [`mise.toml`](mise.toml) installs **only the Trunk CLI** (pinned to match `.trunk/trunk.yaml` `cli.version`). **Go is not installed via mise**; use your normal Go install and the `go` / `toolchain` lines in `go.mod`. When you upgrade Trunk, bump `cli.version` in `.trunk/trunk.yaml` and the `trunk` version in `mise.toml` together.

## Key Commands

- Unit tests: `make test`
- Build: `go build -v ./`
- Generate docs: `go generate ./...`
- Format: `make format`
- Lint: `make lint`

## Project Structure

- `main.go`: provider server entry point (`registry.terraform.io/ubie-oss/n8n`)
- `internal/provider`: provider schema, configuration, resources, data sources, embedded docs, and tests
- `internal/n8n`: n8n Public API HTTP client (`X-N8N-API-KEY`, rate and concurrency limits), plus `models/`, versioned `api/v1/<resource>/` (one HTTP op per file), `services/`, and `controllers/` (Terraform-facing orchestrators; resources and data sources call these, not services directly)
- `examples`: Terraform examples used by docs generation
- `docs`: provider documentation
- `tools`: Go tool dependencies
- `mise.toml`: optional Trunk CLI version for mise users (no Go via mise)

## Development Notes

- Prefer real, deterministic tests over mocks. Do not add a stand-in n8n HTTP server for acceptance tests.
- Acceptance tests need `TF_ACC=1`, `N8N_ENDPOINT`, and `N8N_API_KEY`.
- `make testacc-docker` starts Community n8n from `docker-compose.dev.yml` and mints a key with `scripts/bootstrap-n8n.sh` (`/rest` harness only). Licensed APIs skip on 403 (`feat:projectRole:admin` for team projects, `feat:folders` for folders). Credential CRUD runs on Community (write-only `data` needs Terraform 1.11+; `GNUmakefile` sets `TF_ACC_TERRAFORM_VERSION` from `.terraform-version`). Use `make testacc` against a licensed instance for project/folder CRUD tests. CE capability matrix (API + CLI): [`dev/docs/n8n-ce-public-api-and-cli-limits.md`](dev/docs/n8n-ce-public-api-and-cli-limits.md).
- Git hooks for this repo are **Trunk-only** (`make setup-dev` runs `trunk git-hooks sync`); there is no separate pre-commit install step.
