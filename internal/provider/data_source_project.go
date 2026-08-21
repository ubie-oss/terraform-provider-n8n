package provider

import (
	"context"
	"errors"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/ubie-oss/terraform-provider-n8n/internal/n8n"
	"github.com/ubie-oss/terraform-provider-n8n/internal/n8n/controllers"
	"github.com/ubie-oss/terraform-provider-n8n/internal/n8n/models"
)

var (
	_ datasource.DataSource              = &projectDataSource{}
	_ datasource.DataSourceWithConfigure = &projectDataSource{}
)

type projectDataSource struct {
	client            *n8n.Client
	projectController *controllers.ProjectController
}

type projectDataSourceModel struct {
	ID   types.String `tfsdk:"id"`
	Name types.String `tfsdk:"name"`
	Type types.String `tfsdk:"type"`
}

func NewProjectDataSource() datasource.DataSource {
	return &projectDataSource{}
}

func (d *projectDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_project"
}

func (d *projectDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	markdownDescription, err := readMarkdownDescription(ctx, "internal/provider/docs/data_sources/project.md")
	if err != nil {
		resp.Diagnostics.AddError("Unable to read markdown description", err.Error())
		return
	}

	resp.Schema = schema.Schema{
		MarkdownDescription: markdownDescription,
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Project ID.",
			},
			"name": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Project name.",
			},
			"type": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Project type (`team` or `personal`).",
			},
		},
	}
}

func (d *projectDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
	d.client = client
	d.projectController = controllers.NewProjectController(client)
}

func (d *projectDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config projectDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	got, err := d.projectController.Get(ctx, config.ID.ValueString())
	if err != nil {
		if errors.Is(err, controllers.ErrProjectNotFound) {
			resp.Diagnostics.AddError("Project not found", fmt.Sprintf("No project found with id %q", config.ID.ValueString()))
			return
		}
		resp.Diagnostics.AddError("Error reading project", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, projectDataSourceModelFromAPI(got))...)
}

func projectDataSourceModelFromAPI(p *models.Project) projectDataSourceModel {
	return projectDataSourceModel{
		ID:   types.StringValue(p.ID),
		Name: types.StringValue(p.Name),
		Type: types.StringValue(p.Type),
	}
}
