package credentials

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/ubie-oss/terraform-provider-n8n/internal/n8n"
	"github.com/ubie-oss/terraform-provider-n8n/internal/n8n/models"
)

// UpdateCredentialV1 calls PATCH /api/v1/credentials/{id}. Success is 200 with metadata (no data).
func UpdateCredentialV1(ctx context.Context, c *n8n.Client, id string, in models.CredentialUpdate) (*models.Credential, error) {
	id = strings.TrimSpace(id)
	if id == "" {
		return nil, fmt.Errorf("credential id is empty")
	}
	var out models.Credential
	if err := c.DoJSON(ctx, http.MethodPatch, c.URL("credentials", id), in, &out, credentialResource, id); err != nil {
		return nil, err
	}
	if out.ID == "" {
		return nil, fmt.Errorf("update credential returned empty id")
	}
	return &out, nil
}
