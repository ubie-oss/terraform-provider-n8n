resource "n8n_credential_google_cloud_storage_oauth2_api" "example" {
  name              = "gcs"
  client_id         = "placeholder"
  client_secret     = "placeholder"
  data_version      = 1
  delete_protection = true
}
