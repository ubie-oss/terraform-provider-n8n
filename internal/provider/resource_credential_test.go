package provider

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

func TestCredentialResourceMetadata(t *testing.T) {
	res := NewCredentialResource()

	var resp resource.MetadataResponse
	res.Metadata(context.Background(), resource.MetadataRequest{ProviderTypeName: "n8n"}, &resp)

	if resp.TypeName != "n8n_credential" {
		t.Fatalf("expected n8n_credential, got %q", resp.TypeName)
	}
}

func TestCredentialDataSourceMetadata(t *testing.T) {
	ds := NewCredentialDataSource()

	var resp datasource.MetadataResponse
	ds.Metadata(context.Background(), datasource.MetadataRequest{ProviderTypeName: "n8n"}, &resp)

	if resp.TypeName != "n8n_credential" {
		t.Fatalf("expected n8n_credential, got %q", resp.TypeName)
	}
}

func TestCredentialsDataSourceMetadata(t *testing.T) {
	ds := NewCredentialsDataSource()

	var resp datasource.MetadataResponse
	ds.Metadata(context.Background(), datasource.MetadataRequest{ProviderTypeName: "n8n"}, &resp)

	if resp.TypeName != "n8n_credentials" {
		t.Fatalf("expected n8n_credentials, got %q", resp.TypeName)
	}
}

func TestCredentialSchemaDataSourceMetadata(t *testing.T) {
	ds := NewCredentialSchemaDataSource()

	var resp datasource.MetadataResponse
	ds.Metadata(context.Background(), datasource.MetadataRequest{ProviderTypeName: "n8n"}, &resp)

	if resp.TypeName != "n8n_credential_schema" {
		t.Fatalf("expected n8n_credential_schema, got %q", resp.TypeName)
	}
}
