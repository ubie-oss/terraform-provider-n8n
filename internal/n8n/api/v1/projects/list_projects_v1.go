package projects

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	"github.com/ubie-oss/terraform-provider-n8n/internal/n8n"
	"github.com/ubie-oss/terraform-provider-n8n/internal/n8n/models"
)

const (
	// DefaultProjectListLimit is n8n's documented maximum page size for GET /projects.
	DefaultProjectListLimit = 250
	projectResource         = "project"
)

// ListProjectsV1 calls GET /api/v1/projects?limit=&cursor=.
func ListProjectsV1(ctx context.Context, c *n8n.Client, limit int, cursor string) (*models.ProjectList, error) {
	if limit <= 0 {
		limit = DefaultProjectListLimit
	}
	u, err := url.Parse(c.URL("projects"))
	if err != nil {
		return nil, fmt.Errorf("parse projects URL: %w", err)
	}
	q := u.Query()
	q.Set("limit", strconv.Itoa(limit))
	if cursor != "" {
		q.Set("cursor", cursor)
	}
	u.RawQuery = q.Encode()

	var out models.ProjectList
	if err := c.DoJSON(ctx, http.MethodGet, u.String(), nil, &out, projectResource, ""); err != nil {
		return nil, err
	}
	return &out, nil
}
