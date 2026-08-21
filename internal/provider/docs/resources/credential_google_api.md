Manages an n8n Google Service Account credential (`googleApi`) via the Public API.

`private_key` is write-only. Terraform `impersonate` is sent as n8n `inpersonate` (n8n's spelling). Optional `region`, `delegated_email`, `http_node`, and `scopes` map to the matching n8n data keys.

This resource and generic `n8n_credential` must not manage the same credential id. Terraform 1.11 or later is required for write-only attributes.
