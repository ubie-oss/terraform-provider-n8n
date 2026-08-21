resource "n8n_credential_oauth2_api" "test" {
  name              = "{{NAME}}-renamed"
  grant_type        = "clientCredentials"
  access_token_url  = "https://example.com/oauth/token-v2"
  client_id         = "placeholder-2"
  client_secret     = "placeholder-2"
  authentication    = "body"
  data_version      = 2
  delete_protection = false
}
