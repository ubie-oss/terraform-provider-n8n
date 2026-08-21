package provider

import (
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

func NewCredentialGmailOAuth2Resource() resource.Resource {
	return newTypedCredentialResource(gmailOAuth2Spec())
}

func NewCredentialGoogleSheetsTriggerOAuth2APIResource() resource.Resource {
	return newTypedCredentialResource(googleSheetsTriggerOAuth2APISpec())
}

func NewCredentialGoogleCalendarOAuth2APIResource() resource.Resource {
	return newTypedCredentialResource(googleCalendarOAuth2APISpec())
}

func NewCredentialGoogleCloudStorageOAuth2APIResource() resource.Resource {
	return newTypedCredentialResource(googleCloudStorageOAuth2APISpec())
}

func NewCredentialGoogleDriveOAuth2APIResource() resource.Resource {
	return newTypedCredentialResource(googleDriveOAuth2APISpec())
}

func NewCredentialGoogleAPIResource() resource.Resource {
	return newTypedCredentialResource(googleAPISpec())
}

func NewCredentialGooglePalmAPIResource() resource.Resource {
	return newTypedCredentialResource(googlePalmAPISpec())
}

func NewCredentialGoogleSheetsOAuth2APIResource() resource.Resource {
	return newTypedCredentialResource(googleSheetsOAuth2APISpec())
}

func gmailOAuth2Spec() typedCredentialSpec {
	return googleOAuth2Spec("credential_gmail_oauth2", "gmailOAuth2", "internal/provider/docs/resources/credential_gmail_oauth2.md")
}

func googleSheetsOAuth2APISpec() typedCredentialSpec {
	return googleOAuth2Spec("credential_google_sheets_oauth2_api", "googleSheetsOAuth2Api", "internal/provider/docs/resources/credential_google_sheets_oauth2_api.md")
}

func googleSheetsTriggerOAuth2APISpec() typedCredentialSpec {
	return googleOAuth2Spec("credential_google_sheets_trigger_oauth2_api", "googleSheetsTriggerOAuth2Api", "internal/provider/docs/resources/credential_google_sheets_trigger_oauth2_api.md")
}

func googleCalendarOAuth2APISpec() typedCredentialSpec {
	return googleOAuth2Spec("credential_google_calendar_oauth2_api", "googleCalendarOAuth2Api", "internal/provider/docs/resources/credential_google_calendar_oauth2_api.md")
}

func googleCloudStorageOAuth2APISpec() typedCredentialSpec {
	return googleOAuth2Spec("credential_google_cloud_storage_oauth2_api", "googleCloudStorageOAuth2Api", "internal/provider/docs/resources/credential_google_cloud_storage_oauth2_api.md")
}

func googleDriveOAuth2APISpec() typedCredentialSpec {
	return googleOAuth2Spec("credential_google_drive_oauth2_api", "googleDriveOAuth2Api", "internal/provider/docs/resources/credential_google_drive_oauth2_api.md")
}

func googleOAuth2Spec(suffix, n8nType, docFile string) typedCredentialSpec {
	return typedCredentialSpec{
		TerraformSuffix:     suffix,
		N8nType:             n8nType,
		OAuthPartialDefault: true,
		DocFile:             docFile,
		ExtraAttributes:     googleOAuth2Attributes(),
		BuildData:           googleOAuth2BuildData,
	}
}

func googleAPISpec() typedCredentialSpec {
	return typedCredentialSpec{
		TerraformSuffix: "credential_google_api",
		N8nType:         "googleApi",
		DocFile:         "internal/provider/docs/resources/credential_google_api.md",
		ExtraAttributes: map[string]schema.Attribute{
			"region":          stringState("Google Cloud region or multi-region (for example `global`, `us`, `eu`). n8n defaults to `global` when omitted.", false),
			"email":           stringState("Google service account email.", true),
			"private_key":     stringWriteOnly("Service account private key (PEM). Write-only; never stored in state.", true),
			"impersonate":     boolState("When true, n8n impersonates delegated_email. Sent as n8n data key `inpersonate`."),
			"delegated_email": stringState("User email to impersonate when impersonate is true.", false),
			"http_node":       boolState("When true, configure this credential for the HTTP Request node and send scopes."),
			"scopes":          stringState("OAuth scopes used when http_node is true.", false),
		},
		BuildData: func(bag typedAttrBag) (map[string]any, error) {
			m := map[string]any{}
			if err := bagPutRequiredString(m, "email", bag.strings["email"]); err != nil {
				return nil, err
			}
			if err := bagPutRequiredString(m, "privateKey", bag.strings["private_key"]); err != nil {
				return nil, err
			}
			bagPutString(m, "region", bag.strings["region"])
			bagPutBool(m, "inpersonate", bag.bools["impersonate"])
			bagPutString(m, "delegatedEmail", bag.strings["delegated_email"])
			bagPutBool(m, "httpNode", bag.bools["http_node"])
			bagPutString(m, "scopes", bag.strings["scopes"])
			return m, nil
		},
	}
}

func googlePalmAPISpec() typedCredentialSpec {
	return typedCredentialSpec{
		TerraformSuffix: "credential_google_palm_api",
		N8nType:         "googlePalmApi",
		DocFile:         "internal/provider/docs/resources/credential_google_palm_api.md",
		ExtraAttributes: map[string]schema.Attribute{
			"host":    stringState("Gemini/PaLM API host. n8n defaults to `https://generativelanguage.googleapis.com` when omitted.", false),
			"api_key": stringWriteOnly("Google Gemini API key. Write-only; never stored in state.", true),
		},
		BuildData: func(bag typedAttrBag) (map[string]any, error) {
			m := map[string]any{}
			if err := bagPutRequiredString(m, "apiKey", bag.strings["api_key"]); err != nil {
				return nil, err
			}
			bagPutString(m, "host", bag.strings["host"])
			return m, nil
		},
	}
}
