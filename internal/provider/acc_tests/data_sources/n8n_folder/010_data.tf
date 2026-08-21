resource "n8n_project" "test" {
  name              = "{{NAME}}"
  delete_protection = false
}

resource "n8n_folder" "test" {
  project_id        = n8n_project.test.id
  name              = "{{NAME}}-root"
  delete_protection = false
}

data "n8n_folder" "by_id" {
  project_id = n8n_folder.test.project_id
  folder_id  = n8n_folder.test.folder_id
}

data "n8n_folders" "all" {
  project_id = n8n_project.test.id
}
