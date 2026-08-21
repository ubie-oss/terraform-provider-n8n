resource "n8n_credential_hubspot_app_token" "test" {
  name              = "{{NAME}}-renamed"
  app_token         = "placeholder-2"
  data_version      = 2
  delete_protection = false
}
