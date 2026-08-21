package folders

import (
	"context"
	"net/http"
	"testing"

	"github.com/ubie-oss/terraform-provider-n8n/internal/n8n"
	"github.com/ubie-oss/terraform-provider-n8n/internal/n8n/models"
)

func TestCreateFolderV1(t *testing.T) {
	client := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		assertAPIKey(t, r)
		if r.Method != http.MethodPost || r.URL.Path != "/api/v1/projects/proj-1/folders" {
			t.Errorf("unexpected %s %s", r.Method, r.URL.Path)
		}
		var in models.FolderWrite
		decodeBody(t, r, &in)
		if in.Name != "Alpha" {
			t.Errorf("name=%q", in.Name)
		}
		writeFixture(t, w, http.StatusCreated, "create_folder_201.json")
	})

	got, err := CreateFolderV1(context.Background(), client, "proj-1", models.FolderWrite{Name: "Alpha"})
	if err != nil {
		t.Fatalf("CreateFolderV1: %v", err)
	}
	if got.ID != "fld-1" || got.Name != "Alpha" {
		t.Fatalf("got %+v", got)
	}
}

func TestCreateFolderV1EmptyName(t *testing.T) {
	client := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		t.Error("HTTP should not be called")
	})
	if _, err := CreateFolderV1(context.Background(), client, "proj-1", models.FolderWrite{}); err == nil {
		t.Fatal("expected error")
	}
}

func TestListFoldersV1(t *testing.T) {
	client := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		assertAPIKey(t, r)
		if r.Method != http.MethodGet || r.URL.Path != "/api/v1/projects/proj-1/folders" {
			t.Errorf("unexpected %s %s", r.Method, r.URL.Path)
		}
		if r.URL.Query().Get("skip") != "0" {
			t.Errorf("skip=%q", r.URL.Query().Get("skip"))
		}
		if r.URL.Query().Get("take") != "10" {
			t.Errorf("take=%q", r.URL.Query().Get("take"))
		}
		writeFixture(t, w, http.StatusOK, "list_folders.json")
	})

	list, err := ListFoldersV1(context.Background(), client, "proj-1", 0, DefaultFolderListTake)
	if err != nil {
		t.Fatalf("ListFoldersV1: %v", err)
	}
	if list.Count != 2 || len(list.Data) != 2 || list.Data[1].ID != "fld-2" {
		t.Fatalf("got %+v", list)
	}
}

func TestListFoldersV1FeatureNotLicensed(t *testing.T) {
	client := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		writeFixture(t, w, http.StatusForbidden, "feature_not_licensed_403.json")
	})

	_, err := ListFoldersV1(context.Background(), client, "proj-1", 0, 10)
	if !n8n.IsFoldersUnlicensed(err) {
		t.Fatalf("expected license 403, got %v", err)
	}
}

func TestGetFolderV1(t *testing.T) {
	client := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		assertAPIKey(t, r)
		if r.Method != http.MethodGet || r.URL.Path != "/api/v1/projects/proj-1/folders/fld-1" {
			t.Errorf("unexpected %s %s", r.Method, r.URL.Path)
		}
		writeFixture(t, w, http.StatusOK, "get_folder_200.json")
	})

	got, err := GetFolderV1(context.Background(), client, "proj-1", "fld-1")
	if err != nil {
		t.Fatalf("GetFolderV1: %v", err)
	}
	if got.ID != "fld-1" || got.TotalSubFolders == nil || *got.TotalSubFolders != 0 {
		t.Fatalf("got %+v", got)
	}
}

func TestGetFolderV1404(t *testing.T) {
	client := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})
	_, err := GetFolderV1(context.Background(), client, "proj-1", "missing")
	if !n8n.IsNotFound(err) {
		t.Fatalf("expected NotFoundError, got %T %v", err, err)
	}
}

func TestUpdateFolderV1(t *testing.T) {
	name := "Renamed"
	client := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		assertAPIKey(t, r)
		if r.Method != http.MethodPatch || r.URL.Path != "/api/v1/projects/proj-1/folders/fld-1" {
			t.Errorf("unexpected %s %s", r.Method, r.URL.Path)
		}
		var in models.FolderUpdate
		decodeBody(t, r, &in)
		if in.Name == nil || *in.Name != "Renamed" {
			t.Errorf("body name=%v", in.Name)
		}
		writeFixture(t, w, http.StatusOK, "create_folder_201.json")
	})

	got, err := UpdateFolderV1(context.Background(), client, "proj-1", "fld-1", models.FolderUpdate{Name: &name})
	if err != nil {
		t.Fatalf("UpdateFolderV1: %v", err)
	}
	if got.ID != "fld-1" {
		t.Fatalf("got %+v", got)
	}
}

func TestDeleteFolderV1204(t *testing.T) {
	client := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		assertAPIKey(t, r)
		if r.Method != http.MethodDelete || r.URL.Path != "/api/v1/projects/proj-1/folders/fld-1" {
			t.Errorf("unexpected %s %s", r.Method, r.URL.Path)
		}
		if r.URL.Query().Get("transferToFolderId") != "" {
			t.Errorf("unexpected transfer query %q", r.URL.RawQuery)
		}
		w.WriteHeader(http.StatusNoContent)
	})
	if err := DeleteFolderV1(context.Background(), client, "proj-1", "fld-1", models.DeleteFolderQuery{}); err != nil {
		t.Fatalf("DeleteFolderV1: %v", err)
	}
}

func TestDeleteFolderV1TransferQuery(t *testing.T) {
	target := "fld-root"
	client := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("transferToFolderId") != "fld-root" {
			t.Errorf("transfer=%q", r.URL.Query().Get("transferToFolderId"))
		}
		w.WriteHeader(http.StatusNoContent)
	})
	if err := DeleteFolderV1(context.Background(), client, "proj-1", "fld-1", models.DeleteFolderQuery{TransferToFolderID: &target}); err != nil {
		t.Fatalf("DeleteFolderV1: %v", err)
	}
}

func TestDeleteFolderV1EmptyID(t *testing.T) {
	client := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		t.Error("HTTP should not be called")
	})
	if err := DeleteFolderV1(context.Background(), client, "proj-1", "  ", models.DeleteFolderQuery{}); err == nil {
		t.Fatal("expected error")
	}
}
