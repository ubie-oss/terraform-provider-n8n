resource "n8n_credential_aws" "test" {
  name              = "{{NAME}}-renamed"
  region            = "eu-west-1"
  access_key_id     = "AKIAyyyyyyyy"
  secret_access_key = "placeholder-2"
  data_version      = 2
  delete_protection = false
}
