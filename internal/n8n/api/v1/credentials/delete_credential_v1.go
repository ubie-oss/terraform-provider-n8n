package credentials

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/ubie-oss/terraform-provider-n8n/internal/n8n"
)

// DeleteCredentialV1 calls DELETE /api/v1/credentials/{id}. Success is 200. 404 is NotFoundError.
func DeleteCredentialV1(ctx context.Context, c *n8n.Client, id string) error {
	id = strings.TrimSpace(id)
	if id == "" {
		return fmt.Errorf("credential id is empty")
	}
	return c.DoJSON(ctx, http.MethodDelete, c.URL("credentials", id), nil, nil, credentialResource, id)
}
