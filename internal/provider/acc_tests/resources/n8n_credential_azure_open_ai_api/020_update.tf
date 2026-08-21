resource "n8n_credential_azure_open_ai_api" "test" {
  name              = "{{NAME}}-renamed"
  api_key           = "placeholder-2"
  resource_name     = "tf-acc-aoai-2"
  api_version       = "2024-10-21"
  data_version      = 2
  delete_protection = false
}
