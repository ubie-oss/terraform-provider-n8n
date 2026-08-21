Manages an n8n HubSpot Service Key credential (`hubspotAppToken`) via the Public API.

`app_token` is write-only. This is not the sunset HubSpot API key type (`hubspotApi`) and not HubSpot OAuth2.

This resource and generic `n8n_credential` must not manage the same credential id. Terraform 1.11 or later is required for write-only attributes.
