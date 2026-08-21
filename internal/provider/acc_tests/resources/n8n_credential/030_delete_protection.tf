resource "n8n_credential" "test" {
  name              = "{{NAME}}"
  type              = "httpHeaderAuth"
  data_version      = 1
  delete_protection = true
  data = {
    name  = "X-Test"
    value = "placeholder"
  }
}
