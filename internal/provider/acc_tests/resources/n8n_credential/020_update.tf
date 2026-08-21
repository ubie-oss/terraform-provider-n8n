resource "n8n_credential" "test" {
  name              = "{{NAME}}-renamed"
  type              = "httpHeaderAuth"
  data_version      = 2
  delete_protection = false
  data = {
    name  = "X-Test"
    value = "placeholder-rotated"
  }
}
