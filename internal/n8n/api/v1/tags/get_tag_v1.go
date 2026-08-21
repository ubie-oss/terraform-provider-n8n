package tags

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/ubie-oss/terraform-provider-n8n/internal/n8n"
	"github.com/ubie-oss/terraform-provider-n8n/internal/n8n/models"
)

// GetTagV1 calls GET /api/v1/tags/{id}. Success is 200.
func GetTagV1(ctx context.Context, c *n8n.Client, id string) (*models.Tag, error) {
	id = strings.TrimSpace(id)
	if id == "" {
		return nil, fmt.Errorf("tag id is empty")
	}
	var out models.Tag
	if err := c.DoJSON(ctx, http.MethodGet, c.URL("tags", id), nil, &out, tagResource, id); err != nil {
		return nil, err
	}
	if out.ID == "" {
		return nil, fmt.Errorf("get tag returned empty id")
	}
	return &out, nil
}
