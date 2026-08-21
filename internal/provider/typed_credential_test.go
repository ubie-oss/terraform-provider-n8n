package provider

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func TestTypedHTTPCredentialMetadata(t *testing.T) {
	cases := []struct {
		res      resource.Resource
		typeName string
	}{
		{NewCredentialHTTPHeaderAuthResource(), "n8n_credential_http_header_auth"},
		{NewCredentialHTTPBasicAuthResource(), "n8n_credential_http_basic_auth"},
		{NewCredentialHTTPBearerAuthResource(), "n8n_credential_http_bearer_auth"},
	}
	for _, tc := range cases {
		var resp resource.MetadataResponse
		tc.res.Metadata(context.Background(), resource.MetadataRequest{ProviderTypeName: "n8n"}, &resp)
		if resp.TypeName != tc.typeName {
			t.Fatalf("expected %q, got %q", tc.typeName, resp.TypeName)
		}
		var schemaResp resource.SchemaResponse
		tc.res.Schema(context.Background(), resource.SchemaRequest{}, &schemaResp)
		if schemaResp.Diagnostics.HasError() {
			t.Fatalf("schema: %v", schemaResp.Diagnostics)
		}
	}
}

func TestHTTPHeaderAuthBuildData(t *testing.T) {
	data, err := httpHeaderAuthSpec().BuildData(typedAttrBag{
		strings: map[string]types.String{
			"header_name": types.StringValue("X-Test"),
			"value":       types.StringValue("secret"),
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	if data["name"] != "X-Test" || data["value"] != "secret" {
		t.Fatalf("unexpected data: %#v", data)
	}
}

func TestHTTPBasicAuthBuildData(t *testing.T) {
	data, err := httpBasicAuthSpec().BuildData(typedAttrBag{
		strings: map[string]types.String{
			"user":     types.StringValue("alice"),
			"password": types.StringValue("p"),
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	if data["user"] != "alice" || data["password"] != "p" {
		t.Fatalf("unexpected data: %#v", data)
	}
}

func TestHTTPBearerAuthBuildData(t *testing.T) {
	data, err := httpBearerAuthSpec().BuildData(typedAttrBag{
		strings: map[string]types.String{
			"token": types.StringValue("tok"),
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	if data["token"] != "tok" {
		t.Fatalf("unexpected data: %#v", data)
	}
}

func TestIsCredentialManagedResource(t *testing.T) {
	if !isCredentialManagedResource("n8n_credential") {
		t.Fatal("generic resource should match")
	}
	if !isCredentialManagedResource("n8n_credential_http_header_auth") {
		t.Fatal("typed resource should match")
	}
	if isCredentialManagedResource("n8n_credentials") {
		t.Fatal("data source n8n_credentials must not match")
	}
	if isCredentialManagedResource("n8n_project") {
		t.Fatal("project must not match")
	}
}
