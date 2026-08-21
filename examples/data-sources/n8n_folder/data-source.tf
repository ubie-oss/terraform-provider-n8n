data "n8n_folder" "example" {
  project_id = n8n_folder.example.project_id
  folder_id  = n8n_folder.example.folder_id
}

resource "n8n_project" "example" {
  name              = "platform"
  delete_protection = false
}

resource "n8n_folder" "example" {
  project_id        = n8n_project.example.id
  name              = "ingest"
  delete_protection = false
}
