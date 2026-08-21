resource "n8n_credential_smtp" "test" {
  name              = "{{NAME}}-renamed"
  user              = "ops@example.com"
  password          = "placeholder-2"
  host              = "smtp-updated.example.com"
  port              = 587
  secure            = false
  data_version      = 2
  delete_protection = false
}
