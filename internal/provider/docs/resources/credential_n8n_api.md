Manages an n8n API credential (`n8nApi`) via the Public API.

`api_key` is write-only. `base_url` is the target instance Public API URL (for example `https://example.app.n8n.cloud/api/v1`) and is stored in Terraform state.

This resource and generic `n8n_credential` must not manage the same credential id. Terraform 1.11 or later is required for write-only attributes.
