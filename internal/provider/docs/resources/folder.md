Manages an n8n folder inside a team project via the Public API.

The resource identifier is `projects/{project_id}/folders/{folder_id}`. `project_id` cannot change in place; n8n cannot move a folder between projects. `parent_folder_id` is optional and may be updated to nest the folder under another folder in the same project, or omitted for the project root.

`delete_protection` is required. When set to `true`, Terraform will not destroy the resource. This flag is Terraform-only; n8n has no matching API field. Imported resources default to `delete_protection = true`.

`transfer_to_folder_id` is sent only when destroying, as the Public API `transferToFolderId` query. Apply the value in a prior apply before destroy; a same-apply set-and-destroy uses the previous state value. If you omit it, n8n moves workflows in the folder to the project root and archives them, and it deletes child folders. Folder APIs require an n8n license that includes `feat:folders`.
