Manages an n8n Bearer Auth credential (`httpBearerAuth`) via the Public API.

`token` is write-only and is never stored in state; bump `data_version` to rotate it. n8n GET responses omit secrets, so Terraform cannot detect drift of the token.

This resource and generic `n8n_credential` must not manage the same credential id. Terraform 1.11 or later is required for write-only attributes.
