resource "n8n_credential_salesforce_oauth2_api" "test" {
  name              = "{{NAME}}"
  client_id         = "placeholder"
  client_secret     = "placeholder"
  data_version      = 1
  delete_protection = false
}
