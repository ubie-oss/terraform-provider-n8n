Manages an n8n Google Drive OAuth2 credential (`googleDriveOAuth2Api`) via the Public API.

`client_secret` and optional `oauth_token_data` are write-only. `is_partial_data` defaults to true. The Public API cannot complete a browser OAuth flow.

This resource and generic `n8n_credential` must not manage the same credential id. Terraform 1.11 or later is required for write-only attributes.
