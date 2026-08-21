resource "n8n_credential_ollama_api" "test" {
  name              = "{{NAME}}-renamed"
  base_url          = "http://ollama.example.com:11434"
  api_key           = "placeholder-proxy"
  data_version      = 2
  delete_protection = false
}
