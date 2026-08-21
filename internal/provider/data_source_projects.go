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
	_ datasource.DataSource              = &projectsDataSource{}
	_ datasource.DataSourceWithConfigure = &projectsDataSource{}
)

type projectsDataSource struct {
	client            *n8n.Client
	projectController *controllers.ProjectController
}

type nestedProjectModel struct {
	ID   types.String `tfsdk:"id"`
	Name types.String `tfsdk:"name"`
	Type types.String `tfsdk:"type"`
}

type projectsDataSourceModel struct {
	ID       types.String         `tfsdk:"id"`
	Projects []nestedProjectModel `tfsdk:"projects"`
}

func NewProjectsDataSource() datasource.DataSource {
	return &projectsDataSource{}
}

func (d *projectsDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_projects"
}

func (d *projectsDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	markdownDescription, err := readMarkdownDescription(ctx, "internal/provider/docs/data_sources/projects.md")
	if err != nil {
		resp.Diagnostics.AddError("Unable to read markdown description", err.Error())
		return
	}

	projectAttrs := map[string]schema.Attribute{
		"id": schema.StringAttribute{
			Computed:            true,
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
	}

	resp.Schema = schema.Schema{
		MarkdownDescription: markdownDescription,
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Data source identifier. Always `projects`.",
			},
			"projects": schema.ListNestedAttribute{
				Computed:            true,
				MarkdownDescription: "Projects returned by GET /projects after following pagination.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: projectAttrs,
				},
			},
		},
	}
}

func (d *projectsDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *projectsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	all, err := d.projectController.List(ctx)
	if err != nil {
		resp.Diagnostics.AddError("Error listing projects", err.Error())
		return
	}

	state := projectsDataSourceModel{
		ID:       types.StringValue("projects"),
		Projects: make([]nestedProjectModel, 0, len(all)),
	}
	for _, p := range all {
		state.Projects = append(state.Projects, nestedProjectModel{
			ID:   types.StringValue(p.ID),
			Name: types.StringValue(p.Name),
			Type: types.StringValue(p.Type),
		})
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
