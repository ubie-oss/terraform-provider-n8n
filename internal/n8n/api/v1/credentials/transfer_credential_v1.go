package credentials

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/ubie-oss/terraform-provider-n8n/internal/n8n"
	"github.com/ubie-oss/terraform-provider-n8n/internal/n8n/models"
)

// TransferCredentialV1 calls PUT /api/v1/credentials/{id}/transfer.
func TransferCredentialV1(ctx context.Context, c *n8n.Client, id, destinationProjectID string) error {
	id = strings.TrimSpace(id)
	destinationProjectID = strings.TrimSpace(destinationProjectID)
	if id == "" {
		return fmt.Errorf("credential id is empty")
	}
	if destinationProjectID == "" {
		return fmt.Errorf("destination project id is empty")
	}
	return c.DoJSON(ctx, http.MethodPut, c.URL("credentials", id, "transfer"), models.CredentialTransfer{
		DestinationProjectID: destinationProjectID,
	}, nil, credentialResource, id)
}
