package services

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/ubie-oss/terraform-provider-n8n/internal/n8n"
	credentialsv1 "github.com/ubie-oss/terraform-provider-n8n/internal/n8n/api/v1/credentials"
	"github.com/ubie-oss/terraform-provider-n8n/internal/n8n/models"
)

// CredentialService implements credential operations against GET-by-id and cursor list.
type CredentialService struct {
	client *n8n.Client
}

// NewCredentialService returns a CredentialService using c.
func NewCredentialService(client *n8n.Client) *CredentialService {
	return &CredentialService{client: client}
}

// Create creates a credential.
func (s *CredentialService) Create(ctx context.Context, in models.CredentialWrite) (*models.Credential, error) {
	return credentialsv1.CreateCredentialV1(ctx, s.client, in)
}

// Get returns a credential by id. Secrets are never included.
func (s *CredentialService) Get(ctx context.Context, id string) (*models.Credential, error) {
	return credentialsv1.GetCredentialV1(ctx, s.client, id)
}

// ListAll walks GET /credentials until nextCursor is empty.
func (s *CredentialService) ListAll(ctx context.Context) ([]models.CredentialListItem, error) {
	var all []models.CredentialListItem
	cursor := ""
	seen := make(map[string]struct{})
	for page := 0; page < maxListPages; page++ {
		if cursor != "" {
			if _, ok := seen[cursor]; ok {
				return nil, fmt.Errorf("credential list pagination repeated cursor %q", cursor)
			}
			seen[cursor] = struct{}{}
		}
		list, err := credentialsv1.ListCredentialsV1(ctx, s.client, credentialsv1.DefaultCredentialListLimit, cursor)
		if err != nil {
			return nil, err
		}
		all = append(all, list.Data...)
		if list.NextCursor == nil || strings.TrimSpace(*list.NextCursor) == "" {
			return all, nil
		}
		cursor = *list.NextCursor
	}
	return nil, fmt.Errorf("credential list exceeded %d pages", maxListPages)
}

// Update patches a credential and returns the GET-shaped metadata from the PATCH body.
func (s *CredentialService) Update(ctx context.Context, id string, in models.CredentialUpdate) (*models.Credential, error) {
	return credentialsv1.UpdateCredentialV1(ctx, s.client, id, in)
}

// Transfer moves a credential to destinationProjectID.
func (s *CredentialService) Transfer(ctx context.Context, id, destinationProjectID string) error {
	return credentialsv1.TransferCredentialV1(ctx, s.client, id, destinationProjectID)
}

// Delete removes a credential. Callers treat n8n.NotFoundError as already gone.
func (s *CredentialService) Delete(ctx context.Context, id string) error {
	return credentialsv1.DeleteCredentialV1(ctx, s.client, id)
}

// Schema returns the JSON Schema-like body for a credential type.
func (s *CredentialService) Schema(ctx context.Context, typeName string) (json.RawMessage, error) {
	return credentialsv1.GetCredentialSchemaV1(ctx, s.client, typeName)
}
