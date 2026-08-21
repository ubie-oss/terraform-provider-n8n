package provider

import (
	"context"
	"fmt"
	"strings"
	"unicode/utf8"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/ubie-oss/terraform-provider-n8n/internal/n8n"
	"github.com/ubie-oss/terraform-provider-n8n/internal/n8n/controllers"
	"github.com/ubie-oss/terraform-provider-n8n/internal/n8n/models"
)

var (
	_ resource.Resource                = &tagResource{}
	_ resource.ResourceWithConfigure   = &tagResource{}
	_ resource.ResourceWithImportState = &tagResource{}
)

type tagResource struct {
	tagController *controllers.TagController
}

type tagResourceModel struct {
	ID               types.String `tfsdk:"id"`
	Name             types.String `tfsdk:"name"`
	CreatedAt        types.String `tfsdk:"created_at"`
	UpdatedAt        types.String `tfsdk:"updated_at"`
	DeleteProtection types.Bool   `tfsdk:"delete_protection"`
}

func NewTagResource() resource.Resource {
	return &tagResource{}
}

func (r *tagResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_tag"
}

func (r *tagResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	markdownDescription, err := readMarkdownDescription(ctx, "internal/provider/docs/resources/tag.md")
	if err != nil {
		resp.Diagnostics.AddError("Unable to read markdown description", err.Error())
		return
	}

	resp.Schema = schema.Schema{
		MarkdownDescription: markdownDescription,
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Tag ID assigned by n8n.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Tag name. Must be unique in the instance registry and at most 24 characters.",
			},
			"created_at": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Creation timestamp from n8n.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"updated_at": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Last update timestamp from n8n.",
			},
			"delete_protection": schema.BoolAttribute{
				Required:            true,
				MarkdownDescription: "When set to `true`, prevents Terraform from destroying this tag. This flag is Terraform-only; n8n has no matching API field. Imported resources default to `true`.",
			},
		},
	}
}

func (r *tagResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	client, ok := req.ProviderData.(*n8n.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *n8n.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}
	r.tagController = controllers.NewTagController(client)
}

func (r *tagResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan tagResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	name := strings.TrimSpace(plan.Name.ValueString())
	if errMsg := tagNameAttributeError(name); errMsg != "" {
		resp.Diagnostics.AddAttributeError(path.Root("name"), "Invalid Tag Name", errMsg)
		return
	}

	created, err := r.tagController.Create(ctx, controllers.CreateTagOptions{Name: name})
	if err != nil {
		resp.Diagnostics.AddError("Error creating tag", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, tagModelFromAPI(created, plan.DeleteProtection))...)
}

func (r *tagResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state tagResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	got, err := r.tagController.Get(ctx, state.ID.ValueString())
	if err != nil {
		if n8n.IsNotFound(err) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Error reading tag", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, tagModelFromAPI(got, state.DeleteProtection))...)
}

func (r *tagResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state tagResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	name := strings.TrimSpace(plan.Name.ValueString())
	if errMsg := tagNameAttributeError(name); errMsg != "" {
		resp.Diagnostics.AddAttributeError(path.Root("name"), "Invalid Tag Name", errMsg)
		return
	}

	if name != strings.TrimSpace(state.Name.ValueString()) {
		updated, err := r.tagController.Update(ctx, controllers.UpdateTagOptions{
			ID:   plan.ID.ValueString(),
			Name: name,
		})
		if err != nil {
			resp.Diagnostics.AddError("Error updating tag", err.Error())
			return
		}
		resp.Diagnostics.Append(resp.State.Set(ctx, tagModelFromAPI(updated, plan.DeleteProtection))...)
		return
	}

	state.Name = types.StringValue(name)
	state.DeleteProtection = plan.DeleteProtection
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *tagResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state tagResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.tagController.Delete(ctx, controllers.DeleteTagOptions{
		ID:               state.ID.ValueString(),
		DeleteProtection: state.DeleteProtection.ValueBool(),
	})
	if err != nil {
		resp.Diagnostics.AddError("Error deleting tag", err.Error())
		return
	}
}

func (r *tagResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	got, err := r.tagController.Import(ctx, controllers.ImportTagOptions{ID: req.ID})
	if err != nil {
		resp.Diagnostics.AddError("Error importing tag", err.Error())
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, tagModelFromAPI(got, types.BoolValue(true)))...)
}

func tagModelFromAPI(tag *models.Tag, deleteProtection types.Bool) tagResourceModel {
	return tagResourceModel{
		ID:               types.StringValue(tag.ID),
		Name:             types.StringValue(tag.Name),
		CreatedAt:        types.StringValue(tag.CreatedAt),
		UpdatedAt:        types.StringValue(tag.UpdatedAt),
		DeleteProtection: deleteProtection,
	}
}

func tagNameAttributeError(name string) string {
	if name == "" {
		return "name must not be empty"
	}
	if utf8.RuneCountInString(name) > controllers.MaxTagNameLength {
		return fmt.Sprintf("name must be at most %d characters", controllers.MaxTagNameLength)
	}
	return ""
}
