package folders

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	"github.com/ubie-oss/terraform-provider-n8n/internal/n8n"
	"github.com/ubie-oss/terraform-provider-n8n/internal/n8n/models"
)

const (
	// DefaultFolderListTake is n8n's documented default page size for GET folders.
	DefaultFolderListTake = 10
	// ServiceFolderListTake is the page size FolderService.ListAll uses (larger than the default 10).
	ServiceFolderListTake = 100
)

// ListFoldersV1 calls GET /api/v1/projects/{projectId}/folders?skip=&take= (values encoded as decimal strings).
func ListFoldersV1(ctx context.Context, c *n8n.Client, projectID string, skip, take int) (*models.FolderList, error) {
	projectID, err := requireProjectID(projectID)
	if err != nil {
		return nil, err
	}
	if skip < 0 {
		skip = 0
	}
	if take <= 0 {
		take = DefaultFolderListTake
	}
	u, err := url.Parse(c.URL("projects", projectID, "folders"))
	if err != nil {
		return nil, fmt.Errorf("parse folders URL: %w", err)
	}
	q := u.Query()
	q.Set("skip", strconv.Itoa(skip))
	q.Set("take", strconv.Itoa(take))
	u.RawQuery = q.Encode()

	var out models.FolderList
	if err := c.DoJSON(ctx, http.MethodGet, u.String(), nil, &out, folderResource, ""); err != nil {
		return nil, err
	}
	return &out, nil
}
