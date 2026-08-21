package credentials

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/ubie-oss/terraform-provider-n8n/internal/n8n"
	"github.com/ubie-oss/terraform-provider-n8n/internal/n8n/models"
)

// GetCredentialV1 calls GET /api/v1/credentials/{id}. Success is 200 metadata (no data).
func GetCredentialV1(ctx context.Context, c *n8n.Client, id string) (*models.Credential, error) {
	id = strings.TrimSpace(id)
	if id == "" {
		return nil, fmt.Errorf("credential id is empty")
	}
	var out models.Credential
	if err := c.DoJSON(ctx, http.MethodGet, c.URL("credentials", id), nil, &out, credentialResource, id); err != nil {
		return nil, err
	}
	if out.ID == "" {
		return nil, fmt.Errorf("get credential returned empty id")
	}
	return &out, nil
}
