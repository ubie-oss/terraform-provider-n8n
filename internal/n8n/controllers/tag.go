package controllers

import (
	"context"
	"fmt"
	"sort"
	"strings"
	"unicode/utf8"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/ubie-oss/terraform-provider-n8n/internal/n8n"
	"github.com/ubie-oss/terraform-provider-n8n/internal/n8n/models"
	"github.com/ubie-oss/terraform-provider-n8n/internal/n8n/services"
)

// MaxTagNameLength is the maximum number of Unicode code points n8n accepts for a tag name.
// Longer values return a misleading HTTP 409 ("Tag already exists") from the Public API.
const MaxTagNameLength = 24

// TagController is the Terraform-facing orchestrator for n8n tags.
type TagController struct {
	tags *services.TagService
}

// NewTagController builds a TagController around client.
func NewTagController(client *n8n.Client) *TagController {
	return &TagController{
		tags: services.NewTagService(client),
	}
}

// CreateTagOptions is the input for creating a tag.
type CreateTagOptions struct {
	Name string
}

// UpdateTagOptions is the input for renaming a tag.
type UpdateTagOptions struct {
	ID   string
	Name string
}

// DeleteTagOptions is the input for deleting a tag.
type DeleteTagOptions struct {
	ID               string
	DeleteProtection bool
}

// ImportTagOptions is the input for importing a tag by id.
type ImportTagOptions struct {
	ID string
}

func validateTagName(name string) error {
	if name == "" {
		return fmt.Errorf("tag name is empty")
	}
	if utf8.RuneCountInString(name) > MaxTagNameLength {
		return fmt.Errorf("tag name %q exceeds n8n's %d character limit", name, MaxTagNameLength)
	}
	return nil
}

// Create creates a tag in the global registry.
func (c *TagController) Create(ctx context.Context, options CreateTagOptions) (*models.Tag, error) {
	name := strings.TrimSpace(options.Name)
	tflog.Debug(ctx, "(TagController.Create) creating tag", map[string]any{"name": name})
	if err := validateTagName(name); err != nil {
		return nil, err
	}
	return c.tags.Create(ctx, models.TagWrite{Name: name})
}

// Get returns a tag by id.
func (c *TagController) Get(ctx context.Context, id string) (*models.Tag, error) {
	tflog.Debug(ctx, "(TagController.Get) getting tag", map[string]any{"id": id})
	return c.tags.Get(ctx, id)
}

// List returns every tag visible to the API key, sorted by id.
func (c *TagController) List(ctx context.Context) ([]models.Tag, error) {
	tflog.Debug(ctx, "(TagController.List) listing tags")
	all, err := c.tags.ListAll(ctx)
	if err != nil {
		return nil, err
	}
	sort.Slice(all, func(i, j int) bool { return all[i].ID < all[j].ID })
	return all, nil
}

// Update renames a tag via PUT, then refreshes via GET.
// PUT responses omit createdAt; GET restores the full document.
func (c *TagController) Update(ctx context.Context, options UpdateTagOptions) (*models.Tag, error) {
	name := strings.TrimSpace(options.Name)
	tflog.Debug(ctx, "(TagController.Update) updating tag", map[string]any{"id": options.ID, "name": name})
	if err := validateTagName(name); err != nil {
		return nil, err
	}
	if _, err := c.tags.Update(ctx, options.ID, models.TagWrite{Name: name}); err != nil {
		return nil, err
	}
	return c.tags.Get(ctx, options.ID)
}

// Delete removes a tag. Missing tags are treated as already gone.
func (c *TagController) Delete(ctx context.Context, options DeleteTagOptions) error {
	tflog.Debug(ctx, "(TagController.Delete) deleting tag", map[string]any{
		"id":               options.ID,
		"deleteProtection": options.DeleteProtection,
	})
	if options.DeleteProtection {
		return errDeleteProtected("tag", options.ID)
	}
	err := c.tags.Delete(ctx, options.ID)
	if n8n.IsNotFound(err) {
		return nil
	}
	return err
}

// Import loads a tag by id for terraform import.
func (c *TagController) Import(ctx context.Context, options ImportTagOptions) (*models.Tag, error) {
	tflog.Debug(ctx, "(TagController.Import) importing tag", map[string]any{"id": options.ID})
	return c.Get(ctx, options.ID)
}
