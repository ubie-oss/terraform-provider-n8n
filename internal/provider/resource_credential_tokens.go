package provider

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

func NewCredentialSlackAPIResource() resource.Resource {
	return newTypedCredentialResource(slackAPISpec())
}

func NewCredentialNotionAPIResource() resource.Resource {
	return newTypedCredentialResource(notionAPISpec())
}

func NewCredentialN8nAPIResource() resource.Resource {
	return newTypedCredentialResource(n8nAPISpec())
}

func NewCredentialJiraSoftwareCloudAPIResource() resource.Resource {
	return newTypedCredentialResource(jiraSoftwareCloudAPISpec())
}

func NewCredentialHubspotAppTokenResource() resource.Resource {
	return newTypedCredentialResource(hubspotAppTokenSpec())
}

func NewCredentialSerpAPIResource() resource.Resource {
	return newTypedCredentialResource(serpAPISpec())
}

func NewCredentialOpenAIAPIResource() resource.Resource {
	return newTypedCredentialResource(openAIAPISpec())
}

func NewCredentialJWTAuthResource() resource.Resource {
	return newTypedCredentialResource(jwtAuthSpec())
}

func NewCredentialGitHubAPIResource() resource.Resource {
	return newTypedCredentialResource(githubAPISpec())
}

func NewCredentialSendGridAPIResource() resource.Resource {
	return newTypedCredentialResource(sendGridAPISpec())
}

func NewCredentialStripeAPIResource() resource.Resource {
	return newTypedCredentialResource(stripeAPISpec())
}

func NewCredentialTwilioAPIResource() resource.Resource {
	return newTypedCredentialResource(twilioAPISpec())
}

func slackAPISpec() typedCredentialSpec {
	return typedCredentialSpec{
		TerraformSuffix: "credential_slack_api",
		N8nType:         "slackApi",
		DocFile:         "internal/provider/docs/resources/credential_slack_api.md",
		ExtraAttributes: map[string]schema.Attribute{
			"access_token":     stringWriteOnly("Slack bot access token. Write-only; never stored in state.", true),
			"signature_secret": stringWriteOnly("Slack signing secret used to verify requests. Write-only; never stored in state.", false),
		},
		BuildData: func(bag typedAttrBag) (map[string]any, error) {
			m := map[string]any{}
			if err := bagPutRequiredString(m, "accessToken", bag.strings["access_token"]); err != nil {
				return nil, err
			}
			bagPutString(m, "signatureSecret", bag.strings["signature_secret"])
			return m, nil
		},
	}
}

func notionAPISpec() typedCredentialSpec {
	return typedCredentialSpec{
		TerraformSuffix: "credential_notion_api",
		N8nType:         "notionApi",
		DocFile:         "internal/provider/docs/resources/credential_notion_api.md",
		ExtraAttributes: map[string]schema.Attribute{
			"api_key": stringWriteOnly("Notion Internal Integration Secret. Write-only; never stored in state.", true),
		},
		BuildData: func(bag typedAttrBag) (map[string]any, error) {
			m := map[string]any{}
			if err := bagPutRequiredString(m, "apiKey", bag.strings["api_key"]); err != nil {
				return nil, err
			}
			return m, nil
		},
	}
}

func n8nAPISpec() typedCredentialSpec {
	return typedCredentialSpec{
		TerraformSuffix: "credential_n8n_api",
		N8nType:         "n8nApi",
		DocFile:         "internal/provider/docs/resources/credential_n8n_api.md",
		ExtraAttributes: map[string]schema.Attribute{
			"api_key":  stringWriteOnly("n8n API key. Write-only; never stored in state.", true),
			"base_url": stringState("n8n Public API base URL, for example `https://example.app.n8n.cloud/api/v1`.", true),
		},
		BuildData: func(bag typedAttrBag) (map[string]any, error) {
			m := map[string]any{}
			if err := bagPutRequiredString(m, "apiKey", bag.strings["api_key"]); err != nil {
				return nil, err
			}
			if err := bagPutRequiredString(m, "baseUrl", bag.strings["base_url"]); err != nil {
				return nil, err
			}
			return m, nil
		},
	}
}

