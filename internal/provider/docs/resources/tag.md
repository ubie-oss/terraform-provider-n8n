Manages an n8n global tag via the Public API (`/api/v1/tags`).

Tags live in the instance-wide tag registry. This resource creates, renames, and deletes registry entries. It does **not** attach tags to workflows (that is a separate Public API surface and is out of scope while workflows are unmanaged).

Only `name` is writable. Names must be unique and at most 24 characters (n8n returns a misleading 409 for longer names). `id`, `created_at`, and `updated_at` are assigned by n8n.

`delete_protection` is required. When set to `true`, Terraform will not destroy the resource. This flag is Terraform-only; n8n has no matching API field. Imported resources default to `delete_protection = true`.

Tag CRUD is available on Community Edition.
