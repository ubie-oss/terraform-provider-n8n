# Live contract notes — folders and roles

Probed 2026-08-21 against Community n8n from `docker-compose.dev.yml` (pinned 2.35.5), API key minted by `scripts/bootstrap-n8n.sh`. Auth header `X-N8N-API-KEY`. No writes.

Instance OpenAPI: `GET /api/v1/openapi.yml` (`info.version` 1.1.1).

## Folders

- Paths present: `GET/POST /projects/{projectId}/folders`, `GET/PATCH/DELETE /projects/{projectId}/folders/{folderId}`.
- List pagination: query `skip`/`take` typed as strings; default take 10; envelope `{count, data}`.
- DELETE 204; optional `transferToFolderId`.
- PATCH (not PUT) updates name and/or parent.
- GET-by-id 200 includes `totalSubFolders` / `totalWorkflows`.
- Live `GET /api/v1/projects/personal/folders?take=10&skip=0` → **403** body contains `feat:folders`:
  `Your license does not allow for feat:folders. To enable feat:folders, please upgrade to a license that supports this feature.`

Wave A HTTP ops can use httptest fixtures matching this spec; licensed CRUD is `make testacc`.

## Roles (Wave B gate)

- Instance OpenAPI **does** include `/roles`, but **POST only** (`createRole`, `x-decorator-routed: true`). No `/roles/{slug}`, no GET list, no PUT, no DELETE in the resolved spec.
- Live `GET /api/v1/roles` → **405** `{"message":"GET method not allowed"}`.
- Live `POST /api/v1/roles` → **403** body contains `feat:customRoles` (path is mounted; Community is unlicensed).
- Live PUT/PATCH/DELETE `/api/v1/roles` → **405**.

**Wave B is not started.** GET/PUT/DELETE catalog operations are not proven with API-key auth on this instance. Do not call `/rest/roles`.
