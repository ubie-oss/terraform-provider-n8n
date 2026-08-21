Manages an n8n SerpAPI credential (`serpApi`) via the Public API.

`api_key` is write-only. n8n marks this credential type as deprecated in-source; it remains available on current instances.

This resource and generic `n8n_credential` must not manage the same credential id. Terraform 1.11 or later is required for write-only attributes.
