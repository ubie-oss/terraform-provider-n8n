package folders

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/ubie-oss/terraform-provider-n8n/internal/n8n"
	"github.com/ubie-oss/terraform-provider-n8n/internal/n8n/models"
)

const folderResource = "folder"

// CreateFolderV1 calls POST /api/v1/projects/{projectId}/folders. Success is 201 with an unwrapped folder body.
func CreateFolderV1(ctx context.Context, c *n8n.Client, projectID string, in models.FolderWrite) (*models.Folder, error) {
	projectID, err := requireProjectID(projectID)
	if err != nil {
		return nil, err
	}
	if strings.TrimSpace(in.Name) == "" {
		return nil, fmt.Errorf("folder name is empty")
	}
	in.Name = strings.TrimSpace(in.Name)
	var out models.Folder
	if err := c.DoJSON(ctx, http.MethodPost, c.URL("projects", projectID, "folders"), in, &out, folderResource, ""); err != nil {
		return nil, err
	}
	if out.ID == "" {
		return nil, fmt.Errorf("create folder returned empty id")
	}
	return &out, nil
}
