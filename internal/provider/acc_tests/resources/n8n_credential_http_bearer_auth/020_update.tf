resource "n8n_credential_http_bearer_auth" "test" {
  name              = "{{NAME}}-renamed"
  token             = "placeholder-2"
  data_version      = 2
  delete_protection = false
}
