resource "n8n_credential_github_api" "test" {
  name              = "{{NAME}}-renamed"
  user              = "octocat-2"
  access_token      = "placeholder-2"
  data_version      = 2
  delete_protection = false
}
