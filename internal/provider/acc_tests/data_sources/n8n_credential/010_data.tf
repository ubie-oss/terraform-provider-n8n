resource "n8n_credential" "test" {
  name              = "{{NAME}}"
  type              = "httpHeaderAuth"
  data_version      = 1
  delete_protection = false
  data = {
    name  = "X-Test"
    value = "placeholder"
  }
}

data "n8n_credential" "by_id" {
  id = n8n_credential.test.id
}

data "n8n_credentials" "all" {
  name = n8n_credential.test.name
  type = n8n_credential.test.type
}

data "n8n_credential_schema" "header" {
  type = "httpHeaderAuth"
}
