resource "n8n_credential_twilio_api" "test" {
  name              = "{{NAME}}-renamed"
  auth_type         = "authToken"
  account_sid       = "ACyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyy"
  auth_token        = "placeholder-2"
  data_version      = 2
  delete_protection = false
}
