resource "n8n_credential_mcp_oauth2_api" "test" {
  name                            = "{{NAME}}"
  use_dynamic_client_registration = true
  server_url                      = "https://mcp.example.com"
  data_version                    = 1
  delete_protection               = false
}
