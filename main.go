package main

import (
	"context"
	"flag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"
	"github.com/wangsu/terraform-provider-wangsu/wangsu"
	"log"
)

func main() {
	var debugMode bool

	flag.BoolVar(&debugMode, "debuggable", true, "set to true to run the provider with support for debuggers like delve")
	flag.Parse()

	if debugMode {
		err := plugin.Debug(context.Background(), "registry.terraform.io/wangsustack/wangsu",
			&plugin.ServeOpts{
				ProviderFunc: wangsu.Provider,
				Debug:        debugMode,
			})
		if err != nil {
			log.Println(err.Error())
		}
	} else {
		plugin.Serve(&plugin.ServeOpts{ProviderFunc: wangsu.Provider})
	}
}
