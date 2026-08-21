Manages an n8n MCP OAuth2 API credential (`mcpOAuth2Api`) via the Public API.

This is not MCP Authentication, which is a node parameter that selects Header Auth, Bearer Auth, MCP OAuth2, or Multiple Headers Auth. When `use_dynamic_client_registration` is true, set `server_url`. When false, set client id/secret and token URLs. `is_partial_data` defaults to true. The Public API cannot complete a browser OAuth flow.

This resource and generic `n8n_credential` must not manage the same credential id. Terraform 1.11 or later is required for write-only attributes.
