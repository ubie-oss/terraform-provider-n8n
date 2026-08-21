resource "n8n_credential_http_basic_auth" "test" {
  name              = "{{NAME}}"
  user              = "alice"
  password          = "placeholder"
  data_version      = 1
  delete_protection = false
}
