package provider

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

type typedCredentialAccCase struct {
	spec        typedCredentialSpec
	createAttrs map[string]string
	updateAttrs map[string]string
}

func TestAccN8nCredentialTyped_basic(t *testing.T) {
	if !isIntegrationTestMode() {
		t.Skip("Skipping acceptance test unless TF_ACC=1")
	}

	for _, tc := range typedCredentialAccCases() {
		t.Run(tc.spec.TerraformSuffix, func(t *testing.T) {
			runTypedCredentialAccBasic(t, tc)
		})
	}
}

func TestAccN8nCredentialTyped_deleteProtection(t *testing.T) {
	if !isIntegrationTestMode() {
		t.Skip("Skipping acceptance test unless TF_ACC=1")
	}

	name := acctest.RandomWithPrefix("tf-n8n-cred-dp")

	protectConfig, err := ReadAccTestResource([]string{"resources", "n8n_credential_http_header_auth", "030_delete_protection.tf"})
	if err != nil {
		t.Fatalf("read protect fixture: %v", err)
	}
	allowDestroyConfig, err := ReadAccTestResource([]string{"resources", "n8n_credential_http_header_auth", "040_allow_destroy.tf"})
	if err != nil {
		t.Fatalf("read allow-destroy fixture: %v", err)
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() { testAccPreCheckCredentials(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_11_0),
		},
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckCredentialDestroy,
		Steps: []resource.TestStep{
			{
				Config: getProviderConfig() + replaceAccName(protectConfig, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("n8n_credential_http_header_auth.test", "delete_protection", "true"),
				),
			},
			{
				Config:      getProviderConfig() + replaceAccName(protectConfig, name),
				Destroy:     true,
				ExpectError: regexp.MustCompile("delete protection is enabled"),
			},
			{
				Config: getProviderConfig() + replaceAccName(allowDestroyConfig, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("n8n_credential_http_header_auth.test", "delete_protection", "false"),
				),
			},
		},
	})
}

func runTypedCredentialAccBasic(t *testing.T, tc typedCredentialAccCase) {
	t.Helper()

	name := acctest.RandomWithPrefix("tf-n8n-cred")
	fixtureDir := "n8n_" + tc.spec.TerraformSuffix
	address := "n8n_" + tc.spec.TerraformSuffix + ".test"
	typed := typedResourceFromSpec(t, tc.spec)

	createConfig, err := ReadAccTestResource([]string{"resources", fixtureDir, "010_create.tf"})
	if err != nil {
		t.Fatalf("read create fixture: %v", err)
	}
	updateConfig, err := ReadAccTestResource([]string{"resources", fixtureDir, "020_update.tf"})
	if err != nil {
		t.Fatalf("read update fixture: %v", err)
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() { testAccPreCheckCredentials(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_11_0),
		},
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckCredentialDestroy,
		Steps: []resource.TestStep{
			{
				Config: getProviderConfig() + replaceAccName(createConfig, name),
				Check:  typedCredentialAccChecks(address, name, "1", tc.spec, typed, tc.createAttrs),
			},
			{
				ResourceName:            address,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: typedImportIgnore(t, tc.spec),
			},
			{
				Config: getProviderConfig() + replaceAccName(updateConfig, name),
				Check:  typedCredentialAccChecks(address, name+"-renamed", "2", tc.spec, typed, tc.updateAttrs),
			},
		},
	})
}

