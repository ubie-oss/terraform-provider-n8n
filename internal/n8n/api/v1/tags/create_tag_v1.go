package tags

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/ubie-oss/terraform-provider-n8n/internal/n8n"
	"github.com/ubie-oss/terraform-provider-n8n/internal/n8n/models"
)

const tagResource = "tag"

// CreateTagV1 calls POST /api/v1/tags. Success is 201.
func CreateTagV1(ctx context.Context, c *n8n.Client, in models.TagWrite) (*models.Tag, error) {
	in.Name = strings.TrimSpace(in.Name)
	if in.Name == "" {
		return nil, fmt.Errorf("tag name is empty")
	}
	var out models.Tag
	if err := c.DoJSON(ctx, http.MethodPost, c.URL("tags"), in, &out, tagResource, ""); err != nil {
		return nil, err
	}
	if out.ID == "" {
		return nil, fmt.Errorf("create tag returned empty id")
	}
	return &out, nil
}
