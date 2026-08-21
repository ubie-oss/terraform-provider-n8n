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

func TestTagControllerCreateEmptyName(t *testing.T) {
	c := NewTagController(&n8n.Client{})
	if _, err := c.Create(context.Background(), CreateTagOptions{}); err == nil {
		t.Fatal("expected error")
	}
}

func TestTagControllerCreateNameTooLong(t *testing.T) {
	c := NewTagController(&n8n.Client{})
	long := strings.Repeat("a", MaxTagNameLength+1)
	if _, err := c.Create(context.Background(), CreateTagOptions{Name: long}); err == nil {
		t.Fatal("expected error")
	}
}

func TestTagControllerDeleteProtection(t *testing.T) {
	c := NewTagController(&n8n.Client{})
	err := c.Delete(context.Background(), DeleteTagOptions{ID: "tag-1", DeleteProtection: true})
	if err == nil {
		t.Fatal("expected delete protection error")
	}
	if !strings.Contains(err.Error(), "delete protection is enabled") {
		t.Fatalf("got %v", err)
	}
}

func TestTagControllerListSortsByID(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(models.TagList{
			Data: []models.Tag{
				{ID: "tag-2", Name: "Beta"},
				{ID: "tag-1", Name: "Alpha"},
			},
		})
	}))
	t.Cleanup(srv.Close)

	client, err := n8n.New(srv.URL, "secret", &n8n.Options{RPS: 100})
	if err != nil {
		t.Fatal(err)
	}
	all, err := NewTagController(client).List(context.Background())
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if len(all) != 2 || all[0].ID != "tag-1" || all[1].ID != "tag-2" {
		t.Fatalf("got %+v", all)
	}
}

func TestTagControllerDeleteMissingIsOK(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	t.Cleanup(srv.Close)

	client, err := n8n.New(srv.URL, "secret", &n8n.Options{RPS: 100})
	if err != nil {
		t.Fatal(err)
	}
	if err := NewTagController(client).Delete(context.Background(), DeleteTagOptions{ID: "missing"}); err != nil {
		t.Fatalf("Delete: %v", err)
	}
}

// PUT /tags/{id} omits createdAt; Update must refresh via GET so Terraform state stays consistent.
func TestTagControllerUpdateRefreshesCreatedAt(t *testing.T) {
	var putHits, getHits int
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch r.Method {
		case http.MethodPut:
			putHits++
			_ = json.NewEncoder(w).Encode(models.Tag{
				ID:        "tag-1",
				Name:      "Staging",
				UpdatedAt: "2026-01-02T00:00:00.000Z",
			})
		case http.MethodGet:
			getHits++
			_ = json.NewEncoder(w).Encode(models.Tag{
				ID:        "tag-1",
				Name:      "Staging",
				CreatedAt: "2026-01-01T00:00:00.000Z",
				UpdatedAt: "2026-01-02T00:00:00.000Z",
			})
		default:
			t.Errorf("unexpected %s %s", r.Method, r.URL.Path)
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	}))
	t.Cleanup(srv.Close)

	client, err := n8n.New(srv.URL, "secret", &n8n.Options{RPS: 100})
	if err != nil {
		t.Fatal(err)
	}
	got, err := NewTagController(client).Update(context.Background(), UpdateTagOptions{ID: "tag-1", Name: "Staging"})
	if err != nil {
		t.Fatalf("Update: %v", err)
	}
	if putHits != 1 || getHits != 1 {
		t.Fatalf("put=%d get=%d", putHits, getHits)
	}
	if got.CreatedAt != "2026-01-01T00:00:00.000Z" {
		t.Fatalf("createdAt=%q want prior create timestamp from GET", got.CreatedAt)
	}
}
