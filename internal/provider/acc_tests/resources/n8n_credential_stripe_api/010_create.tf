resource "n8n_credential_stripe_api" "test" {
  name              = "{{NAME}}"
  secret_key        = "placeholder"
  data_version      = 1
  delete_protection = false
}
