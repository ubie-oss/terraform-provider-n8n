package folders

import (
	"context"
	"fmt"
	"net/http"

	"github.com/ubie-oss/terraform-provider-n8n/internal/n8n"
	"github.com/ubie-oss/terraform-provider-n8n/internal/n8n/models"
)

// UpdateFolderV1 calls PATCH /api/v1/projects/{projectId}/folders/{folderId}. Success is 200 with a folder body.
func UpdateFolderV1(ctx context.Context, c *n8n.Client, projectID, folderID string, in models.FolderUpdate) (*models.Folder, error) {
	projectID, folderID, err := requireProjectFolderIDs(projectID, folderID)
	if err != nil {
		return nil, err
	}
	var out models.Folder
	if err := c.DoJSON(ctx, http.MethodPatch, c.URL("projects", projectID, "folders", folderID), in, &out, folderResource, folderID); err != nil {
		return nil, err
	}
	if out.ID == "" {
		return nil, fmt.Errorf("update folder returned empty id")
	}
	return &out, nil
}
