package credentials

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
	// DefaultCredentialListLimit is n8n's documented maximum page size for GET /credentials.
	DefaultCredentialListLimit = 250
)

// ListCredentialsV1 calls GET /api/v1/credentials?limit=&cursor=.
func ListCredentialsV1(ctx context.Context, c *n8n.Client, limit int, cursor string) (*models.CredentialList, error) {
	if limit <= 0 {
		limit = DefaultCredentialListLimit
	}
	u, err := url.Parse(c.URL("credentials"))
	if err != nil {
		return nil, fmt.Errorf("parse credentials URL: %w", err)
	}
	q := u.Query()
	q.Set("limit", strconv.Itoa(limit))
	if cursor != "" {
		q.Set("cursor", cursor)
	}
	u.RawQuery = q.Encode()

	var out models.CredentialList
	if err := c.DoJSON(ctx, http.MethodGet, u.String(), nil, &out, credentialResource, ""); err != nil {
		return nil, err
	}
	return &out, nil
}
