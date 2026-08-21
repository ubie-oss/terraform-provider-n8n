Manages a generic n8n OAuth2 API credential (`oAuth2Api`) for HTTP Request and other generic OAuth2 use.

Configure `grant_type`, token/authorization URLs, client id/secret, and optional scope/authentication settings. Dynamic client registration is not exposed (n8n keeps it hidden/disabled on this type). Browser OAuth is not available over the Public API; supply `oauth_token_data` when needed. Terraform 1.11 or later is required.
