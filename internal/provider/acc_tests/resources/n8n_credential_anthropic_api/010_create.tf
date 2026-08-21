resource "n8n_credential_anthropic_api" "test" {
  name              = "{{NAME}}"
  api_key           = "placeholder"
  data_version      = 1
  delete_protection = false
}
