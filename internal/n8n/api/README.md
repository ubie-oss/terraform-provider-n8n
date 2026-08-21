# `internal/n8n/api`

Versioned n8n Public API HTTP operations. One directory per API version and resource, one file per HTTP verb — the same split Lightdash uses under `internal/lightdash/api/v1` and `api/v2`, with an extra resource folder so later n8n surfaces (credentials, variables, tags) do not share one package.

```text
api/
  v1/
    projects/          # GET/POST /projects, PUT/DELETE /projects/{projectId}
      list_projects_v1.go
      create_project_v1.go
      update_project_v1.go
      delete_project_v1.go
    folders/           # GET/POST /projects/{projectId}/folders, GET/PATCH/DELETE .../{folderId}
      list_folders_v1.go
      create_folder_v1.go
      get_folder_v1.go
      update_folder_v1.go
      delete_folder_v1.go
    credentials/       # GET/POST /credentials, GET/PATCH/DELETE /credentials/{id}, schema, transfer
      list_credentials_v1.go
      create_credential_v1.go
      get_credential_v1.go
      update_credential_v1.go
      delete_credential_v1.go
      get_credential_schema_v1.go
      transfer_credential_v1.go
```

Callers import the resource package (typically as `projectsv1`) and use `*V1` function names. Shared HTTP (`DoJSON`, `X-N8N-API-KEY`) stays on `internal/n8n.Client`.
