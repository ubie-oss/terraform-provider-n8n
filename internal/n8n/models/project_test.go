package models

import (
	"encoding/json"
	"testing"
)

func TestProjectUnmarshalCreate201(t *testing.T) {
	const raw = `{
		"id": "abc123",
		"name": "some-project",
		"type": "team",
		"icon": null,
		"description": null,
		"creatorId": "user-1",
		"customTelemetryTags": [],
		"createdAt": "2024-01-01T00:00:00.000Z",
		"updatedAt": "2024-01-01T00:00:00.000Z",
		"role": "project:admin",
		"scopes": ["project:create"]
	}`

	var p Project
	if err := json.Unmarshal([]byte(raw), &p); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if p.ID != "abc123" || p.Name != "some-project" || p.Type != ProjectTypeTeam {
		t.Fatalf("got %+v", p)
	}
	if !p.IsTeam() {
		t.Fatal("expected team project")
	}
	if p.Role != "project:admin" || p.CreatorID != "user-1" {
		t.Fatalf("extra fields: role=%q creator=%q", p.Role, p.CreatorID)
	}
}

func TestProjectListUnmarshal(t *testing.T) {
	const raw = `{
		"data": [
			{"id": "p1", "name": "Personal", "type": "personal"},
			{"id": "t1", "name": "platform", "type": "team"}
		],
		"nextCursor": null
	}`

	var list ProjectList
	if err := json.Unmarshal([]byte(raw), &list); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if len(list.Data) != 2 {
		t.Fatalf("len=%d", len(list.Data))
	}
	if list.NextCursor != nil {
		t.Fatalf("expected null nextCursor, got %v", list.NextCursor)
	}
	if list.Data[0].IsTeam() {
		t.Fatal("personal project should not be team")
	}
	if !list.Data[1].IsTeam() {
		t.Fatal("expected second project to be team")
	}
}

func TestProjectWriteMarshal(t *testing.T) {
	b, err := json.Marshal(ProjectWrite{Name: "platform"})
	if err != nil {
		t.Fatal(err)
	}
	if string(b) != `{"name":"platform"}` {
		t.Fatalf("got %s", b)
	}
}
