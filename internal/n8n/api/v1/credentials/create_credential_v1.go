package credentials

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/ubie-oss/terraform-provider-n8n/internal/n8n"
	"github.com/ubie-oss/terraform-provider-n8n/internal/n8n/models"
)

const credentialResource = "credential"

// CreateCredentialV1 calls POST /api/v1/credentials. Success is 200 with metadata (no data).
func CreateCredentialV1(ctx context.Context, c *n8n.Client, in models.CredentialWrite) (*models.Credential, error) {
	in.Name = strings.TrimSpace(in.Name)
	in.Type = strings.TrimSpace(in.Type)
	if in.Name == "" {
		return nil, fmt.Errorf("credential name is empty")
	}
	if in.Type == "" {
		return nil, fmt.Errorf("credential type is empty")
	}
	if in.Data == nil {
		return nil, fmt.Errorf("credential data is empty")
	}
	var out models.Credential
	if err := c.DoJSON(ctx, http.MethodPost, c.URL("credentials"), in, &out, credentialResource, ""); err != nil {
		return nil, err
	}
	if out.ID == "" {
		return nil, fmt.Errorf("create credential returned empty id")
	}
	return &out, nil
}