func jiraSoftwareCloudAPISpec() typedCredentialSpec {
	return typedCredentialSpec{
		TerraformSuffix: "credential_jira_software_cloud_api",
		N8nType:         "jiraSoftwareCloudApi",
		DocFile:         "internal/provider/docs/resources/credential_jira_software_cloud_api.md",
		ExtraAttributes: map[string]schema.Attribute{
			"email":     stringState("Atlassian account email.", true),
			"api_token": stringWriteOnly("Atlassian API token. Write-only; never stored in state.", true),
			"domain":    stringState("Jira Cloud domain, for example `https://example.atlassian.net`.", true),
		},
		BuildData: func(bag typedAttrBag) (map[string]any, error) {
			m := map[string]any{}
			if err := bagPutRequiredString(m, "email", bag.strings["email"]); err != nil {
				return nil, err
			}
			if err := bagPutRequiredString(m, "apiToken", bag.strings["api_token"]); err != nil {
				return nil, err
			}
			if err := bagPutRequiredString(m, "domain", bag.strings["domain"]); err != nil {
				return nil, err
			}
			return m, nil
		},
	}
}

func hubspotAppTokenSpec() typedCredentialSpec {
	return typedCredentialSpec{
		TerraformSuffix: "credential_hubspot_app_token",
		N8nType:         "hubspotAppToken",
		DocFile:         "internal/provider/docs/resources/credential_hubspot_app_token.md",
		ExtraAttributes: map[string]schema.Attribute{
			"app_token": stringWriteOnly("HubSpot private app service key. Write-only; never stored in state.", true),
		},
		BuildData: func(bag typedAttrBag) (map[string]any, error) {
			m := map[string]any{}
			if err := bagPutRequiredString(m, "appToken", bag.strings["app_token"]); err != nil {
				return nil, err
			}
			return m, nil
		},
	}
}

func serpAPISpec() typedCredentialSpec {
	return typedCredentialSpec{
		TerraformSuffix: "credential_serp_api",
		N8nType:         "serpApi",
		DocFile:         "internal/provider/docs/resources/credential_serp_api.md",
		ExtraAttributes: map[string]schema.Attribute{
			"api_key": stringWriteOnly("SerpAPI key. Write-only; never stored in state.", true),
		},
		BuildData: func(bag typedAttrBag) (map[string]any, error) {
			m := map[string]any{}
			if err := bagPutRequiredString(m, "apiKey", bag.strings["api_key"]); err != nil {
				return nil, err
			}
			return m, nil
		},
	}
}

func openAIAPISpec() typedCredentialSpec {
	return typedCredentialSpec{
		TerraformSuffix: "credential_open_ai_api",
		N8nType:         "openAiApi",
		DocFile:         "internal/provider/docs/resources/credential_open_ai_api.md",
		ExtraAttributes: map[string]schema.Attribute{
			"api_key":         stringWriteOnly("OpenAI API key. Write-only; never stored in state.", true),
			"organization_id": stringState("OpenAI organization id. Required only if you belong to multiple organizations.", false),
			"url":             stringState("OpenAI API base URL. Defaults in n8n to `https://api.openai.com/v1` when omitted.", false),
			"header":          boolState("When true, send the custom header_name and header_value."),
			"header_name":     stringState("Custom header name used when header is true.", false),
			"header_value":    stringWriteOnly("Custom header value used when header is true. Write-only; never stored in state.", false),
		},
		BuildData: func(bag typedAttrBag) (map[string]any, error) {
			m := map[string]any{}
			if err := bagPutRequiredString(m, "apiKey", bag.strings["api_key"]); err != nil {
				return nil, err
			}
			bagPutString(m, "organizationId", bag.strings["organization_id"])
			bagPutString(m, "url", bag.strings["url"])
			bagPutBool(m, "header", bag.bools["header"])
			bagPutString(m, "headerName", bag.strings["header_name"])
			bagPutString(m, "headerValue", bag.strings["header_value"])
			return m, nil
		},
	}
}

func jwtAuthSpec() typedCredentialSpec {
	return typedCredentialSpec{
		TerraformSuffix: "credential_jwt_auth",
		N8nType:         "jwtAuth",
		DocFile:         "internal/provider/docs/resources/credential_jwt_auth.md",
		ExtraAttributes: map[string]schema.Attribute{
			"key_type":    stringState("Key type: `passphrase` or `pemKey`.", true),
			"secret":      stringWriteOnly("HMAC secret used when key_type is passphrase. Write-only; never stored in state.", false),
			"private_key": stringWriteOnly("PEM private key used when key_type is pemKey. Write-only; never stored in state.", false),
			"public_key":  stringWriteOnly("PEM public key used when key_type is pemKey. Write-only; never stored in state.", false),
			"algorithm":   stringState("JWT algorithm (for example HS256, RS256). n8n defaults to HS256 when omitted.", false),
		},
		BuildData: func(bag typedAttrBag) (map[string]any, error) {
			m := map[string]any{}
			if err := bagPutRequiredString(m, "keyType", bag.strings["key_type"]); err != nil {
				return nil, err
			}
			bagPutString(m, "secret", bag.strings["secret"])
			bagPutString(m, "privateKey", bag.strings["private_key"])
			bagPutString(m, "publicKey", bag.strings["public_key"])
			bagPutString(m, "algorithm", bag.strings["algorithm"])
			return m, nil
		},
	}
}

