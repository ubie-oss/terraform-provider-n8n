resource "n8n_project" "test" {
  name              = "{{NAME}}"
  delete_protection = false
}

data "n8n_project" "by_id" {
  id = n8n_project.test.id
}

data "n8n_projects" "all" {}
