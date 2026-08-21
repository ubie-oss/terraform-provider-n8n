resource "n8n_credential_http_multiple_headers_auth" "test" {
  name = "{{NAME}}-renamed"
  headers = {
    values = [
      {
        name  = "X-Api-Key"
        value = "placeholder-2"
      },
      {
        name  = "X-Extra"
        value = "extra"
      }
    ]
  }
  data_version      = 2
  delete_protection = false
}
