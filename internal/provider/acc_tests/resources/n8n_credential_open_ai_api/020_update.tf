resource "n8n_credential_open_ai_api" "test" {
  name              = "{{NAME}}-renamed"
  api_key           = "placeholder-2"
  data_version      = 2
  delete_protection = false
}
