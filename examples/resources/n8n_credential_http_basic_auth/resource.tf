resource "n8n_credential_http_basic_auth" "example" {
  name              = "http-basic"
  user              = "alice"
  password          = "placeholder"
  data_version      = 1
  delete_protection = false
}
