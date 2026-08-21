Manages an n8n Anthropic credential (`anthropicApi`) via the Public API.

`api_key` is write-only. Optional `url` (n8n defaults to `https://api.anthropic.com`) and a custom header (`header`, `header_name`, `header_value`) map to n8n `data` keys. `header_value` is write-only.

This resource and generic `n8n_credential` must not manage the same credential id. Terraform 1.11 or later is required for write-only attributes.
