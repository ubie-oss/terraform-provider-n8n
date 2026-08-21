Manages an n8n Ollama credential (`ollamaApi`) via the Public API.

`base_url` is required (n8n UI defaults to `http://localhost:11434`). Optional write-only `api_key` is a Bearer token for authenticated proxies such as Open WebUI; omit it for a default local Ollama install.

This resource and generic `n8n_credential` must not manage the same credential id. Terraform 1.11 or later is required for write-only attributes.
