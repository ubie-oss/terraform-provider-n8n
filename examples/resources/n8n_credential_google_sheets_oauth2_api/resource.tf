resource "n8n_credential_google_sheets_oauth2_api" "example" {
  name              = "sheets"
  client_id         = "placeholder"
  client_secret     = "placeholder"
  data_version      = 1
  delete_protection = true
}
