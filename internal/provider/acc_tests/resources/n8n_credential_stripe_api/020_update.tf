resource "n8n_credential_stripe_api" "test" {
  name              = "{{NAME}}-renamed"
  secret_key        = "placeholder-2"
  data_version      = 2
  delete_protection = false
}
