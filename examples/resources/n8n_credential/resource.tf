resource "n8n_credential" "example" {
  name              = "http-header"
  type              = "httpHeaderAuth"
  data_version      = 1
  delete_protection = false
  data = {
    name  = "X-Test"
    value = "placeholder"
  }
}
