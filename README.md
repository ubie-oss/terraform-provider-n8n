# Terraform Provider for n8n

Manage [n8n](https://n8n.io/) **team projects**, **folders**, **tags**, and **credentials** as code with Terraform, using the [n8n Public API](https://docs.n8n.io/connect/n8n-api).

Team-project APIs require an n8n license that includes `feat:projectRole:admin`. Folder APIs require `feat:folders`. Community self-hosted n8n returns HTTP 403 for those operations. Tag and credential CRUD are available on Community Edition. Write-only secret attributes require Terraform 1.11 or later.

Prefer a typed credential resource when one exists (`n8n_credential_http_header_auth`, `n8n_credential_slack_api`, `n8n_credential_aws`, `n8n_credential_github_api`, `n8n_credential_open_ai_api`, `n8n_credential_anthropic_api`, and others). Generic `n8n_credential` remains the escape hatch. Do not manage the same credential id with both. There is no `n8n_credential_mcp_authentication` resource: MCP Authentication is a node option that selects Header Auth, Bearer Auth, MCP OAuth2, or Multiple Headers Auth. Do not confuse `githubApi` with `githubOAuth2Api`, `googleSheetsOAuth2Api` with the Sheets Trigger type, `aws` with `awsAssumeRole`, or `openAiApi` with `azureOpenAiApi`. Azure Entra ID for Cognitive Services is a separate n8n type (`azureEntraCognitiveServicesOAuth2Api`), not an auth mode on `n8n_credential_azure_open_ai_api`.

## Example Usage

```hcl
terraform {
  required_providers {
    n8n = {
      source = "ubie-oss/n8n"
    }
  }
}

provider "n8n" {
  endpoint = "https://n8n.example.com"
  api_key  = var.api_key

  # Or export N8N_ENDPOINT and N8N_API_KEY instead.
  # max_concurrent_requests = 5
  # requests_per_second     = 20
}

resource "n8n_project" "platform" {
  name              = "platform"
  delete_protection = false
}

resource "n8n_folder" "ingest" {
  project_id        = n8n_project.platform.id
  name              = "ingest"
  delete_protection = false
}

resource "n8n_credential_http_header_auth" "header" {
  name              = "http-header"
  header_name       = "X-Test"
  value             = "placeholder"
  data_version      = 1
  delete_protection = false
}

data "n8n_project" "platform" {
  id = n8n_project.platform.id
}

data "n8n_projects" "all" {}

data "n8n_folders" "in_platform" {
  project_id = n8n_project.platform.id
}

data "n8n_credential_schema" "header" {
  type = "httpHeaderAuth"
}
```

## Authentication

Create an API key in n8n under **Settings → n8n API**. The provider sends it as the `X-N8N-API-KEY` header.

## Development

Optional: with [mise](https://mise.jdx.dev/), run `mise trust` if prompted, then `mise install` to get the pinned Trunk CLI from `mise.toml` (Go is not installed via mise).

```shell
make test
go build -v ./
go generate ./...
```

Use `internal/provider` for Terraform wiring and `internal/n8n` for the Public API client (`models`, versioned `api/v1/<resource>`, `services`, `controllers`).
