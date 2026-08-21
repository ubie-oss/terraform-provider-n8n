package n8n

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNewRequiresAPIKey(t *testing.T) {
	_, err := New("https://n8n.example.com", "", nil)
	if err == nil {
		t.Fatal("expected empty API key to return an error")
	}
}

func TestNewRequiresEndpoint(t *testing.T) {
	_, err := New("", "test-key", nil)
	if err == nil {
		t.Fatal("expected empty endpoint to return an error")
	}
}

func TestNormalizeEndpoint(t *testing.T) {
	tests := []struct {
		in   string
		want string
	}{
		{"https://n8n.example.com", "https://n8n.example.com/api/v1"},
		{"https://n8n.example.com/", "https://n8n.example.com/api/v1"},
		{"https://n8n.example.com/api/v1", "https://n8n.example.com/api/v1"},
		{"https://n8n.example.com/api/v1/", "https://n8n.example.com/api/v1"},
		{"https://my.app.n8n.cloud/api/v1", "https://my.app.n8n.cloud/api/v1"},
	}
	for _, tt := range tests {
		got, err := NormalizeEndpoint(tt.in)
		if err != nil {
			t.Fatalf("NormalizeEndpoint(%q): %v", tt.in, err)
		}
		if got != tt.want {
			t.Fatalf("NormalizeEndpoint(%q)=%q, want %q", tt.in, got, tt.want)
		}
	}
}

func TestNewSetsAPIKeyHeader(t *testing.T) {
	var gotHeader string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotHeader = r.Header.Get("X-N8N-API-KEY")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{}`))
	}))
	t.Cleanup(srv.Close)

	client, err := New(srv.URL, "secret-key", &Options{MaxConcurrent: 5, RPS: 100})
	if err != nil {
		t.Fatalf("New: %v", err)
	}

	req, err := http.NewRequest(http.MethodGet, client.Endpoint+"/projects", nil)
	if err != nil {
		t.Fatalf("NewRequest: %v", err)
	}
	resp, err := client.HTTP.Do(req)
	if err != nil {
		t.Fatalf("Do: %v", err)
	}
	_ = resp.Body.Close()

	if gotHeader != "secret-key" {
		t.Fatalf("expected X-N8N-API-KEY=secret-key, got %q", gotHeader)
	}
	if auth := req.Header.Get("Authorization"); auth != "" {
		t.Fatalf("did not expect Authorization on request clone input, got %q", auth)
	}
}

func TestNewRejectsInvalidLimits(t *testing.T) {
	_, err := New("https://n8n.example.com", "k", &Options{MaxConcurrent: -1, RPS: 1})
	if err == nil {
		t.Fatal("expected error for negative max concurrent")
	}
	_, err = New("https://n8n.example.com", "k", &Options{MaxConcurrent: 1, RPS: -0.5})
	if err == nil {
		t.Fatal("expected error for negative RPS")
	}
}
