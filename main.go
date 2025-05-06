package main

import (
	"context"
	"flag"
	"log"

	"github.com/cofide/terraform-provider-cofide/internal"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
)

// Format terraform and generate docs:
//go:generate terraform fmt -recursive ./examples/
//go:generate go run github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs

var (
	// These will be set by the goreleaser configuration
	// to appropriate values for the compiled binary.
	version string = "dev"
)

func main() {
	var debug bool

	flag.BoolVar(&debug, "debug", false, "set to true to run the provider with support for debuggers like delve")
	flag.Parse()

	opts := providerserver.ServeOpts{
		Address: "registry.terraform.io/cofide/cofide",
		Debug:   debug,
	}

	err := providerserver.Serve(context.Background(), internal.NewProvider(version), opts)
	if err != nil {
		log.Fatal(err.Error())
	}
}
