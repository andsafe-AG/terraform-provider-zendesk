// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package main

import (
	"flag"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6/tf6server"
	"log"
	"terraform-provider-zendesk/internal/combined_provider"
	"terraform-provider-zendesk/internal/provider"
)

// Run "go generate" to format example terraform files and generate the docs for the registry/website

// If you do not have terraform installed, you can remove the formatting command, but its suggested to
// ensure the documentation is formatted properly.
//go:generate terraform fmt -recursive ./examples/

// Run the docs generation tool, check its repository for more information on how it works and how docs
// can be customized.
//go:generate go run github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs generate -provider-name zendesk

var (
	// these will be set by the goreleaser configuration
	// to appropriate values for the compiled binary.
	version string = "dev"

	// goreleaser can pass other information to the main package, such as the specific commit
	// https://goreleaser.com/cookbooks/using-main.version/
)

func main() {
	var debug bool
	var err error

	flag.BoolVar(&debug, "debug", false, "set to true to run the provider with support for debuggers like delve")
	flag.Parse()

	var serveOpts []tf6server.ServeOpt

	if debug {
		serveOpts = append(serveOpts, tf6server.WithManagedDebug())
	}

	err = tf6server.Serve(
		"registry.terraform.io/andsafe-AG/zendesk",
		func() tfprotov6.ProviderServer {
			muxServer, err2 := combined_provider.BuildMuxProviderServer(provider.New(version)())
			if err2 != nil {
				log.Fatal(err2)
			}
			return *muxServer
		},
		serveOpts...,
	)

	if err != nil {
		log.Fatal(err)
	}
}
