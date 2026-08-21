resource "n8n_credential_jwt_auth" "test" {
  name              = "{{NAME}}-renamed"
  key_type          = "passphrase"
  secret            = "placeholder-2"
  algorithm         = "HS512"
  data_version      = 2
  delete_protection = false
}
