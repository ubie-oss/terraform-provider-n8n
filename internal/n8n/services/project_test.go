package services

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/ubie-oss/terraform-provider-n8n/internal/n8n"
	"github.com/ubie-oss/terraform-provider-n8n/internal/n8n/models"
)

func testdata(t *testing.T, name string) []byte {
	t.Helper()
	b, err := os.ReadFile(filepath.Join("..", "api", "v1", "projects", "testdata", name))
	if err != nil {
		t.Fatal(err)
	}
	return b
}

func TestProjectServiceListAllPaginates(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v1/projects" {
			t.Errorf("path=%s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		switch r.URL.Query().Get("cursor") {
		case "":
			_, _ = w.Write(testdata(t, "list_projects_page1.json"))
		case "cur2":
			_, _ = w.Write(testdata(t, "list_projects_page2.json"))
		default:
			t.Errorf("unexpected cursor %q", r.URL.Query().Get("cursor"))
			w.WriteHeader(http.StatusInternalServerError)
		}
	}))
	t.Cleanup(srv.Close)

	client, err := n8n.New(srv.URL, "secret", &n8n.Options{RPS: 100})
	if err != nil {
		t.Fatal(err)
	}
	svc := NewProjectService(client)
	all, err := svc.ListAll(context.Background())
	if err != nil {
		t.Fatalf("ListAll: %v", err)
	}
	if len(all) != 3 {
		t.Fatalf("len=%d", len(all))
	}
	if all[2].ID != "p3" {
		t.Fatalf("last=%+v", all[2])
	}

	got, err := svc.GetByID(context.Background(), "p2")
	if err != nil {
		t.Fatalf("GetByID: %v", err)
	}
	if got.Name != "two" {
		t.Fatalf("name=%q", got.Name)
	}

	_, err = svc.GetByID(context.Background(), "missing")
	if !errors.Is(err, ErrProjectNotFound) {
		t.Fatalf("expected ErrProjectNotFound, got %v", err)
	}
}

func TestProjectServiceGetByIDEmpty(t *testing.T) {
	svc := NewProjectService(&n8n.Client{})
	if _, err := svc.GetByID(context.Background(), "  "); err == nil {
		t.Fatal("expected error")
	}
}

func TestProjectServiceCreateUpdateDelete(t *testing.T) {
	projects := []models.Project{{
		ID:   "abc123",
		Name: "some-project",
		Type: models.ProjectTypeTeam,
	}}

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch {
		case r.Method == http.MethodPost && r.URL.Path == "/api/v1/projects":
			w.WriteHeader(http.StatusCreated)
			_, _ = w.Write(testdata(t, "create_project_201.json"))
		case r.Method == http.MethodPut && r.URL.Path == "/api/v1/projects/abc123":
			var body struct {
				Name string `json:"name"`
			}
			if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
				t.Errorf("decode: %v", err)
			}
			projects[0].Name = body.Name
			w.WriteHeader(http.StatusNoContent)
		case r.Method == http.MethodGet && r.URL.Path == "/api/v1/projects":
			_ = json.NewEncoder(w).Encode(models.ProjectList{Data: projects})
		case r.Method == http.MethodDelete && r.URL.Path == "/api/v1/projects/abc123":
			w.WriteHeader(http.StatusNoContent)
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
	svc := NewProjectService(client)

	created, err := svc.Create(context.Background(), "some-project")
	if err != nil {
		t.Fatalf("Create: %v", err)
	}
	if created.ID != "abc123" {
		t.Fatalf("id=%q", created.ID)
	}

	updated, err := svc.Update(context.Background(), "abc123", "renamed")
	if err != nil {
		t.Fatalf("Update: %v", err)
	}
	if updated.Name != "renamed" {
		t.Fatalf("name=%q", updated.Name)
	}

	if err := svc.Delete(context.Background(), "abc123"); err != nil {
		t.Fatalf("Delete: %v", err)
	}
}

func TestProjectServiceListAllLicense403(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusForbidden)
		_, _ = w.Write(testdata(t, "feature_not_licensed_403.json"))
	}))
	t.Cleanup(srv.Close)

	client, err := n8n.New(srv.URL, "secret", &n8n.Options{RPS: 100})
	if err != nil {
		t.Fatal(err)
	}
	_, err = NewProjectService(client).ListAll(context.Background())
	if !n8n.IsProjectRoleAdminUnlicensed(err) {
		t.Fatalf("expected license error, got %v", err)
	}
}