func typedCredentialAccCases() []typedCredentialAccCase {
	oauthClientCreate := map[string]string{"client_id": "placeholder"}
	oauthClientUpdate := map[string]string{"client_id": "placeholder-2"}

	return []typedCredentialAccCase{
		{
			spec:        httpHeaderAuthSpec(),
			createAttrs: map[string]string{"header_name": "X-Test"},
			updateAttrs: map[string]string{"header_name": "X-Test-Updated"},
		},
		{
			spec:        httpBasicAuthSpec(),
			createAttrs: map[string]string{"user": "alice"},
			updateAttrs: map[string]string{"user": "bob"},
		},
		{spec: httpBearerAuthSpec()},
		{spec: slackAPISpec()},
		{spec: notionAPISpec()},
		{
			spec:        n8nAPISpec(),
			createAttrs: map[string]string{"base_url": "https://example.app.n8n.cloud/api/v1"},
			updateAttrs: map[string]string{"base_url": "https://example.app.n8n.cloud/api/v1-updated"},
		},
		{
			spec: jiraSoftwareCloudAPISpec(),
			createAttrs: map[string]string{
				"email":  "user@example.com",
				"domain": "https://example.atlassian.net",
			},
			updateAttrs: map[string]string{
				"email":  "updated@example.com",
				"domain": "https://updated.atlassian.net",
			},
		},
		{spec: hubspotAppTokenSpec()},
		{spec: serpAPISpec()},
		{spec: openAIAPISpec()},
		{
			spec: anthropicAPISpec(),
			updateAttrs: map[string]string{
				"url": "https://api.anthropic.com",
			},
		},
		{
			spec: azureOpenAiAPISpec(),
			createAttrs: map[string]string{
				"resource_name": "tf-acc-aoai",
				"api_version":   "2025-03-01-preview",
			},
			updateAttrs: map[string]string{
				"resource_name": "tf-acc-aoai-2",
				"api_version":   "2024-10-21",
			},
		},
		{
			spec: ollamaAPISpec(),
			createAttrs: map[string]string{
				"base_url": "http://127.0.0.1:11434",
			},
			updateAttrs: map[string]string{
				"base_url": "http://ollama.example.com:11434",
			},
		},
		{
			spec: jwtAuthSpec(),
			createAttrs: map[string]string{
				"key_type":  "passphrase",
				"algorithm": "HS256",
			},
			updateAttrs: map[string]string{
				"key_type":  "passphrase",
				"algorithm": "HS512",
			},
		},
		{
			spec:        googleAPISpec(),
			createAttrs: map[string]string{"email": "sa@example.iam.gserviceaccount.com"},
			updateAttrs: map[string]string{"email": "sa-updated@example.iam.gserviceaccount.com"},
		},
		{
			spec: googlePalmAPISpec(),
			createAttrs: map[string]string{
				"host": "https://generativelanguage.googleapis.com",
			},
			updateAttrs: map[string]string{
				"host": "https://generativelanguage.googleapis.com",
			},
		},
		{spec: gmailOAuth2Spec(), createAttrs: oauthClientCreate, updateAttrs: oauthClientUpdate},
		{spec: googleSheetsTriggerOAuth2APISpec(), createAttrs: oauthClientCreate, updateAttrs: oauthClientUpdate},
		{spec: googleCalendarOAuth2APISpec(), createAttrs: oauthClientCreate, updateAttrs: oauthClientUpdate},
		{spec: googleCloudStorageOAuth2APISpec(), createAttrs: oauthClientCreate, updateAttrs: oauthClientUpdate},
		{spec: googleDriveOAuth2APISpec(), createAttrs: oauthClientCreate, updateAttrs: oauthClientUpdate},
		{spec: githubOAuth2APISpec(), createAttrs: oauthClientCreate, updateAttrs: oauthClientUpdate},
		{
			spec:        salesforceOAuth2APISpec(),
			createAttrs: map[string]string{"client_id": "placeholder"},
			updateAttrs: map[string]string{"client_id": "placeholder-2", "environment": "production"},
		},
		{spec: twitterOAuth2APISpec(), createAttrs: oauthClientCreate, updateAttrs: oauthClientUpdate},
		{spec: microsoftGraphSecurityOAuth2APISpec(), createAttrs: oauthClientCreate, updateAttrs: oauthClientUpdate},
		{
			spec: mcpOAuth2APISpec(),
			createAttrs: map[string]string{
				"use_dynamic_client_registration": "true",
				"server_url":                      "https://mcp.example.com",
			},
			updateAttrs: map[string]string{
				"use_dynamic_client_registration": "true",
				"server_url":                      "https://mcp-updated.example.com",
			},
		},
		{
			spec: githubAPISpec(),
			createAttrs: map[string]string{
				"user": "octocat",
			},
			updateAttrs: map[string]string{
				"user": "octocat-2",
			},
		},
		{spec: sendGridAPISpec()},
		{spec: stripeAPISpec()},
		{
			spec: twilioAPISpec(),
			createAttrs: map[string]string{
				"auth_type":   "authToken",
				"account_sid": "ACxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx",
			},
			updateAttrs: map[string]string{
				"auth_type":   "authToken",
				"account_sid": "ACyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyy",
			},
		},
		{
			spec: smtpSpec(),
			createAttrs: map[string]string{
				"user":   "alerts@example.com",
				"host":   "smtp.example.com",
				"port":   "465",
				"secure": "true",
			},
			updateAttrs: map[string]string{
				"user":   "ops@example.com",
				"host":   "smtp-updated.example.com",
				"port":   "587",
				"secure": "false",
			},
		},
		{
			spec: awsSpec(),
			createAttrs: map[string]string{
				"region":        "us-east-1",
				"access_key_id": "AKIAxxxxxxxx",
			},
			updateAttrs: map[string]string{
				"region":        "eu-west-1",
				"access_key_id": "AKIAyyyyyyyy",
			},
		},
		{spec: googleSheetsOAuth2APISpec(), createAttrs: oauthClientCreate, updateAttrs: oauthClientUpdate},
		{
			spec: oAuth2APISpec(),
			createAttrs: map[string]string{
				"grant_type":       "clientCredentials",
				"access_token_url": "https://example.com/oauth/token",
				"client_id":        "placeholder",
				"authentication":   "header",
			},
			updateAttrs: map[string]string{
				"grant_type":       "clientCredentials",
				"access_token_url": "https://example.com/oauth/token-v2",
				"client_id":        "placeholder-2",
				"authentication":   "body",
			},
		},
		{spec: httpMultipleHeadersAuthSpec()},
	}
}

