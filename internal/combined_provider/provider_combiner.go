package combined_provider

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-mux/tf5to6server"
	"github.com/hashicorp/terraform-plugin-mux/tf6muxserver"
	"github.com/nukosuke/terraform-provider-zendesk/zendesk"
	"log"
)

func BuildMuxProviderServer(pluginFrameworkProvider provider.Provider) (*tfprotov6.ProviderServer, error) {
	ctx := context.Background()

	// upgrade the zendesk provider to the Terraform Plugin Framework (version 6.0)
	upgradedNukosukeZendeskProvider, err2 := tf5to6server.UpgradeServer(
		ctx,
		zendesk.Provider().GRPCProvider,
	)

	if err2 != nil {
		log.Fatal(err2)

	}

	providers := []func() tfprotov6.ProviderServer{
		// andsafe Zendesk provider, developed in the Terraform Plugin Framework (version 6.0)
		providerserver.NewProtocol6(pluginFrameworkProvider),
		// upgraded provider
		func() tfprotov6.ProviderServer {
			return upgradedNukosukeZendeskProvider
		},
	}

	// Mix the two providers
	muxServer, err := tf6muxserver.NewMuxServer(ctx, providers...)
	if err != nil {
		return nil, err
	}
	server := muxServer.ProviderServer()
	return &server, err
}
