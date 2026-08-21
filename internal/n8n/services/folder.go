package services

import (
	"context"
	"fmt"

	"github.com/ubie-oss/terraform-provider-n8n/internal/n8n"
	foldersv1 "github.com/ubie-oss/terraform-provider-n8n/internal/n8n/api/v1/folders"
	"github.com/ubie-oss/terraform-provider-n8n/internal/n8n/models"
)

// FolderService implements folder operations against GET-by-id and skip/take list.
type FolderService struct {
	client *n8n.Client
}

// NewFolderService returns a FolderService using c.
func NewFolderService(client *n8n.Client) *FolderService {
	return &FolderService{client: client}
}

// Create creates a folder in projectID.
func (s *FolderService) Create(ctx context.Context, projectID string, in models.FolderWrite) (*models.Folder, error) {
	return foldersv1.CreateFolderV1(ctx, s.client, projectID, in)
}

// Get returns a folder by project and folder id.
func (s *FolderService) Get(ctx context.Context, projectID, folderID string) (*models.Folder, error) {
	return foldersv1.GetFolderV1(ctx, s.client, projectID, folderID)
}

// ListAll walks GET folders with skip/take until skip >= count or a page is empty.
func (s *FolderService) ListAll(ctx context.Context, projectID string) ([]models.Folder, error) {
	var all []models.Folder
	skip := 0
	for page := 0; page < maxListPages; page++ {
		list, err := foldersv1.ListFoldersV1(ctx, s.client, projectID, skip, foldersv1.ServiceFolderListTake)
		if err != nil {
			return nil, err
		}
		all = append(all, list.Data...)
		skip += len(list.Data)
		if len(list.Data) == 0 || skip >= list.Count {
			return all, nil
		}
	}
	return nil, fmt.Errorf("folder list exceeded %d pages", maxListPages)
}

// Update patches name and/or parent folder.
func (s *FolderService) Update(ctx context.Context, projectID, folderID string, in models.FolderUpdate) (*models.Folder, error) {
	return foldersv1.UpdateFolderV1(ctx, s.client, projectID, folderID, in)
}

// Delete removes a folder. Callers treat n8n.NotFoundError as already gone.
func (s *FolderService) Delete(ctx context.Context, projectID, folderID string, q models.DeleteFolderQuery) error {
	return foldersv1.DeleteFolderV1(ctx, s.client, projectID, folderID, q)
}
