package tags

import (
	"context"
	"errors"
	"net/http"
	"testing"

	"github.com/ubie-oss/terraform-provider-n8n/internal/n8n"
	"github.com/ubie-oss/terraform-provider-n8n/internal/n8n/models"
)

func TestCreateTagV1(t *testing.T) {
	client := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		assertAPIKey(t, r)
		if r.Method != http.MethodPost || r.URL.Path != "/api/v1/tags" {
			t.Errorf("unexpected %s %s", r.Method, r.URL.Path)
		}
		var in models.TagWrite
		decodeBody(t, r, &in)
		if in.Name != "Production" {
			t.Errorf("in=%+v", in)
		}
		writeFixture(t, w, http.StatusCreated, "tag_200.json")
	})

	got, err := CreateTagV1(context.Background(), client, models.TagWrite{Name: "Production"})
	if err != nil {
		t.Fatalf("CreateTagV1: %v", err)
	}
	if got.ID != "tag-1" || got.Name != "Production" {
		t.Fatalf("got %+v", got)
	}
}

func TestCreateTagV1EmptyName(t *testing.T) {
	client := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		t.Error("HTTP should not be called")
	})
	if _, err := CreateTagV1(context.Background(), client, models.TagWrite{}); err == nil {
		t.Fatal("expected error")
	}
}

func TestCreateTagV1409(t *testing.T) {
	client := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		writeFixture(t, w, http.StatusConflict, "conflict_409.json")
	})
	_, err := CreateTagV1(context.Background(), client, models.TagWrite{Name: "Production"})
	var apiErr *n8n.APIError
	if err == nil || !errors.As(err, &apiErr) || apiErr.StatusCode != http.StatusConflict {
		t.Fatalf("expected 409 APIError, got %T %v", err, err)
	}
}

func TestListTagsV1(t *testing.T) {
	client := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		assertAPIKey(t, r)
		if r.Method != http.MethodGet || r.URL.Path != "/api/v1/tags" {
			t.Errorf("unexpected %s %s", r.Method, r.URL.Path)
		}
		if r.URL.Query().Get("limit") != "250" {
			t.Errorf("limit=%q", r.URL.Query().Get("limit"))
		}
		if r.URL.Query().Get("cursor") != "abc" {
			t.Errorf("cursor=%q", r.URL.Query().Get("cursor"))
		}
		writeFixture(t, w, http.StatusOK, "list_tags.json")
	})

	list, err := ListTagsV1(context.Background(), client, 0, "abc")
	if err != nil {
		t.Fatalf("ListTagsV1: %v", err)
	}
	if len(list.Data) != 1 || list.Data[0].ID != "tag-1" {
		t.Fatalf("got %+v", list.Data)
	}
	if list.NextCursor != nil {
		t.Fatalf("nextCursor=%v", list.NextCursor)
	}
}

func TestGetTagV1(t *testing.T) {
	client := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		assertAPIKey(t, r)
		if r.Method != http.MethodGet || r.URL.Path != "/api/v1/tags/tag-1" {
			t.Errorf("unexpected %s %s", r.Method, r.URL.Path)
		}
		writeFixture(t, w, http.StatusOK, "tag_200.json")
	})
	got, err := GetTagV1(context.Background(), client, "tag-1")
	if err != nil {
		t.Fatalf("GetTagV1: %v", err)
	}
	if got.ID != "tag-1" {
		t.Fatalf("got %+v", got)
	}
}

func TestGetTagV1404(t *testing.T) {
	client := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})
	_, err := GetTagV1(context.Background(), client, "missing")
	if !n8n.IsNotFound(err) {
		t.Fatalf("expected NotFoundError, got %T %v", err, err)
	}
}

func TestUpdateTagV1(t *testing.T) {
	client := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		assertAPIKey(t, r)
		if r.Method != http.MethodPut || r.URL.Path != "/api/v1/tags/tag-1" {
			t.Errorf("unexpected %s %s", r.Method, r.URL.Path)
		}
		var in models.TagWrite
		decodeBody(t, r, &in)
		if in.Name != "Staging" {
			t.Errorf("body=%+v", in)
		}
		// Live n8n PUT /tags/{id} omits createdAt; keep the fixture honest.
		writeFixture(t, w, http.StatusOK, "tag_update_200.json")
	})
	got, err := UpdateTagV1(context.Background(), client, "tag-1", models.TagWrite{Name: "Staging"})
	if err != nil {
		t.Fatalf("UpdateTagV1: %v", err)
	}
	if got.ID != "tag-1" || got.Name != "Staging" {
		t.Fatalf("got %+v", got)
	}
	if got.CreatedAt != "" {
		t.Fatalf("expected empty createdAt on PUT body, got %q", got.CreatedAt)
	}
}

func TestDeleteTagV1(t *testing.T) {
	client := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		assertAPIKey(t, r)
		if r.Method != http.MethodDelete || r.URL.Path != "/api/v1/tags/tag-1" {
			t.Errorf("unexpected %s %s", r.Method, r.URL.Path)
		}
		w.WriteHeader(http.StatusOK)
	})
	if err := DeleteTagV1(context.Background(), client, "tag-1"); err != nil {
		t.Fatalf("DeleteTagV1: %v", err)
	}
}
