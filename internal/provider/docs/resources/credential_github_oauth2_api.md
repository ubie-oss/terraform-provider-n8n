Manages an n8n GitHub OAuth2 credential (`githubOAuth2Api`) via the Public API.

This is not GitHub PAT (`githubApi`) or GitHub App (`githubAppApi`). Optional `server` selects GitHub Enterprise. `client_secret` and optional `oauth_token_data` are write-only. `is_partial_data` defaults to true.

This resource and generic `n8n_credential` must not manage the same credential id. Terraform 1.11 or later is required for write-only attributes.
