data "n8n_credentials" "all" {
}

data "n8n_credentials" "headers" {
  type = "httpHeaderAuth"
}
