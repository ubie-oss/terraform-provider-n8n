Manages an n8n Azure OpenAI API-key credential (`azureOpenAiApi`) via the Public API.

This is the API-key auth path only. Azure Entra ID (OAuth2) uses a separate n8n type (`azureEntraCognitiveServicesOAuth2Api`) and is not managed by this resource.

`api_key` is write-only. Required state fields are `resource_name` and `api_version`. Optional `endpoint` overrides the regional Cognitive Services host. Deployment names belong on Azure OpenAI nodes as the model name, not on this credential.

This resource and generic `n8n_credential` must not manage the same credential id. Terraform 1.11 or later is required for write-only attributes.
