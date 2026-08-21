package models

import (
	"encoding/json"
	"testing"
)

func TestFolderUnmarshalCreate201(t *testing.T) {
	const raw = `{
		"id":"fld-1",
		"name":"Alpha",
		"parentFolderId":null,
		"createdAt":"2026-01-01T00:00:00.000Z",
		"updatedAt":"2026-01-01T00:00:00.000Z"
	}`

	var f Folder
	if err := json.Unmarshal([]byte(raw), &f); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if f.ID != "fld-1" || f.Name != "Alpha" {
		t.Fatalf("got %+v", f)
	}
	if f.ParentFolderID != nil {
		t.Fatalf("parent=%v", f.ParentFolderID)
	}
}

func TestFolderUnmarshalGet200Counts(t *testing.T) {
	const raw = `{
		"id":"fld-1",
		"name":"Alpha",
		"parentFolderId":"fld-0",
		"createdAt":"2026-01-01T00:00:00.000Z",
		"updatedAt":"2026-01-01T00:00:00.000Z",
		"totalSubFolders":2,
		"totalWorkflows":3
	}`

	var f Folder
	if err := json.Unmarshal([]byte(raw), &f); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if f.ParentFolderID == nil || *f.ParentFolderID != "fld-0" {
		t.Fatalf("parent=%v", f.ParentFolderID)
	}
	if f.TotalSubFolders == nil || *f.TotalSubFolders != 2 {
		t.Fatalf("sub=%v", f.TotalSubFolders)
	}
	if f.TotalWorkflows == nil || *f.TotalWorkflows != 3 {
		t.Fatalf("workflows=%v", f.TotalWorkflows)
	}
}

func TestFolderListUnmarshal(t *testing.T) {
	const raw = `{"count":1,"data":[{"id":"fld-1","name":"Alpha","parentFolderId":null}]}`

	var list FolderList
	if err := json.Unmarshal([]byte(raw), &list); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if list.Count != 1 || len(list.Data) != 1 || list.Data[0].ID != "fld-1" {
		t.Fatalf("got %+v", list)
	}
}

func TestFolderWriteMarshal(t *testing.T) {
	b, err := json.Marshal(FolderWrite{Name: "Alpha"})
	if err != nil {
		t.Fatal(err)
	}
	if string(b) != `{"name":"Alpha"}` {
		t.Fatalf("got %s", b)
	}
}

func TestFolderUpdateMarshalMoveToRoot(t *testing.T) {
	b, err := json.Marshal(FolderUpdate{MoveToRoot: true})
	if err != nil {
		t.Fatal(err)
	}
	if string(b) != `{"parentFolderId":null}` {
		t.Fatalf("got %s", b)
	}
}
