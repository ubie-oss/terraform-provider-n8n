package provider

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/resource"
)

func TestProjectResourceMetadata(t *testing.T) {
	res := NewProjectResource()

	var resp resource.MetadataResponse
	res.Metadata(context.Background(), resource.MetadataRequest{ProviderTypeName: "n8n"}, &resp)

	if resp.TypeName != "n8n_project" {
		t.Fatalf("expected n8n_project, got %q", resp.TypeName)
	}
}
