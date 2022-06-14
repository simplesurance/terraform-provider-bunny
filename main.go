package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"

	"github.com/simplesurance/terraform-provider-bunny/internal/provider"
)

// Run "go generate" to format example terraform files and generate the docs for the registry/website

// If you do not have terraform installed, you can remove the formatting command, but its suggested to
// ensure the documentation is formatted properly.
//go:generate terraform fmt -recursive ./examples/

// Run the docs generation tool, check its repository for more information on how it works and how docs
// can be customized.
//go:generate go run github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs generate --provider-name bunny

func main() {
	var debugMode bool
	var showVersion bool

	flag.BoolVar(&debugMode, "debug", false, "set to true to run the provider with support for debuggers like delve")
	flag.BoolVar(&showVersion, "version", false, "print the version and exit")
	flag.Parse()

	if showVersion {
		fmt.Printf("terraform provider bunny, version %s (%s)\n", provider.Version, provider.Commit)
		os.Exit(0)
	}

	opts := &plugin.ServeOpts{
		ProviderFunc: provider.New,
		Debug:        debugMode,
	}

	plugin.Serve(opts)
}
