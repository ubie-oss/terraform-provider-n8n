package tags

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/ubie-oss/terraform-provider-n8n/internal/n8n"
)

// DeleteTagV1 calls DELETE /api/v1/tags/{id}. Success is 200. 404 is NotFoundError.
func DeleteTagV1(ctx context.Context, c *n8n.Client, id string) error {
	id = strings.TrimSpace(id)
	if id == "" {
		return fmt.Errorf("tag id is empty")
	}
	return c.DoJSON(ctx, http.MethodDelete, c.URL("tags", id), nil, nil, tagResource, id)
}
