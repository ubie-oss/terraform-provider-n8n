package provider

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/ubie-oss/terraform-provider-n8n/internal/n8n"
	"github.com/ubie-oss/terraform-provider-n8n/internal/n8n/controllers"
	"github.com/ubie-oss/terraform-provider-n8n/internal/n8n/models"
)

var (
	_ resource.Resource                = &credentialResource{}
	_ resource.ResourceWithConfigure   = &credentialResource{}
	_ resource.ResourceWithImportState = &credentialResource{}
)

type credentialResource struct {
	credentialController *controllers.CredentialController
}

type credentialResourceModel struct {
	ID                      types.String  `tfsdk:"id"`
	Name                    types.String  `tfsdk:"name"`
	Type                    types.String  `tfsdk:"type"`
	Data                    types.Dynamic `tfsdk:"data"`
	DataVersion             types.Int64   `tfsdk:"data_version"`
	IsPartialData           types.Bool    `tfsdk:"is_partial_data"`
	ProjectID               types.String  `tfsdk:"project_id"`
	IsResolvable            types.Bool    `tfsdk:"is_resolvable"`
	IsGlobal                types.Bool    `tfsdk:"is_global"`
	DeleteProtection        types.Bool    `tfsdk:"delete_protection"`
	IsManaged               types.Bool    `tfsdk:"is_managed"`
	ResolvableAllowFallback types.Bool    `tfsdk:"resolvable_allow_fallback"`
	ResolverID              types.String  `tfsdk:"resolver_id"`
	CreatedAt               types.String  `tfsdk:"created_at"`
	UpdatedAt               types.String  `tfsdk:"updated_at"`
}

func NewCredentialResource() resource.Resource {
	return &credentialResource{}
}

func (r *credentialResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_credential"
}

func (r *credentialResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	markdownDescription, err := readMarkdownDescription(ctx, "internal/provider/docs/resources/credential.md")
	if err != nil {
		resp.Diagnostics.AddError("Unable to read markdown description", err.Error())
		return
	}

	resp.Schema = schema.Schema{
		MarkdownDescription: markdownDescription,
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Credential ID assigned by n8n.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Credential name.",
			},
			"type": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "n8n credential type name (for example `httpHeaderAuth`). Changing this forces replacement.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"data": schema.DynamicAttribute{
				Required:            true,
				WriteOnly:           true,
				Sensitive:           true,
				MarkdownDescription: "Credential payload for this type. Write-only; never stored in state. Change `data_version` to push an update.",
			},
			"data_version": schema.Int64Attribute{
				Required:            true,
				MarkdownDescription: "Keeper for write-only `data`. Create always sends `data`. Update sends `data` only when this value changes.",
			},
			"is_partial_data": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
				MarkdownDescription: "When true, n8n merges `data` into the stored secret object on update. When false, `data` replaces the entire object. Sent only when `data_version` changes.",
			},
			"project_id": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Project that owns the credential. Omit to use the API key owner's personal project. Changing a previously set value transfers the credential; setting it for the first time after import adopts without transfer.",
			},
			"is_resolvable": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "Whether this credential has resolvable fields.",
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"is_global": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "Whether this credential is available globally. Applied after create via update. Community n8n returns 403 when set to true.",
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"delete_protection": schema.BoolAttribute{
				Required:            true,
				MarkdownDescription: "When set to `true`, prevents Terraform from destroying this credential. This flag is Terraform-only; n8n has no matching API field. Imported resources default to `true`.",
			},
			"is_managed": schema.BoolAttribute{
				Computed:            true,
				MarkdownDescription: "Whether n8n manages this credential. Managed credentials cannot be edited or deleted via the API.",
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"resolvable_allow_fallback": schema.BoolAttribute{
				Computed:            true,
				MarkdownDescription: "Whether a credential resolver may fall back to static data.",
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"resolver_id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Dynamic credential resolver id, if any.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
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
		},
	}
}

func (r *credentialResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
	r.credentialController = controllers.NewCredentialController(client)
}

func (r *credentialResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan, config credentialResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	name := strings.TrimSpace(plan.Name.ValueString())
	if name == "" {
		resp.Diagnostics.AddAttributeError(path.Root("name"), "Invalid Credential Name", "name must not be empty")
		return
	}
	typ := strings.TrimSpace(plan.Type.ValueString())
	if typ == "" {
		resp.Diagnostics.AddAttributeError(path.Root("type"), "Invalid Credential Type", "type must not be empty")
		return
	}
	data, err := dynamicObjectToMap(config.Data)
	if err != nil {
		resp.Diagnostics.AddAttributeError(path.Root("data"), "Invalid Credential Data", err.Error())
		return
	}

	created, err := r.credentialController.Create(ctx, controllers.CreateCredentialOptions{
		Name:         name,
		Type:         typ,
		Data:         data,
		ProjectID:    optionalStringPointer(plan.ProjectID),
		IsResolvable: optionalBoolPointer(plan.IsResolvable),
	})
	if err != nil {
		resp.Diagnostics.AddError("Error creating credential", err.Error())
		return
	}

	if wantGlobal := optionalBoolPointer(plan.IsGlobal); wantGlobal != nil && *wantGlobal != created.IsGlobal {
		created, err = r.credentialController.Update(ctx, controllers.UpdateCredentialOptions{
			ID:       created.ID,
			Name:     created.Name,
			IsGlobal: wantGlobal,
		})
		if err != nil {
			resp.Diagnostics.AddError("Error setting credential is_global after create", err.Error())
			return
		}
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, credentialResourceModelFromAPI(created, plan))...)
}

