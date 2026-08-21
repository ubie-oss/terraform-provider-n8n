Manages an n8n JWT Auth credential (`jwtAuth`) via the Public API.

Set `key_type` to `passphrase` or `pemKey`. `secret`, `private_key`, and `public_key` are write-only. Optional `algorithm` defaults in n8n to HS256 when omitted.

This resource and generic `n8n_credential` must not manage the same credential id. Terraform 1.11 or later is required for write-only attributes.
