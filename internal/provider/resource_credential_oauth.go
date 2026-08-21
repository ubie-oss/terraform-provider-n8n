package provider

import (
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

func NewCredentialMCPOAuth2APIResource() resource.Resource {
	return newTypedCredentialResource(mcpOAuth2APISpec())
}

func NewCredentialGitHubOAuth2APIResource() resource.Resource {
	return newTypedCredentialResource(githubOAuth2APISpec())
}

func NewCredentialSalesforceOAuth2APIResource() resource.Resource {
	return newTypedCredentialResource(salesforceOAuth2APISpec())
}

func NewCredentialTwitterOAuth2APIResource() resource.Resource {
	return newTypedCredentialResource(twitterOAuth2APISpec())
}

func NewCredentialMicrosoftGraphSecurityOAuth2APIResource() resource.Resource {
	return newTypedCredentialResource(microsoftGraphSecurityOAuth2APISpec())
}

func NewCredentialOAuth2APIResource() resource.Resource {
	return newTypedCredentialResource(oAuth2APISpec())
}

func mcpOAuth2APISpec() typedCredentialSpec {
	attrs := oauthClientAttributes()
	attrs["client_id"] = stringState("OAuth client ID. Required when use_dynamic_client_registration is false.", false)
	attrs["client_secret"] = stringWriteOnly("OAuth client secret. Write-only. Required when use_dynamic_client_registration is false.", false)
	attrs["use_dynamic_client_registration"] = boolState("When true, n8n uses OAuth 2.0 dynamic client registration. n8n defaults to true.")
	attrs["resource_url"] = stringState("Optional protected resource URL required by the MCP server.", false)
	attrs["server_url"] = stringState("Authorization server URL used when use_dynamic_client_registration is true.", false)
	attrs["grant_type"] = stringState("OAuth grant type: `authorizationCode`, `clientCredentials`, or `pkce`. Used when dynamic client registration is false.", false)
	attrs["auth_url"] = stringState("Authorization URL used when dynamic client registration is false.", false)
	attrs["access_token_url"] = stringState("Access token URL used when dynamic client registration is false.", false)
	attrs["scope"] = stringState("OAuth scope used when dynamic client registration is false.", false)
	return typedCredentialSpec{
		TerraformSuffix:     "credential_mcp_oauth2_api",
		N8nType:             "mcpOAuth2Api",
		OAuthPartialDefault: true,
		DocFile:             "internal/provider/docs/resources/credential_mcp_oauth2_api.md",
		ExtraAttributes:     attrs,
		BuildData: func(bag typedAttrBag) (map[string]any, error) {
			m := map[string]any{}
			bagPutBool(m, "useDynamicClientRegistration", bag.bools["use_dynamic_client_registration"])
			bagPutString(m, "resourceUrl", bag.strings["resource_url"])
			bagPutString(m, "serverUrl", bag.strings["server_url"])
			bagPutString(m, "grantType", bag.strings["grant_type"])
			bagPutString(m, "authUrl", bag.strings["auth_url"])
			bagPutString(m, "accessTokenUrl", bag.strings["access_token_url"])
			bagPutString(m, "clientId", bag.strings["client_id"])
			bagPutString(m, "clientSecret", bag.strings["client_secret"])
			bagPutString(m, "scope", bag.strings["scope"])
			bagPutBool(m, "ignoreSSLIssues", bag.bools["ignore_ssl_issues"])
			token, ok, err := bagDynamicValue(bag.dynamics["oauth_token_data"], "oauth_token_data")
			if err != nil {
				return nil, err
			}
			if ok {
				m["oauthTokenData"] = token
			}
			return m, nil
		},
	}
}

func githubOAuth2APISpec() typedCredentialSpec {
	attrs := oauthClientAttributes()
	attrs["server"] = stringState("GitHub API server. Defaults in n8n to `https://api.github.com`. Set this for GitHub Enterprise.", false)
	return typedCredentialSpec{
		TerraformSuffix:     "credential_github_oauth2_api",
		N8nType:             "githubOAuth2Api",
		OAuthPartialDefault: true,
		DocFile:             "internal/provider/docs/resources/credential_github_oauth2_api.md",
		ExtraAttributes:     attrs,
		BuildData: func(bag typedAttrBag) (map[string]any, error) {
			m, err := oauthBuildData(bag, false, true)
			if err != nil {
				return nil, err
			}
			bagPutString(m, "server", bag.strings["server"])
			return m, nil
		},
	}
}

func salesforceOAuth2APISpec() typedCredentialSpec {
	attrs := oauthClientAttributes()
	attrs["environment"] = stringState("Salesforce environment: `production` or `sandbox`.", false)
	return typedCredentialSpec{
		TerraformSuffix:     "credential_salesforce_oauth2_api",
		N8nType:             "salesforceOAuth2Api",
		OAuthPartialDefault: true,
		DocFile:             "internal/provider/docs/resources/credential_salesforce_oauth2_api.md",
		ExtraAttributes:     attrs,
		BuildData: func(bag typedAttrBag) (map[string]any, error) {
			m, err := oauthBuildData(bag, false, true)
			if err != nil {
				return nil, err
			}
			bagPutString(m, "environment", bag.strings["environment"])
			return m, nil
		},
	}
}

func twitterOAuth2APISpec() typedCredentialSpec {
	return typedCredentialSpec{
		TerraformSuffix:     "credential_twitter_oauth2_api",
		N8nType:             "twitterOAuth2Api",
		OAuthPartialDefault: true,
		DocFile:             "internal/provider/docs/resources/credential_twitter_oauth2_api.md",
		ExtraAttributes:     oauthClientAttributes(),
		BuildData: func(bag typedAttrBag) (map[string]any, error) {
			return oauthBuildData(bag, false, true)
		},
	}
}

func microsoftGraphSecurityOAuth2APISpec() typedCredentialSpec {
	attrs := googleOAuth2Attributes()
	attrs["client_secret"] = stringWriteOnly("OAuth client secret. Write-only. Required when client_credential_type is clientSecret.", false)
	attrs["client_credential_type"] = stringState("How n8n authenticates to Entra: `clientSecret` or `certificate`.", false)
	attrs["private_key"] = stringWriteOnly("PEM private key used when client_credential_type is certificate. Write-only; never stored in state.", false)
	attrs["certificate"] = stringWriteOnly("PEM public certificate used when client_credential_type is certificate. Write-only; never stored in state.", false)
	attrs["auth_url"] = stringState("Microsoft authorization URL. n8n defaults to the common v2 authorize endpoint.", false)
	attrs["access_token_url"] = stringState("Microsoft token URL. n8n defaults to the common v2 token endpoint.", false)
	attrs["graph_api_base_url"] = stringState("Microsoft Graph base URL (global, US Government, DoD, or China).", false)
	return typedCredentialSpec{
		TerraformSuffix:     "credential_microsoft_graph_security_oauth2_api",
		N8nType:             "microsoftGraphSecurityOAuth2Api",
		OAuthPartialDefault: true,
		DocFile:             "internal/provider/docs/resources/credential_microsoft_graph_security_oauth2_api.md",
		ExtraAttributes:     attrs,
		BuildData: func(bag typedAttrBag) (map[string]any, error) {
			m, err := oauthBuildData(bag, true, false)
			if err != nil {
				return nil, err
			}
			bagPutString(m, "clientCredentialType", bag.strings["client_credential_type"])
			bagPutString(m, "privateKey", bag.strings["private_key"])
			bagPutString(m, "certificate", bag.strings["certificate"])
			bagPutString(m, "authUrl", bag.strings["auth_url"])
			bagPutString(m, "accessTokenUrl", bag.strings["access_token_url"])
			bagPutString(m, "graphApiBaseUrl", bag.strings["graph_api_base_url"])
			return m, nil
		},
	}
}

func oAuth2APISpec() typedCredentialSpec {
	return typedCredentialSpec{
		TerraformSuffix:     "credential_oauth2_api",
		N8nType:             "oAuth2Api",
		OAuthPartialDefault: true,
		DocFile:             "internal/provider/docs/resources/credential_oauth2_api.md",
		ExtraAttributes: map[string]schema.Attribute{
			"grant_type":                      stringState("OAuth grant type: `authorizationCode`, `clientCredentials`, or `pkce`.", true),
			"auth_url":                        stringState("Authorization URL used for authorizationCode and pkce grants.", false),
			"access_token_url":                stringState("Access token URL.", true),
			"client_id":                       stringState("OAuth client ID.", true),
			"client_secret":                   stringWriteOnly("OAuth client secret. Write-only; never stored in state.", true),
			"scope":                           stringState("OAuth scopes.", false),
			"auth_query_parameters":           stringState("Extra authorization URI query parameters, for example `access_type=offline`.", false),
			"authentication":                  stringState("How to send client credentials to the token endpoint: `header` or `body`.", false),
			"send_additional_body_properties": boolState("When true with clientCredentials + body auth, send additional_body_properties."),
			"additional_body_properties":      stringState("JSON object of extra body properties for clientCredentials body auth.", false),
			"ignore_ssl_issues":               boolState("Whether to ignore TLS certificate issues when talking to the token endpoint."),
			"token_expired_status_code":       int64State("HTTP status that marks an expired token. n8n defaults to 401.", false),
			"oauth_token_data": dynamicWriteOnly(
				"OAuth token payload (`oauthTokenData`), typically including access_token and refresh_token. Write-only. The Public API cannot complete a browser OAuth flow.",
			),
		},
		BuildData: func(bag typedAttrBag) (map[string]any, error) {
			m := map[string]any{}
			if err := bagPutRequiredString(m, "grantType", bag.strings["grant_type"]); err != nil {
				return nil, err
			}
			if err := bagPutRequiredString(m, "accessTokenUrl", bag.strings["access_token_url"]); err != nil {
				return nil, err
			}
			if err := bagPutRequiredString(m, "clientId", bag.strings["client_id"]); err != nil {
				return nil, err
			}
			if err := bagPutRequiredString(m, "clientSecret", bag.strings["client_secret"]); err != nil {
				return nil, err
			}
			bagPutString(m, "authUrl", bag.strings["auth_url"])
			bagPutString(m, "scope", bag.strings["scope"])
			bagPutString(m, "authQueryParameters", bag.strings["auth_query_parameters"])
			bagPutString(m, "authentication", bag.strings["authentication"])
			bagPutBool(m, "sendAdditionalBodyProperties", bag.bools["send_additional_body_properties"])
			bagPutString(m, "additionalBodyProperties", bag.strings["additional_body_properties"])
			bagPutBool(m, "ignoreSSLIssues", bag.bools["ignore_ssl_issues"])
			bagPutInt64(m, "tokenExpiredStatusCode", bag.int64s["token_expired_status_code"])
			token, ok, err := bagDynamicValue(bag.dynamics["oauth_token_data"], "oauth_token_data")
			if err != nil {
				return nil, err
			}
			if ok {
				m["oauthTokenData"] = token
			}
			return m, nil
		},
	}
}
