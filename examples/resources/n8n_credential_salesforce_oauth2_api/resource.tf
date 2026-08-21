resource "n8n_credential_salesforce_oauth2_api" "example" {
  name              = "salesforce"
  client_id         = "placeholder"
  client_secret     = "placeholder"
  environment       = "production"
  data_version      = 1
  delete_protection = true
}
