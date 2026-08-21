package credentials

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/ubie-oss/terraform-provider-n8n/internal/n8n"
)

// GetCredentialSchemaV1 calls GET /api/v1/credentials/schema/{credentialTypeName}.
// Success is 200 with a JSON Schema-like object. Unknown types are 404.
func GetCredentialSchemaV1(ctx context.Context, c *n8n.Client, typeName string) (json.RawMessage, error) {
	typeName = strings.TrimSpace(typeName)
	if typeName == "" {
		return nil, fmt.Errorf("credential type is empty")
	}
	var out json.RawMessage
	if err := c.DoJSON(ctx, http.MethodGet, c.URL("credentials", "schema", typeName), nil, &out, "credential schema", typeName); err != nil {
		return nil, err
	}
	if len(out) == 0 {
		return nil, fmt.Errorf("get credential schema returned empty body")
	}
	return out, nil
}
