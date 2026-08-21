resource "n8n_credential_mcp_oauth2_api" "example" {
  name                            = "mcp-oauth"
  use_dynamic_client_registration = true
  server_url                      = "https://mcp.example.com"
  data_version                    = 1
  delete_protection               = true
}
