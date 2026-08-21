resource "n8n_credential_twilio_api" "test" {
  name              = "{{NAME}}"
  auth_type         = "authToken"
  account_sid       = "ACxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"
  auth_token        = "placeholder"
  data_version      = 1
  delete_protection = false
}
