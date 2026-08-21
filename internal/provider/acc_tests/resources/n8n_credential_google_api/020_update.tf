resource "n8n_credential_google_api" "test" {
  name              = "{{NAME}}-renamed"
  email             = "sa-updated@example.iam.gserviceaccount.com"
  private_key       = "placeholder-2"
  data_version      = 2
  delete_protection = false
}
