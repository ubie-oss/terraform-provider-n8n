provider "n8n" {
  endpoint = "https://n8n.example.com"
  api_key  = "your-api-key"

  # Or set N8N_ENDPOINT and N8N_API_KEY environment variables instead.
  # Optional: tune HTTP client (defaults: 10 concurrent, 10 RPS)
  # max_concurrent_requests = 5
  # requests_per_second     = 20.5
}
