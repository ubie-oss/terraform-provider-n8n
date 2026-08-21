resource "n8n_credential_http_basic_auth" "test" {
  name              = "{{NAME}}-renamed"
  user              = "bob"
  password          = "placeholder-2"
  data_version      = 2
  delete_protection = false
}
