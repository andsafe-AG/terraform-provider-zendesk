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
	//"terraform-provider-zendesk/internal/provider"
)

func BuildMuxProviderServer(pluginFrameworkProvider provider.Provider) (*tfprotov6.ProviderServer, error) {
	ctx := context.Background()

	upgradedNukosukeZendeskProvider, err2 := tf5to6server.UpgradeServer(
		ctx,
		zendesk.Provider().GRPCProvider,
	)

	if err2 != nil {
		log.Fatal(err2)

	}

	providers := []func() tfprotov6.ProviderServer{
		providerserver.NewProtocol6(pluginFrameworkProvider), // Example terraform-plugin-framework provider
		func() tfprotov6.ProviderServer {
			return upgradedNukosukeZendeskProvider
		},
	}

	muxServer, err := tf6muxserver.NewMuxServer(ctx, providers...)
	if err != nil {
		return nil, err
	}
	server := muxServer.ProviderServer()
	return &server, err
}
