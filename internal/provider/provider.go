// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"strings"
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
	Account types.String `tfsdk:"account"`
	Email   types.String `tfsdk:"email"`
	Token   types.String `tfsdk:"token"`
}

func (p *zendeskProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "zendesk"
	resp.Version = p.version
}

// Schema defines the provider-level schema for configuration data.
func (p *zendeskProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "",
		Attributes: map[string]schema.Attribute{
			"account": schema.StringAttribute{
				Description: "Account name of your Zendesk instance.",
				Optional:    true,
			},
			"email": schema.StringAttribute{
				Description: "Email address of agent user who have permission to access the API.",
				Optional:    true,
			},
			"token": schema.StringAttribute{
				Description: "[API token](https://developer.zendesk.com/rest_api/docs/support/introduction#api-token) for your Zendesk instance.",
				Optional:    true,
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

	if config.Account.IsUnknown() {
		errorSummary := "Unknown Zendesk API Account"
		tflog.Error(ctx, errorSummary)
		resp.Diagnostics.AddAttributeError(
			path.Root("account"),
			errorSummary,
			"The provider cannot create the Zendesk API client as there is an unknown configuration value for the Zendesk API account URL. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the ZENDESK_ACCOUNT environment variable.",
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

	if config.Token.IsUnknown() {
		errorSummary := "Unknown Zendesk API Token"
		tflog.Error(ctx, errorSummary)
		resp.Diagnostics.AddAttributeError(
			path.Root("token"),
			errorSummary,
			"The provider cannot create the Zendesk API client as there is an unknown configuration value for the Zendesk API Token. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the ZENDESK_TOKEN environment variable.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	var hostUrl, email, apiToken string

	if !config.Account.IsNull() {
		baseURLFormat := "https://%s.zendesk.com"
		if strings.Contains(config.Account.ValueString(), "localhost") {
			hostUrl = "http://localhost:"
			if strings.Contains(config.Account.ValueString(), "Port") {
				hostUrl += strings.Replace(config.Account.ValueString(), "localhostPort", "", 1)
			} else {
				hostUrl += "8080"
			}
		} else {
			hostUrl = fmt.Sprintf(baseURLFormat, config.Account.ValueString())
		}

		tflog.Debug(ctx, "Used account from configuration")
	}

	if !config.Email.IsNull() {
		email = config.Email.ValueString()
		tflog.Debug(ctx, "Used email from configuration")
	}

	if !config.Token.IsNull() {
		apiToken = config.Token.ValueString()
		tflog.Debug(ctx, "Used token from configuration")
	}

	// If any of the expected configurations are missing, return
	// errors with provider-specific guidance.

	if hostUrl == "" {
		tflog.Error(ctx, "Missing Zendesk API Account")
		resp.Diagnostics.AddAttributeError(
			path.Root("account"),
			"Missing Zendesk API Account",
			"The provider cannot create the Zendesk API client as there is a missing or empty value for the Zendesk API account. "+
				"Set the account value in the configuration or use the ZENDESK_ACCOUNT environment variable. "+
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
		tflog.Error(ctx, "Missing Zendesk API Token")
		resp.Diagnostics.AddAttributeError(
			path.Root("token"),
			"Missing Zendesk API Token",
			"The provider cannot create the Zendesk API client as there is a missing or empty value for the Zendesk API token. "+
				"Set the token value in the configuration or use the ZENDESK_TOKEN environment variable. "+
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
