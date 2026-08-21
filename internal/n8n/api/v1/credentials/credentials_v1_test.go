package credentials

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"testing"

	"github.com/ubie-oss/terraform-provider-n8n/internal/n8n"
	"github.com/ubie-oss/terraform-provider-n8n/internal/n8n/models"
)

func TestCreateCredentialV1(t *testing.T) {
	client := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		assertAPIKey(t, r)
		if r.Method != http.MethodPost || r.URL.Path != "/api/v1/credentials" {
			t.Errorf("unexpected %s %s", r.Method, r.URL.Path)
		}
		var in models.CredentialWrite
		decodeBody(t, r, &in)
		if in.Name != "Header" || in.Type != "httpHeaderAuth" {
			t.Errorf("in=%+v", in)
		}
		if in.Data["name"] != "X-Test" {
			t.Errorf("data=%v", in.Data)
		}
		writeFixture(t, w, http.StatusOK, "create_credential_200.json")
	})

	got, err := CreateCredentialV1(context.Background(), client, models.CredentialWrite{
		Name: "Header",
		Type: "httpHeaderAuth",
		Data: map[string]any{"name": "X-Test"},
	})
	if err != nil {
		t.Fatalf("CreateCredentialV1: %v", err)
	}
	if got.ID != "cred-1" || got.Type != "httpHeaderAuth" {
		t.Fatalf("got %+v", got)
	}
}

func TestCreateCredentialV1EmptyName(t *testing.T) {
	client := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		t.Error("HTTP should not be called")
	})
	if _, err := CreateCredentialV1(context.Background(), client, models.CredentialWrite{
		Type: "httpHeaderAuth",
		Data: map[string]any{},
	}); err == nil {
		t.Fatal("expected error")
	}
}

func TestListCredentialsV1(t *testing.T) {
	client := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		assertAPIKey(t, r)
		if r.Method != http.MethodGet || r.URL.Path != "/api/v1/credentials" {
			t.Errorf("unexpected %s %s", r.Method, r.URL.Path)
		}
		if r.URL.Query().Get("limit") != "250" {
			t.Errorf("limit=%q", r.URL.Query().Get("limit"))
		}
		if r.URL.Query().Get("cursor") != "abc" {
			t.Errorf("cursor=%q", r.URL.Query().Get("cursor"))
		}
		writeFixture(t, w, http.StatusOK, "list_credentials.json")
	})

	list, err := ListCredentialsV1(context.Background(), client, 0, "abc")
	if err != nil {
		t.Fatalf("ListCredentialsV1: %v", err)
	}
	if len(list.Data) != 1 || list.Data[0].ID != "cred-1" {
		t.Fatalf("got %+v", list.Data)
	}
	if list.NextCursor != nil {
		t.Fatalf("nextCursor=%v", list.NextCursor)
	}
}

func TestListCredentialsV1403(t *testing.T) {
	client := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		writeFixture(t, w, http.StatusForbidden, "list_forbidden_403.json")
	})
	_, err := ListCredentialsV1(context.Background(), client, 100, "")
	var apiErr *n8n.APIError
	if err == nil || !errors.As(err, &apiErr) || apiErr.StatusCode != http.StatusForbidden {
		t.Fatalf("expected 403 APIError, got %T %v", err, err)
	}
}

func TestGetCredentialV1(t *testing.T) {
	client := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		assertAPIKey(t, r)
		if r.Method != http.MethodGet || r.URL.Path != "/api/v1/credentials/cred-1" {
			t.Errorf("unexpected %s %s", r.Method, r.URL.Path)
		}
		writeFixture(t, w, http.StatusOK, "get_credential_200.json")
	})
	got, err := GetCredentialV1(context.Background(), client, "cred-1")
	if err != nil {
		t.Fatalf("GetCredentialV1: %v", err)
	}
	if got.ID != "cred-1" {
		t.Fatalf("got %+v", got)
	}
}

