resource "n8n_credential_gmail_oauth2" "test" {
  name              = "{{NAME}}-renamed"
  client_id         = "placeholder-2"
  client_secret     = "placeholder-2"
  data_version      = 2
  delete_protection = false
}
