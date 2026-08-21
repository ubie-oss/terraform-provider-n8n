Manages an n8n Header Auth credential (`httpHeaderAuth`) via the Public API.

`header_name` maps to n8n `data.name` (the HTTP header name). It is stored in Terraform state. `value` is write-only and is never stored in state; bump `data_version` to rotate it. n8n GET responses omit secrets, so Terraform cannot detect drift of the header value.

This resource and generic `n8n_credential` must not manage the same credential id. Terraform 1.11 or later is required for write-only attributes.
