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
	_ datasource.DataSource              = &tagsDataSource{}
	_ datasource.DataSourceWithConfigure = &tagsDataSource{}
)

type tagsDataSource struct {
	tagController *controllers.TagController
}

type tagsDataSourceModel struct {
	ID   types.String         `tfsdk:"id"`
	Tags []tagDataSourceModel `tfsdk:"tags"`
}

func NewTagsDataSource() datasource.DataSource {
	return &tagsDataSource{}
}

func (d *tagsDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_tags"
}

func (d *tagsDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	markdownDescription, err := readMarkdownDescription(ctx, "internal/provider/docs/data_sources/tags.md")
	if err != nil {
		resp.Diagnostics.AddError("Unable to read markdown description", err.Error())
		return
	}

	tagAttrs := map[string]schema.Attribute{
		"id": schema.StringAttribute{
			Computed:            true,
			MarkdownDescription: "Tag ID.",
		},
		"name": schema.StringAttribute{
			Computed:            true,
			MarkdownDescription: "Tag name.",
		},
		"created_at": schema.StringAttribute{
			Computed:            true,
			MarkdownDescription: "Creation timestamp from n8n.",
		},
		"updated_at": schema.StringAttribute{
			Computed:            true,
			MarkdownDescription: "Last update timestamp from n8n.",
		},
	}

	resp.Schema = schema.Schema{
		MarkdownDescription: markdownDescription,
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Data source identifier. Always `tags`.",
			},
			"tags": schema.ListNestedAttribute{
				Computed:            true,
				MarkdownDescription: "Tags returned by GET /tags after following pagination.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: tagAttrs,
				},
			},
		},
	}
}

func (d *tagsDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
	d.tagController = controllers.NewTagController(client)
}

func (d *tagsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	all, err := d.tagController.List(ctx)
	if err != nil {
		resp.Diagnostics.AddError("Error listing tags", err.Error())
		return
	}

	state := tagsDataSourceModel{
		ID:   types.StringValue("tags"),
		Tags: make([]tagDataSourceModel, 0, len(all)),
	}
	for _, tag := range all {
		state.Tags = append(state.Tags, tagDataSourceModelFromAPI(&tag))
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
