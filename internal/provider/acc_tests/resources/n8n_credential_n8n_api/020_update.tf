resource "n8n_credential_n8n_api" "test" {
  name              = "{{NAME}}-renamed"
  api_key           = "placeholder-2"
  base_url          = "https://example.app.n8n.cloud/api/v1-updated"
  data_version      = 2
  delete_protection = false
}
