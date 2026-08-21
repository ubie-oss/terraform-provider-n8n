package tags

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/ubie-oss/terraform-provider-n8n/internal/n8n"
	"github.com/ubie-oss/terraform-provider-n8n/internal/n8n/models"
)

// UpdateTagV1 calls PUT /api/v1/tags/{id}. Success is 200.
func UpdateTagV1(ctx context.Context, c *n8n.Client, id string, in models.TagWrite) (*models.Tag, error) {
	id = strings.TrimSpace(id)
	if id == "" {
		return nil, fmt.Errorf("tag id is empty")
	}
	in.Name = strings.TrimSpace(in.Name)
	if in.Name == "" {
		return nil, fmt.Errorf("tag name is empty")
	}
	var out models.Tag
	if err := c.DoJSON(ctx, http.MethodPut, c.URL("tags", id), in, &out, tagResource, id); err != nil {
		return nil, err
	}
	if out.ID == "" {
		return nil, fmt.Errorf("update tag returned empty id")
	}
	return &out, nil
}
