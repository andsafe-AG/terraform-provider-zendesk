// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"terraform-provider-zendesk/zendesk_api"
)

// Ensure zendeskProvider satisfies various provider interfaces.
var _ provider.Provider = &zendeskProvider{}

// zendeskProvider defines the provider implementation.
type zendeskProvider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
}

// zendeskProviderModel maps provider schema data to a Go type.
type zendeskProviderModel struct {
	HostUrl  types.String `tfsdk:"host_url"`
	Email    types.String `tfsdk:"email"`
	ApiToken types.String `tfsdk:"api_token"`
}

func (p *zendeskProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "zendesk"
	resp.Version = p.version
}

// Schema defines the provider-level schema for configuration data.
func (p *zendeskProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "The Zendesk provider allows you to interact with the Zendesk API.",
		Attributes: map[string]schema.Attribute{
			"host_url": schema.StringAttribute{
				Description: "The base URL of your Zendesk instance.",
				Optional:    false,
				Required:    true,
			},
			"email": schema.StringAttribute{
				Description: "The email address of the user to authenticate with. It will be masked.",
				Optional:    false,
				Required:    true,
				Sensitive:   true,
			},
			"api_token": schema.StringAttribute{
				Description: "The API token to authenticate with.",
				Optional:    false,
				Required:    true,
				Sensitive:   true,
			},
		},
	}
}

func (p *zendeskProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	// Retrieve provider data from configuration
	tflog.Info(ctx, "Configuring Zendesk Provider")

	var config zendeskProviderModel
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// If practitioner provided a configuration value for any of the
	// attributes, it must be a known value.

	if config.HostUrl.IsUnknown() {
		errorSummary := "Unknown Zendesk API HostUrl"
		tflog.Error(ctx, errorSummary)
		resp.Diagnostics.AddAttributeError(
			path.Root("host_url"),
			errorSummary,
			"The provider cannot create the Zendesk API client as there is an unknown configuration value for the Zendesk API host_url URL. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the ZENDESK_HOST_URL environment variable.",
		)
	}

	if config.Email.IsUnknown() {
		errorSummary := "Unknown Zendesk API Email"
		tflog.Error(ctx, errorSummary)
		resp.Diagnostics.AddAttributeError(
			path.Root("email"),
			errorSummary,
			"The provider cannot create the Zendesk API client as there is an unknown configuration value for the Zendesk API Email Address. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the ZENDESK_EMAIL environment variable.",
		)
	}

	if config.ApiToken.IsUnknown() {
		errorSummary := "Unknown Zendesk API ApiToken"
		tflog.Error(ctx, errorSummary)
		resp.Diagnostics.AddAttributeError(
			path.Root("api_token"),
			errorSummary,
			"The provider cannot create the Zendesk API client as there is an unknown configuration value for the Zendesk API ApiToken. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the ZENDESK_API_TOKEN environment variable.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	var hostUrl, email, apiToken string

	if !config.HostUrl.IsNull() {
		hostUrl = config.HostUrl.ValueString()
		tflog.Debug(ctx, "Used host_url from configuration")
	}

	if !config.Email.IsNull() {
		email = config.Email.ValueString()
		tflog.Debug(ctx, "Used email from configuration")
	}

	if !config.ApiToken.IsNull() {
		apiToken = config.ApiToken.ValueString()
		tflog.Debug(ctx, "Used api_token from configuration")
	}

	// If any of the expected configurations are missing, return
	// errors with provider-specific guidance.

	if hostUrl == "" {
		tflog.Error(ctx, "Missing Zendesk API HostUrl")
		resp.Diagnostics.AddAttributeError(
			path.Root("host_url"),
			"Missing Zendesk API HostUrl",
			"The provider cannot create the Zendesk API client as there is a missing or empty value for the Zendesk API host_url. "+
				"Set the host_url value in the configuration or use the ZENDESK_HOST_URL environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
	}

	if email == "" {
		tflog.Error(ctx, "Missing Zendesk API Email")
		resp.Diagnostics.AddAttributeError(
			path.Root("email"),
			"Missing Zendesk API Email",
			"The provider cannot create the Zendesk API client as there is a missing or empty value for the Zendesk API email. "+
				"Set the email value in the configuration or use the ZENDESK_EMAIL environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
	}

	if apiToken == "" {
		tflog.Error(ctx, "Missing Zendesk API ApiToken")
		resp.Diagnostics.AddAttributeError(
			path.Root("api_token"),
			"Missing Zendesk API ApiToken",
			"The provider cannot create the Zendesk API client as there is a missing or empty value for the Zendesk API api_token. "+
				"Set the api_token value in the configuration or use the ZENDESK_API_TOKEN environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "Creating Zendesk API client")
	// Make the Zendesk client available during DataSource and Resource
	// type Configure methods.*/
	client :=
		zendesk_api.NewSupportApi(
			hostUrl,
			email,
			apiToken,
		)

	resp.DataSourceData = client
	resp.ResourceData = client
	tflog.Info(ctx, "Zendesk Provider configured successfully")
}

func (p *zendeskProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewCustomStatusResource,
	}
}

func (p *zendeskProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	//TODO
	return nil
}

// New is a helper function to simplify provider server and testing implementation.
func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &zendeskProvider{
			version: version,
		}
	}
}
