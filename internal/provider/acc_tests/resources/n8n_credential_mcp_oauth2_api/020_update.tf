resource "n8n_credential_mcp_oauth2_api" "test" {
  name                            = "{{NAME}}-renamed"
  use_dynamic_client_registration = true
  server_url                      = "https://mcp-updated.example.com"
  data_version                    = 2
  delete_protection               = false
}
