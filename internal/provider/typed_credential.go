package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/ubie-oss/terraform-provider-n8n/internal/n8n"
	"github.com/ubie-oss/terraform-provider-n8n/internal/n8n/controllers"
	"github.com/ubie-oss/terraform-provider-n8n/internal/n8n/models"
)

var (
	_ resource.Resource                = &typedCredentialResource{}
	_ resource.ResourceWithConfigure   = &typedCredentialResource{}
	_ resource.ResourceWithImportState = &typedCredentialResource{}
)

type typedAttrKind int

const (
	typedAttrString typedAttrKind = iota
	typedAttrBool
	typedAttrInt64
	typedAttrDynamic
)

type typedCredentialSpec struct {
	TerraformSuffix     string
	N8nType             string
	OAuthPartialDefault bool
	DocFile             string
	ExtraAttributes     map[string]schema.Attribute
	BuildData           func(typedAttrBag) (map[string]any, error)
}

type typedCredentialResource struct {
	spec                 typedCredentialSpec
	secretKeys           []string
	stateKeys            []string
	kinds                map[string]typedAttrKind
	credentialController *controllers.CredentialController
}

type typedAttrBag struct {
	strings  map[string]types.String
	bools    map[string]types.Bool
	int64s   map[string]types.Int64
	dynamics map[string]types.Dynamic
}

type attributeReader interface {
	GetAttribute(ctx context.Context, attrPath path.Path, target interface{}) diag.Diagnostics
}

func newTypedCredentialResource(spec typedCredentialSpec) resource.Resource {
	r := &typedCredentialResource{
		spec:  spec,
		kinds: map[string]typedAttrKind{},
	}
	for name, attrSchema := range spec.ExtraAttributes {
		kind, writeOnly, ok := classifyTypedAttribute(attrSchema)
		if !ok {
			panic(fmt.Sprintf("typed credential %q: unsupported schema type for attribute %q", spec.TerraformSuffix, name))
		}
		r.kinds[name] = kind
		if writeOnly {
			r.secretKeys = append(r.secretKeys, name)
		} else {
			r.stateKeys = append(r.stateKeys, name)
		}
	}
	return r
}

func classifyTypedAttribute(a schema.Attribute) (typedAttrKind, bool, bool) {
	switch v := a.(type) {
	case schema.StringAttribute:
		return typedAttrString, v.WriteOnly, true
	case schema.BoolAttribute:
		return typedAttrBool, v.WriteOnly, true
	case schema.Int64Attribute:
		return typedAttrInt64, v.WriteOnly, true
	case schema.DynamicAttribute:
		return typedAttrDynamic, v.WriteOnly, true
	default:
		return 0, false, false
	}
}

func (r *typedCredentialResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + r.spec.TerraformSuffix
}

func (r *typedCredentialResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	markdownDescription, err := readMarkdownDescription(ctx, r.spec.DocFile)
	if err != nil {
		resp.Diagnostics.AddError("Unable to read markdown description", err.Error())
		return
	}

	attrs := typedCredentialEnvelopeAttributes(r.spec.N8nType, r.spec.OAuthPartialDefault)
	for name, extra := range r.spec.ExtraAttributes {
		if _, exists := attrs[name]; exists {
			resp.Diagnostics.AddError(
				"Invalid typed credential spec",
				fmt.Sprintf("extra attribute %q collides with the shared credential envelope", name),
			)
			return
		}
		if _, _, ok := classifyTypedAttribute(extra); !ok {
			resp.Diagnostics.AddError(
				"Invalid typed credential spec",
				fmt.Sprintf("extra attribute %q has an unsupported schema type", name),
			)
			return
		}
		attrs[name] = extra
	}

	resp.Schema = schema.Schema{
		MarkdownDescription: markdownDescription,
		Attributes:          attrs,
	}
}

func typedCredentialEnvelopeAttributes(n8nType string, oauthPartialDefault bool) map[string]schema.Attribute {
	return map[string]schema.Attribute{
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
			Computed:            true,
			Default:             stringdefault.StaticString(n8nType),
			MarkdownDescription: fmt.Sprintf("n8n credential type (`%s`).", n8nType),
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
			},
		},
		"data_version": schema.Int64Attribute{
			Required:            true,
			MarkdownDescription: "Keeper for write-only secrets. Create always sends credential data. Update sends data when this value changes or when a non-secret payload attribute changes.",
		},
		"is_partial_data": schema.BoolAttribute{
			Optional:            true,
			Computed:            true,
			Default:             booldefault.StaticBool(oauthPartialDefault),
			MarkdownDescription: "When true, n8n merges payload fields into the stored secret object on update. When false, the payload replaces the entire object. OAuth types default to true so omitted oauth_token_data does not wipe UI-obtained tokens.",
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
	}
}

