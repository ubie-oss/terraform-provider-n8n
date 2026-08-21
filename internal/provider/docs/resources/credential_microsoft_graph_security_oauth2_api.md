Manages an n8n Microsoft Graph Security OAuth2 credential (`microsoftGraphSecurityOAuth2Api`) via the Public API.

Supports `client_credential_type` `clientSecret` or `certificate`. Certificate mode uses write-only `private_key` and `certificate`. `is_partial_data` defaults to true. The Public API cannot complete a browser OAuth flow.

This resource and generic `n8n_credential` must not manage the same credential id. Terraform 1.11 or later is required for write-only attributes.
