// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"terraform-provider-zendesk/internal/combined_provider"
)

// testAccProtoV6ProviderFactories are used to instantiate a provider during
// acceptance testing. The factory function will be invoked for every Terraform
// CLI command executed to create a provider server to which the CLI can
// reattach.
var testAccProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
	"zendesk": func() (tfprotov6.ProviderServer, error) {
		server, err := combined_provider.BuildMuxProviderServer(New("test")())
		return *server, err
	},
}

/*
var oldProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
	"zendesk": providerserver.NewProtocol6WithError(New("test")()),
}*/

const (
	// providerConfig is a shared configuration to combine with the actual
	// test configuration so the Zendesk client is properly configured.
	providerConfig = `
provider "zendesk" {
  email = "education"
  token = "test123"
  account     = "localhost"
}
`
)
