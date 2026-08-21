resource "n8n_credential_github_oauth2_api" "example" {
  name              = "github-oauth"
  client_id         = "placeholder"
  client_secret     = "placeholder"
  data_version      = 1
  delete_protection = true
}
