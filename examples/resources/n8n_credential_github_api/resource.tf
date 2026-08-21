resource "n8n_credential_github_api" "example" {
  name              = "github-pat"
  user              = "octocat"
  access_token      = "placeholder"
  data_version      = 1
  delete_protection = false
}
