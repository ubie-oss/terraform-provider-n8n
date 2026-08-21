Manages an n8n Gmail OAuth2 credential (`gmailOAuth2`) via the Public API.

The n8n type name has no `Api` suffix. `client_secret` and optional `oauth_token_data` are write-only. `is_partial_data` defaults to true so omitted tokens are not wiped on update. The Public API cannot complete a browser OAuth flow.

This resource and generic `n8n_credential` must not manage the same credential id. Terraform 1.11 or later is required for write-only attributes.
