resource "n8n_project" "test" {
  name              = "{{NAME}}"
  delete_protection = false
}

resource "n8n_folder" "test" {
  project_id        = n8n_project.test.id
  name              = "{{NAME}}-root"
  delete_protection = false
}

resource "n8n_folder" "child" {
  project_id        = n8n_project.test.id
  parent_folder_id  = n8n_folder.test.folder_id
  name              = "{{NAME}}-child"
  delete_protection = false
}
