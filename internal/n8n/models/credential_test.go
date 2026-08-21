package models

import (
	"encoding/json"
	"testing"
)

func TestCredentialUnmarshalCreate200(t *testing.T) {
	const raw = `{
		"id":"cred-1",
		"name":"Header",
		"type":"httpHeaderAuth",
		"isManaged":false,
		"isGlobal":false,
		"isResolvable":false,
		"resolvableAllowFallback":false,
		"resolverId":null,
		"createdAt":"2026-01-01T00:00:00.000Z",
		"updatedAt":"2026-01-01T00:00:00.000Z"
	}`

	var c Credential
	if err := json.Unmarshal([]byte(raw), &c); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if c.ID != "cred-1" || c.Name != "Header" || c.Type != "httpHeaderAuth" {
		t.Fatalf("got %+v", c)
	}
	if c.ResolverID != nil {
		t.Fatalf("resolver=%v", c.ResolverID)
	}
	if c.IsManaged {
		t.Fatal("expected unmanaged")
	}
}

func TestCredentialListUnmarshalLiveShape(t *testing.T) {
	const raw = `{
		"data":[{
			"id":"cred-1",
			"name":"Header",
			"type":"httpHeaderAuth",
			"createdAt":"2026-01-01T00:00:00.000Z",
			"updatedAt":"2026-01-01T00:00:00.000Z",
			"shared":[{"id":"proj-1","name":"Personal","role":"credential:owner"}]
		}],
		"nextCursor":null
	}`

	var list CredentialList
	if err := json.Unmarshal([]byte(raw), &list); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if len(list.Data) != 1 || list.Data[0].ID != "cred-1" {
		t.Fatalf("got %+v", list)
	}
	if list.NextCursor != nil {
		t.Fatalf("nextCursor=%v", list.NextCursor)
	}
	if len(list.Data[0].Shared) != 1 || list.Data[0].Shared[0].Role != "credential:owner" {
		t.Fatalf("shared=%+v", list.Data[0].Shared)
	}
}

func TestCredentialWriteMarshalOmitsEmptyProject(t *testing.T) {
	b, err := json.Marshal(CredentialWrite{
		Name: "Header",
		Type: "httpHeaderAuth",
		Data: map[string]any{"name": "X-Test"},
	})
	if err != nil {
		t.Fatal(err)
	}
	var m map[string]any
	if err := json.Unmarshal(b, &m); err != nil {
		t.Fatal(err)
	}
	if _, ok := m["projectId"]; ok {
		t.Fatalf("unexpected projectId in %s", b)
	}
	if m["type"] != "httpHeaderAuth" {
		t.Fatalf("got %s", b)
	}
}

func TestCredentialUpdateMarshalNameOnly(t *testing.T) {
	name := "Renamed"
	b, err := json.Marshal(CredentialUpdate{Name: &name})
	if err != nil {
		t.Fatal(err)
	}
	if string(b) != `{"name":"Renamed"}` {
		t.Fatalf("got %s", b)
	}
}

func TestCredentialUpdateMarshalSendDataPartialFalse(t *testing.T) {
	partial := false
	b, err := json.Marshal(CredentialUpdate{
		SendData:      true,
		Data:          map[string]any{"value": "x"},
		IsPartialData: &partial,
	})
	if err != nil {
		t.Fatal(err)
	}
	var m map[string]any
	if err := json.Unmarshal(b, &m); err != nil {
		t.Fatal(err)
	}
	if m["isPartialData"] != false {
		t.Fatalf("isPartialData omitted or wrong: %s", b)
	}
	data, ok := m["data"].(map[string]any)
	if !ok || data["value"] != "x" {
		t.Fatalf("data=%v", m["data"])
	}
}

func TestCredentialUpdateEmptyError(t *testing.T) {
	if _, err := json.Marshal(CredentialUpdate{}); err == nil {
		t.Fatal("expected empty update error")
	}
}
