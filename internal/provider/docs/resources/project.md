Manages an n8n team project via the Public API.

Only the project name is writable in n8n. The identifier and type are assigned by n8n; created projects have type `team`. Personal projects cannot be created, updated, imported, or destroyed by this resource.

`delete_protection` is required. When set to `true`, Terraform will not destroy the resource. This flag is Terraform-only; n8n has no matching API field. Imported resources default to `delete_protection = true`.

There is no GET-by-id endpoint: Terraform refreshes by listing projects until the id matches. Team-project APIs require an n8n license that includes `feat:projectRole:admin`.
