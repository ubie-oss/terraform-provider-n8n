resource "n8n_credential_google_sheets_trigger_oauth2_api" "example" {
  name              = "sheets-trigger"
  client_id         = "placeholder"
  client_secret     = "placeholder"
  data_version      = 1
  delete_protection = true
}
