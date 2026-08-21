package services

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/ubie-oss/terraform-provider-n8n/internal/n8n"
	"github.com/ubie-oss/terraform-provider-n8n/internal/n8n/models"
)

func credentialTestdata(t *testing.T, name string) []byte {
	t.Helper()
	b, err := os.ReadFile(filepath.Join("..", "api", "v1", "credentials", "testdata", name))
	if err != nil {
		t.Fatal(err)
	}
	return b
}

func TestCredentialServiceListAllPaginates(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v1/credentials" {
			t.Errorf("path=%s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		switch r.URL.Query().Get("cursor") {
		case "":
			_, _ = w.Write(credentialTestdata(t, "list_credentials_page1.json"))
		case "cur2":
			_, _ = w.Write(credentialTestdata(t, "list_credentials_page2.json"))
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
	all, err := NewCredentialService(client).ListAll(context.Background())
	if err != nil {
		t.Fatalf("ListAll: %v", err)
	}
	if len(all) != 2 || all[0].ID != "cred-1" || all[1].ID != "cred-2" {
		t.Fatalf("got %+v", all)
	}
}

func TestCredentialServiceCreateGetDelete(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch {
		case r.Method == http.MethodPost && r.URL.Path == "/api/v1/credentials":
			_, _ = w.Write(credentialTestdata(t, "create_credential_200.json"))
		case r.Method == http.MethodGet && r.URL.Path == "/api/v1/credentials/cred-1":
			_, _ = w.Write(credentialTestdata(t, "get_credential_200.json"))
		case r.Method == http.MethodDelete && r.URL.Path == "/api/v1/credentials/cred-1":
			w.WriteHeader(http.StatusOK)
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
	svc := NewCredentialService(client)
	created, err := svc.Create(context.Background(), models.CredentialWrite{
		Name: "Header",
		Type: "httpHeaderAuth",
		Data: map[string]any{"name": "X-Test"},
	})
	if err != nil {
		t.Fatalf("Create: %v", err)
	}
	if created.ID != "cred-1" {
		t.Fatalf("id=%q", created.ID)
	}
	got, err := svc.Get(context.Background(), "cred-1")
	if err != nil {
		t.Fatalf("Get: %v", err)
	}
	if got.Name != "Header" {
		t.Fatalf("name=%q", got.Name)
	}
	if err := svc.Delete(context.Background(), "cred-1"); err != nil {
		t.Fatalf("Delete: %v", err)
	}
}

func TestCredentialServiceSchema(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v1/credentials/schema/httpHeaderAuth" {
			t.Errorf("path=%s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write(credentialTestdata(t, "schema_httpHeaderAuth.json"))
	}))
	t.Cleanup(srv.Close)

	client, err := n8n.New(srv.URL, "secret", &n8n.Options{RPS: 100})
	if err != nil {
		t.Fatal(err)
	}
	raw, err := NewCredentialService(client).Schema(context.Background(), "httpHeaderAuth")
	if err != nil {
		t.Fatal(err)
	}
	var m map[string]any
	if err := json.Unmarshal(raw, &m); err != nil {
		t.Fatal(err)
	}
	if m["type"] != "object" {
		t.Fatalf("schema=%s", raw)
	}
}