func githubAPISpec() typedCredentialSpec {
	return typedCredentialSpec{
		TerraformSuffix: "credential_github_api",
		N8nType:         "githubApi",
		DocFile:         "internal/provider/docs/resources/credential_github_api.md",
		ExtraAttributes: map[string]schema.Attribute{
			"server":       stringState("GitHub API server. Defaults in n8n to `https://api.github.com`. Set this for GitHub Enterprise.", false),
			"user":         stringState("GitHub username associated with the access token.", false),
			"access_token": stringWriteOnly("GitHub personal access token. Write-only; never stored in state.", true),
		},
		BuildData: func(bag typedAttrBag) (map[string]any, error) {
			m := map[string]any{}
			if err := bagPutRequiredString(m, "accessToken", bag.strings["access_token"]); err != nil {
				return nil, err
			}
			bagPutString(m, "server", bag.strings["server"])
			bagPutString(m, "user", bag.strings["user"])
			return m, nil
		},
	}
}

func sendGridAPISpec() typedCredentialSpec {
	return typedCredentialSpec{
		TerraformSuffix: "credential_send_grid_api",
		N8nType:         "sendGridApi",
		DocFile:         "internal/provider/docs/resources/credential_send_grid_api.md",
		ExtraAttributes: map[string]schema.Attribute{
			"api_key": stringWriteOnly("SendGrid API key. Write-only; never stored in state.", true),
		},
		BuildData: func(bag typedAttrBag) (map[string]any, error) {
			m := map[string]any{}
			if err := bagPutRequiredString(m, "apiKey", bag.strings["api_key"]); err != nil {
				return nil, err
			}
			return m, nil
		},
	}
}

func stripeAPISpec() typedCredentialSpec {
	return typedCredentialSpec{
		TerraformSuffix: "credential_stripe_api",
		N8nType:         "stripeApi",
		DocFile:         "internal/provider/docs/resources/credential_stripe_api.md",
		ExtraAttributes: map[string]schema.Attribute{
			"secret_key":       stringWriteOnly("Stripe secret key (`sk_live_` or `sk_test_`). Write-only; never stored in state.", true),
			"signature_secret": stringWriteOnly("Stripe webhook signing secret (`whsec_`). Write-only; never stored in state.", false),
		},
		BuildData: func(bag typedAttrBag) (map[string]any, error) {
			m := map[string]any{}
			if err := bagPutRequiredString(m, "secretKey", bag.strings["secret_key"]); err != nil {
				return nil, err
			}
			bagPutString(m, "signatureSecret", bag.strings["signature_secret"])
			return m, nil
		},
	}
}

func twilioAPISpec() typedCredentialSpec {
	return typedCredentialSpec{
		TerraformSuffix: "credential_twilio_api",
		N8nType:         "twilioApi",
		DocFile:         "internal/provider/docs/resources/credential_twilio_api.md",
		ExtraAttributes: map[string]schema.Attribute{
			"auth_type":      stringState("Auth type: `authToken` or `apiKey`.", true),
			"account_sid":    stringState("Twilio Account SID.", true),
			"auth_token":     stringWriteOnly("Twilio Auth Token used when auth_type is authToken. Write-only; never stored in state.", false),
			"api_key_sid":    stringWriteOnly("Twilio API Key SID used when auth_type is apiKey. Write-only; never stored in state.", false),
			"api_key_secret": stringWriteOnly("Twilio API Key Secret used when auth_type is apiKey. Write-only; never stored in state.", false),
		},
		BuildData: func(bag typedAttrBag) (map[string]any, error) {
			m := map[string]any{}
			if err := bagPutRequiredString(m, "authType", bag.strings["auth_type"]); err != nil {
				return nil, err
			}
			if err := bagPutRequiredString(m, "accountSid", bag.strings["account_sid"]); err != nil {
				return nil, err
			}
			authType := bag.strings["auth_type"].ValueString()
			switch authType {
			case "authToken":
				if err := bagPutRequiredString(m, "authToken", bag.strings["auth_token"]); err != nil {
					return nil, err
				}
			case "apiKey":
				if err := bagPutRequiredString(m, "apiKeySid", bag.strings["api_key_sid"]); err != nil {
					return nil, err
				}
				if err := bagPutRequiredString(m, "apiKeySecret", bag.strings["api_key_secret"]); err != nil {
					return nil, err
				}
			default:
				return nil, fmt.Errorf("authType must be authToken or apiKey")
			}
			return m, nil
		},
	}
}
