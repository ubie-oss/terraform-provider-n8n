Manages an n8n AWS (IAM) credential (`aws`) via the Public API.

Provide `region`, `access_key_id`, and write-only `secret_access_key`. Optional STS `temporary_credentials` / `session_token` and VPC `custom_endpoints` with per-service endpoint URLs. This is not AWS Assume Role (`awsAssumeRole`). Terraform 1.11 or later is required.
