package provider

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
)

func TestTagDataSourceMetadata(t *testing.T) {
	ds := NewTagDataSource()

	var resp datasource.MetadataResponse
	ds.Metadata(context.Background(), datasource.MetadataRequest{ProviderTypeName: "n8n"}, &resp)

	if resp.TypeName != "n8n_tag" {
		t.Fatalf("expected n8n_tag, got %q", resp.TypeName)
	}
}

func TestTagsDataSourceMetadata(t *testing.T) {
	ds := NewTagsDataSource()

	var resp datasource.MetadataResponse
	ds.Metadata(context.Background(), datasource.MetadataRequest{ProviderTypeName: "n8n"}, &resp)

	if resp.TypeName != "n8n_tags" {
		t.Fatalf("expected n8n_tags, got %q", resp.TypeName)
	}
}
