package provider

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/resource"
)

func TestFolderResourceMetadata(t *testing.T) {
	res := NewFolderResource()

	var resp resource.MetadataResponse
	res.Metadata(context.Background(), resource.MetadataRequest{ProviderTypeName: "n8n"}, &resp)

	if resp.TypeName != "n8n_folder" {
		t.Fatalf("expected n8n_folder, got %q", resp.TypeName)
	}
}

func TestFormatFolderResourceID(t *testing.T) {
	got := formatFolderResourceID("proj-1", "fld-1")
	if got != "projects/proj-1/folders/fld-1" {
		t.Fatalf("got %q", got)
	}
}

func TestParseFolderResourceID(t *testing.T) {
	projectID, folderID, err := parseFolderResourceID("projects/proj-1/folders/fld-1")
	if err != nil {
		t.Fatal(err)
	}
	if projectID != "proj-1" || folderID != "fld-1" {
		t.Fatalf("got %q %q", projectID, folderID)
	}

	if _, _, err := parseFolderResourceID("fld-1"); err == nil {
		t.Fatal("expected error for bare id")
	}
	if _, _, err := parseFolderResourceID("projects/proj-1"); err == nil {
		t.Fatal("expected error for missing folder")
	}
	if _, _, err := parseFolderResourceID("projects//folders/fld-1"); err == nil {
		t.Fatal("expected error for empty project")
	}
}
