// Copyright 2026 Ubie, inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package provider

import (
	"context"
	"os"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/ubie-oss/terraform-provider-n8n/internal/n8n"
)

// Ensure n8nProvider satisfies the provider interface.
var _ provider.Provider = &n8nProvider{}

// n8nProvider defines the provider implementation.
type n8nProvider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
}

// n8nProviderModel describes the provider data model.
type n8nProviderModel struct {
	Endpoint              types.String  `tfsdk:"endpoint"`
	APIKey                types.String  `tfsdk:"api_key"`
	MaxConcurrentRequests types.Int64   `tfsdk:"max_concurrent_requests"`
	RequestsPerSecond     types.Float64 `tfsdk:"requests_per_second"`
}

func (p *n8nProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "n8n"
	resp.Version = p.version
}

func (p *n8nProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manage n8n team projects, folders, and credentials via the n8n Public API.",
		Attributes: map[string]schema.Attribute{
			"endpoint": schema.StringAttribute{
				MarkdownDescription: "n8n instance base URL (for example `https://n8n.example.com` or `https://<subdomain>.app.n8n.cloud`). `/api/v1` is appended when missing. May also be set via the `N8N_ENDPOINT` environment variable.",
				Optional:            true,
			},
			"api_key": schema.StringAttribute{
				MarkdownDescription: "n8n API key sent as the `X-N8N-API-KEY` header. May also be set via the `N8N_API_KEY` environment variable.",
				Optional:            true,
				Sensitive:           true,
			},
			"max_concurrent_requests": schema.Int64Attribute{
				MarkdownDescription: "Maximum number of concurrent HTTP requests to the API. Defaults to 10 when omitted.",
				Optional:            true,
			},
			"requests_per_second": schema.Float64Attribute{
				MarkdownDescription: "Average sustained API requests per second (token bucket). Defaults to 10 when omitted.",
				Optional:            true,
			},
		},
	}
}

func (p *n8nProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var config n8nProviderModel
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if config.Endpoint.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("endpoint"),
			"Unknown API Endpoint",
			"The provider cannot configure the API client when `endpoint` is unknown.",
		)
		return
	}
	if config.APIKey.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("api_key"),
			"Unknown API Key",
			"The provider cannot configure the API client when `api_key` is unknown.",
		)
		return
	}

	endpoint := os.Getenv("N8N_ENDPOINT")
	if !config.Endpoint.IsNull() && config.Endpoint.ValueString() != "" {
		endpoint = config.Endpoint.ValueString()
	}
	apiKey := os.Getenv("N8N_API_KEY")
	if !config.APIKey.IsNull() && config.APIKey.ValueString() != "" {
		apiKey = config.APIKey.ValueString()
	}

	if endpoint == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("endpoint"),
			"Missing Endpoint",
			"Set the `endpoint` provider attribute or the `N8N_ENDPOINT` environment variable.",
		)
		return
	}
	if apiKey == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("api_key"),
			"Missing API Key",
			"Set the `api_key` provider attribute or the `N8N_API_KEY` environment variable.",
		)
		return
	}

	opts := &n8n.Options{}
	if !config.MaxConcurrentRequests.IsNull() && !config.MaxConcurrentRequests.IsUnknown() {
		opts.MaxConcurrent = config.MaxConcurrentRequests.ValueInt64()
	}
	if !config.RequestsPerSecond.IsNull() && !config.RequestsPerSecond.IsUnknown() {
		opts.RPS = config.RequestsPerSecond.ValueFloat64()
	}

	client, err := n8n.New(endpoint, apiKey, opts)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Configure API Client",
			err.Error(),
		)
		return
	}
	resp.DataSourceData = client
	resp.ResourceData = client
}

func (p *n8nProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewProjectResource,
		NewFolderResource,
		NewCredentialResource,
		NewCredentialHTTPHeaderAuthResource,
		NewCredentialHTTPBasicAuthResource,
		NewCredentialHTTPBearerAuthResource,
		NewCredentialSlackAPIResource,
		NewCredentialNotionAPIResource,
		NewCredentialN8nAPIResource,
		NewCredentialJiraSoftwareCloudAPIResource,
		NewCredentialHubspotAppTokenResource,
		NewCredentialSerpAPIResource,
		NewCredentialOpenAIAPIResource,
		NewCredentialAnthropicAPIResource,
		NewCredentialAzureOpenAIAPIResource,
		NewCredentialOllamaAPIResource,
		NewCredentialJWTAuthResource,
		NewCredentialGmailOAuth2Resource,
		NewCredentialGoogleSheetsTriggerOAuth2APIResource,
		NewCredentialGoogleCalendarOAuth2APIResource,
		NewCredentialGoogleCloudStorageOAuth2APIResource,
		NewCredentialGoogleDriveOAuth2APIResource,
		NewCredentialGoogleAPIResource,
		NewCredentialGooglePalmAPIResource,
		NewCredentialMCPOAuth2APIResource,
		NewCredentialGitHubOAuth2APIResource,
		NewCredentialSalesforceOAuth2APIResource,
		NewCredentialTwitterOAuth2APIResource,
		NewCredentialMicrosoftGraphSecurityOAuth2APIResource,
		NewCredentialGitHubAPIResource,
		NewCredentialSendGridAPIResource,
		NewCredentialStripeAPIResource,
		NewCredentialTwilioAPIResource,
		NewCredentialSMTPResource,
		NewCredentialAWSResource,
		NewCredentialGoogleSheetsOAuth2APIResource,
		NewCredentialOAuth2APIResource,
		NewCredentialHTTPMultipleHeadersAuthResource,
	}
}

func (p *n8nProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewProjectDataSource,
		NewProjectsDataSource,
		NewFolderDataSource,
		NewFoldersDataSource,
		NewCredentialDataSource,
		NewCredentialsDataSource,
		NewCredentialSchemaDataSource,
	}
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &n8nProvider{
			version: version,
		}
	}
}
