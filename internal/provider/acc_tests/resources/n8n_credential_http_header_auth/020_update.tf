resource "n8n_credential_http_header_auth" "test" {
  name              = "{{NAME}}-renamed"
  header_name       = "X-Test-Updated"
  value             = "placeholder-2"
  data_version      = 2
  delete_protection = false
}
