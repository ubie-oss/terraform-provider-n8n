package controllers

import (
	"context"
	"errors"
	"fmt"
	"sort"
	"strings"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/ubie-oss/terraform-provider-n8n/internal/n8n"
	"github.com/ubie-oss/terraform-provider-n8n/internal/n8n/models"
	"github.com/ubie-oss/terraform-provider-n8n/internal/n8n/services"
)

// ErrProjectNotFound is returned when a project id is not present in GET /projects.
var ErrProjectNotFound = services.ErrProjectNotFound

// ProjectController is the Terraform-facing orchestrator for n8n projects.
// Resources and data sources call this layer; it uses ProjectService for HTTP.
type ProjectController struct {
	projects *services.ProjectService
}

// NewProjectController builds a ProjectController around client.
func NewProjectController(client *n8n.Client) *ProjectController {
	return &ProjectController{
		projects: services.NewProjectService(client),
	}
}

// CreateProjectOptions is the input for creating a team project.
type CreateProjectOptions struct {
	Name string
}

// UpdateProjectOptions is the input for renaming a team project.
type UpdateProjectOptions struct {
	ID   string
	Name string
}

// DeleteProjectOptions is the input for deleting a team project.
type DeleteProjectOptions struct {
	ID               string
	DeleteProtection bool
}

// ImportProjectOptions is the input for importing a team project by id.
type ImportProjectOptions struct {
	ID string
}

// Create creates a team project. Personal projects are rejected.
func (c *ProjectController) Create(ctx context.Context, options CreateProjectOptions) (*models.Project, error) {
	name := strings.TrimSpace(options.Name)
	tflog.Debug(ctx, "(ProjectController.Create) creating project", map[string]any{"name": name})
	if name == "" {
		return nil, fmt.Errorf("project name is empty")
	}
	created, err := c.projects.Create(ctx, name)
	if err != nil {
		return nil, err
	}
	if err := errIfNotTeamProject(created); err != nil {
		return nil, err
	}
	return created, nil
}

// Get returns a project of any type (team or personal). Missing ids return ErrProjectNotFound.
func (c *ProjectController) Get(ctx context.Context, id string) (*models.Project, error) {
	tflog.Debug(ctx, "(ProjectController.Get) getting project", map[string]any{"id": id})
	return c.projects.GetByID(ctx, id)
}

// GetTeam returns a team project. Personal projects are rejected.
func (c *ProjectController) GetTeam(ctx context.Context, id string) (*models.Project, error) {
	tflog.Debug(ctx, "(ProjectController.GetTeam) getting team project", map[string]any{"id": id})
	got, err := c.projects.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if err := errIfNotTeamProject(got); err != nil {
		return nil, err
	}
	return got, nil
}

// List returns every project visible to the API key, sorted by id.
func (c *ProjectController) List(ctx context.Context) ([]models.Project, error) {
	tflog.Debug(ctx, "(ProjectController.List) listing projects")
	all, err := c.projects.ListAll(ctx)
	if err != nil {
		return nil, err
	}
	sort.Slice(all, func(i, j int) bool { return all[i].ID < all[j].ID })
	return all, nil
}

// Update renames a team project and refreshes it.
func (c *ProjectController) Update(ctx context.Context, options UpdateProjectOptions) (*models.Project, error) {
	name := strings.TrimSpace(options.Name)
	tflog.Debug(ctx, "(ProjectController.Update) updating project", map[string]any{"id": options.ID, "name": name})
	if name == "" {
		return nil, fmt.Errorf("project name is empty")
	}
	if _, err := c.GetTeam(ctx, options.ID); err != nil {
		return nil, err
	}
	updated, err := c.projects.Update(ctx, options.ID, name)
	if err != nil {
		return nil, err
	}
	if err := errIfNotTeamProject(updated); err != nil {
		return nil, err
	}
	return updated, nil
}

// Delete removes a team project. Missing projects are treated as already gone.
func (c *ProjectController) Delete(ctx context.Context, options DeleteProjectOptions) error {
	tflog.Debug(ctx, "(ProjectController.Delete) deleting project", map[string]any{
		"id":               options.ID,
		"deleteProtection": options.DeleteProtection,
	})
	if options.DeleteProtection {
		return errDeleteProtected("project", options.ID)
	}
	got, err := c.projects.GetByID(ctx, options.ID)
	if err != nil {
		if errors.Is(err, ErrProjectNotFound) {
			return nil
		}
		return err
	}
	if err := errIfNotTeamProject(got); err != nil {
		return fmt.Errorf("refusing to delete personal project: %w", err)
	}
	err = c.projects.Delete(ctx, options.ID)
	if n8n.IsNotFound(err) {
		return nil
	}
	return err
}

// Import loads a team project by id for terraform import.
func (c *ProjectController) Import(ctx context.Context, options ImportProjectOptions) (*models.Project, error) {
	tflog.Debug(ctx, "(ProjectController.Import) importing project", map[string]any{"id": options.ID})
	return c.GetTeam(ctx, options.ID)
}

func errIfNotTeamProject(p *models.Project) error {
	if p == nil {
		return fmt.Errorf("project is nil")
	}
	if p.IsTeam() {
		return nil
	}
	return fmt.Errorf("n8n personal projects cannot be managed by this resource (id %q, type %q); only team projects are supported", p.ID, p.Type)
}
