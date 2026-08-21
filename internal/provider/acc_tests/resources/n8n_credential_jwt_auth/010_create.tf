resource "n8n_credential_jwt_auth" "test" {
  name              = "{{NAME}}"
  key_type          = "passphrase"
  secret            = "placeholder"
  algorithm         = "HS256"
  data_version      = 1
  delete_protection = false
}
