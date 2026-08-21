package services

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strconv"
	"testing"

	"github.com/ubie-oss/terraform-provider-n8n/internal/n8n"
	"github.com/ubie-oss/terraform-provider-n8n/internal/n8n/models"
)

func folderTestdata(t *testing.T, name string) []byte {
	t.Helper()
	b, err := os.ReadFile(filepath.Join("..", "api", "v1", "folders", "testdata", name))
	if err != nil {
		t.Fatal(err)
	}
	return b
}

func TestFolderServiceListAllPaginates(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v1/projects/proj-1/folders" {
			t.Errorf("path=%s", r.URL.Path)
		}
		if r.URL.Query().Get("take") == "" {
			t.Error("missing take")
		}
		w.Header().Set("Content-Type", "application/json")
		skip, _ := strconv.Atoi(r.URL.Query().Get("skip"))
		switch skip {
		case 0:
			_, _ = w.Write([]byte(`{"count":2,"data":[{"id":"fld-1","name":"Alpha","parentFolderId":null}]}`))
		case 1:
			_, _ = w.Write([]byte(`{"count":2,"data":[{"id":"fld-2","name":"Beta","parentFolderId":"fld-1"}]}`))
		default:
			t.Errorf("unexpected skip %d", skip)
			w.WriteHeader(http.StatusInternalServerError)
		}
	}))
	t.Cleanup(srv.Close)

	client, err := n8n.New(srv.URL, "secret", &n8n.Options{RPS: 100})
	if err != nil {
		t.Fatal(err)
	}
	all, err := NewFolderService(client).ListAll(context.Background(), "proj-1")
	if err != nil {
		t.Fatalf("ListAll: %v", err)
	}
	if len(all) != 2 || all[0].ID != "fld-1" || all[1].ID != "fld-2" {
		t.Fatalf("got %+v", all)
	}
}

func TestFolderServiceCreateGetDelete(t *testing.T) {
	var deletedWithTransfer string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch {
		case r.Method == http.MethodPost && r.URL.Path == "/api/v1/projects/proj-1/folders":
			w.WriteHeader(http.StatusCreated)
			_, _ = w.Write(folderTestdata(t, "create_folder_201.json"))
		case r.Method == http.MethodGet && r.URL.Path == "/api/v1/projects/proj-1/folders/fld-1":
			_, _ = w.Write(folderTestdata(t, "get_folder_200.json"))
		case r.Method == http.MethodDelete && r.URL.Path == "/api/v1/projects/proj-1/folders/fld-1":
			deletedWithTransfer = r.URL.Query().Get("transferToFolderId")
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
	svc := NewFolderService(client)

	created, err := svc.Create(context.Background(), "proj-1", models.FolderWrite{Name: "Alpha"})
	if err != nil {
		t.Fatalf("Create: %v", err)
	}
	if created.ID != "fld-1" {
		t.Fatalf("id=%q", created.ID)
	}

	got, err := svc.Get(context.Background(), "proj-1", "fld-1")
	if err != nil {
		t.Fatalf("Get: %v", err)
	}
	if got.Name != "Alpha" {
		t.Fatalf("name=%q", got.Name)
	}

	if err := svc.Delete(context.Background(), "proj-1", "fld-1", models.DeleteFolderQuery{}); err != nil {
		t.Fatalf("Delete: %v", err)
	}
	if deletedWithTransfer != "" {
		t.Fatalf("unexpected transfer %q", deletedWithTransfer)
	}

	target := "fld-root"
	if err := svc.Delete(context.Background(), "proj-1", "fld-1", models.DeleteFolderQuery{TransferToFolderID: &target}); err != nil {
		t.Fatalf("Delete with transfer: %v", err)
	}
	if deletedWithTransfer != "fld-root" {
		t.Fatalf("transfer=%q", deletedWithTransfer)
	}
}

func TestFolderServiceListAllLicense403(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusForbidden)
		_, _ = w.Write(folderTestdata(t, "feature_not_licensed_403.json"))
	}))
	t.Cleanup(srv.Close)

	client, err := n8n.New(srv.URL, "secret", &n8n.Options{RPS: 100})
	if err != nil {
		t.Fatal(err)
	}
	_, err = NewFolderService(client).ListAll(context.Background(), "proj-1")
	if !n8n.IsFoldersUnlicensed(err) {
		t.Fatalf("expected license error, got %v", err)
	}
}

func TestFolderServiceListAllEmptyProjectID(t *testing.T) {
	svc := NewFolderService(&n8n.Client{})
	if _, err := svc.ListAll(context.Background(), "  "); err == nil {
		t.Fatal("expected error")
	}
}

func TestFolderServiceUpdate(t *testing.T) {
	name := "Renamed"
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPatch {
			t.Errorf("method=%s", r.Method)
		}
		var body models.FolderUpdate
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Errorf("decode: %v", err)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write(folderTestdata(t, "create_folder_201.json"))
	}))
	t.Cleanup(srv.Close)

	client, err := n8n.New(srv.URL, "secret", &n8n.Options{RPS: 100})
	if err != nil {
		t.Fatal(err)
	}
	got, err := NewFolderService(client).Update(context.Background(), "proj-1", "fld-1", models.FolderUpdate{Name: &name})
	if err != nil {
		t.Fatalf("Update: %v", err)
	}
	if got.ID != "fld-1" {
		t.Fatalf("got %+v", got)
	}
}
