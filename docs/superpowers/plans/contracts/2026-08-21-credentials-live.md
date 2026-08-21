# Live contract notes — credentials

Probed 2026-08-21 against Community n8n from `docker-compose.dev.yml` (pinned 2.35.5), API key minted by `scripts/bootstrap-n8n.sh`. Auth header `X-N8N-API-KEY`. Throwaway `httpHeaderAuth` credentials were created and deleted. Secret `data` values are not recorded here.

Instance OpenAPI: `GET /api/v1/openapi.yml` (`info.version` 1.1.1). Paths present: `GET/POST /credentials`, `GET/PATCH/DELETE /credentials/{id}`, `GET /credentials/schema/{credentialTypeName}`, `PUT /credentials/{id}/transfer`, `POST /credentials/{id}/test`.

## Schema

- `GET /credentials/schema/httpHeaderAuth` → 200 JSON Schema-like object: `additionalProperties: false`, `type: object`, properties `name`, `value` (string), plus `useCustomAuth` (`type: notice`), `allowedHttpRequestDomains`, `allowedDomains`. **`required` is `[]`** (not `name`/`value`).
- `httpBasicAuth` similar (`user`, `password`). `slackOAuth2Api` includes n8n-specific types (`notice`, `json`) and `allOf` conditionals. This is **not** strict JSON Schema — do not add plan-time validation in v1.
- Unknown type → 404 `{"message":"Not Found"}`.
- Acc-test type: **`httpHeaderAuth`** with `data = { name = "X-Test", value = "placeholder" }`.

## Create / read

- `POST /credentials` → **200** (not 201). Body matches OpenAPI `create-credential-response`: `id`, `name`, `type`, `isManaged`, `isGlobal`, `isResolvable`, `resolvableAllowFallback`, `resolverId`, `createdAt`, `updatedAt`. **No `data`. No `projectId`. No `shared`.**
- `isResolvable: true` on create is honored. `isGlobal: true` on create is **ignored** (response stays `false`).
- `GET /credentials/{id}` → same metadata shape as create. **Never includes `data`.**
- `GET /credentials?limit=` → `{ data, nextCursor }`. List **items are slimmer than OpenAPI `credentialListItem`**: `id`, `name`, `type`, `createdAt`, `updatedAt`, `shared`. No `isManaged` / `isGlobal` / `isResolvable` on list items. No `data`.
- `shared[]` on list: `{ id, name, role, createdAt, updatedAt }`. Owner entry is `role: credential:owner`. On CE the `id` is the **personal project UUID** (not the alias `personal`); `name` is the personal-project display name.

## Update / delete / transfer / global

- `PATCH` name-only → 200, name changes.
- `PATCH` `data` with `isPartialData` true or false → 200. GET still omits `data` (cannot observe merge vs replace from GET).
- `PATCH` `isGlobal: true` → **403** `You are not licensed for sharing credentials` on Community. Do not set `is_global` in CE acc tests.
- `PUT /credentials/{id}/transfer` with `destinationProjectId: "personal"` → **404** (alias is not accepted). Fake UUID → 404. Transfer needs a real team project id (licensed).
- `DELETE` → **200** with metadata plus extra `usageScope` (`"project"`). **No `id` in the delete body.** Ignore the delete payload.

Live wins over OpenAPI where they disagree (list item fields, `isGlobal` on create, delete body, schema `required: []`).
