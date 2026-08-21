resource "n8n_credential_http_bearer_auth" "example" {
  name              = "http-bearer"
  token             = "placeholder"
  data_version      = 1
  delete_protection = false
}
