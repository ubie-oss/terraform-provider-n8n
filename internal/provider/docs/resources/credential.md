Manages an n8n credential of any type via the Public API.

Prefer a typed resource when one exists for the credential type (for example `n8n_credential_http_header_auth` or `n8n_credential_slack_api`). Use this generic resource as the escape hatch for types that are not specialized. Do not manage the same credential id with both a typed resource and `n8n_credential`.

`type` is the n8n credential type name for this instance (for example `httpHeaderAuth` or `githubApi`). There is no Public API catalog of types; use the `n8n_credential_schema` data source to inspect the payload keys for a type. Changing `type` replaces the resource.

`data` is write-only and is never stored in Terraform state. n8n GET and list responses omit secrets. On create, the provider always sends `data`. On update, it sends `data` only when `data_version` changes. Pair `data` with a secrets manager or Terraform ephemeral values. Terraform 1.11 or later is required for write-only attributes.

OAuth credential types still need tokens you obtained outside Terraform; the Public API cannot complete a browser OAuth flow.

`project_id` is optional. If omitted, n8n creates the credential in the API key owner's personal project. Changing `project_id` after it was already set in state transfers the credential (`credential:move`). Setting `project_id` for the first time after import (when state had no project) adopts the value without transferring—GET does not return ownership, so Terraform cannot detect the real project. Clearing `project_id` after create is not supported. Transfer needs a real team project id (the path alias `personal` is not accepted).

`is_global` is not part of create; when set in config the provider applies it with a follow-up update after create. Community n8n returns HTTP 403 (`You are not licensed for sharing credentials`) if you set it to true.

`delete_protection` is required. When set to `true`, Terraform will not destroy the resource. This flag is Terraform-only; n8n has no matching API field. Imported resources default to `delete_protection = true`. After import you must set `data` and `data_version` in configuration; import cannot recover secrets. Use the same `data_version` as state (import sets `1`) until you intend to rotate.

Managed credentials (`is_managed`) cannot be updated or deleted via the API. Prefer not to import managed credentials.
