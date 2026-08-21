package provider

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
)

func TestFolderDataSourceMetadata(t *testing.T) {
	ds := NewFolderDataSource()

	var resp datasource.MetadataResponse
	ds.Metadata(context.Background(), datasource.MetadataRequest{ProviderTypeName: "n8n"}, &resp)

	if resp.TypeName != "n8n_folder" {
		t.Fatalf("expected n8n_folder, got %q", resp.TypeName)
	}
}

func TestFoldersDataSourceMetadata(t *testing.T) {
	ds := NewFoldersDataSource()

	var resp datasource.MetadataResponse
	ds.Metadata(context.Background(), datasource.MetadataRequest{ProviderTypeName: "n8n"}, &resp)

	if resp.TypeName != "n8n_folders" {
		t.Fatalf("expected n8n_folders, got %q", resp.TypeName)
	}
}
