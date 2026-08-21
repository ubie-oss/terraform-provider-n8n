Manages an n8n Notion API credential (`notionApi`) via the Public API.

`api_key` is the Notion Internal Integration Secret and is write-only. Bump `data_version` to rotate it. This is not Notion OAuth2 (`notionOAuth2Api`).

This resource and generic `n8n_credential` must not manage the same credential id. Terraform 1.11 or later is required for write-only attributes.
