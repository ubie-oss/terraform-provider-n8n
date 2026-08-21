# Community Edition: Public API and n8n CLI limits

Developer reference for what **Community** n8n from [`docker-compose.dev.yml`](../../docker-compose.dev.yml) can and cannot exercise. Use this when choosing `make testacc-docker` vs licensed `make testacc`, or when probing with [n8n CLI](https://docs.n8n.io/connect/n8n-cli).

## Disclaimer

- The Terraform provider contract is the **Public REST API** (`/api/v1` + `X-N8N-API-KEY`) and instance/GitHub OpenAPI — **not** CLI help or CLI command shapes.
- n8n CLI wraps the Public API. CLI success on CE means the underlying Public API is available on that edition; it is not a license to invent Terraform schemas from CLI output.
- See also: [api-surfaces.md](../../.claude/skills/investigate-n8n-rest-api/references/api-surfaces.md).

## Probe pins

| Item      | Value                                                                                                 |
| --------- | ----------------------------------------------------------------------------------------------------- |
| Date      | 2026-08-21                                                                                            |
| n8n image | `docker.n8n.io/n8nio/n8n:2.35.5` ([compose](../../docker-compose.dev.yml))                            |
| CLI       | `@n8n/cli` `0.15.0` (`npx @n8n/cli`)                                                                  |
| Base URL  | `http://127.0.0.1:5678`                                                                               |
| Auth      | API key from [`scripts/bootstrap-n8n.sh`](../../scripts/bootstrap-n8n.sh) (`N8N_URL` / `N8N_API_KEY`) |

Re-verify after bumping the compose image pin or a major CLI release.

## Works on Community Edition

Exercised successfully via n8n CLI (Public API underneath):

| Area                                            | Commands / ops                                                                                     | Notes                                                                        |
| ----------------------------------------------- | -------------------------------------------------------------------------------------------------- | ---------------------------------------------------------------------------- |
| **workflow**                                    | `list`, `get`, `create`, `update`, `delete`, `activate`, `deactivate`, `tags`                      | `activate` needs a trigger, webhook, or polling node (manual-only fails).    |
| **credential**                                  | `list`, `get`, `schema`, `create`, `delete`                                                        | `get` returns metadata, not secrets.                                         |
| **tag**                                         | `list`, `create`, `update`, `delete`                                                               | Full CRUD.                                                                   |
| **user**                                        | `list`, `get`                                                                                      |                                                                              |
| **execution**                                   | `list`, `get`, `retry`, `stop`, `delete`                                                           | `retry` only for failed runs; `stop` cancels running executions.             |
| **data-table**                                  | `list`, `get`, `create`, `delete`, `rows`, `add-rows`, `update-rows`, `upsert-rows`, `delete-rows` | Tables land under a **personal** `projectId`. See payload notes below.       |
| **package**                                     | `export`, `import`                                                                                 | Preview. Conflict policies: `fail`, `skip`, `new-version` (not `overwrite`). |
| **audit**                                       | top-level `audit`                                                                                  |                                                                              |
| **config** / **skill** / **login** / **logout** | local CLI config and skill install                                                                 | Not instance features; no CE license gate.                                   |

## Blocked on Community Edition (license)

| Area               | Commands                                                                       | Typical error                                                                    |
| ------------------ | ------------------------------------------------------------------------------ | -------------------------------------------------------------------------------- |
| **project**        | `list`, `create`, `update`, `delete`, `members`, `add-member`, `remove-member` | `feat:projectRole:admin` not licensed                                            |
| **project**        | `get`                                                                          | `GET method not allowed` (even for personal project id observed via data tables) |
| **variable**       | `list`, `create`, `update`, `delete`                                           | `feat:variables` not licensed                                                    |
| **source-control** | `pull`                                                                         | `Source Control feature is not licensed`                                         |

Provider resources that need the same gates (documented elsewhere): team projects → `feat:projectRole:admin`; folders → `feat:folders`; custom-role writes → Enterprise / `feat:customRoles`.

## Folders and roles (CLI probe, 2026-08-21)

`@n8n/cli` **0.15.0 has no** `folder`, `folders`, `role`, or `roles` topics (`command … not found`). Folder and role behavior is only reachable indirectly.

### Folders

| Surface                                                          | Result on CE                                                                                                                                                              |
| ---------------------------------------------------------------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| CLI topic `folder` / `folders`                                   | **Missing** — not implemented in this CLI version                                                                                                                         |
| Public API `GET/POST …/projects/{id}/folders` (incl. `personal`) | **403** `feat:folders`                                                                                                                                                    |
| `package export --folder-id=…`                                   | Calls `POST /api/v1/n8n-packages/export`; with a fake/missing id → **400** “folder(s) not found or not accessible” (cannot create folders on CE to exercise a happy path) |
| `package import --folder-id=…`                                   | **Error** “Folder not found in target project” for a fake id                                                                                                              |
| `package import --folder-conflict-policy=…`                      | Flag accepted; with a workflow-only package (no folders) import can still succeed (`skip` conflict policy)                                                                |

**Conclusion:** CLI cannot CRUD folders. Package flags assume folders already exist (licensed). Provider folder work stays on Public API + licensed `make testacc` (see [folders/roles live contract](../../docs/superpowers/plans/contracts/2026-08-21-folders-roles-live.md)).

### Roles

| Surface                                                   | Result on CE                                                                                           |
| --------------------------------------------------------- | ------------------------------------------------------------------------------------------------------ |
| CLI topic `role` / `roles`                                | **Missing** — not implemented in this CLI version                                                      |
| Public API `GET /api/v1/roles`                            | **405** method not allowed                                                                             |
| Public API `POST /api/v1/roles`                           | **403** `feat:customRoles`                                                                             |
| `project add-member --role=project:editor\|viewer\|admin` | **403** `feat:projectRole:admin` (same for `members` / `remove-member`, including personal project id) |
| `project create` / `list` / …                             | **403** `feat:projectRole:admin`                                                                       |

**Conclusion:** The only CLI “role” knob is **project membership role strings** on `project add-member`, and those ops are license-gated with team projects. Custom role catalog CRUD is not exposed as CLI commands and is unlicensed on CE via Public API.

## Reachable but not useful on CE: transfer

`workflow transfer` and `credential transfer` call the Public API. On CE, resources already live in the **personal** project and there is no other team project to target:

| Attempt                                     | Result                                 |
| ------------------------------------------- | -------------------------------------- |
| Transfer into the same personal `projectId` | Error: already belongs to that project |
| Transfer into a missing / fake project id   | 404                                    |
| Create a team project first, then transfer  | Blocked by `feat:projectRole:admin`    |

Treat transfer as **unusable on CE** for practical testing.

## Data-table CLI payload notes (CE)

Handy when re-probing; still not a Terraform schema:

- `data-table create --columns='[...]'` expects **inline JSON**, not a file path.
- `update-rows` / `upsert-rows` `--file` / `--stdin` body: `{ "filter": { "filters": [ { "columnName", "condition", "value" } ] }, "data": { ... } }`.
- `delete-rows --filter` needs the nested `filters` object shape above (a bare `{"name":"..."}` is rejected).

## Acceptance-test implications

| Goal                                               | Use                                                               |
| -------------------------------------------------- | ----------------------------------------------------------------- |
| Public API smoke, bootstrap, CE-available surfaces | `make testacc-docker` (Community compose + bootstrap)             |
| Team-project CRUD (`n8n_project`, …)               | Licensed instance + `make testacc`                                |
| Folder / custom-role CRUD                          | Licensed features (`feat:folders`, custom roles) + `make testacc` |

Do **not** add a mock Public API server. Licensed tests should **skip** on CE 403, not fail.

## Related

- [CONTRIBUTING.md](../../CONTRIBUTING.md) — how to run acceptance tests
- [AGENTS.md](../../AGENTS.md) — short agent notes
- [n8n CLI docs](https://docs.n8n.io/connect/n8n-cli)
- [n8n Public API docs](https://docs.n8n.io/connect/n8n-api.md)
