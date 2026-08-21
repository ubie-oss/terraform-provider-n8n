resource "n8n_credential_oauth2_api" "example" {
  name              = "generic-oauth2"
  grant_type        = "clientCredentials"
  access_token_url  = "https://example.com/oauth/token"
  client_id         = "placeholder"
  client_secret     = "placeholder"
  authentication    = "header"
  data_version      = 1
  delete_protection = true
}
