Manages an n8n Basic Auth credential (`httpBasicAuth`) via the Public API.

`user` is stored in Terraform state. `password` is write-only and is never stored in state; bump `data_version` to rotate it. n8n GET responses omit secrets, so Terraform cannot detect drift of the password.

This resource and generic `n8n_credential` must not manage the same credential id. Terraform 1.11 or later is required for write-only attributes.
