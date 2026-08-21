resource "n8n_credential_http_header_auth" "test" {
  name              = "{{NAME}}"
  header_name       = "X-Test"
  value             = "placeholder"
  data_version      = 1
  delete_protection = false
}
