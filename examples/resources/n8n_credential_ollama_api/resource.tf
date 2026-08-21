resource "n8n_credential_ollama_api" "example" {
  name              = "ollama"
  base_url          = "http://localhost:11434"
  data_version      = 1
  delete_protection = false
}
