resource "n8n_credential_jira_software_cloud_api" "test" {
  name              = "{{NAME}}-renamed"
  email             = "updated@example.com"
  api_token         = "placeholder-2"
  domain            = "https://updated.atlassian.net"
  data_version      = 2
  delete_protection = false
}
