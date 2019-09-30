package main

import (
	"github.com/hashicorp/terraform-plugin-sdk/plugin"
	"github.com/terraform-providers/terraform-provider-rancher/rancher"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: rancher.Provider})
}
