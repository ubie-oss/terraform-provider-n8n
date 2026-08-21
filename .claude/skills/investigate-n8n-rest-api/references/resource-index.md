# Public API resource index (tags → path prefixes)

Lean map from OpenAPI `tags` to common `/api/v1` path prefixes. Use this to
locate operations in a resolved instance OpenAPI or the GitHub index. **No
schemas here** — always fetch the live OpenAPI for fields and scopes.

Tag names match `info`/paths tagging in recent Public API OpenAPI (`info.version`
historically `1.1.1`). Instance editions may omit licensed features.

| OpenAPI tag      | Primary path prefixes                                                           | Notes                                             |
| ---------------- | ------------------------------------------------------------------------------- | ------------------------------------------------- |
| Audit            | `/audit`                                                                        | POST generate audit                               |
| CommunityPackage | `/community-packages`                                                           | Install/list/update/uninstall                     |
| Credential       | `/credentials`, `/credentials/{id}`, `/credentials/schema/{credentialTypeName}` | Secrets often omitted on read                     |
| DataTable        | `/data-tables`, `/data-tables/{dataTableId}`, `.../rows`, `.../columns`         | Table + row + column ops                          |
| Discover         | `/discover`                                                                     | Capability map; good live verify                  |
| Evaluation       | `/workflows/{id}/test-runs`, `.../test-cases`                                   | Workflow evaluation runs                          |
| Execution        | `/executions`, `/executions/{id}`, retry/stop/tags                              | List/read/retry/stop                              |
| Folders          | `/projects/{projectId}/folders`, `.../{folderId}`                               | Project-scoped                                    |
| Insights         | `/insights/summary`                                                             | Metrics summary                                   |
| LogStreaming     | `/settings/log-streaming/...`                                                   | Destinations and event types                      |
| N8nPackage       | `/n8n-packages/export`, `/n8n-packages/import`                                  | Beta — may break without major bump               |
| Projects         | `/projects`, `/projects/{projectId}`, `.../users`                               | Projects and membership                           |
| Role             | `/roles`, `/roles/{slug}`                                                       | Instance roles                                    |
| RoleMappingRule  | `/role-mapping-rules`                                                           | IdP role mapping                                  |
| SecurityPolicy   | `/settings/security-policy`                                                     | Instance security policy                          |
| SettingsLdap     | `/settings/ldap`, `/settings/ldap/sync`                                         | LDAP config                                       |
| SettingsOtel     | `/settings/otel`, `/settings/otel/test-trace`                                   | OpenTelemetry                                     |
| SettingsSsoOidc  | `/settings/sso/oidc`                                                            | OIDC SSO                                          |
| SettingsSsoSaml  | `/settings/sso/saml`                                                            | SAML SSO                                          |
| SourceControl    | `/source-control/pull`                                                          | Pull from connected repo                          |
| Tags             | `/tags`, `/tags/{id}`                                                           | Global tag registry                               |
| User             | `/users`, `/users/{id}`, `/users/{id}/role`                                     | Users and global role                             |
| Variables        | `/variables`, `/variables/{id}`                                                 | Instance variables                                |
| Workflow         | `/workflows`, `/workflows/{id}`, `/workflows/{workflowId}`                      | See path-overlap note in [sources.md](sources.md) |

## Workflow path reminders

- Core collection: `GET/POST /workflows`
- By id: `/workflows/{id}` plus activate, deactivate, publish, unpublish, archive, unarchive, transfer, version id
- Overlap: `/workflows/{workflowId}` (get/update-style decorator routes), history, tags

Always re-check the fetched OpenAPI; this index can lag new tags.
