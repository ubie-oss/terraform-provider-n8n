resource "n8n_credential_slack_api" "example" {
  name              = "slack-bot"
  access_token      = "placeholder"
  data_version      = 1
  delete_protection = false
}
