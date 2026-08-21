package projects

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/ubie-oss/terraform-provider-n8n/internal/n8n"
	"github.com/ubie-oss/terraform-provider-n8n/internal/n8n/models"
)

// UpdateProjectV1 calls PUT /api/v1/projects/{projectId}. Success is 204 with an empty body.
func UpdateProjectV1(ctx context.Context, c *n8n.Client, id string, in models.ProjectWrite) error {
	id = strings.TrimSpace(id)
	if id == "" {
		return fmt.Errorf("project id is empty")
	}
	if in.Name == "" {
		return fmt.Errorf("project name is empty")
	}
	return c.DoJSON(ctx, http.MethodPut, c.URL("projects", id), in, nil, projectResource, id)
}
