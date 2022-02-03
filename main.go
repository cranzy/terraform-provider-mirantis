package main

import (
  "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
  "github.com/hashicorp/terraform-plugin-sdk/v2/plugin"

  "github.com/Mirantis/terraform-provider-mirantis/mirantis"
)

func main() {
  plugin.Serve(&plugin.ServeOpts{
    ProviderFunc: func() *schema.Provider {
      return mirantis.Provider()
    },
  })
}