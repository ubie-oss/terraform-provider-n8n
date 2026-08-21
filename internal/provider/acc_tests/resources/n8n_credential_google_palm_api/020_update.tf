resource "n8n_credential_google_palm_api" "test" {
  name              = "{{NAME}}-renamed"
  host              = "https://generativelanguage.googleapis.com"
  api_key           = "placeholder-2"
  data_version      = 2
  delete_protection = false
}
