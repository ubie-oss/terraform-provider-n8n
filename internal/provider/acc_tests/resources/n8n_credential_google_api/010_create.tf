resource "n8n_credential_google_api" "test" {
  name              = "{{NAME}}"
  email             = "sa@example.iam.gserviceaccount.com"
  private_key       = "placeholder"
  data_version      = 1
  delete_protection = false
}
