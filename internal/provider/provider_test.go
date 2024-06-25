// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
)

// testAccProtoV6ProviderFactories are used to instantiate a provider during
// acceptance testing. The factory function will be invoked for every Terraform
// CLI command executed to create a provider server to which the CLI can
// reattach.
var testAccProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
	"zendesk": providerserver.NewProtocol6WithError(New("test")()),
}

const (
	// providerConfig is a shared configuration to combine with the actual
	// test configuration so the Zendesk client is properly configured.
	providerConfig = `
provider "zendesk" {
  email = "education"
  api_token = "test123"
  host_url     = "http://localhost:8080"
}
`
)
