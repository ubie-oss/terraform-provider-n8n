resource "n8n_credential_anthropic_api" "test" {
  name              = "{{NAME}}-renamed"
  api_key           = "placeholder-2"
  url               = "https://api.anthropic.com"
  data_version      = 2
  delete_protection = false
}
