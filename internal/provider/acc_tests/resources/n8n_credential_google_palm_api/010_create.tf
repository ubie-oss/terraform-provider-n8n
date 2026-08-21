resource "n8n_credential_google_palm_api" "test" {
  name              = "{{NAME}}"
  host              = "https://generativelanguage.googleapis.com"
  api_key           = "placeholder"
  data_version      = 1
  delete_protection = false
}
