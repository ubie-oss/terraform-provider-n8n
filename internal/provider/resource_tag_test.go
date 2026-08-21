package provider

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/resource"
)

func TestTagResourceMetadata(t *testing.T) {
	res := NewTagResource()

	var resp resource.MetadataResponse
	res.Metadata(context.Background(), resource.MetadataRequest{ProviderTypeName: "n8n"}, &resp)

	if resp.TypeName != "n8n_tag" {
		t.Fatalf("expected n8n_tag, got %q", resp.TypeName)
	}
}
