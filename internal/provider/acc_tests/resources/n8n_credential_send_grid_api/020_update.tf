resource "n8n_credential_send_grid_api" "test" {
  name              = "{{NAME}}-renamed"
  api_key           = "placeholder-2"
  data_version      = 2
  delete_protection = false
}