func TestGetCredentialV1404(t *testing.T) {
	client := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})
	_, err := GetCredentialV1(context.Background(), client, "missing")
	if !n8n.IsNotFound(err) {
		t.Fatalf("expected NotFoundError, got %T %v", err, err)
	}
}

func TestUpdateCredentialV1(t *testing.T) {
	name := "Renamed"
	client := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		assertAPIKey(t, r)
		if r.Method != http.MethodPatch || r.URL.Path != "/api/v1/credentials/cred-1" {
			t.Errorf("unexpected %s %s", r.Method, r.URL.Path)
		}
		var in map[string]any
		decodeBody(t, r, &in)
		if in["name"] != "Renamed" {
			t.Errorf("body=%v", in)
		}
		writeFixture(t, w, http.StatusOK, "get_credential_200.json")
	})
	got, err := UpdateCredentialV1(context.Background(), client, "cred-1", models.CredentialUpdate{Name: &name})
	if err != nil {
		t.Fatalf("UpdateCredentialV1: %v", err)
	}
	if got.ID != "cred-1" {
		t.Fatalf("got %+v", got)
	}
}

func TestDeleteCredentialV1(t *testing.T) {
	client := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		assertAPIKey(t, r)
		if r.Method != http.MethodDelete || r.URL.Path != "/api/v1/credentials/cred-1" {
			t.Errorf("unexpected %s %s", r.Method, r.URL.Path)
		}
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"name":"Header","type":"httpHeaderAuth"}`))
	})
	if err := DeleteCredentialV1(context.Background(), client, "cred-1"); err != nil {
		t.Fatalf("DeleteCredentialV1: %v", err)
	}
}

func TestDeleteCredentialV1EmptyID(t *testing.T) {
	client := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		t.Error("HTTP should not be called")
	})
	if err := DeleteCredentialV1(context.Background(), client, "  "); err == nil {
		t.Fatal("expected error")
	}
}

func TestGetCredentialSchemaV1(t *testing.T) {
	client := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v1/credentials/schema/httpHeaderAuth" {
			t.Errorf("path=%s", r.URL.Path)
		}
		writeFixture(t, w, http.StatusOK, "schema_httpHeaderAuth.json")
	})
	raw, err := GetCredentialSchemaV1(context.Background(), client, "httpHeaderAuth")
	if err != nil {
		t.Fatalf("GetCredentialSchemaV1: %v", err)
	}
	var m map[string]any
	if err := json.Unmarshal(raw, &m); err != nil {
		t.Fatal(err)
	}
	if m["type"] != "object" {
		t.Fatalf("schema=%s", raw)
	}
}

func TestGetCredentialSchemaV1404(t *testing.T) {
	client := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		writeFixture(t, w, http.StatusNotFound, "not_found_404.json")
	})
	_, err := GetCredentialSchemaV1(context.Background(), client, "missing")
	if !n8n.IsNotFound(err) {
		t.Fatalf("expected NotFoundError, got %T %v", err, err)
	}
}

func TestTransferCredentialV1(t *testing.T) {
	client := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut || r.URL.Path != "/api/v1/credentials/cred-1/transfer" {
			t.Errorf("unexpected %s %s", r.Method, r.URL.Path)
		}
		var in models.CredentialTransfer
		decodeBody(t, r, &in)
		if in.DestinationProjectID != "proj-2" {
			t.Errorf("dest=%q", in.DestinationProjectID)
		}
		w.WriteHeader(http.StatusOK)
	})
	if err := TransferCredentialV1(context.Background(), client, "cred-1", "proj-2"); err != nil {
		t.Fatalf("TransferCredentialV1: %v", err)
	}
}

func TestTransferCredentialV1EmptyDest(t *testing.T) {
	client := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		t.Error("HTTP should not be called")
	})
	if err := TransferCredentialV1(context.Background(), client, "cred-1", " "); err == nil {
		t.Fatal("expected error")
	}
}
