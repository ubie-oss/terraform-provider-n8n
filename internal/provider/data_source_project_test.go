package provider

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
)

func TestProjectDataSourceMetadata(t *testing.T) {
	ds := NewProjectDataSource()

	var resp datasource.MetadataResponse
	ds.Metadata(context.Background(), datasource.MetadataRequest{ProviderTypeName: "n8n"}, &resp)

	if resp.TypeName != "n8n_project" {
		t.Fatalf("expected n8n_project, got %q", resp.TypeName)
	}
}

func TestProjectsDataSourceMetadata(t *testing.T) {
	ds := NewProjectsDataSource()

	var resp datasource.MetadataResponse
	ds.Metadata(context.Background(), datasource.MetadataRequest{ProviderTypeName: "n8n"}, &resp)

	if resp.TypeName != "n8n_projects" {
		t.Fatalf("expected n8n_projects, got %q", resp.TypeName)
	}
}
