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
	_ datasource.DataSource              = &foldersDataSource{}
	_ datasource.DataSourceWithConfigure = &foldersDataSource{}
)

type foldersDataSource struct {
	folderController *controllers.FolderController
}

type nestedFolderModel struct {
	FolderID       types.String `tfsdk:"folder_id"`
	Name           types.String `tfsdk:"name"`
	ParentFolderID types.String `tfsdk:"parent_folder_id"`
}

type foldersDataSourceModel struct {
	ID             types.String        `tfsdk:"id"`
	ProjectID      types.String        `tfsdk:"project_id"`
	ParentFolderID types.String        `tfsdk:"parent_folder_id"`
	Name           types.String        `tfsdk:"name"`
	Folders        []nestedFolderModel `tfsdk:"folders"`
}

func NewFoldersDataSource() datasource.DataSource {
	return &foldersDataSource{}
}

func (d *foldersDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_folders"
}

func (d *foldersDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	markdownDescription, err := readMarkdownDescription(ctx, "internal/provider/docs/data_sources/folders.md")
	if err != nil {
		resp.Diagnostics.AddError("Unable to read markdown description", err.Error())
		return
	}

	folderAttrs := map[string]schema.Attribute{
		"folder_id": schema.StringAttribute{
			Computed:            true,
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
	}

	resp.Schema = schema.Schema{
		MarkdownDescription: markdownDescription,
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Data source identifier. Always `folders:{project_id}`.",
			},
			"project_id": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Project whose folders are listed.",
			},
			"parent_folder_id": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "When set, only folders with this parent are returned. Use an empty string for the project root.",
			},
			"name": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "When set, only folders with this exact name are returned.",
			},
			"folders": schema.ListNestedAttribute{
				Computed:            true,
				MarkdownDescription: "Folders returned after following skip/take pagination.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: folderAttrs,
				},
			},
		},
	}
}

func (d *foldersDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *foldersDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config foldersDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	opts := controllers.ListFoldersOptions{
		ProjectID: config.ProjectID.ValueString(),
		Name:      config.Name.ValueString(),
	}
	// Keep empty parent_folder_id as a pointer so List can mean "root only".
	// optionalStringPointer would collapse "" to nil (no filter).
	if !config.ParentFolderID.IsNull() && !config.ParentFolderID.IsUnknown() {
		parent := config.ParentFolderID.ValueString()
		opts.ParentFolderID = &parent
	}

	all, err := d.folderController.List(ctx, opts)
	if err != nil {
		resp.Diagnostics.AddError("Error listing folders", err.Error())
		return
	}

	state := foldersDataSourceModel{
		ID:             types.StringValue("folders:" + config.ProjectID.ValueString()),
		ProjectID:      config.ProjectID,
		ParentFolderID: config.ParentFolderID,
		Name:           config.Name,
		Folders:        make([]nestedFolderModel, 0, len(all)),
	}
	for _, f := range all {
		state.Folders = append(state.Folders, nestedFolderModel{
			FolderID:       types.StringValue(f.ID),
			ParentFolderID: optionalStringValue(f.ParentFolderID),
			Name:           types.StringValue(f.Name),
		})
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
