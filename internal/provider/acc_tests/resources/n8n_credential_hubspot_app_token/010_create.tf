resource "n8n_credential_hubspot_app_token" "test" {
  name              = "{{NAME}}"
  app_token         = "placeholder"
  data_version      = 1
  delete_protection = false
}
