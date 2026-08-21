package projects

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/ubie-oss/terraform-provider-n8n/internal/n8n"
)

// DeleteProjectV1 calls DELETE /api/v1/projects/{projectId}. Success is 204. 404 is NotFoundError.
func DeleteProjectV1(ctx context.Context, c *n8n.Client, id string) error {
	id = strings.TrimSpace(id)
	if id == "" {
		return fmt.Errorf("project id is empty")
	}
	return c.DoJSON(ctx, http.MethodDelete, c.URL("projects", id), nil, nil, projectResource, id)
}
