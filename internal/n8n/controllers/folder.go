package controllers

import (
	"context"
	"fmt"
	"sort"
	"strings"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/ubie-oss/terraform-provider-n8n/internal/n8n"
	"github.com/ubie-oss/terraform-provider-n8n/internal/n8n/models"
	"github.com/ubie-oss/terraform-provider-n8n/internal/n8n/services"
)

const personalProjectPath = "personal"

// FolderController is the Terraform-facing orchestrator for n8n folders.
type FolderController struct {
	folders  *services.FolderService
	projects *ProjectController
}

// NewFolderController builds a FolderController around client.
func NewFolderController(client *n8n.Client) *FolderController {
	return &FolderController{
		folders:  services.NewFolderService(client),
		projects: NewProjectController(client),
	}
}

// CreateFolderOptions is the input for creating a folder in a team project.
type CreateFolderOptions struct {
	ProjectID      string
	Name           string
	ParentFolderID *string
}

// UpdateFolderOptions is the input for renaming or moving a folder.
type UpdateFolderOptions struct {
	ProjectID      string
	FolderID       string
	Name           string
	ParentFolderID *string
	MoveToRoot     bool
}

// DeleteFolderOptions is the input for deleting a folder.
type DeleteFolderOptions struct {
	ProjectID          string
	FolderID           string
	DeleteProtection   bool
	TransferToFolderID *string
}

// ImportFolderOptions is the input for importing a folder by project and folder id.
type ImportFolderOptions struct {
	ProjectID string
	FolderID  string
}

// ListFoldersOptions filters a project's folders after walking skip/take pagination.
type ListFoldersOptions struct {
	ProjectID      string
	ParentFolderID *string
	Name           string
}

// Create creates a folder in a team project. Path alias "personal" is rejected.
func (c *FolderController) Create(ctx context.Context, options CreateFolderOptions) (*models.Folder, error) {
	name := strings.TrimSpace(options.Name)
	tflog.Debug(ctx, "(FolderController.Create) creating folder", map[string]any{
		"projectId": options.ProjectID,
		"name":      name,
	})
	if name == "" {
		return nil, fmt.Errorf("folder name is empty")
	}
	if err := c.requireTeamProject(ctx, options.ProjectID); err != nil {
		return nil, err
	}
	return c.folders.Create(ctx, options.ProjectID, models.FolderWrite{
		Name:           name,
		ParentFolderID: options.ParentFolderID,
	})
}

// Get returns a folder of any parent project type, including personal.
func (c *FolderController) Get(ctx context.Context, projectID, folderID string) (*models.Folder, error) {
	tflog.Debug(ctx, "(FolderController.Get) getting folder", map[string]any{
		"projectId": projectID,
		"folderId":  folderID,
	})
	return c.folders.Get(ctx, projectID, folderID)
}

// List returns folders in a project, optionally filtered by parent and exact name.
func (c *FolderController) List(ctx context.Context, options ListFoldersOptions) ([]models.Folder, error) {
	tflog.Debug(ctx, "(FolderController.List) listing folders", map[string]any{"projectId": options.ProjectID})
	all, err := c.folders.ListAll(ctx, options.ProjectID)
	if err != nil {
		return nil, err
	}
	name := strings.TrimSpace(options.Name)
	var out []models.Folder
	for _, f := range all {
		if name != "" && f.Name != name {
			continue
		}
		if options.ParentFolderID != nil {
			want := strings.TrimSpace(*options.ParentFolderID)
			got := ""
			if f.ParentFolderID != nil {
				got = *f.ParentFolderID
			}
			if want != got {
				continue
			}
		}
		out = append(out, f)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].ID < out[j].ID })
	return out, nil
}

// Update patches a folder in a team project.
func (c *FolderController) Update(ctx context.Context, options UpdateFolderOptions) (*models.Folder, error) {
	name := strings.TrimSpace(options.Name)
	tflog.Debug(ctx, "(FolderController.Update) updating folder", map[string]any{
		"projectId": options.ProjectID,
		"folderId":  options.FolderID,
		"name":      name,
	})
	if name == "" {
		return nil, fmt.Errorf("folder name is empty")
	}
	if err := c.requireTeamProject(ctx, options.ProjectID); err != nil {
		return nil, err
	}
	in := models.FolderUpdate{Name: &name}
	if options.MoveToRoot {
		in.MoveToRoot = true
	} else if options.ParentFolderID != nil {
		in.ParentFolderID = options.ParentFolderID
	}
	return c.folders.Update(ctx, options.ProjectID, options.FolderID, in)
}

// Delete removes a folder in a team project. Missing folders are treated as already gone.
func (c *FolderController) Delete(ctx context.Context, options DeleteFolderOptions) error {
	tflog.Debug(ctx, "(FolderController.Delete) deleting folder", map[string]any{
		"projectId":        options.ProjectID,
		"folderId":         options.FolderID,
		"deleteProtection": options.DeleteProtection,
	})
	if options.DeleteProtection {
		return errDeleteProtected("folder", options.FolderID)
	}
	if err := c.requireTeamProject(ctx, options.ProjectID); err != nil {
		return err
	}
	err := c.folders.Delete(ctx, options.ProjectID, options.FolderID, models.DeleteFolderQuery{
		TransferToFolderID: options.TransferToFolderID,
	})
	if n8n.IsNotFound(err) {
		return nil
	}
	return err
}

// Import loads a folder in a team project for terraform import.
func (c *FolderController) Import(ctx context.Context, options ImportFolderOptions) (*models.Folder, error) {
	tflog.Debug(ctx, "(FolderController.Import) importing folder", map[string]any{
		"projectId": options.ProjectID,
		"folderId":  options.FolderID,
	})
	if err := c.requireTeamProject(ctx, options.ProjectID); err != nil {
		return nil, err
	}
	return c.folders.Get(ctx, options.ProjectID, options.FolderID)
}

func (c *FolderController) requireTeamProject(ctx context.Context, projectID string) error {
	if strings.TrimSpace(projectID) == personalProjectPath {
		return fmt.Errorf("n8n personal projects cannot be managed by this resource (id %q); only team projects are supported", personalProjectPath)
	}
	_, err := c.projects.GetTeam(ctx, projectID)
	return err
}