func typedResourceFromSpec(t *testing.T, spec typedCredentialSpec) *typedCredentialResource {
	t.Helper()
	typed, ok := newTypedCredentialResource(spec).(*typedCredentialResource)
	if !ok {
		t.Fatalf("newTypedCredentialResource(%q) did not return *typedCredentialResource", spec.TerraformSuffix)
	}
	return typed
}

func typedImportIgnore(t *testing.T, spec typedCredentialSpec) []string {
	t.Helper()
	typed := typedResourceFromSpec(t, spec)
	ignore := make([]string, 0, len(typed.secretKeys)+len(typed.stateKeys)+3)
	ignore = append(ignore, typed.secretKeys...)
	ignore = append(ignore, typed.stateKeys...)
	return append(ignore, "data_version", "is_partial_data", "delete_protection")
}

func typedCredentialAccChecks(
	address, name, dataVersion string,
	spec typedCredentialSpec,
	typed *typedCredentialResource,
	extra map[string]string,
) resource.TestCheckFunc {
	checks := []resource.TestCheckFunc{
		resource.TestCheckResourceAttr(address, "name", name),
		resource.TestCheckResourceAttr(address, "type", spec.N8nType),
		resource.TestCheckResourceAttr(address, "delete_protection", "false"),
		resource.TestCheckResourceAttr(address, "data_version", dataVersion),
		resource.TestCheckResourceAttrSet(address, "id"),
	}
	if spec.OAuthPartialDefault {
		checks = append(checks, resource.TestCheckResourceAttr(address, "is_partial_data", "true"))
	}
	for _, key := range typed.secretKeys {
		checks = append(checks, resource.TestCheckNoResourceAttr(address, key))
	}
	for key, value := range extra {
		checks = append(checks, resource.TestCheckResourceAttr(address, key, value))
	}
	return resource.ComposeTestCheckFunc(checks...)
}
