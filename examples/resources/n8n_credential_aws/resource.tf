resource "n8n_credential_aws" "example" {
  name              = "aws-iam"
  region            = "us-east-1"
  access_key_id     = "AKIAxxxxxxxx"
  secret_access_key = "placeholder"
  data_version      = 1
  delete_protection = false
}
