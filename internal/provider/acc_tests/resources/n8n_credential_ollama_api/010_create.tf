resource "n8n_credential_ollama_api" "test" {
  name              = "{{NAME}}"
  base_url          = "http://127.0.0.1:11434"
  data_version      = 1
  delete_protection = false
}
