// Package apitest provides shared httptest helpers for Public API v1 package tests.
package apitest

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/ubie-oss/terraform-provider-n8n/internal/n8n"
)

// Testdata reads a fixture from the caller's package testdata directory.
func Testdata(t *testing.T, name string) []byte {
	t.Helper()
	root, err := os.OpenRoot("testdata")
	if err != nil {
		t.Fatal(err)
	}
	defer root.Close()
	b, err := root.ReadFile(name)
	if err != nil {
		t.Fatal(err)
	}
	return b
}

// TestClient builds an n8n.Client pointed at an httptest server with API key "secret".
func TestClient(t *testing.T, handler http.HandlerFunc) *n8n.Client {
	t.Helper()
	srv := httptest.NewServer(handler)
	t.Cleanup(srv.Close)
	client, err := n8n.New(srv.URL, "secret", &n8n.Options{RPS: 100})
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	return client
}

// DecodeBody unmarshals the request body into out.
func DecodeBody(t *testing.T, r *http.Request, out any) {
	t.Helper()
	b, err := io.ReadAll(r.Body)
	if err != nil {
		t.Fatalf("read body: %v", err)
	}
	if err := json.Unmarshal(b, out); err != nil {
		t.Fatalf("decode body %s: %v", b, err)
	}
}

// AssertAPIKey checks the httptest request carries X-N8N-API-KEY=secret.
func AssertAPIKey(t *testing.T, r *http.Request) {
	t.Helper()
	if r.Header.Get("X-N8N-API-KEY") != "secret" {
		t.Errorf("missing API key header")
	}
}
