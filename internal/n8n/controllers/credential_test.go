package controllers

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/ubie-oss/terraform-provider-n8n/internal/n8n"
	"github.com/ubie-oss/terraform-provider-n8n/internal/n8n/models"
)

func TestCredentialControllerCreateEmptyName(t *testing.T) {
	c := NewCredentialController(&n8n.Client{})
	if _, err := c.Create(context.Background(), CreateCredentialOptions{
		Type: "httpHeaderAuth",
		Data: map[string]any{"name": "X-Test"},
	}); err == nil {
		t.Fatal("expected error")
	}
}

func TestCredentialControllerDeleteProtection(t *testing.T) {
	c := NewCredentialController(&n8n.Client{})
	err := c.Delete(context.Background(), DeleteCredentialOptions{ID: "cred-1", DeleteProtection: true})
	if err == nil {
		t.Fatal("expected delete protection error")
	}
	if !strings.Contains(err.Error(), "delete protection is enabled") {
		t.Fatalf("got %v", err)
	}
}

func TestCredentialControllerUpdateRejectsManaged(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(models.Credential{
			ID:        "cred-1",
			Name:      "Managed",
			Type:      "httpHeaderAuth",
			IsManaged: true,
		})
	}))
	t.Cleanup(srv.Close)

	client, err := n8n.New(srv.URL, "secret", &n8n.Options{RPS: 100})
	if err != nil {
		t.Fatal(err)
	}
	_, err = NewCredentialController(client).Update(context.Background(), UpdateCredentialOptions{
		ID:   "cred-1",
		Name: "Managed",
	})
	if err == nil || !strings.Contains(err.Error(), "managed by n8n") {
		t.Fatalf("got %v", err)
	}
}

func TestCredentialControllerUpdateRejectsClearingProjectID(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(models.Credential{ID: "cred-1", Name: "Header", Type: "httpHeaderAuth"})
	}))
	t.Cleanup(srv.Close)

	client, err := n8n.New(srv.URL, "secret", &n8n.Options{RPS: 100})
	if err != nil {
		t.Fatal(err)
	}
	prior := "proj-1"
	_, err = NewCredentialController(client).Update(context.Background(), UpdateCredentialOptions{
		ID:           "cred-1",
		Name:         "Header",
		PriorProject: &prior,
	})
	if err == nil || !strings.Contains(err.Error(), "cannot clear project_id") {
		t.Fatalf("got %v", err)
	}
}

func TestCredentialControllerListFilters(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(models.CredentialList{
			Data: []models.CredentialListItem{
				{ID: "cred-2", Name: "Beta", Type: "httpBasicAuth"},
				{ID: "cred-1", Name: "Alpha", Type: "httpHeaderAuth"},
			},
		})
	}))
	t.Cleanup(srv.Close)

	client, err := n8n.New(srv.URL, "secret", &n8n.Options{RPS: 100})
	if err != nil {
		t.Fatal(err)
	}
	c := NewCredentialController(client)
	all, err := c.List(context.Background(), ListCredentialsOptions{})
	if err != nil {
		t.Fatal(err)
	}
	if len(all) != 2 || all[0].ID != "cred-1" {
		t.Fatalf("sorted list=%+v", all)
	}
	filtered, err := c.List(context.Background(), ListCredentialsOptions{Type: "httpHeaderAuth"})
	if err != nil {
		t.Fatal(err)
	}
	if len(filtered) != 1 || filtered[0].ID != "cred-1" {
		t.Fatalf("type filter=%+v", filtered)
	}
}

func TestCredentialControllerUpdateTransfersThenPatches(t *testing.T) {
	var transferredTo string
	var patched bool
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch {
		case r.Method == http.MethodGet && r.URL.Path == "/api/v1/credentials/cred-1":
			_ = json.NewEncoder(w).Encode(models.Credential{ID: "cred-1", Name: "Header", Type: "httpHeaderAuth"})
		case r.Method == http.MethodPut && r.URL.Path == "/api/v1/credentials/cred-1/transfer":
			var in models.CredentialTransfer
			_ = json.NewDecoder(r.Body).Decode(&in)
			transferredTo = in.DestinationProjectID
			w.WriteHeader(http.StatusOK)
		case r.Method == http.MethodPatch && r.URL.Path == "/api/v1/credentials/cred-1":
			patched = true
			_ = json.NewEncoder(w).Encode(models.Credential{ID: "cred-1", Name: "Renamed", Type: "httpHeaderAuth"})
		default:
			t.Errorf("unexpected %s %s", r.Method, r.URL.Path)
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	t.Cleanup(srv.Close)

	client, err := n8n.New(srv.URL, "secret", &n8n.Options{RPS: 100})
	if err != nil {
		t.Fatal(err)
	}
	next := "proj-2"
	prior := "proj-1"
	got, err := NewCredentialController(client).Update(context.Background(), UpdateCredentialOptions{
		ID:           "cred-1",
		Name:         "Renamed",
		ProjectID:    &next,
		PriorProject: &prior,
	})
	if err != nil {
		t.Fatal(err)
	}
	if transferredTo != "proj-2" || !patched || got.Name != "Renamed" {
		t.Fatalf("transfer=%q patched=%v got=%+v", transferredTo, patched, got)
	}
}

func TestCredentialControllerUpdateAdoptsProjectWithoutTransfer(t *testing.T) {
	var transferred bool
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch {
		case r.Method == http.MethodGet && r.URL.Path == "/api/v1/credentials/cred-1":
			_ = json.NewEncoder(w).Encode(models.Credential{ID: "cred-1", Name: "Header", Type: "httpHeaderAuth"})
		case r.Method == http.MethodPut && strings.HasSuffix(r.URL.Path, "/transfer"):
			transferred = true
			w.WriteHeader(http.StatusOK)
		case r.Method == http.MethodPatch && r.URL.Path == "/api/v1/credentials/cred-1":
			_ = json.NewEncoder(w).Encode(models.Credential{ID: "cred-1", Name: "Header", Type: "httpHeaderAuth"})
		default:
			t.Errorf("unexpected %s %s", r.Method, r.URL.Path)
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	t.Cleanup(srv.Close)

	client, err := n8n.New(srv.URL, "secret", &n8n.Options{RPS: 100})
	if err != nil {
		t.Fatal(err)
	}
	next := "proj-1"
	if _, err := NewCredentialController(client).Update(context.Background(), UpdateCredentialOptions{
		ID:        "cred-1",
		Name:      "Header",
		ProjectID: &next,
	}); err != nil {
		t.Fatal(err)
	}
	if transferred {
		t.Fatal("nil prior project_id must adopt without Transfer")
	}
}
