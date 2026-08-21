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

func TestFolderControllerCreateEmptyName(t *testing.T) {
	c := NewFolderController(&n8n.Client{})
	if _, err := c.Create(context.Background(), CreateFolderOptions{ProjectID: "proj-1", Name: "  "}); err == nil {
		t.Fatal("expected error")
	}
}

func TestFolderControllerCreateRejectsPersonalPath(t *testing.T) {
	c := NewFolderController(&n8n.Client{})
	if _, err := c.Create(context.Background(), CreateFolderOptions{ProjectID: "personal", Name: "Alpha"}); err == nil {
		t.Fatal("expected personal reject")
	}
}

func TestFolderControllerDeleteProtection(t *testing.T) {
	c := NewFolderController(&n8n.Client{})
	err := c.Delete(context.Background(), DeleteFolderOptions{
		ProjectID:        "proj-1",
		FolderID:         "fld-1",
		DeleteProtection: true,
	})
	if err == nil {
		t.Fatal("expected delete protection error")
	}
	if !strings.Contains(err.Error(), "delete protection is enabled") {
		t.Fatalf("got %v", err)
	}
}

func TestFolderControllerUpdateRejectsPersonal(t *testing.T) {
	c := NewFolderController(&n8n.Client{})
	if _, err := c.Update(context.Background(), UpdateFolderOptions{
		ProjectID: "personal",
		FolderID:  "fld-1",
		Name:      "Alpha",
	}); err == nil {
		t.Fatal("expected personal reject")
	}
}

func TestFolderControllerImportRejectsPersonalProject(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(models.ProjectList{Data: []models.Project{{
			ID:   "p1",
			Name: "Personal",
			Type: models.ProjectTypePersonal,
		}}})
	}))
	t.Cleanup(srv.Close)

	client, err := n8n.New(srv.URL, "secret", &n8n.Options{RPS: 100})
	if err != nil {
		t.Fatal(err)
	}
	c := NewFolderController(client)
	if _, err := c.Import(context.Background(), ImportFolderOptions{ProjectID: "p1", FolderID: "fld-1"}); err == nil {
		t.Fatal("expected refuse personal project")
	}
}

func TestFolderControllerGetSkipsTeamProjectCheck(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v1/projects/p1/folders/fld-1" {
			t.Errorf("path=%s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(models.Folder{ID: "fld-1", Name: "Inbox"})
	}))
	t.Cleanup(srv.Close)

	client, err := n8n.New(srv.URL, "secret", &n8n.Options{RPS: 100})
	if err != nil {
		t.Fatal(err)
	}
	got, err := NewFolderController(client).Get(context.Background(), "p1", "fld-1")
	if err != nil {
		t.Fatalf("Get should skip team-project checks: %v", err)
	}
	if got.Name != "Inbox" {
		t.Fatalf("name=%q", got.Name)
	}
}

func TestFolderControllerListFilters(t *testing.T) {
	parent := "fld-1"
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(models.FolderList{
			Count: 2,
			Data: []models.Folder{
				{ID: "fld-2", Name: "Beta", ParentFolderID: &parent},
				{ID: "fld-1", Name: "Alpha"},
			},
		})
	}))
	t.Cleanup(srv.Close)

	client, err := n8n.New(srv.URL, "secret", &n8n.Options{RPS: 100})
	if err != nil {
		t.Fatal(err)
	}
	c := NewFolderController(client)
	all, err := c.List(context.Background(), ListFoldersOptions{ProjectID: "proj-1"})
	if err != nil {
		t.Fatal(err)
	}
	if len(all) != 2 || all[0].ID != "fld-1" {
		t.Fatalf("sorted list=%+v", all)
	}
	filtered, err := c.List(context.Background(), ListFoldersOptions{ProjectID: "proj-1", Name: "Beta"})
	if err != nil {
		t.Fatal(err)
	}
	if len(filtered) != 1 || filtered[0].ID != "fld-2" {
		t.Fatalf("name filter=%+v", filtered)
	}
}
