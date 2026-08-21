package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/ubie-oss/terraform-provider-n8n/internal/n8n"
	"github.com/ubie-oss/terraform-provider-n8n/internal/n8n/controllers"
)

var (
	_ datasource.DataSource              = &credentialSchemaDataSource{}
	_ datasource.DataSourceWithConfigure = &credentialSchemaDataSource{}
)

type credentialSchemaDataSource struct {
	credentialController *controllers.CredentialController
}

type credentialSchemaDataSourceModel struct {
	ID         types.String `tfsdk:"id"`
	Type       types.String `tfsdk:"type"`
	SchemaJSON types.String `tfsdk:"schema_json"`
}

func NewCredentialSchemaDataSource() datasource.DataSource {
	return &credentialSchemaDataSource{}
}

func (d *credentialSchemaDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_credential_schema"
}

func (d *credentialSchemaDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	markdownDescription, err := readMarkdownDescription(ctx, "internal/provider/docs/data_sources/credential_schema.md")
	if err != nil {
		resp.Diagnostics.AddError("Unable to read markdown description", err.Error())
		return
	}

	resp.Schema = schema.Schema{
		MarkdownDescription: markdownDescription,
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Data source identifier. Always `schema:{type}`.",
			},
			"type": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "n8n credential type name (for example `httpHeaderAuth`).",
			},
			"schema_json": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "JSON Schema-like object returned by GET /credentials/schema/{type}.",
			},
		},
	}
}

func (d *credentialSchemaDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	client, ok := req.ProviderData.(*n8n.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *n8n.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}
	d.credentialController = controllers.NewCredentialController(client)
}

func (d *credentialSchemaDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config credentialSchemaDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	typ := config.Type.ValueString()
	raw, err := d.credentialController.Schema(ctx, typ)
	if err != nil {
		resp.Diagnostics.AddError("Error reading credential schema", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, credentialSchemaDataSourceModel{
		ID:         types.StringValue("schema:" + typ),
		Type:       types.StringValue(typ),
		SchemaJSON: types.StringValue(string(raw)),
	})...)
}
