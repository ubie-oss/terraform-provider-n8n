package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"sort"
	"strings"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/ubie-oss/terraform-provider-n8n/internal/n8n"
	"github.com/ubie-oss/terraform-provider-n8n/internal/n8n/models"
	"github.com/ubie-oss/terraform-provider-n8n/internal/n8n/services"
)

// CredentialController is the Terraform-facing orchestrator for n8n credentials.
type CredentialController struct {
	credentials *services.CredentialService
}

// NewCredentialController builds a CredentialController around client.
func NewCredentialController(client *n8n.Client) *CredentialController {
	return &CredentialController{
		credentials: services.NewCredentialService(client),
	}
}

// CreateCredentialOptions is the input for creating a credential.
type CreateCredentialOptions struct {
	Name         string
	Type         string
	Data         map[string]any
	ProjectID    *string
	IsResolvable *bool
}

// UpdateCredentialOptions is the input for patching a credential.
type UpdateCredentialOptions struct {
	ID            string
	Name          string
	Data          map[string]any
	SendData      bool
	IsPartialData bool
	ProjectID     *string
	PriorProject  *string
	IsResolvable  *bool
	IsGlobal      *bool
}

// DeleteCredentialOptions is the input for deleting a credential.
type DeleteCredentialOptions struct {
	ID               string
	DeleteProtection bool
}

// ImportCredentialOptions is the input for importing a credential by id.
type ImportCredentialOptions struct {
	ID string
}

// ListCredentialsOptions filters credentials after walking cursor pagination.
type ListCredentialsOptions struct {
	Name string
	Type string
}

// Create creates a credential. Data is required.
func (c *CredentialController) Create(ctx context.Context, options CreateCredentialOptions) (*models.Credential, error) {
	name := strings.TrimSpace(options.Name)
	typ := strings.TrimSpace(options.Type)
	tflog.Debug(ctx, "(CredentialController.Create) creating credential", map[string]any{
		"name": name,
		"type": typ,
	})
	if name == "" {
		return nil, fmt.Errorf("credential name is empty")
	}
	if typ == "" {
		return nil, fmt.Errorf("credential type is empty")
	}
	if options.Data == nil {
		return nil, fmt.Errorf("credential data is empty")
	}
	return c.credentials.Create(ctx, models.CredentialWrite{
		Name:         name,
		Type:         typ,
		Data:         options.Data,
		ProjectID:    trimOptionalString(options.ProjectID),
		IsResolvable: options.IsResolvable,
	})
}

// Get returns credential metadata by id. Secrets are never included.
func (c *CredentialController) Get(ctx context.Context, id string) (*models.Credential, error) {
	tflog.Debug(ctx, "(CredentialController.Get) getting credential", map[string]any{"id": id})
	return c.credentials.Get(ctx, id)
}

// List returns credentials, optionally filtered by exact name and type.
func (c *CredentialController) List(ctx context.Context, options ListCredentialsOptions) ([]models.CredentialListItem, error) {
	tflog.Debug(ctx, "(CredentialController.List) listing credentials")
	all, err := c.credentials.ListAll(ctx)
	if err != nil {
		return nil, err
	}
	name := strings.TrimSpace(options.Name)
	typ := strings.TrimSpace(options.Type)
	var out []models.CredentialListItem
	for _, item := range all {
		if name != "" && item.Name != name {
			continue
		}
		if typ != "" && item.Type != typ {
			continue
		}
		out = append(out, item)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].ID < out[j].ID })
	return out, nil
}

// Update patches name/flags and optionally data. Managed credentials cannot be edited.
// Changing project_id transfers the credential; clearing project_id is rejected.
func (c *CredentialController) Update(ctx context.Context, options UpdateCredentialOptions) (*models.Credential, error) {
	name := strings.TrimSpace(options.Name)
	tflog.Debug(ctx, "(CredentialController.Update) updating credential", map[string]any{
		"id":       options.ID,
		"name":     name,
		"sendData": options.SendData,
	})
	got, err := c.credentials.Get(ctx, options.ID)
	if err != nil {
		return nil, err
	}
	if got.IsManaged {
		return nil, fmt.Errorf("credential %s is managed by n8n and cannot be edited via the API", options.ID)
	}
	if name == "" {
		return nil, fmt.Errorf("credential name is empty")
	}

	prior := trimOptionalString(options.PriorProject)
	next := trimOptionalString(options.ProjectID)
	if prior != nil && next == nil {
		return nil, fmt.Errorf("cannot clear project_id on credential %s; transfer to personal is unsupported", options.ID)
	}
	transferred := false
	if projectIDChanged(prior, next) {
		if err := c.credentials.Transfer(ctx, options.ID, *next); err != nil {
			return nil, err
		}
		transferred = true
	}

	in := models.CredentialUpdate{Name: &name, IsResolvable: options.IsResolvable, IsGlobal: options.IsGlobal}
	if options.SendData {
		in.SendData = true
		in.Data = options.Data
		partial := options.IsPartialData
		in.IsPartialData = &partial
	}
	updated, err := c.credentials.Update(ctx, options.ID, in)
	if err != nil {
		if transferred {
			return nil, fmt.Errorf("credential %s was transferred to project %s, but the follow-up update failed: %w", options.ID, *next, err)
		}
		return nil, err
	}
	return updated, nil
}

// Delete removes a credential. Missing credentials are treated as already gone.
func (c *CredentialController) Delete(ctx context.Context, options DeleteCredentialOptions) error {
	tflog.Debug(ctx, "(CredentialController.Delete) deleting credential", map[string]any{
		"id":               options.ID,
		"deleteProtection": options.DeleteProtection,
	})
	if options.DeleteProtection {
		return errDeleteProtected("credential", options.ID)
	}
	got, err := c.credentials.Get(ctx, options.ID)
	if err != nil {
		if n8n.IsNotFound(err) {
			return nil
		}
		return err
	}
	if got.IsManaged {
		return fmt.Errorf("credential %s is managed by n8n and cannot be deleted via the API", options.ID)
	}
	err = c.credentials.Delete(ctx, options.ID)
	if n8n.IsNotFound(err) {
		return nil
	}
	return err
}

// Import loads a credential by id for terraform import.
func (c *CredentialController) Import(ctx context.Context, options ImportCredentialOptions) (*models.Credential, error) {
	tflog.Debug(ctx, "(CredentialController.Import) importing credential", map[string]any{"id": options.ID})
	return c.credentials.Get(ctx, options.ID)
}

// Schema returns the JSON Schema-like body for a credential type.
func (c *CredentialController) Schema(ctx context.Context, typeName string) (json.RawMessage, error) {
	tflog.Debug(ctx, "(CredentialController.Schema) fetching credential schema", map[string]any{"type": typeName})
	return c.credentials.Schema(ctx, typeName)
}

func trimOptionalString(v *string) *string {
	if v == nil {
		return nil
	}
	s := strings.TrimSpace(*v)
	if s == "" {
		return nil
	}
	return &s
}

// projectIDChanged reports whether ownership should move via Transfer.
// A nil prior means ownership is unknown (import or never set in Terraform);
// adopting a project_id into state must not Transfer to a project the credential
// may already belong to.
func projectIDChanged(prior, next *string) bool {
	if prior == nil || next == nil {
		return false
	}
	return *prior != *next
}
