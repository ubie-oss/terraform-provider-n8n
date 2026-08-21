package folders

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/ubie-oss/terraform-provider-n8n/internal/n8n"
	"github.com/ubie-oss/terraform-provider-n8n/internal/n8n/models"
)

// DeleteFolderV1 calls DELETE /api/v1/projects/{projectId}/folders/{folderId}. Success is 204.
func DeleteFolderV1(ctx context.Context, c *n8n.Client, projectID, folderID string, q models.DeleteFolderQuery) error {
	projectID, folderID, err := requireProjectFolderIDs(projectID, folderID)
	if err != nil {
		return err
	}
	u, err := url.Parse(c.URL("projects", projectID, "folders", folderID))
	if err != nil {
		return fmt.Errorf("parse folder URL: %w", err)
	}
	if q.TransferToFolderID != nil {
		id := strings.TrimSpace(*q.TransferToFolderID)
		if id != "" {
			query := u.Query()
			query.Set("transferToFolderId", id)
			u.RawQuery = query.Encode()
		}
	}
	return c.DoJSON(ctx, http.MethodDelete, u.String(), nil, nil, folderResource, folderID)
}
