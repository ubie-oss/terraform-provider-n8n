package folders

import (
	"context"
	"fmt"
	"net/http"

	"github.com/ubie-oss/terraform-provider-n8n/internal/n8n"
	"github.com/ubie-oss/terraform-provider-n8n/internal/n8n/models"
)

// GetFolderV1 calls GET /api/v1/projects/{projectId}/folders/{folderId}. Success is 200.
func GetFolderV1(ctx context.Context, c *n8n.Client, projectID, folderID string) (*models.Folder, error) {
	projectID, folderID, err := requireProjectFolderIDs(projectID, folderID)
	if err != nil {
		return nil, err
	}
	var out models.Folder
	if err := c.DoJSON(ctx, http.MethodGet, c.URL("projects", projectID, "folders", folderID), nil, &out, folderResource, folderID); err != nil {
		return nil, err
	}
	if out.ID == "" {
		return nil, fmt.Errorf("get folder returned empty id")
	}
	return &out, nil
}
