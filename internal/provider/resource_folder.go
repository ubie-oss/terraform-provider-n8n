package provider

import (
	"context"
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
	_ resource.Resource                = &folderResource{}
	_ resource.ResourceWithConfigure   = &folderResource{}
	_ resource.ResourceWithImportState = &folderResource{}
)

type folderResource struct {
	folderController *controllers.FolderController
}

type folderResourceModel struct {
	ID                 types.String `tfsdk:"id"`
	ProjectID          types.String `tfsdk:"project_id"`
	FolderID           types.String `tfsdk:"folder_id"`
	Name               types.String `tfsdk:"name"`
	ParentFolderID     types.String `tfsdk:"parent_folder_id"`
	DeleteProtection   types.Bool   `tfsdk:"delete_protection"`
	TransferToFolderID types.String `tfsdk:"transfer_to_folder_id"`
}

func NewFolderResource() resource.Resource {
	return &folderResource{}
}

func (r *folderResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_folder"
}

func (r *folderResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	markdownDescription, err := readMarkdownDescription(ctx, "internal/provider/docs/resources/folder.md")
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
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"project_id": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Team project that owns the folder. Changing this forces replacement; n8n cannot move a folder between projects.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"folder_id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Folder ID assigned by n8n.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Folder name.",
			},
			"parent_folder_id": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Parent folder ID in the same project. Omit for the project root.",
			},
			"delete_protection": schema.BoolAttribute{
				Required:            true,
				MarkdownDescription: "When set to `true`, prevents Terraform from destroying this folder. This flag is Terraform-only; n8n has no matching API field. Imported resources default to `true`.",
			},
			"transfer_to_folder_id": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "When destroying, n8n query `transferToFolderId` that receives workflows and child folders. If omitted, workflows are moved to the project root and archived, and child folders are deleted. This value is Terraform-only. Apply the value in a prior apply before destroy; a same-apply set-and-destroy uses the previous state value.",
			},
		},
	}
}

func (r *folderResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
	r.folderController = controllers.NewFolderController(client)
}

func (r *folderResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan folderResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	name := strings.TrimSpace(plan.Name.ValueString())
	if name == "" {
		resp.Diagnostics.AddAttributeError(path.Root("name"), "Invalid Folder Name", "name must not be empty")
		return
	}

	created, err := r.folderController.Create(ctx, controllers.CreateFolderOptions{
		ProjectID:      plan.ProjectID.ValueString(),
		Name:           name,
		ParentFolderID: optionalStringPointer(plan.ParentFolderID),
	})
	if err != nil {
		resp.Diagnostics.AddError("Error creating folder", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, folderResourceModelFromAPI(plan.ProjectID.ValueString(), created, plan.DeleteProtection, plan.TransferToFolderID))...)
}

func (r *folderResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state folderResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	got, err := r.folderController.Get(ctx, state.ProjectID.ValueString(), state.FolderID.ValueString())
	if err != nil {
		if n8n.IsNotFound(err) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Error reading folder", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, folderResourceModelFromAPI(state.ProjectID.ValueString(), got, state.DeleteProtection, state.TransferToFolderID))...)
}

func (r *folderResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state folderResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	name := strings.TrimSpace(plan.Name.ValueString())
	if name == "" {
		resp.Diagnostics.AddAttributeError(path.Root("name"), "Invalid Folder Name", "name must not be empty")
		return
	}

	nameChanged := name != strings.TrimSpace(state.Name.ValueString())
	parentChanged, moveToRoot, parent := parentFolderChange(plan.ParentFolderID, state.ParentFolderID)
	if !nameChanged && !parentChanged {
		state.DeleteProtection = plan.DeleteProtection
		state.TransferToFolderID = plan.TransferToFolderID
		resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
		return
	}

	updated, err := r.folderController.Update(ctx, controllers.UpdateFolderOptions{
		ProjectID:      plan.ProjectID.ValueString(),
		FolderID:       plan.FolderID.ValueString(),
		Name:           name,
		ParentFolderID: parent,
		MoveToRoot:     moveToRoot,
	})
	if err != nil {
		resp.Diagnostics.AddError("Error updating folder", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, folderResourceModelFromAPI(plan.ProjectID.ValueString(), updated, plan.DeleteProtection, plan.TransferToFolderID))...)
}

func (r *folderResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state folderResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.folderController.Delete(ctx, controllers.DeleteFolderOptions{
		ProjectID:          state.ProjectID.ValueString(),
		FolderID:           state.FolderID.ValueString(),
		DeleteProtection:   state.DeleteProtection.ValueBool(),
		TransferToFolderID: optionalStringPointer(state.TransferToFolderID),
	})
	if err != nil {
		resp.Diagnostics.AddError("Error deleting folder", err.Error())
		return
	}
}

func (r *folderResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	projectID, folderID, err := parseFolderResourceID(req.ID)
	if err != nil {
		resp.Diagnostics.AddError("Error importing folder", err.Error())
		return
	}
	got, err := r.folderController.Import(ctx, controllers.ImportFolderOptions{
		ProjectID: projectID,
		FolderID:  folderID,
	})
	if err != nil {
		resp.Diagnostics.AddError("Error importing folder", err.Error())
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, folderResourceModelFromAPI(projectID, got, types.BoolValue(true), types.StringNull()))...)
}

func folderResourceModelFromAPI(projectID string, f *models.Folder, deleteProtection types.Bool, transfer types.String) folderResourceModel {
	return folderResourceModel{
		ID:                 types.StringValue(formatFolderResourceID(projectID, f.ID)),
		ProjectID:          types.StringValue(projectID),
		FolderID:           types.StringValue(f.ID),
		Name:               types.StringValue(f.Name),
		ParentFolderID:     optionalStringValue(f.ParentFolderID),
		DeleteProtection:   deleteProtection,
		TransferToFolderID: transfer,
	}
}

func parentFolderChange(plan, state types.String) (changed bool, moveToRoot bool, parent *string) {
	planParent := optionalStringPointer(plan)
	stateParent := optionalStringPointer(state)
	switch {
	case planParent == nil && stateParent == nil:
		return false, false, nil
	case planParent == nil:
		return true, true, nil
	case stateParent == nil:
		return true, false, planParent
	case *planParent != *stateParent:
		return true, false, planParent
	default:
		return false, false, nil
	}
}
