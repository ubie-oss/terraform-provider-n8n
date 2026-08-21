package provider

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/ubie-oss/terraform-provider-n8n/internal/n8n"
	"github.com/ubie-oss/terraform-provider-n8n/internal/n8n/controllers"
	"github.com/ubie-oss/terraform-provider-n8n/internal/n8n/models"
)

var (
	_ datasource.DataSource              = &credentialDataSource{}
	_ datasource.DataSourceWithConfigure = &credentialDataSource{}
)

type credentialDataSource struct {
	credentialController *controllers.CredentialController
}

type credentialDataSourceModel struct {
	ID                      types.String `tfsdk:"id"`
	Name                    types.String `tfsdk:"name"`
	Type                    types.String `tfsdk:"type"`
	IsManaged               types.Bool   `tfsdk:"is_managed"`
	IsGlobal                types.Bool   `tfsdk:"is_global"`
	IsResolvable            types.Bool   `tfsdk:"is_resolvable"`
	ResolvableAllowFallback types.Bool   `tfsdk:"resolvable_allow_fallback"`
	ResolverID              types.String `tfsdk:"resolver_id"`
	CreatedAt               types.String `tfsdk:"created_at"`
	UpdatedAt               types.String `tfsdk:"updated_at"`
}

func NewCredentialDataSource() datasource.DataSource {
	return &credentialDataSource{}
}

func (d *credentialDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_credential"
}

func (d *credentialDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	markdownDescription, err := readMarkdownDescription(ctx, "internal/provider/docs/data_sources/credential.md")
	if err != nil {
		resp.Diagnostics.AddError("Unable to read markdown description", err.Error())
		return
	}

	resp.Schema = schema.Schema{
		MarkdownDescription: markdownDescription,
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Credential ID.",
			},
			"name": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Credential name.",
			},
			"type": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "n8n credential type name.",
			},
			"is_managed": schema.BoolAttribute{
				Computed:            true,
				MarkdownDescription: "Whether n8n manages this credential.",
			},
			"is_global": schema.BoolAttribute{
				Computed:            true,
				MarkdownDescription: "Whether this credential is available globally.",
			},
			"is_resolvable": schema.BoolAttribute{
				Computed:            true,
				MarkdownDescription: "Whether this credential has resolvable fields.",
			},
			"resolvable_allow_fallback": schema.BoolAttribute{
				Computed:            true,
				MarkdownDescription: "Whether a credential resolver may fall back to static data.",
			},
			"resolver_id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Dynamic credential resolver id, if any.",
			},
			"created_at": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Creation timestamp from n8n.",
			},
			"updated_at": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Last update timestamp from n8n.",
			},
		},
	}
}

func (d *credentialDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *credentialDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config credentialDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	got, err := d.credentialController.Get(ctx, config.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error reading credential", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, credentialDataSourceModelFromAPI(got))...)
}

func credentialDataSourceModelFromAPI(c *models.Credential) credentialDataSourceModel {
	resolver := types.StringNull()
	if c.ResolverID != nil && strings.TrimSpace(*c.ResolverID) != "" {
		resolver = types.StringValue(*c.ResolverID)
	}
	return credentialDataSourceModel{
		ID:                      types.StringValue(c.ID),
		Name:                    types.StringValue(c.Name),
		Type:                    types.StringValue(c.Type),
		IsManaged:               types.BoolValue(c.IsManaged),
		IsGlobal:                types.BoolValue(c.IsGlobal),
		IsResolvable:            types.BoolValue(c.IsResolvable),
		ResolvableAllowFallback: types.BoolValue(c.ResolvableAllowFallback),
		ResolverID:              resolver,
		CreatedAt:               types.StringValue(c.CreatedAt),
		UpdatedAt:               types.StringValue(c.UpdatedAt),
	}
}
