Manages an n8n Jira Software Cloud API credential (`jiraSoftwareCloudApi`) via the Public API.

`email` and `domain` are stored in Terraform state. `api_token` is write-only. This is the email-and-token type, not Jira Cloud OAuth2.

This resource and generic `n8n_credential` must not manage the same credential id. Terraform 1.11 or later is required for write-only attributes.
