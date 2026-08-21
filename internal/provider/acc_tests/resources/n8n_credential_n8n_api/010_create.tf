resource "n8n_credential_n8n_api" "test" {
  name              = "{{NAME}}"
  api_key           = "placeholder"
  base_url          = "https://example.app.n8n.cloud/api/v1"
  data_version      = 1
  delete_protection = false
}
