package services

import (
	"context"
	"fmt"
	"strings"

	"github.com/ubie-oss/terraform-provider-n8n/internal/n8n"
	tagsv1 "github.com/ubie-oss/terraform-provider-n8n/internal/n8n/api/v1/tags"
	"github.com/ubie-oss/terraform-provider-n8n/internal/n8n/models"
)

// TagService implements tag operations against GET-by-id and cursor list.
type TagService struct {
	client *n8n.Client
}

// NewTagService returns a TagService using c.
func NewTagService(client *n8n.Client) *TagService {
	return &TagService{client: client}
}

// Create creates a tag.
func (s *TagService) Create(ctx context.Context, in models.TagWrite) (*models.Tag, error) {
	return tagsv1.CreateTagV1(ctx, s.client, in)
}

// Get returns a tag by id.
func (s *TagService) Get(ctx context.Context, id string) (*models.Tag, error) {
	return tagsv1.GetTagV1(ctx, s.client, id)
}

// ListAll walks GET /tags until nextCursor is empty.
func (s *TagService) ListAll(ctx context.Context) ([]models.Tag, error) {
	var all []models.Tag
	cursor := ""
	seen := make(map[string]struct{})
	for page := 0; page < maxListPages; page++ {
		if cursor != "" {
			if _, ok := seen[cursor]; ok {
				return nil, fmt.Errorf("tag list pagination repeated cursor %q", cursor)
			}
			seen[cursor] = struct{}{}
		}
		list, err := tagsv1.ListTagsV1(ctx, s.client, tagsv1.DefaultTagListLimit, cursor)
		if err != nil {
			return nil, err
		}
		all = append(all, list.Data...)
		if list.NextCursor == nil || strings.TrimSpace(*list.NextCursor) == "" {
			return all, nil
		}
		cursor = *list.NextCursor
	}
	return nil, fmt.Errorf("tag list exceeded %d pages", maxListPages)
}

// Update renames a tag via PUT.
func (s *TagService) Update(ctx context.Context, id string, in models.TagWrite) (*models.Tag, error) {
	return tagsv1.UpdateTagV1(ctx, s.client, id, in)
}

// Delete removes a tag. Callers treat n8n.NotFoundError as already gone.
func (s *TagService) Delete(ctx context.Context, id string) error {
	return tagsv1.DeleteTagV1(ctx, s.client, id)
}
