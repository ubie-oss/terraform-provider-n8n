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
