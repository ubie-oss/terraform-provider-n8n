Manages an n8n Salesforce OAuth2 credential (`salesforceOAuth2Api`) via the Public API.

Optional `environment` is `production` or `sandbox`. This is not Salesforce JWT (`salesforceJwtApi`). `is_partial_data` defaults to true. The Public API cannot complete a browser OAuth flow.

This resource and generic `n8n_credential` must not manage the same credential id. Terraform 1.11 or later is required for write-only attributes.
