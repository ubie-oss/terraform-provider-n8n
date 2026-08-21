package provider

import (
	"context"
	"errors"
	"fmt"
	"strings"

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
	_ resource.Resource                = &projectResource{}
	_ resource.ResourceWithConfigure   = &projectResource{}
	_ resource.ResourceWithImportState = &projectResource{}
)

type projectResource struct {
	client            *n8n.Client
	projectController *controllers.ProjectController
}

type projectResourceModel struct {
	ID               types.String `tfsdk:"id"`
	Name             types.String `tfsdk:"name"`
	Type             types.String `tfsdk:"type"`
	DeleteProtection types.Bool   `tfsdk:"delete_protection"`
}

func NewProjectResource() resource.Resource {
	return &projectResource{}
}

func (r *projectResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_project"
}

func (r *projectResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	markdownDescription, err := readMarkdownDescription(ctx, "internal/provider/docs/resources/project.md")
	if err != nil {
		resp.Diagnostics.AddError("Unable to read markdown description", err.Error())
		return
	}

	resp.Schema = schema.Schema{
		MarkdownDescription: markdownDescription,
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Project ID assigned by n8n.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Project name.",
			},
			"type": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Project type from n8n. Managed resources are always `team`.",
			},
			"delete_protection": schema.BoolAttribute{
				Required:            true,
				MarkdownDescription: "When set to `true`, prevents Terraform from destroying this project. This flag is Terraform-only; n8n has no matching API field. Imported resources default to `true`.",
			},
		},
	}
}

func (r *projectResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
	r.client = client
	r.projectController = controllers.NewProjectController(client)
}

func (r *projectResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan projectResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	name := strings.TrimSpace(plan.Name.ValueString())
	if name == "" {
		resp.Diagnostics.AddAttributeError(path.Root("name"), "Invalid Project Name", "name must not be empty")
		return
	}

	created, err := r.projectController.Create(ctx, controllers.CreateProjectOptions{Name: name})
	if err != nil {
		resp.Diagnostics.AddError("Error creating project", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, projectModelFromAPI(created, plan.DeleteProtection))...)
}

func (r *projectResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state projectResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	got, err := r.projectController.GetTeam(ctx, state.ID.ValueString())
	if err != nil {
		if errors.Is(err, controllers.ErrProjectNotFound) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Error reading project", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, projectModelFromAPI(got, state.DeleteProtection))...)
}

func (r *projectResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state projectResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	name := strings.TrimSpace(plan.Name.ValueString())
	if name == "" {
		resp.Diagnostics.AddAttributeError(path.Root("name"), "Invalid Project Name", "name must not be empty")
		return
	}

	if name != strings.TrimSpace(state.Name.ValueString()) {
		updated, err := r.projectController.Update(ctx, controllers.UpdateProjectOptions{
			ID:   plan.ID.ValueString(),
			Name: name,
		})
		if err != nil {
			resp.Diagnostics.AddError("Error updating project", err.Error())
			return
		}
		resp.Diagnostics.Append(resp.State.Set(ctx, projectModelFromAPI(updated, plan.DeleteProtection))...)
		return
	}

	state.Name = types.StringValue(name)
	state.DeleteProtection = plan.DeleteProtection
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *projectResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state projectResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.projectController.Delete(ctx, controllers.DeleteProjectOptions{
		ID:               state.ID.ValueString(),
		DeleteProtection: state.DeleteProtection.ValueBool(),
	})
	if err != nil {
		resp.Diagnostics.AddError("Error deleting project", err.Error())
		return
	}
}

func (r *projectResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	got, err := r.projectController.Import(ctx, controllers.ImportProjectOptions{ID: req.ID})
	if err != nil {
		resp.Diagnostics.AddError("Error importing project", err.Error())
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, projectModelFromAPI(got, types.BoolValue(true)))...)
}

func projectModelFromAPI(p *models.Project, deleteProtection types.Bool) projectResourceModel {
	return projectResourceModel{
		ID:               types.StringValue(p.ID),
		Name:             types.StringValue(p.Name),
		Type:             types.StringValue(p.Type),
		DeleteProtection: deleteProtection,
	}
}
