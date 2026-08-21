package controllers

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/ubie-oss/terraform-provider-n8n/internal/n8n"
	"github.com/ubie-oss/terraform-provider-n8n/internal/n8n/models"
)

func TestErrIfNotTeamProject(t *testing.T) {
	if err := errIfNotTeamProject(&models.Project{ID: "t1", Type: models.ProjectTypeTeam}); err != nil {
		t.Fatalf("team should be allowed: %v", err)
	}
	if err := errIfNotTeamProject(&models.Project{ID: "p1", Type: models.ProjectTypePersonal}); err == nil {
		t.Fatal("personal project should be rejected")
	}
	if err := errIfNotTeamProject(nil); err == nil {
		t.Fatal("nil project should be rejected")
	}
}

func TestProjectControllerCreateEmptyName(t *testing.T) {
	c := NewProjectController(&n8n.Client{})
	if _, err := c.Create(context.Background(), CreateProjectOptions{Name: "  "}); err == nil {
		t.Fatal("expected error")
	}
}

func TestProjectControllerGetTeamRejectsPersonal(t *testing.T) {
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
	c := NewProjectController(client)

	got, err := c.Get(context.Background(), "p1")
	if err != nil {
		t.Fatalf("Get should allow personal: %v", err)
	}
	if got.Type != models.ProjectTypePersonal {
		t.Fatalf("type=%q", got.Type)
	}

	if _, err := c.GetTeam(context.Background(), "p1"); err == nil {
		t.Fatal("GetTeam should reject personal")
	}
}

func TestProjectControllerDeleteMissingIsSuccess(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(models.ProjectList{Data: []models.Project{}})
	}))
	t.Cleanup(srv.Close)

	client, err := n8n.New(srv.URL, "secret", &n8n.Options{RPS: 100})
	if err != nil {
		t.Fatal(err)
	}
	c := NewProjectController(client)
	if err := c.Delete(context.Background(), DeleteProjectOptions{ID: "gone"}); err != nil {
		t.Fatalf("missing delete should succeed: %v", err)
	}
}

func TestProjectControllerDeleteProtection(t *testing.T) {
	c := NewProjectController(&n8n.Client{})
	err := c.Delete(context.Background(), DeleteProjectOptions{ID: "p1", DeleteProtection: true})
	if err == nil {
		t.Fatal("expected delete protection error")
	}
	if !strings.Contains(err.Error(), "delete protection is enabled") {
		t.Fatalf("got %v", err)
	}
}

func TestProjectControllerDeleteRejectsPersonal(t *testing.T) {
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
	c := NewProjectController(client)
	if err := c.Delete(context.Background(), DeleteProjectOptions{ID: "p1"}); err == nil {
		t.Fatal("expected refuse personal")
	}
}

func TestProjectControllerGetNotFound(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(models.ProjectList{Data: []models.Project{}})
	}))
	t.Cleanup(srv.Close)

	client, err := n8n.New(srv.URL, "secret", &n8n.Options{RPS: 100})
	if err != nil {
		t.Fatal(err)
	}
	_, err = NewProjectController(client).Get(context.Background(), "missing")
	if !errors.Is(err, ErrProjectNotFound) {
		t.Fatalf("expected ErrProjectNotFound, got %v", err)
	}
}

func TestProjectControllerListSortsByID(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(models.ProjectList{Data: []models.Project{
			{ID: "b", Name: "two", Type: models.ProjectTypeTeam},
			{ID: "a", Name: "one", Type: models.ProjectTypePersonal},
		}})
	}))
	t.Cleanup(srv.Close)

	client, err := n8n.New(srv.URL, "secret", &n8n.Options{RPS: 100})
	if err != nil {
		t.Fatal(err)
	}
	all, err := NewProjectController(client).List(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if len(all) != 2 || all[0].ID != "a" || all[1].ID != "b" {
		t.Fatalf("got %+v", all)
	}
}
