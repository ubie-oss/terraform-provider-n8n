resource "n8n_project" "test" {
  name              = "{{NAME}}"
  delete_protection = false
}

resource "n8n_folder" "test" {
  project_id        = n8n_project.test.id
  name              = "{{NAME}}-root"
  delete_protection = false
}
