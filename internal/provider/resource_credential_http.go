package provider

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

func NewCredentialHTTPHeaderAuthResource() resource.Resource {
	return newTypedCredentialResource(httpHeaderAuthSpec())
}

func NewCredentialHTTPBasicAuthResource() resource.Resource {
	return newTypedCredentialResource(httpBasicAuthSpec())
}

func NewCredentialHTTPBearerAuthResource() resource.Resource {
	return newTypedCredentialResource(httpBearerAuthSpec())
}

func NewCredentialHTTPMultipleHeadersAuthResource() resource.Resource {
	return newTypedCredentialResource(httpMultipleHeadersAuthSpec())
}

func httpHeaderAuthSpec() typedCredentialSpec {
	return typedCredentialSpec{
		TerraformSuffix: "credential_http_header_auth",
		N8nType:         "httpHeaderAuth",
		DocFile:         "internal/provider/docs/resources/credential_http_header_auth.md",
		ExtraAttributes: map[string]schema.Attribute{
			"header_name": stringState("HTTP header name sent as n8n data key `name`.", true),
			"value":       stringWriteOnly("HTTP header value. Write-only; never stored in state.", true),
		},
		BuildData: func(bag typedAttrBag) (map[string]any, error) {
			m := map[string]any{}
			if err := bagPutRequiredString(m, "name", bag.strings["header_name"]); err != nil {
				return nil, err
			}
			if err := bagPutRequiredString(m, "value", bag.strings["value"]); err != nil {
				return nil, err
			}
			return m, nil
		},
	}
}

func httpBasicAuthSpec() typedCredentialSpec {
	return typedCredentialSpec{
		TerraformSuffix: "credential_http_basic_auth",
		N8nType:         "httpBasicAuth",
		DocFile:         "internal/provider/docs/resources/credential_http_basic_auth.md",
		ExtraAttributes: map[string]schema.Attribute{
			"user":     stringState("Basic auth username.", true),
			"password": stringWriteOnly("Basic auth password. Write-only; never stored in state.", true),
		},
		BuildData: func(bag typedAttrBag) (map[string]any, error) {
			m := map[string]any{}
			if err := bagPutRequiredString(m, "user", bag.strings["user"]); err != nil {
				return nil, err
			}
			if err := bagPutRequiredString(m, "password", bag.strings["password"]); err != nil {
				return nil, err
			}
			return m, nil
		},
	}
}

func httpBearerAuthSpec() typedCredentialSpec {
	return typedCredentialSpec{
		TerraformSuffix: "credential_http_bearer_auth",
		N8nType:         "httpBearerAuth",
		DocFile:         "internal/provider/docs/resources/credential_http_bearer_auth.md",
		ExtraAttributes: map[string]schema.Attribute{
			"token": stringWriteOnly("Bearer token. Write-only; never stored in state.", true),
		},
		BuildData: func(bag typedAttrBag) (map[string]any, error) {
			m := map[string]any{}
			if err := bagPutRequiredString(m, "token", bag.strings["token"]); err != nil {
				return nil, err
			}
			return m, nil
		},
	}
}

func httpMultipleHeadersAuthSpec() typedCredentialSpec {
	return typedCredentialSpec{
		TerraformSuffix: "credential_http_multiple_headers_auth",
		N8nType:         "httpMultipleHeadersAuth",
		DocFile:         "internal/provider/docs/resources/credential_http_multiple_headers_auth.md",
		ExtraAttributes: map[string]schema.Attribute{
			"headers": schema.DynamicAttribute{
				Required:            true,
				WriteOnly:           true,
				Sensitive:           true,
				MarkdownDescription: "Header collection matching n8n `headers` (`{ values = [{ name, value }] }`). Write-only; never stored in state.",
			},
		},
		BuildData: func(bag typedAttrBag) (map[string]any, error) {
			m := map[string]any{}
			headers, ok, err := bagDynamicValue(bag.dynamics["headers"], "headers")
			if err != nil {
				return nil, err
			}
			if !ok {
				return nil, fmt.Errorf("headers is required")
			}
			m["headers"] = headers
			return m, nil
		},
	}
}
