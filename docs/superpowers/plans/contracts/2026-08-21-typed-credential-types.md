# Typed credential catalog — n8n@2.35.5

Source clone: `/Users/yu/tmp/n8n-io-n8n` at tag `n8n@2.35.5` (`602f737`).
Credential classes: `packages/nodes-base/credentials/*.credentials.ts` and
`packages/@n8n/nodes-langchain/credentials/*.credentials.ts`.
Public API still one object: `POST/PATCH /api/v1/credentials` with `type` + `data`.
Schema/validation: `GET /credentials/schema/{type}` and create/update middleware
use `CredentialsHelper.getCredentialsProperties(type)` then **drop `type: hidden`**.
`toJsonSchema` sets `additionalProperties: false`. Hidden defaults (grant URLs,
scopes) are applied by n8n when decrypting, not by sending those keys.

GET never returns `data`. Browser OAuth is not a Public API operation. Sending
`oauthTokenData` is the only API way to attach tokens. `PATCH` without
`isPartialData: true` **replaces** the whole blob and can wipe tokens.

Do not treat `/rest`, MCP, or CLI as the Terraform contract.

## Not a credential type

**MCP Authentication** is a **node parameter** on MCP Client / MCP Client Tool
(`authentication`: `none` | `bearerAuth` | `headerAuth` | `mcpOAuth2Api` |
`multipleHeadersAuth`). It selects one of:

| Node option | Credential `type` | Display name |
| --- | --- | --- |
| `headerAuth` | `httpHeaderAuth` | Header Auth |
| `bearerAuth` | `httpBearerAuth` | Bearer Auth |
| `mcpOAuth2Api` | `mcpOAuth2Api` | MCP OAuth2 API |
| `multipleHeadersAuth` | `httpMultipleHeadersAuth` | Multiple Headers Auth |

There is no `mcpAuthentication` / `mcpAuth` ICredentialType. Do not add
`n8n_credential_mcp_authentication`.

## Shared Terraform envelope (every typed resource)

Same as generic `n8n_credential`: `name`, write-only secrets, required
`data_version`, optional `is_partial_data` (default false), optional
`project_id`, optional `is_resolvable` / `is_global`, required
`delete_protection` (import default true). Hard-code `type`. Mutual exclusion:
one Terraform type owns a given credential id (generic **or** typed).

Skip `notice` fields in Terraform (schema type `notice`; UI-only).

## OAuth2 inherited writable keys

Children of `oAuth2Api` inherit (after hidden filter) roughly:

- `clientId`, `clientSecret` (required when not using DCR)
- `grantType` (`authorizationCode` | `clientCredentials` | `pkce`) unless the child hides it
- `serverUrl` when `useDynamicClientRegistration` is true
- `ignoreSSLIssues`, `tokenExpiredStatusCode`, `jweEnabled`, `jwksUri`, `inlineJwks`
- `sendAdditionalBodyProperties`, `additionalBodyProperties` (client-credentials + body auth)
- `oauthTokenData` (`json`) — injected on the `oAuth2Api` parent and merged into children

Google family (`extends: googleOAuth2Api` → `oAuth2Api`) hide grant/URLs/auth
method. Extra user fields: `customScopes` (bool), `enabledScopes` (string).
Microsoft Graph Security also exposes `clientCredentialType`, `authUrl`,
`accessTokenUrl`, `privateKey`, `certificate`, `graphApiBaseUrl`.

## Catalog

Resource names: `n8n_credential_<snake_case(type)>`.

