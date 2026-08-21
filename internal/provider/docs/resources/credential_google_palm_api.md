Manages an n8n Google Gemini (PaLM) API credential (`googlePalmApi`) via the Public API.

The n8n type name is `googlePalmApi`, not `googleGeminiApi`. `api_key` is write-only. Optional `host` defaults in n8n to `https://generativelanguage.googleapis.com` when omitted.

This resource and generic `n8n_credential` must not manage the same credential id. Terraform 1.11 or later is required for write-only attributes.
