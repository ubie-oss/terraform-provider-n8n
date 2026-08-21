resource "n8n_tag" "test" {
  name              = "{{NAME}}"
  delete_protection = false
}

data "n8n_tag" "by_id" {
  id = n8n_tag.test.id
}

data "n8n_tags" "all" {}
