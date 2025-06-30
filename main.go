package main

import (
	"context"
	"flag"
	"log"

	provider "github.com/controlplane-com/terraform-provider-cpln/internal/provider"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
)

var (
	// version holds the provider version. It defaults to "dev" for local development
	// but is set by the goreleaser configuration during release builds.
	version string = "dev"

	// Additional metadata, like commit hash or build date, can also be set by goreleaser.
	// This is managed via goreleaser configuration: https://goreleaser.com/cookbooks/using-main.version/
)

func main() {
	var debug bool

	// Define a command-line flag to enable debug mode
	flag.BoolVar(&debug, "debug", false, "set to true to run the provider with support for debuggers like delve")
	flag.Parse()

	// Configure server options, setting the provider address for registry use
	// and enabling debug mode if the debug flag is set.
	opts := providerserver.ServeOpts{
		Address: "registry.terraform.io/controlplane/cpln",
		Debug:   debug,
	}

	// Initialize and serve the provider. The provider's New function takes the version
	// to ensure versioning information is available to users and tools.
	err := providerserver.Serve(context.Background(), provider.New(version), opts)

	// Log a fatal error if the provider server fails to start, terminating execution.
	if err != nil {
		log.Fatal(err.Error())
	}
}