| Display name (user list) | `type` (POST) | Suggested Terraform resource | Kind | Source |
| --- | --- | --- | --- | --- |
| Slack API | `slackApi` | `n8n_credential_slack_api` | token | `SlackApi.credentials.ts` |
| MCP OAuth2 API | `mcpOAuth2Api` | `n8n_credential_mcp_oauth2_api` | OAuth2 + DCR | `@n8n/nodes-langchain/.../McpOAuth2Api.credentials.ts` |
| Gmail OAuth2 API | **`gmailOAuth2`** (no `Api` suffix) | `n8n_credential_gmail_oauth2` | Google OAuth2 | `GmailOAuth2Api.credentials.ts` |
| Google Sheets Trigger OAuth2 API | `googleSheetsTriggerOAuth2Api` | `n8n_credential_google_sheets_trigger_oauth2_api` | Google OAuth2 | `GoogleSheetsTriggerOAuth2Api.credentials.ts` |
| Google Calendar OAuth2 API | `googleCalendarOAuth2Api` | `n8n_credential_google_calendar_oauth2_api` | Google OAuth2 | `GoogleCalendarOAuth2Api.credentials.ts` |
| Google Cloud Storage OAuth2 API | `googleCloudStorageOAuth2Api` | `n8n_credential_google_cloud_storage_oauth2_api` | Google OAuth2 | `GoogleCloudStorageOAuth2Api.credentials.ts` |
| GitHub OAuth2 API | `githubOAuth2Api` | `n8n_credential_github_oauth2_api` | OAuth2 | `GithubOAuth2Api.credentials.ts` |
| Header Auth | `httpHeaderAuth` | `n8n_credential_http_header_auth` | generic HTTP | `HttpHeaderAuth.credentials.ts` |
| Notion API | `notionApi` | `n8n_credential_notion_api` | token | `NotionApi.credentials.ts` |
| Google Service Account API | `googleApi` | `n8n_credential_google_api` | service account | `GoogleApi.credentials.ts` |
| JWT Auth | `jwtAuth` | `n8n_credential_jwt_auth` | JWT | `JwtAuth.credentials.ts` |
| Google Drive OAuth2 API | `googleDriveOAuth2Api` | `n8n_credential_google_drive_oauth2_api` | Google OAuth2 | `GoogleDriveOAuth2Api.credentials.ts` |
| Salesforce OAuth2 API | `salesforceOAuth2Api` | `n8n_credential_salesforce_oauth2_api` | OAuth2 PKCE | `SalesforceOAuth2Api.credentials.ts` |
| Basic Auth | `httpBasicAuth` | `n8n_credential_http_basic_auth` | generic HTTP | `HttpBasicAuth.credentials.ts` |
| n8n API | `n8nApi` | `n8n_credential_n8n_api` | token | `N8nApi.credentials.ts` |
| Jira SW Cloud API | `jiraSoftwareCloudApi` | `n8n_credential_jira_software_cloud_api` | basic+token | `JiraSoftwareCloudApi.credentials.ts` |
| X OAuth2 API | **`twitterOAuth2Api`** | `n8n_credential_twitter_oauth2_api` | OAuth2 PKCE | `TwitterOAuth2Api.credentials.ts` |
| HubSpot Service Key | **`hubspotAppToken`** | `n8n_credential_hubspot_app_token` | token | `HubspotAppToken.credentials.ts` |
| SerpAPI | `serpApi` | `n8n_credential_serp_api` | token (deprecated node) | `@n8n/nodes-langchain/.../SerpApi.credentials.ts` |
| Google Gemini Api | **`googlePalmApi`** | `n8n_credential_google_palm_api` | token | `@n8n/nodes-langchain/.../GooglePalmApi.credentials.ts` |
| Microsoft Graph Security OAuth2 API | `microsoftGraphSecurityOAuth2Api` | `n8n_credential_microsoft_graph_security_oauth2_api` | Microsoft OAuth2 | `MicrosoftGraphSecurityOAuth2Api.credentials.ts` |
| MCP Authentication | — | **do not implement** | node option | `nodes/mcp/shared/descriptions.ts` |
| Bearer Auth | `httpBearerAuth` | `n8n_credential_http_bearer_auth` | generic HTTP | `HttpBearerAuth.credentials.ts` |
| OpenAI | `openAiApi` | `n8n_credential_open_ai_api` | token | `OpenAiApi.credentials.ts` |

Nearby types that are **not** the listed display names (do not substitute):

- Slack OAuth2 = `slackOAuth2Api` (user asked Slack **API** = bot token)
- Google Sheets (non-trigger) = `googleSheetsOAuth2Api`
- GitHub PAT = `githubApi`; GitHub App = `githubAppApi`
- HubSpot API key (sunset) = `hubspotApi`; HubSpot OAuth = `hubspotOAuth2Api`
- Notion OAuth2 = `notionOAuth2Api`
- Jira Cloud OAuth2 = `jiraSoftwareCloudOAuth2Api`
- Salesforce JWT = `salesforceJwtApi`

## Type-specific `data` keys (non-notice)

Secrets marked `password` should be Terraform write-only.

### `slackApi`

- `accessToken` (string, required, password)
- `signatureSecret` (string, optional, password)

### `mcpOAuth2Api` (`extends: oAuth2Api`)

- `useDynamicClientRegistration` (bool, default true)
- `resourceUrl` (string, optional)
- plus OAuth2 inherited: when DCR true → `serverUrl` required; when DCR false → `grantType`, `authUrl`/`accessTokenUrl`, `clientId`, `clientSecret`, `scope`, `oauthTokenData`

### Google OAuth2 family

Shared Terraform attrs: `client_id`, `client_secret`, `custom_scopes`, `enabled_scopes`, `ignore_ssl_issues`, write-only `oauth_token_data` (object/json).

