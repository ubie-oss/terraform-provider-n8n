package services

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/ubie-oss/terraform-provider-n8n/internal/n8n"
	projectsv1 "github.com/ubie-oss/terraform-provider-n8n/internal/n8n/api/v1/projects"
	"github.com/ubie-oss/terraform-provider-n8n/internal/n8n/models"
)

// ErrProjectNotFound is returned when GetByID exhausts GET /projects without a match.
// The Public API has no GET /projects/{id}.
var ErrProjectNotFound = errors.New("project not found")

const maxListPages = 1000

// ProjectService implements project operations, including Read-via-list.
type ProjectService struct {
	client *n8n.Client
}

// NewProjectService returns a ProjectService using c.
func NewProjectService(client *n8n.Client) *ProjectService {
	return &ProjectService{client: client}
}

// ListAll walks GET /projects until nextCursor is empty.
func (s *ProjectService) ListAll(ctx context.Context) ([]models.Project, error) {
	var all []models.Project
	cursor := ""
	seen := make(map[string]struct{})
	for page := 0; page < maxListPages; page++ {
		if cursor != "" {
			if _, ok := seen[cursor]; ok {
				return nil, fmt.Errorf("project list pagination repeated cursor %q", cursor)
			}
			seen[cursor] = struct{}{}
		}
		list, err := projectsv1.ListProjectsV1(ctx, s.client, projectsv1.DefaultProjectListLimit, cursor)
		if err != nil {
			return nil, err
		}
		all = append(all, list.Data...)
		if list.NextCursor == nil || strings.TrimSpace(*list.NextCursor) == "" {
			return all, nil
		}
		cursor = *list.NextCursor
	}
	return nil, fmt.Errorf("project list exceeded %d pages", maxListPages)
}

// GetByID scans ListAll for id. Missing projects return ErrProjectNotFound.
func (s *ProjectService) GetByID(ctx context.Context, id string) (*models.Project, error) {
	id = strings.TrimSpace(id)
	if id == "" {
		return nil, fmt.Errorf("project id is empty")
	}
	all, err := s.ListAll(ctx)
	if err != nil {
		return nil, err
	}
	for i := range all {
		if all[i].ID == id {
			p := all[i]
			return &p, nil
		}
	}
	return nil, fmt.Errorf("%w: %q", ErrProjectNotFound, id)
}

// Create creates a team project by name.
func (s *ProjectService) Create(ctx context.Context, name string) (*models.Project, error) {
	return projectsv1.CreateProjectV1(ctx, s.client, models.ProjectWrite{Name: strings.TrimSpace(name)})
}

// Update renames a project (PUT 204) then refreshes via GetByID.
func (s *ProjectService) Update(ctx context.Context, id, name string) (*models.Project, error) {
	if err := projectsv1.UpdateProjectV1(ctx, s.client, id, models.ProjectWrite{Name: strings.TrimSpace(name)}); err != nil {
		return nil, err
	}
	return s.GetByID(ctx, id)
}

// Delete removes a project. Callers treat n8n.NotFoundError as already gone.
func (s *ProjectService) Delete(ctx context.Context, id string) error {
	return projectsv1.DeleteProjectV1(ctx, s.client, id)
}