func (r *typedCredentialResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *typedCredentialResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	name, diags := readStringAttr(ctx, req.Plan, "name")
	resp.Diagnostics.Append(diags...)
	projectID, diags := readStringAttr(ctx, req.Plan, "project_id")
	resp.Diagnostics.Append(diags...)
	isResolvable, diags := readBoolAttr(ctx, req.Plan, "is_resolvable")
	resp.Diagnostics.Append(diags...)
	isGlobal, diags := readBoolAttr(ctx, req.Plan, "is_global")
	resp.Diagnostics.Append(diags...)
	dataVersion, diags := readInt64Attr(ctx, req.Plan, "data_version")
	resp.Diagnostics.Append(diags...)
	isPartialData, diags := readBoolAttr(ctx, req.Plan, "is_partial_data")
	resp.Diagnostics.Append(diags...)
	deleteProtection, diags := readBoolAttr(ctx, req.Plan, "delete_protection")
	resp.Diagnostics.Append(diags...)
	bag, diags := r.collectBag(ctx, req.Config, req.Plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	trimmed := strings.TrimSpace(name.ValueString())
	if trimmed == "" {
		resp.Diagnostics.AddAttributeError(path.Root("name"), "Invalid Credential Name", "name must not be empty")
		return
	}
	data, err := r.spec.BuildData(bag)
	if err != nil {
		resp.Diagnostics.AddError("Invalid Credential Data", err.Error())
		return
	}

	created, err := r.credentialController.Create(ctx, controllers.CreateCredentialOptions{
		Name:         trimmed,
		Type:         r.spec.N8nType,
		Data:         data,
		ProjectID:    optionalStringPointer(projectID),
		IsResolvable: optionalBoolPointer(isResolvable),
	})
	if err != nil {
		resp.Diagnostics.AddError("Error creating credential", err.Error())
		return
	}

	if wantGlobal := optionalBoolPointer(isGlobal); wantGlobal != nil && *wantGlobal != created.IsGlobal {
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

	resp.Diagnostics.Append(r.writeState(ctx, &resp.State, created, typedStatePersist{
		DataVersion:      dataVersion,
		IsPartialData:    isPartialData,
		ProjectID:        projectID,
		DeleteProtection: deleteProtection,
		Bag:              bag,
	})...)
}

func (r *typedCredentialResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	id, diags := readStringAttr(ctx, req.State, "id")
	resp.Diagnostics.Append(diags...)
	dataVersion, diags := readInt64Attr(ctx, req.State, "data_version")
	resp.Diagnostics.Append(diags...)
	isPartialData, diags := readBoolAttr(ctx, req.State, "is_partial_data")
	resp.Diagnostics.Append(diags...)
	projectID, diags := readStringAttr(ctx, req.State, "project_id")
	resp.Diagnostics.Append(diags...)
	deleteProtection, diags := readBoolAttr(ctx, req.State, "delete_protection")
	resp.Diagnostics.Append(diags...)
	bag, diags := r.readBag(ctx, req.State, r.stateKeys)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	got, err := r.credentialController.Get(ctx, id.ValueString())
	if err != nil {
		if n8n.IsNotFound(err) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Error reading credential", err.Error())
		return
	}

	resp.Diagnostics.Append(r.writeState(ctx, &resp.State, got, typedStatePersist{
		DataVersion:      dataVersion,
		IsPartialData:    isPartialData,
		ProjectID:        projectID,
		DeleteProtection: deleteProtection,
		Bag:              bag,
	})...)
}

func (r *typedCredentialResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	id, diags := readStringAttr(ctx, req.Plan, "id")
	resp.Diagnostics.Append(diags...)
	name, diags := readStringAttr(ctx, req.Plan, "name")
	resp.Diagnostics.Append(diags...)
	planProject, diags := readStringAttr(ctx, req.Plan, "project_id")
	resp.Diagnostics.Append(diags...)
	stateProject, diags := readStringAttr(ctx, req.State, "project_id")
	resp.Diagnostics.Append(diags...)
	planResolvable, diags := readBoolAttr(ctx, req.Plan, "is_resolvable")
	resp.Diagnostics.Append(diags...)
	stateResolvable, diags := readBoolAttr(ctx, req.State, "is_resolvable")
	resp.Diagnostics.Append(diags...)
	planGlobal, diags := readBoolAttr(ctx, req.Plan, "is_global")
	resp.Diagnostics.Append(diags...)
	stateGlobal, diags := readBoolAttr(ctx, req.State, "is_global")
	resp.Diagnostics.Append(diags...)
	planVersion, diags := readInt64Attr(ctx, req.Plan, "data_version")
	resp.Diagnostics.Append(diags...)
	stateVersion, diags := readInt64Attr(ctx, req.State, "data_version")
	resp.Diagnostics.Append(diags...)
	isPartialData, diags := readBoolAttr(ctx, req.Plan, "is_partial_data")
	resp.Diagnostics.Append(diags...)
	deleteProtection, diags := readBoolAttr(ctx, req.Plan, "delete_protection")
	resp.Diagnostics.Append(diags...)
	planBag, diags := r.collectBag(ctx, req.Config, req.Plan)
	resp.Diagnostics.Append(diags...)
	stateBag, diags := r.readBag(ctx, req.State, r.stateKeys)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	trimmed := strings.TrimSpace(name.ValueString())
	if trimmed == "" {
		resp.Diagnostics.AddAttributeError(path.Root("name"), "Invalid Credential Name", "name must not be empty")
		return
	}

	sendData := !planVersion.Equal(stateVersion) || !stateKeysEqual(r.stateKeys, r.kinds, planBag, stateBag)
	var data map[string]any
	if sendData {
		var err error
		data, err = r.spec.BuildData(planBag)
		if err != nil {
			resp.Diagnostics.AddError("Invalid Credential Data", err.Error())
			return
		}
	}

	stateName, diags := readStringAttr(ctx, req.State, "name")
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	nameChanged := trimmed != strings.TrimSpace(stateName.ValueString())
	projectChanged := !optionalStringsEqual(planProject, stateProject)
	resolvableChanged := !planResolvable.Equal(stateResolvable)
	globalChanged := !planGlobal.Equal(stateGlobal)

	if !nameChanged && !sendData && !projectChanged && !resolvableChanged && !globalChanged {
		diags := diag.Diagnostics{}
		diags.Append(resp.State.SetAttribute(ctx, path.Root("data_version"), planVersion)...)
		diags.Append(resp.State.SetAttribute(ctx, path.Root("is_partial_data"), isPartialData)...)
		diags.Append(resp.State.SetAttribute(ctx, path.Root("project_id"), planProject)...)
		diags.Append(resp.State.SetAttribute(ctx, path.Root("delete_protection"), deleteProtection)...)
		for _, key := range r.stateKeys {
			diags.Append(setBagAttr(ctx, &resp.State, key, r.kinds[key], planBag)...)
		}
		resp.Diagnostics.Append(diags...)
		return
	}

	var isResolvable, isGlobal *bool
	if resolvableChanged {
		isResolvable = optionalBoolPointer(planResolvable)
	}
	if globalChanged {
		isGlobal = optionalBoolPointer(planGlobal)
	}

	updated, err := r.credentialController.Update(ctx, controllers.UpdateCredentialOptions{
		ID:            id.ValueString(),
		Name:          trimmed,
		Data:          data,
		SendData:      sendData,
		IsPartialData: isPartialData.ValueBool(),
		ProjectID:     optionalStringPointer(planProject),
		PriorProject:  optionalStringPointer(stateProject),
		IsResolvable:  isResolvable,
		IsGlobal:      isGlobal,
	})
	if err != nil {
		resp.Diagnostics.AddError("Error updating credential", err.Error())
		return
	}

	resp.Diagnostics.Append(r.writeState(ctx, &resp.State, updated, typedStatePersist{
		DataVersion:      planVersion,
		IsPartialData:    isPartialData,
		ProjectID:        planProject,
		DeleteProtection: deleteProtection,
		Bag:              planBag,
	})...)
}

func (r *typedCredentialResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	id, diags := readStringAttr(ctx, req.State, "id")
	resp.Diagnostics.Append(diags...)
	deleteProtection, diags := readBoolAttr(ctx, req.State, "delete_protection")
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.credentialController.Delete(ctx, controllers.DeleteCredentialOptions{
		ID:               id.ValueString(),
		DeleteProtection: deleteProtection.ValueBool(),
	})
	if err != nil {
		resp.Diagnostics.AddError("Error deleting credential", err.Error())
	}
}

func (r *typedCredentialResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	got, err := r.credentialController.Import(ctx, controllers.ImportCredentialOptions{ID: req.ID})
	if err != nil {
		resp.Diagnostics.AddError("Error importing credential", err.Error())
		return
	}
	if got.Type != r.spec.N8nType {
		resp.Diagnostics.AddError(
			"Credential type mismatch",
			fmt.Sprintf("imported credential %s has type %q, expected %q for this resource", req.ID, got.Type, r.spec.N8nType),
		)
		return
	}

	resp.Diagnostics.Append(r.writeState(ctx, &resp.State, got, typedStatePersist{
		DataVersion:      types.Int64Value(1),
		IsPartialData:    types.BoolValue(r.spec.OAuthPartialDefault),
		ProjectID:        types.StringNull(),
		DeleteProtection: types.BoolValue(true),
		Bag:              typedAttrBag{},
	})...)
}

type typedStatePersist struct {
	DataVersion      types.Int64
	IsPartialData    types.Bool
	ProjectID        types.String
	DeleteProtection types.Bool
	Bag              typedAttrBag
}

func (r *typedCredentialResource) writeState(ctx context.Context, state *tfsdk.State, c *models.Credential, persist typedStatePersist) diag.Diagnostics {
	var diags diag.Diagnostics
	resolver := types.StringNull()
	if c.ResolverID != nil && strings.TrimSpace(*c.ResolverID) != "" {
		resolver = types.StringValue(*c.ResolverID)
	}
	diags.Append(state.SetAttribute(ctx, path.Root("id"), types.StringValue(c.ID))...)
	diags.Append(state.SetAttribute(ctx, path.Root("name"), types.StringValue(c.Name))...)
	diags.Append(state.SetAttribute(ctx, path.Root("type"), types.StringValue(c.Type))...)
	diags.Append(state.SetAttribute(ctx, path.Root("data_version"), persist.DataVersion)...)
	diags.Append(state.SetAttribute(ctx, path.Root("is_partial_data"), persist.IsPartialData)...)
	diags.Append(state.SetAttribute(ctx, path.Root("project_id"), persist.ProjectID)...)
	diags.Append(state.SetAttribute(ctx, path.Root("is_resolvable"), types.BoolValue(c.IsResolvable))...)
	diags.Append(state.SetAttribute(ctx, path.Root("is_global"), types.BoolValue(c.IsGlobal))...)
	diags.Append(state.SetAttribute(ctx, path.Root("delete_protection"), persist.DeleteProtection)...)
	diags.Append(state.SetAttribute(ctx, path.Root("is_managed"), types.BoolValue(c.IsManaged))...)
	diags.Append(state.SetAttribute(ctx, path.Root("resolvable_allow_fallback"), types.BoolValue(c.ResolvableAllowFallback))...)
	diags.Append(state.SetAttribute(ctx, path.Root("resolver_id"), resolver)...)
	diags.Append(state.SetAttribute(ctx, path.Root("created_at"), types.StringValue(c.CreatedAt))...)
	diags.Append(state.SetAttribute(ctx, path.Root("updated_at"), types.StringValue(c.UpdatedAt))...)
	for _, key := range r.stateKeys {
		diags.Append(setBagAttr(ctx, state, key, r.kinds[key], persist.Bag)...)
	}
	for _, key := range r.secretKeys {
		diags.Append(setNullSecret(ctx, state, key, r.kinds[key])...)
	}
	return diags
}

func setBagAttr(ctx context.Context, state *tfsdk.State, key string, kind typedAttrKind, bag typedAttrBag) diag.Diagnostics {
	switch kind {
	case typedAttrString:
		v, ok := bag.strings[key]
		if !ok {
			v = types.StringNull()
		}
		return state.SetAttribute(ctx, path.Root(key), v)
	case typedAttrBool:
		v, ok := bag.bools[key]
		if !ok {
			v = types.BoolNull()
		}
		return state.SetAttribute(ctx, path.Root(key), v)
	case typedAttrInt64:
		v, ok := bag.int64s[key]
		if !ok {
			v = types.Int64Null()
		}
		return state.SetAttribute(ctx, path.Root(key), v)
	case typedAttrDynamic:
		v, ok := bag.dynamics[key]
		if !ok {
			v = types.DynamicNull()
		}
		return state.SetAttribute(ctx, path.Root(key), v)
	default:
		return nil
	}
}

func setNullSecret(ctx context.Context, state *tfsdk.State, key string, kind typedAttrKind) diag.Diagnostics {
	switch kind {
	case typedAttrString:
		return state.SetAttribute(ctx, path.Root(key), types.StringNull())
	case typedAttrBool:
		return state.SetAttribute(ctx, path.Root(key), types.BoolNull())
	case typedAttrInt64:
		return state.SetAttribute(ctx, path.Root(key), types.Int64Null())
	case typedAttrDynamic:
		return state.SetAttribute(ctx, path.Root(key), types.DynamicNull())
	default:
		return nil
	}
}

func (r *typedCredentialResource) collectBag(ctx context.Context, config, plan attributeReader) (typedAttrBag, diag.Diagnostics) {
	var diags diag.Diagnostics
	bag, secretDiags := r.readBag(ctx, config, r.secretKeys)
	diags.Append(secretDiags...)
	stateBag, stateDiags := r.readBag(ctx, plan, r.stateKeys)
	diags.Append(stateDiags...)
	mergeTypedAttrBag(&bag, stateBag)
	return bag, diags
}

func (r *typedCredentialResource) readBag(ctx context.Context, src attributeReader, keys []string) (typedAttrBag, diag.Diagnostics) {
	var diags diag.Diagnostics
	bag := newTypedAttrBag()
	if src == nil {
		return bag, diags
	}
	for _, key := range keys {
		kind := r.kinds[key]
		switch kind {
		case typedAttrString:
			v, d := readStringAttr(ctx, src, key)
			diags.Append(d...)
			bag.strings[key] = v
		case typedAttrBool:
			v, d := readBoolAttr(ctx, src, key)
			diags.Append(d...)
			bag.bools[key] = v
		case typedAttrInt64:
			v, d := readInt64Attr(ctx, src, key)
			diags.Append(d...)
			bag.int64s[key] = v
		case typedAttrDynamic:
			v, d := readDynamicAttr(ctx, src, key)
			diags.Append(d...)
			bag.dynamics[key] = v
		}
	}
	return bag, diags
}

func newTypedAttrBag() typedAttrBag {
	return typedAttrBag{
		strings:  map[string]types.String{},
		bools:    map[string]types.Bool{},
		int64s:   map[string]types.Int64{},
		dynamics: map[string]types.Dynamic{},
	}
}

func mergeTypedAttrBag(dst *typedAttrBag, src typedAttrBag) {
	for k, v := range src.strings {
		dst.strings[k] = v
	}
	for k, v := range src.bools {
		dst.bools[k] = v
	}
	for k, v := range src.int64s {
		dst.int64s[k] = v
	}
	for k, v := range src.dynamics {
		dst.dynamics[k] = v
	}
}

func stateKeysEqual(keys []string, kinds map[string]typedAttrKind, plan, state typedAttrBag) bool {
	for _, key := range keys {
		switch kinds[key] {
		case typedAttrString:
			if !plan.strings[key].Equal(state.strings[key]) {
				return false
			}
		case typedAttrBool:
			if !plan.bools[key].Equal(state.bools[key]) {
				return false
			}
		case typedAttrInt64:
			if !plan.int64s[key].Equal(state.int64s[key]) {
				return false
			}
		case typedAttrDynamic:
			if !plan.dynamics[key].Equal(state.dynamics[key]) {
				return false
			}
		}
	}
	return true
}

func readStringAttr(ctx context.Context, src attributeReader, name string) (types.String, diag.Diagnostics) {
	var v types.String
	return v, src.GetAttribute(ctx, path.Root(name), &v)
}

func readBoolAttr(ctx context.Context, src attributeReader, name string) (types.Bool, diag.Diagnostics) {
	var v types.Bool
	return v, src.GetAttribute(ctx, path.Root(name), &v)
}

func readInt64Attr(ctx context.Context, src attributeReader, name string) (types.Int64, diag.Diagnostics) {
	var v types.Int64
	return v, src.GetAttribute(ctx, path.Root(name), &v)
}

func readDynamicAttr(ctx context.Context, src attributeReader, name string) (types.Dynamic, diag.Diagnostics) {
	var v types.Dynamic
	return v, src.GetAttribute(ctx, path.Root(name), &v)
}

func bagPutString(m map[string]any, apiKey string, v types.String) {
	if v.IsNull() || v.IsUnknown() {
		return
	}
	m[apiKey] = v.ValueString()
}

func bagPutBool(m map[string]any, apiKey string, v types.Bool) {
	if v.IsNull() || v.IsUnknown() {
		return
	}
	m[apiKey] = v.ValueBool()
}

func bagPutInt64(m map[string]any, apiKey string, v types.Int64) {
	if v.IsNull() || v.IsUnknown() {
		return
	}
	m[apiKey] = v.ValueInt64()
}

func bagPutRequiredString(m map[string]any, apiKey string, v types.String) error {
	if v.IsNull() || v.IsUnknown() || strings.TrimSpace(v.ValueString()) == "" {
		return fmt.Errorf("%s is required", apiKey)
	}
	m[apiKey] = v.ValueString()
	return nil
}

func bagDynamicValue(v types.Dynamic, attrName string) (any, bool, error) {
	if v.IsNull() || v.IsUnknown() {
		return nil, false, nil
	}
	if v.IsUnderlyingValueNull() || v.IsUnderlyingValueUnknown() {
		return nil, false, nil
	}
	raw, err := attrValueToGo(v.UnderlyingValue())
	if err != nil {
		return nil, false, err
	}
	if s, ok := raw.(string); ok {
		trimmed := strings.TrimSpace(s)
		if trimmed == "" {
			return nil, false, nil
		}
		var decoded any
		if err := json.Unmarshal([]byte(trimmed), &decoded); err != nil {
			return nil, false, fmt.Errorf("%s must be a JSON object: %w", attrName, err)
		}
		return decoded, true, nil
	}
	return raw, true, nil
}

func stringWriteOnly(desc string, required bool) schema.StringAttribute {
	return schema.StringAttribute{
		Required:            required,
		Optional:            !required,
		WriteOnly:           true,
		Sensitive:           true,
		MarkdownDescription: desc,
	}
}

func stringState(desc string, required bool) schema.StringAttribute {
	return schema.StringAttribute{
		Required:            required,
		Optional:            !required,
		MarkdownDescription: desc,
	}
}

func boolState(desc string) schema.BoolAttribute {
	return schema.BoolAttribute{
		Optional:            true,
		MarkdownDescription: desc,
	}
}

func int64State(desc string, required bool) schema.Int64Attribute {
	return schema.Int64Attribute{
		Required:            required,
		Optional:            !required,
		MarkdownDescription: desc,
	}
}

func dynamicWriteOnly(desc string) schema.DynamicAttribute {
	return schema.DynamicAttribute{
		Optional:            true,
		WriteOnly:           true,
		Sensitive:           true,
		MarkdownDescription: desc,
	}
}

func oauthClientAttributes() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"client_id":     stringState("OAuth client ID.", true),
		"client_secret": stringWriteOnly("OAuth client secret. Write-only; never stored in state.", true),
		"oauth_token_data": dynamicWriteOnly(
			"OAuth token payload (`oauthTokenData`), typically including access_token and refresh_token. Write-only. The Public API cannot complete a browser OAuth flow.",
		),
		"ignore_ssl_issues": boolState("Whether to ignore TLS certificate issues when talking to the token endpoint."),
	}
}

