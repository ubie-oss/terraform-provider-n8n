Manages an n8n GitHub API credential (`githubApi`) via the Public API.

Uses a personal access token (`access_token`), optional `user`, and optional `server` for GitHub Enterprise. This is not GitHub OAuth2 (`githubOAuth2Api`) or GitHub App (`githubAppApi`).

Write-only secrets require Terraform 1.11+. Do not manage the same credential id with generic `n8n_credential`.
