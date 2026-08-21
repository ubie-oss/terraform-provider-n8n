resource "n8n_credential_github_oauth2_api" "test" {
  name              = "{{NAME}}-renamed"
  client_id         = "placeholder-2"
  client_secret     = "placeholder-2"
  data_version      = 2
  delete_protection = false
}
