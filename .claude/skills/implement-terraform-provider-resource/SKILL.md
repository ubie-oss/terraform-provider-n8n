---
name: implement-terraform-provider-resource
description: Implement Terraform resources and data sources in `internal/provider` with testing aligned to this provider.
---

# Implement Terraform provider resource

## Description

Implement a new Terraform resource or data source in `internal/provider` using `hashicorp/terraform-plugin-framework`, following patterns in this repository (for example `resource_project.go`, `data_source_project.go`, and the API client in `internal/n8n`).

## Input

The user should provide:

1. **Type**: Resource or data source.
2. **Terraform type name suffix** (e.g., `tag` → `n8n_tag` when `ProviderTypeName` is `n8n`).
3. **Schema**: Attributes (name, type, required / optional / computed).
4. **Client operations**: HTTP helpers in `internal/n8n/api/v1/<resource>` and/or a service in `internal/n8n/services`. Terraform CRUD goes through `internal/n8n/controllers`, not services.

## Workflow

### 1. Implementation

- **Files**: Add `internal/provider/resource_<name>.go` or `internal/provider/data_source_<name>.go`.
- **Model**: Go struct with `tfsdk` tags; use `types.String`, `types.Bool`, etc.
- **Interfaces**: Implement `resource.Resource` (+ `Configure`, `ImportState` as needed) or `datasource.DataSource` (+ `Configure` as needed).
- **Schema**: Define in `Schema`; use clear descriptions.
- **Configure**: Read `*n8n.Client` from `req.ProviderData`, matching `internal/provider/provider.go` `Configure`, then construct the controller (`controllers.NewProjectController(client)`).
- **CRUD / read**: Call the controller; surface errors with `resp.Diagnostics`. Do not construct `services.ProjectService` in the resource.

### 2. Unit tests

- Add `internal/provider/resource_<name>_test.go` or `data_source_<name>_test.go` for metadata, helpers, and behavior that does not need Terraform CLI.

### 3. Acceptance tests (optional)

- Use `isIntegrationTestMode()` and `testAccPreCheck(t)`.
- Use `testAccProtoV6ProviderFactories` with provider key `n8n`.
- Prefer small `.tf` fixtures under `internal/provider/acc_tests/` when you add them.
- Use `resource.Test` from `github.com/hashicorp/terraform-plugin-testing/helper/resource`.

## Example patterns

### Resource structure (same package as provider)

```go
package provider

import (
    "context"

    "github.com/ubie-oss/terraform-provider-n8n/internal/n8n"
    "github.com/hashicorp/terraform-plugin-framework/resource"
    "github.com/hashicorp/terraform-plugin-framework/types"
)

var _ resource.Resource = &widgetResource{}

type widgetResource struct {
    client *n8n.Client
}

type widgetResourceModel struct {
    ID   types.String `tfsdk:"id"`
    Name types.String `tfsdk:"name"`
}
```

### Acceptance test snippet

```go
func TestAccWidget_basic(t *testing.T) {
    if !isIntegrationTestMode() {
        t.Skip("acceptance tests require TF_ACC=1")
    }
    resource.Test(t, resource.TestCase{
        PreCheck:                 func() { testAccPreCheck(t) },
        ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
        Steps: []resource.TestStep{
            {
                Config: testAccWidgetConfigBasic,
                Check: resource.ComposeTestCheckFunc(
                    resource.TestCheckResourceAttr("n8n_widget.test", "name", "value"),
                ),
            },
        },
    })
}
```

## Reference

- [Terraform Plugin Framework](https://developer.hashicorp.com/terraform/plugin/framework)
- [Provider implementation rules](.cursor/rules/provider-implementation.mdc)
- [Project structure](.cursor/rules/project-structure.mdc)
- [Testing setup](references/testing_setup.md)

## Assets

- [Resource boilerplate](assets/resource_boilerplate.go)
- [Data source boilerplate](assets/data_source_boilerplate.go)
- [Acceptance test boilerplate](assets/test_boilerplate.go)
