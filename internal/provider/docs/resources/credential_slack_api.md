Manages an n8n Slack API credential (`slackApi`) via the Public API.

`access_token` and optional `signature_secret` are write-only. Bump `data_version` to rotate secrets. This is the bot-token type, not Slack OAuth2 (`slackOAuth2Api`).

This resource and generic `n8n_credential` must not manage the same credential id. Terraform 1.11 or later is required for write-only attributes.
