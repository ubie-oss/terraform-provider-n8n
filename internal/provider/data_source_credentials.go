package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/ubie-oss/terraform-provider-n8n/internal/n8n"
	"github.com/ubie-oss/terraform-provider-n8n/internal/n8n/controllers"
	"github.com/ubie-oss/terraform-provider-n8n/internal/n8n/models"
)

var (
	_ datasource.DataSource              = &credentialsDataSource{}
	_ datasource.DataSourceWithConfigure = &credentialsDataSource{}
)

type credentialsDataSource struct {
	credentialController *controllers.CredentialController
}

type nestedCredentialSharedModel struct {
	ID        types.String `tfsdk:"id"`
	Name      types.String `tfsdk:"name"`
	Role      types.String `tfsdk:"role"`
	CreatedAt types.String `tfsdk:"created_at"`
	UpdatedAt types.String `tfsdk:"updated_at"`
}

type nestedCredentialModel struct {
	ID        types.String                  `tfsdk:"id"`
	Name      types.String                  `tfsdk:"name"`
	Type      types.String                  `tfsdk:"type"`
	CreatedAt types.String                  `tfsdk:"created_at"`
	UpdatedAt types.String                  `tfsdk:"updated_at"`
	Shared    []nestedCredentialSharedModel `tfsdk:"shared"`
}

type credentialsDataSourceModel struct {
	ID          types.String            `tfsdk:"id"`
	Name        types.String            `tfsdk:"name"`
	Type        types.String            `tfsdk:"type"`
	Credentials []nestedCredentialModel `tfsdk:"credentials"`
}

func NewCredentialsDataSource() datasource.DataSource {
	return &credentialsDataSource{}
}

func (d *credentialsDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_credentials"
}

func (d *credentialsDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	markdownDescription, err := readMarkdownDescription(ctx, "internal/provider/docs/data_sources/credentials.md")
	if err != nil {
		resp.Diagnostics.AddError("Unable to read markdown description", err.Error())
		return
	}

	sharedAttrs := map[string]schema.Attribute{
		"id": schema.StringAttribute{
			Computed:            true,
			MarkdownDescription: "Project ID this credential is shared with.",
		},
		"name": schema.StringAttribute{
			Computed:            true,
			MarkdownDescription: "Project name.",
		},
		"role": schema.StringAttribute{
			Computed:            true,
			MarkdownDescription: "Sharing role (for example `credential:owner`).",
		},
		"created_at": schema.StringAttribute{
			Computed:            true,
			MarkdownDescription: "When the credential was shared with this project.",
		},
		"updated_at": schema.StringAttribute{
			Computed:            true,
			MarkdownDescription: "When the sharing was last updated.",
		},
	}

	credentialAttrs := map[string]schema.Attribute{
		"id": schema.StringAttribute{
			Computed:            true,
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
		"created_at": schema.StringAttribute{
			Computed:            true,
			MarkdownDescription: "Creation timestamp from n8n.",
		},
		"updated_at": schema.StringAttribute{
			Computed:            true,
			MarkdownDescription: "Last update timestamp from n8n.",
		},
		"shared": schema.ListNestedAttribute{
			Computed:            true,
			MarkdownDescription: "Projects this credential is shared with.",
			NestedObject: schema.NestedAttributeObject{
				Attributes: sharedAttrs,
			},
		},
	}

	resp.Schema = schema.Schema{
		MarkdownDescription: markdownDescription,
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Data source identifier. Always `credentials`.",
			},
			"name": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "When set, only credentials with this exact name are returned.",
			},
			"type": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "When set, only credentials with this exact type are returned.",
			},
			"credentials": schema.ListNestedAttribute{
				Computed:            true,
				MarkdownDescription: "Credentials returned after following cursor pagination.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: credentialAttrs,
				},
			},
		},
	}
}

func (d *credentialsDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *credentialsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config credentialsDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	all, err := d.credentialController.List(ctx, controllers.ListCredentialsOptions{
		Name: config.Name.ValueString(),
		Type: config.Type.ValueString(),
	})
	if err != nil {
		resp.Diagnostics.AddError("Error listing credentials", err.Error())
		return
	}

	state := credentialsDataSourceModel{
		ID:          types.StringValue("credentials"),
		Name:        config.Name,
		Type:        config.Type,
		Credentials: make([]nestedCredentialModel, 0, len(all)),
	}
	for _, item := range all {
		state.Credentials = append(state.Credentials, nestedCredentialModelFromAPI(item))
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func nestedCredentialModelFromAPI(item models.CredentialListItem) nestedCredentialModel {
	shared := make([]nestedCredentialSharedModel, 0, len(item.Shared))
	for _, s := range item.Shared {
		shared = append(shared, nestedCredentialSharedModel{
			ID:        types.StringValue(s.ID),
			Name:      types.StringValue(s.Name),
			Role:      types.StringValue(s.Role),
			CreatedAt: types.StringValue(s.CreatedAt),
			UpdatedAt: types.StringValue(s.UpdatedAt),
		})
	}
	return nestedCredentialModel{
		ID:        types.StringValue(item.ID),
		Name:      types.StringValue(item.Name),
		Type:      types.StringValue(item.Type),
		CreatedAt: types.StringValue(item.CreatedAt),
		UpdatedAt: types.StringValue(item.UpdatedAt),
		Shared:    shared,
	}
}
