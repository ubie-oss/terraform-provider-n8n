package projects

import (
	"context"
	"net/http"
	"testing"

	"github.com/ubie-oss/terraform-provider-n8n/internal/n8n"
	"github.com/ubie-oss/terraform-provider-n8n/internal/n8n/models"
)

func TestListProjectsV1(t *testing.T) {
	client := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		assertAPIKey(t, r)
		if r.Method != http.MethodGet || r.URL.Path != "/api/v1/projects" {
			t.Errorf("unexpected %s %s", r.Method, r.URL.Path)
		}
		if r.URL.Query().Get("limit") != "250" {
			t.Errorf("limit=%q", r.URL.Query().Get("limit"))
		}
		if r.URL.Query().Get("cursor") != "abc" {
			t.Errorf("cursor=%q", r.URL.Query().Get("cursor"))
		}
		writeFixture(t, w, http.StatusOK, "list_projects.json")
	})

	list, err := ListProjectsV1(context.Background(), client, 0, "abc")
	if err != nil {
		t.Fatalf("ListProjectsV1: %v", err)
	}
	if len(list.Data) != 2 || list.Data[1].ID != "abc123" {
		t.Fatalf("got %+v", list.Data)
	}
	if list.NextCursor != nil {
		t.Fatalf("nextCursor=%v", list.NextCursor)
	}
}

func TestListProjectsV1FeatureNotLicensed(t *testing.T) {
	client := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		writeFixture(t, w, http.StatusForbidden, "feature_not_licensed_403.json")
	})

	_, err := ListProjectsV1(context.Background(), client, 100, "")
	if !n8n.IsProjectRoleAdminUnlicensed(err) {
		t.Fatalf("expected license 403, got %v", err)
	}
}

func TestCreateProjectV1(t *testing.T) {
	client := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		assertAPIKey(t, r)
		if r.Method != http.MethodPost || r.URL.Path != "/api/v1/projects" {
			t.Errorf("unexpected %s %s", r.Method, r.URL.Path)
		}
		var in models.ProjectWrite
		decodeBody(t, r, &in)
		if in.Name != "some-project" {
			t.Errorf("name=%q", in.Name)
		}
		writeFixture(t, w, http.StatusCreated, "create_project_201.json")
	})

	got, err := CreateProjectV1(context.Background(), client, models.ProjectWrite{Name: "some-project"})
	if err != nil {
		t.Fatalf("CreateProjectV1: %v", err)
	}
	if got.ID != "abc123" || got.Type != models.ProjectTypeTeam || got.Name != "some-project" {
		t.Fatalf("got %+v", got)
	}
}

func TestCreateProjectV1EmptyName(t *testing.T) {
	client := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		t.Error("HTTP should not be called")
	})
	if _, err := CreateProjectV1(context.Background(), client, models.ProjectWrite{}); err == nil {
		t.Fatal("expected error")
	}
}

func TestUpdateProjectV1204(t *testing.T) {
	client := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		assertAPIKey(t, r)
		if r.Method != http.MethodPut || r.URL.Path != "/api/v1/projects/abc123" {
			t.Errorf("unexpected %s %s", r.Method, r.URL.Path)
		}
		var in models.ProjectWrite
		decodeBody(t, r, &in)
		if in.Name != "renamed" {
			t.Errorf("name=%q", in.Name)
		}
		w.WriteHeader(http.StatusNoContent)
	})

	if err := UpdateProjectV1(context.Background(), client, "abc123", models.ProjectWrite{Name: "renamed"}); err != nil {
		t.Fatalf("UpdateProjectV1: %v", err)
	}
}

func TestUpdateProjectV1404(t *testing.T) {
	client := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte(`{"message":"not found"}`))
	})
	err := UpdateProjectV1(context.Background(), client, "missing", models.ProjectWrite{Name: "x"})
	if !n8n.IsNotFound(err) {
		t.Fatalf("expected NotFoundError, got %T %v", err, err)
	}
}

func TestDeleteProjectV1204(t *testing.T) {
	client := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		assertAPIKey(t, r)
		if r.Method != http.MethodDelete || r.URL.Path != "/api/v1/projects/abc123" {
			t.Errorf("unexpected %s %s", r.Method, r.URL.Path)
		}
		w.WriteHeader(http.StatusNoContent)
	})
	if err := DeleteProjectV1(context.Background(), client, "abc123"); err != nil {
		t.Fatalf("DeleteProjectV1: %v", err)
	}
}

func TestDeleteProjectV1404(t *testing.T) {
	client := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})
	err := DeleteProjectV1(context.Background(), client, "gone")
	if !n8n.IsNotFound(err) {
		t.Fatalf("expected NotFoundError, got %T %v", err, err)
	}
}

func TestDeleteProjectV1EmptyID(t *testing.T) {
	client := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		t.Error("HTTP should not be called")
	})
	if err := DeleteProjectV1(context.Background(), client, "  "); err == nil {
		t.Fatal("expected error")
	}
}
