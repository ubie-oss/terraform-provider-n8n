Manages an n8n Google Sheets OAuth2 credential (`googleSheetsOAuth2Api`) via the Public API.

Same client fields as other Google OAuth typed resources. This is the Sheets node credential, not Sheets Trigger (`googleSheetsTriggerOAuth2Api`). Browser OAuth is not available over the Public API; supply `oauth_token_data` or complete Connect in the UI and use `is_partial_data`. Terraform 1.11 or later is required.
