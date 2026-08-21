resource "n8n_credential_microsoft_graph_security_oauth2_api" "example" {
  name              = "graph-security"
  client_id         = "placeholder"
  client_secret     = "placeholder"
  data_version      = 1
  delete_protection = true
}
