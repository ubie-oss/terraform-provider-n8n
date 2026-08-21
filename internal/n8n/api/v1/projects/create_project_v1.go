package projects

import (
	"context"
	"fmt"
	"net/http"

	"github.com/ubie-oss/terraform-provider-n8n/internal/n8n"
	"github.com/ubie-oss/terraform-provider-n8n/internal/n8n/models"
)

// CreateProjectV1 calls POST /api/v1/projects. Live n8n returns 201 with an unwrapped project body.
func CreateProjectV1(ctx context.Context, c *n8n.Client, in models.ProjectWrite) (*models.Project, error) {
	if in.Name == "" {
		return nil, fmt.Errorf("project name is empty")
	}
	var out models.Project
	if err := c.DoJSON(ctx, http.MethodPost, c.URL("projects"), in, &out, projectResource, ""); err != nil {
		return nil, err
	}
	if out.ID == "" {
		return nil, fmt.Errorf("create project returned empty id")
	}
	return &out, nil
}
