Manages an n8n Stripe API credential (`stripeApi`) via the Public API.

`secret_key` is required (use `sk_live_` / `sk_test_`, not publishable `pk_` keys). Optional `signature_secret` (`whsec_`) verifies Stripe Trigger webhooks. Both are write-only. Terraform 1.11 or later is required.
