package tags

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
	// DefaultTagListLimit is n8n's documented maximum page size for GET /tags.
	DefaultTagListLimit = 250
)

// ListTagsV1 calls GET /api/v1/tags?limit=&cursor=.
func ListTagsV1(ctx context.Context, c *n8n.Client, limit int, cursor string) (*models.TagList, error) {
	if limit <= 0 {
		limit = DefaultTagListLimit
	}
	u, err := url.Parse(c.URL("tags"))
	if err != nil {
		return nil, fmt.Errorf("parse tags URL: %w", err)
	}
	q := u.Query()
	q.Set("limit", strconv.Itoa(limit))
	if cursor != "" {
		q.Set("cursor", cursor)
	}
	u.RawQuery = q.Encode()

	var out models.TagList
	if err := c.DoJSON(ctx, http.MethodGet, u.String(), nil, &out, tagResource, ""); err != nil {
		return nil, err
	}
	return &out, nil
}
