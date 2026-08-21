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
	_ datasource.DataSource              = &folderDataSource{}
	_ datasource.DataSourceWithConfigure = &folderDataSource{}
)

type folderDataSource struct {
	folderController *controllers.FolderController
}

type folderDataSourceModel struct {
	ID              types.String `tfsdk:"id"`
	ProjectID       types.String `tfsdk:"project_id"`
	FolderID        types.String `tfsdk:"folder_id"`
	Name            types.String `tfsdk:"name"`
	ParentFolderID  types.String `tfsdk:"parent_folder_id"`
	TotalSubFolders types.Int64  `tfsdk:"total_sub_folders"`
	TotalWorkflows  types.Int64  `tfsdk:"total_workflows"`
}

func NewFolderDataSource() datasource.DataSource {
	return &folderDataSource{}
}

func (d *folderDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_folder"
}

func (d *folderDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	markdownDescription, err := readMarkdownDescription(ctx, "internal/provider/docs/data_sources/folder.md")
	if err != nil {
		resp.Diagnostics.AddError("Unable to read markdown description", err.Error())
		return
	}

	resp.Schema = schema.Schema{
		MarkdownDescription: markdownDescription,
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Composite identifier `projects/{project_id}/folders/{folder_id}`.",
			},
			"project_id": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Project that owns the folder.",
			},
			"folder_id": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Folder ID.",
			},
			"name": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Folder name.",
			},
			"parent_folder_id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Parent folder ID, or null at the project root.",
			},
			"total_sub_folders": schema.Int64Attribute{
				Computed:            true,
				MarkdownDescription: "Recursive count of child folders.",
			},
			"total_workflows": schema.Int64Attribute{
				Computed:            true,
				MarkdownDescription: "Recursive count of workflows in the folder tree.",
			},
		},
	}
}

func (d *folderDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
	d.folderController = controllers.NewFolderController(client)
}

func (d *folderDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config folderDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	got, err := d.folderController.Get(ctx, config.ProjectID.ValueString(), config.FolderID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error reading folder", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, folderDataSourceModelFromAPI(config.ProjectID.ValueString(), got))...)
}

func folderDataSourceModelFromAPI(projectID string, f *models.Folder) folderDataSourceModel {
	return folderDataSourceModel{
		ID:              types.StringValue(formatFolderResourceID(projectID, f.ID)),
		ProjectID:       types.StringValue(projectID),
		FolderID:        types.StringValue(f.ID),
		Name:            types.StringValue(f.Name),
		ParentFolderID:  optionalStringValue(f.ParentFolderID),
		TotalSubFolders: optionalInt64Value(f.TotalSubFolders),
		TotalWorkflows:  optionalInt64Value(f.TotalWorkflows),
	}
}

func optionalInt64Value(v *int) types.Int64 {
	if v == nil {
		return types.Int64Null()
	}
	return types.Int64Value(int64(*v))
}
