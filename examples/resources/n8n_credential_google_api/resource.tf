resource "n8n_credential_google_api" "example" {
  name              = "google-sa"
  email             = "sa@example.iam.gserviceaccount.com"
  private_key       = "placeholder"
  region            = "global"
  data_version      = 1
  delete_protection = false
}
