resource "n8n_credential_smtp" "test" {
  name              = "{{NAME}}"
  user              = "alerts@example.com"
  password          = "placeholder"
  host              = "smtp.example.com"
  port              = 465
  secure            = true
  data_version      = 1
  delete_protection = false
}
