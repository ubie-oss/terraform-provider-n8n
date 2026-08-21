package models

import (
	"encoding/json"
	"fmt"
)

// Folder is an n8n Public API folder document.
// GET-by-id includes totalSubFolders and totalWorkflows; list/create/update may omit them.
type Folder struct {
	ID              string  `json:"id"`
	Name            string  `json:"name"`
	ParentFolderID  *string `json:"parentFolderId"`
	CreatedAt       string  `json:"createdAt,omitempty"`
	UpdatedAt       string  `json:"updatedAt,omitempty"`
	TotalSubFolders *int    `json:"totalSubFolders,omitempty"`
	TotalWorkflows  *int    `json:"totalWorkflows,omitempty"`
}

// FolderList is the GET /projects/{projectId}/folders envelope (skip/take pagination).
type FolderList struct {
	Count int      `json:"count"`
	Data  []Folder `json:"data"`
}

// FolderWrite is the request body for POST /projects/{projectId}/folders.
type FolderWrite struct {
	Name           string  `json:"name"`
	ParentFolderID *string `json:"parentFolderId,omitempty"`
}

// FolderUpdate is the request body for PATCH /projects/{projectId}/folders/{folderId}.
// MoveToRoot sends parentFolderId: null (omitempty would omit a nil pointer).
type FolderUpdate struct {
	Name           *string `json:"name,omitempty"`
	ParentFolderID *string `json:"parentFolderId,omitempty"`
	MoveToRoot     bool    `json:"-"`
}

// MarshalJSON encodes a PATCH body. At least one field must be set.
func (u FolderUpdate) MarshalJSON() ([]byte, error) {
	m := map[string]any{}
	if u.Name != nil {
		m["name"] = *u.Name
	}
	switch {
	case u.MoveToRoot:
		m["parentFolderId"] = nil
	case u.ParentFolderID != nil:
		m["parentFolderId"] = *u.ParentFolderID
	}
	if len(m) == 0 {
		return nil, fmt.Errorf("folder update is empty")
	}
	return json.Marshal(m)
}

// DeleteFolderQuery is the optional query for DELETE /projects/{projectId}/folders/{folderId}.
type DeleteFolderQuery struct {
	TransferToFolderID *string
}
