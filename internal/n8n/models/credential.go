package models

import (
	"encoding/json"
	"fmt"
)

// Credential is an n8n Public API credential document (GET-by-id / create / PATCH).
// Secret `data` is never returned. List items omit several of these flags; see CredentialListItem.
type Credential struct {
	ID                      string  `json:"id"`
	Name                    string  `json:"name"`
	Type                    string  `json:"type"`
	IsManaged               bool    `json:"isManaged"`
	IsGlobal                bool    `json:"isGlobal"`
	IsResolvable            bool    `json:"isResolvable"`
	ResolvableAllowFallback bool    `json:"resolvableAllowFallback"`
	ResolverID              *string `json:"resolverId"`
	CreatedAt               string  `json:"createdAt,omitempty"`
	UpdatedAt               string  `json:"updatedAt,omitempty"`
	UsageScope              string  `json:"usageScope,omitempty"`
}

// CredentialShared is a project-sharing row from GET /credentials list items.
type CredentialShared struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Role      string `json:"role"`
	CreatedAt string `json:"createdAt,omitempty"`
	UpdatedAt string `json:"updatedAt,omitempty"`
}

// CredentialListItem is one element of GET /credentials. Live n8n omits isManaged and related flags.
type CredentialListItem struct {
	ID        string             `json:"id"`
	Name      string             `json:"name"`
	Type      string             `json:"type"`
	CreatedAt string             `json:"createdAt,omitempty"`
	UpdatedAt string             `json:"updatedAt,omitempty"`
	Shared    []CredentialShared `json:"shared"`
}

// CredentialList is the GET /credentials cursor envelope.
type CredentialList struct {
	Data       []CredentialListItem `json:"data"`
	NextCursor *string              `json:"nextCursor"`
}

// CredentialWrite is the request body for POST /credentials.
type CredentialWrite struct {
	Name         string         `json:"name"`
	Type         string         `json:"type"`
	Data         map[string]any `json:"data"`
	ProjectID    *string        `json:"projectId,omitempty"`
	IsResolvable *bool          `json:"isResolvable,omitempty"`
}

// CredentialUpdate is the request body for PATCH /credentials/{id}.
// SendData controls whether `data` (and optional isPartialData) are marshaled.
type CredentialUpdate struct {
	Name          *string
	Data          map[string]any
	SendData      bool
	IsPartialData *bool
	IsGlobal      *bool
	IsResolvable  *bool
}

// MarshalJSON encodes a PATCH body. At least one field must be set.
func (u CredentialUpdate) MarshalJSON() ([]byte, error) {
	m := map[string]any{}
	if u.Name != nil {
		m["name"] = *u.Name
	}
	if u.SendData {
		data := u.Data
		if data == nil {
			data = map[string]any{}
		}
		m["data"] = data
		if u.IsPartialData != nil {
			m["isPartialData"] = *u.IsPartialData
		}
	}
	if u.IsGlobal != nil {
		m["isGlobal"] = *u.IsGlobal
	}
	if u.IsResolvable != nil {
		m["isResolvable"] = *u.IsResolvable
	}
	if len(m) == 0 {
		return nil, fmt.Errorf("credential update is empty")
	}
	return json.Marshal(m)
}

// CredentialTransfer is the request body for PUT /credentials/{id}/transfer.
type CredentialTransfer struct {
	DestinationProjectID string `json:"destinationProjectId"`
}