func googleOAuth2Attributes() map[string]schema.Attribute {
	attrs := oauthClientAttributes()
	attrs["custom_scopes"] = boolState("When true, send enabled_scopes instead of the n8n default scope list.")
	attrs["enabled_scopes"] = stringState("Space-separated OAuth scopes used when custom_scopes is true.", false)
	return attrs
}

func googleOAuth2BuildData(bag typedAttrBag) (map[string]any, error) {
	return oauthBuildData(bag, true, true)
}

func oauthBuildData(bag typedAttrBag, withCustomScopes bool, requireClientSecret bool) (map[string]any, error) {
	m := map[string]any{}
	if err := bagPutRequiredString(m, "clientId", bag.strings["client_id"]); err != nil {
		return nil, err
	}
	if requireClientSecret {
		if err := bagPutRequiredString(m, "clientSecret", bag.strings["client_secret"]); err != nil {
			return nil, err
		}
	} else {
		bagPutString(m, "clientSecret", bag.strings["client_secret"])
	}
	bagPutBool(m, "ignoreSSLIssues", bag.bools["ignore_ssl_issues"])
	if withCustomScopes {
		bagPutBool(m, "customScopes", bag.bools["custom_scopes"])
		bagPutString(m, "enabledScopes", bag.strings["enabled_scopes"])
	}
	token, ok, err := bagDynamicValue(bag.dynamics["oauth_token_data"], "oauth_token_data")
	if err != nil {
		return nil, err
	}
	if ok {
		m["oauthTokenData"] = token
	}
	return m, nil
}

func isCredentialManagedResource(typ string) bool {
	return typ == "n8n_credential" || strings.HasPrefix(typ, "n8n_credential_")
}
