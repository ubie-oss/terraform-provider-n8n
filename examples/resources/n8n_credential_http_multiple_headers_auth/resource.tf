resource "n8n_credential_http_multiple_headers_auth" "example" {
  name = "multi-headers"
  headers = {
    values = [
      {
        name  = "X-Api-Key"
        value = "placeholder"
      }
    ]
  }
  data_version      = 1
  delete_protection = false
}
