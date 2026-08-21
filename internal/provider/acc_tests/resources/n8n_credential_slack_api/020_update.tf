resource "n8n_credential_slack_api" "test" {
  name              = "{{NAME}}-renamed"
  access_token      = "placeholder-2"
  data_version      = 2
  delete_protection = false
}