func (r *credentialResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state credentialResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	got, err := r.credentialController.Get(ctx, state.ID.ValueString())
	if err != nil {
		if n8n.IsNotFound(err) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Error reading credential", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, credentialResourceModelFromAPI(got, state))...)
}

func (r *credentialResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state, config credentialResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	name := strings.TrimSpace(plan.Name.ValueString())
	if name == "" {
		resp.Diagnostics.AddAttributeError(path.Root("name"), "Invalid Credential Name", "name must not be empty")
		return
	}

	sendData := plan.DataVersion.ValueInt64() != state.DataVersion.ValueInt64()
	var data map[string]any
	if sendData {
		var err error
		data, err = dynamicObjectToMap(config.Data)
		if err != nil {
			resp.Diagnostics.AddAttributeError(path.Root("data"), "Invalid Credential Data", err.Error())
			return
		}
	}

	nameChanged := name != strings.TrimSpace(state.Name.ValueString())
	projectChanged := !optionalStringsEqual(plan.ProjectID, state.ProjectID)
	resolvableChanged := !plan.IsResolvable.Equal(state.IsResolvable)
	globalChanged := !plan.IsGlobal.Equal(state.IsGlobal)

	if !nameChanged && !sendData && !projectChanged && !resolvableChanged && !globalChanged {
		state.DeleteProtection = plan.DeleteProtection
		state.DataVersion = plan.DataVersion
		state.IsPartialData = plan.IsPartialData
		resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
		return
	}

	var isResolvable, isGlobal *bool
	if resolvableChanged {
		isResolvable = optionalBoolPointer(plan.IsResolvable)
	}
	if globalChanged {
		isGlobal = optionalBoolPointer(plan.IsGlobal)
	}

	updated, err := r.credentialController.Update(ctx, controllers.UpdateCredentialOptions{
		ID:            plan.ID.ValueString(),
		Name:          name,
		Data:          data,
		SendData:      sendData,
		IsPartialData: plan.IsPartialData.ValueBool(),
		ProjectID:     optionalStringPointer(plan.ProjectID),
		PriorProject:  optionalStringPointer(state.ProjectID),
		IsResolvable:  isResolvable,
		IsGlobal:      isGlobal,
	})
	if err != nil {
		resp.Diagnostics.AddError("Error updating credential", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, credentialResourceModelFromAPI(updated, plan))...)
}

func (r *credentialResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state credentialResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.credentialController.Delete(ctx, controllers.DeleteCredentialOptions{
		ID:               state.ID.ValueString(),
		DeleteProtection: state.DeleteProtection.ValueBool(),
	})
	if err != nil {
		resp.Diagnostics.AddError("Error deleting credential", err.Error())
		return
	}
}

func (r *credentialResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	got, err := r.credentialController.Import(ctx, controllers.ImportCredentialOptions{ID: req.ID})
	if err != nil {
		resp.Diagnostics.AddError("Error importing credential", err.Error())
		return
	}
	plan := credentialResourceModel{
		DataVersion:      types.Int64Value(1),
		IsPartialData:    types.BoolValue(false),
		DeleteProtection: types.BoolValue(true),
		ProjectID:        types.StringNull(),
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, credentialResourceModelFromAPI(got, plan))...)
}

func credentialResourceModelFromAPI(c *models.Credential, persist credentialResourceModel) credentialResourceModel {
	resolver := types.StringNull()
	if c.ResolverID != nil && strings.TrimSpace(*c.ResolverID) != "" {
		resolver = types.StringValue(*c.ResolverID)
	}
	return credentialResourceModel{
		ID:                      types.StringValue(c.ID),
		Name:                    types.StringValue(c.Name),
		Type:                    types.StringValue(c.Type),
		Data:                    types.DynamicNull(),
		DataVersion:             persist.DataVersion,
		IsPartialData:           persist.IsPartialData,
		ProjectID:               persist.ProjectID,
		IsResolvable:            types.BoolValue(c.IsResolvable),
		IsGlobal:                types.BoolValue(c.IsGlobal),
		DeleteProtection:        persist.DeleteProtection,
		IsManaged:               types.BoolValue(c.IsManaged),
		ResolvableAllowFallback: types.BoolValue(c.ResolvableAllowFallback),
		ResolverID:              resolver,
		CreatedAt:               types.StringValue(c.CreatedAt),
		UpdatedAt:               types.StringValue(c.UpdatedAt),
	}
}
