package provider

import (
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

func NewCredentialAnthropicAPIResource() resource.Resource {
	return newTypedCredentialResource(anthropicAPISpec())
}

func NewCredentialAzureOpenAIAPIResource() resource.Resource {
	return newTypedCredentialResource(azureOpenAiAPISpec())
}

func NewCredentialOllamaAPIResource() resource.Resource {
	return newTypedCredentialResource(ollamaAPISpec())
}

func anthropicAPISpec() typedCredentialSpec {
	return typedCredentialSpec{
		TerraformSuffix: "credential_anthropic_api",
		N8nType:         "anthropicApi",
		DocFile:         "internal/provider/docs/resources/credential_anthropic_api.md",
		ExtraAttributes: map[string]schema.Attribute{
			"api_key":      stringWriteOnly("Anthropic API key. Write-only; never stored in state.", true),
			"url":          stringState("Anthropic API base URL. Defaults in n8n to `https://api.anthropic.com` when omitted.", false),
			"header":       boolState("When true, send the custom header_name and header_value."),
			"header_name":  stringState("Custom header name used when header is true.", false),
			"header_value": stringWriteOnly("Custom header value used when header is true. Write-only; never stored in state.", false),
		},
		BuildData: func(bag typedAttrBag) (map[string]any, error) {
			m := map[string]any{}
			if err := bagPutRequiredString(m, "apiKey", bag.strings["api_key"]); err != nil {
				return nil, err
			}
			bagPutString(m, "url", bag.strings["url"])
			bagPutBool(m, "header", bag.bools["header"])
			bagPutString(m, "headerName", bag.strings["header_name"])
			bagPutString(m, "headerValue", bag.strings["header_value"])
			return m, nil
		},
	}
}

func azureOpenAiAPISpec() typedCredentialSpec {
	return typedCredentialSpec{
		TerraformSuffix: "credential_azure_open_ai_api",
		N8nType:         "azureOpenAiApi",
		DocFile:         "internal/provider/docs/resources/credential_azure_open_ai_api.md",
		ExtraAttributes: map[string]schema.Attribute{
			"api_key":       stringWriteOnly("Azure OpenAI API key (for example Key 1). Write-only; never stored in state.", true),
			"resource_name": stringState("Azure OpenAI resource name.", true),
			"api_version":   stringState("Azure OpenAI API version. n8n defaults to a preview version when omitted in the UI; Terraform requires an explicit value.", true),
			"endpoint":      stringState("Optional custom Azure OpenAI endpoint, for example `https://westeurope.api.cognitive.microsoft.com`.", false),
		},
		BuildData: func(bag typedAttrBag) (map[string]any, error) {
			m := map[string]any{}
			if err := bagPutRequiredString(m, "apiKey", bag.strings["api_key"]); err != nil {
				return nil, err
			}
			if err := bagPutRequiredString(m, "resourceName", bag.strings["resource_name"]); err != nil {
				return nil, err
			}
			if err := bagPutRequiredString(m, "apiVersion", bag.strings["api_version"]); err != nil {
				return nil, err
			}
			bagPutString(m, "endpoint", bag.strings["endpoint"])
			return m, nil
		},
	}
}

func ollamaAPISpec() typedCredentialSpec {
	return typedCredentialSpec{
		TerraformSuffix: "credential_ollama_api",
		N8nType:         "ollamaApi",
		DocFile:         "internal/provider/docs/resources/credential_ollama_api.md",
		ExtraAttributes: map[string]schema.Attribute{
			"base_url": stringState("Ollama base URL. Defaults in n8n to `http://localhost:11434` when omitted in the UI; Terraform requires an explicit value.", true),
			"api_key":  stringWriteOnly("Optional Bearer token for authenticated Ollama proxies (for example Open WebUI). Write-only; never stored in state.", false),
		},
		BuildData: func(bag typedAttrBag) (map[string]any, error) {
			m := map[string]any{}
			if err := bagPutRequiredString(m, "baseUrl", bag.strings["base_url"]); err != nil {
				return nil, err
			}
			bagPutString(m, "apiKey", bag.strings["api_key"])
			return m, nil
		},
	}
}
