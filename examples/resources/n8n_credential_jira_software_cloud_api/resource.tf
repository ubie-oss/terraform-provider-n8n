resource "n8n_credential_jira_software_cloud_api" "example" {
  name              = "jira-cloud"
  email             = "user@example.com"
  api_token         = "placeholder"
  domain            = "https://example.atlassian.net"
  data_version      = 1
  delete_protection = false
}
