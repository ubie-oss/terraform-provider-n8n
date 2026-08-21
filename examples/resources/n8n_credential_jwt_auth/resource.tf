resource "n8n_credential_jwt_auth" "example" {
  name              = "jwt"
  key_type          = "passphrase"
  secret            = "placeholder"
  algorithm         = "HS256"
  data_version      = 1
  delete_protection = false
}