Default scopes (space-joined hidden `scope` unless `customScopes`):

- `gmailOAuth2`: gmail.labels, gmail addons compose/action, `https://mail.google.com/`, gmail.modify, gmail.compose
- `googleSheetsTriggerOAuth2Api`: drive, drive.file, spreadsheets, drive.metadata
- `googleCalendarOAuth2Api`: calendar, calendar.events
- `googleCloudStorageOAuth2Api`: cloud-platform, cloud-platform.read-only, devstorage.full_control / read_only / read_write
- `googleDriveOAuth2Api`: drive, drive.appdata, drive.photos.readonly

### `githubOAuth2Api` (`extends: oAuth2Api`)

- `server` (string, default `https://api.github.com`) — GitHub Enterprise
- `clientId`, `clientSecret`, `oauthTokenData`, `ignoreSSLIssues`
- grant/URLs/scopes hidden (authorizationCode; broad default repo/admin/workflow scopes)

### `httpHeaderAuth`

- `name`, `value` (password). Schema `required: []` live. Extra live keys: `useCustomAuth` (notice), `allowedHttpRequestDomains`, `allowedDomains`.

### `notionApi`

- `apiKey` (password) — Internal Integration Secret

### `googleApi` (service account)

- `region` (options, default `global`)
- `email` (required) — service account email
- `privateKey` (required, password)
- `inpersonate` (bool, **typo in n8n**, not `impersonate`)
- `delegatedEmail` (when `inpersonate`)
- `httpNode` (bool), `scopes` (when `httpNode`)

### `jwtAuth`

- `keyType`: `passphrase` | `pemKey`
- `secret` (password, passphrase)
- `privateKey`, `publicKey` (password, pemKey)
- `algorithm`: HS256/384/512, RS*, ES*, PS*, `none` (default HS256)

### `salesforceOAuth2Api`

- `environment`: `production` | `sandbox`
- `clientId`, `clientSecret`, `oauthTokenData`
- hidden grant `pkce`, scope `full refresh_token`

### `httpBasicAuth`

- `user`, `password` (password). Both `resolvableField`.

### `n8nApi`

- `apiKey` (password), `baseUrl` (e.g. `https://<name>.app.n8n.cloud/api/v1`)

### `jiraSoftwareCloudApi`

- `email`, `apiToken` (password), `domain` (`https://example.atlassian.net`)

### `twitterOAuth2Api` (X)

- `clientId`, `clientSecret`, `oauthTokenData`
- hidden grant `pkce`; auth `https://x.com/i/oauth2/authorize`; token `https://api.x.com/2/oauth2/token`

### `hubspotAppToken`

- `appToken` (password) — Service Key

### `serpApi`

- `apiKey` (password). Credential/node marked deprecated in-source.

### `googlePalmApi` (Gemini)

- `host` (required, default `https://generativelanguage.googleapis.com`)
- `apiKey` (required, password)

### `microsoftGraphSecurityOAuth2Api` (`extends: microsoftOAuth2Api` → `oAuth2Api`)

- `customScopes`, `enabledScopes` (default `SecurityEvents.ReadWrite.All offline_access`)
- `clientCredentialType`: `clientSecret` | `certificate`
- `clientId`, `clientSecret` or `privateKey`+`certificate`
- `authUrl` / `accessTokenUrl` (defaults `.../common/oauth2/v2.0/...`)
- `graphApiBaseUrl`: graph.microsoft.com | graph.microsoft.us | dod-graph.microsoft.us | microsoftgraph.chinacloudapi.cn
- `oauthTokenData`

### `httpBearerAuth`

- `token` (password, resolvable)

### `openAiApi`

- `apiKey` (required, password)
- `organizationId` (optional)
- `url` (default `https://api.openai.com/v1`)
- `header` (bool), `headerName`, `headerValue` (password)

## Implementation notes

1. Typed resources are thin schemas that build `models.CredentialWrite{Type, Data}` and call the existing `CredentialController`. Do not add per-type HTTP clients.
2. OAuth resources are **client-credentials + imported tokens**, not Connect-in-UI. Document that `oauth_token_data` is required for a usable OAuth credential unless the user completes OAuth in n8n and then does not rotate `data` (use `is_partial_data` or omit `data` on update).
3. Acc tests on CE: prefer non-OAuth types (`httpHeaderAuth`, `httpBasicAuth`, `httpBearerAuth`, `n8nApi` against the same instance, `jwtAuth`). Licensed/OAuth types skip or use placeholder tokens that will not pass `POST .../test`.
4. Keep generic `n8n_credential` as escape hatch for types not in this list.
