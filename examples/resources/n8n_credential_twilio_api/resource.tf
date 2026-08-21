resource "n8n_credential_twilio_api" "example" {
  name              = "twilio"
  auth_type         = "authToken"
  account_sid       = "ACxxxxxxxx"
  auth_token        = "placeholder"
  data_version      = 1
  delete_protection = false
}
