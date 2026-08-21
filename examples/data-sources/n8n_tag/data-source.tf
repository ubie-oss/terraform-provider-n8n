data "n8n_tag" "example" {
  id = n8n_tag.example.id
}

resource "n8n_tag" "example" {
  name              = "production"
  delete_protection = false
}
