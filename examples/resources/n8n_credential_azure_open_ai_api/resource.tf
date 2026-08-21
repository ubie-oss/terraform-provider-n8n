resource "n8n_credential_azure_open_ai_api" "example" {
  name              = "azure-openai"
  api_key           = "placeholder"
  resource_name     = "my-aoai-resource"
  api_version       = "2025-03-01-preview"
  data_version      = 1
  delete_protection = false
}
