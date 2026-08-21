package services

import (
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/ubie-oss/terraform-provider-n8n/internal/n8n"
	"github.com/ubie-oss/terraform-provider-n8n/internal/n8n/models"
)

func tagTestdata(t *testing.T, name string) []byte {
	t.Helper()
	b, err := os.ReadFile(filepath.Join("..", "api", "v1", "tags", "testdata", name))
	if err != nil {
		t.Fatal(err)
	}
	return b
}

func TestTagServiceListAllPaginates(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v1/tags" {
			t.Errorf("path=%s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		switch r.URL.Query().Get("cursor") {
		case "":
			_, _ = w.Write(tagTestdata(t, "list_tags_page1.json"))
		case "cur2":
			_, _ = w.Write(tagTestdata(t, "list_tags_page2.json"))
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
	all, err := NewTagService(client).ListAll(context.Background())
	if err != nil {
		t.Fatalf("ListAll: %v", err)
	}
	if len(all) != 2 || all[0].ID != "tag-1" || all[1].ID != "tag-2" {
		t.Fatalf("got %+v", all)
	}
}

func TestTagServiceCreateGetDelete(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch {
		case r.Method == http.MethodPost && r.URL.Path == "/api/v1/tags":
			w.WriteHeader(http.StatusCreated)
			_, _ = w.Write(tagTestdata(t, "tag_200.json"))
		case r.Method == http.MethodGet && r.URL.Path == "/api/v1/tags/tag-1":
			_, _ = w.Write(tagTestdata(t, "tag_200.json"))
		case r.Method == http.MethodDelete && r.URL.Path == "/api/v1/tags/tag-1":
			_, _ = w.Write(tagTestdata(t, "tag_200.json"))
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
	svc := NewTagService(client)
	created, err := svc.Create(context.Background(), models.TagWrite{Name: "Production"})
	if err != nil {
		t.Fatalf("Create: %v", err)
	}
	if created.ID != "tag-1" {
		t.Fatalf("created %+v", created)
	}
	got, err := svc.Get(context.Background(), "tag-1")
	if err != nil {
		t.Fatalf("Get: %v", err)
	}
	if got.Name != "Production" {
		t.Fatalf("got %+v", got)
	}
	if err := svc.Delete(context.Background(), "tag-1"); err != nil {
		t.Fatalf("Delete: %v", err)
	}
}
