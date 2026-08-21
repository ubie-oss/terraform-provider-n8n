Manages an n8n OpenAI credential (`openAiApi`) via the Public API.

`api_key` is write-only. Optional `organization_id`, `url`, and a custom header (`header`, `header_name`, `header_value`) map to n8n `data` keys. `header_value` is write-only.

This resource and generic `n8n_credential` must not manage the same credential id. Terraform 1.11 or later is required for write-only attributes.
