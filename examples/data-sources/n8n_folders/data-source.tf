data "n8n_folders" "all" {
  project_id = n8n_project.example.id
}

resource "n8n_project" "example" {
  name              = "platform"
  delete_protection = false
}
