package msr

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/Mirantis/terraform-provider-mirantis/mirantis/mke/client"
)

func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"endpoint": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("MKE_ENDPOINT", nil),
			},
			"username": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("MKE_USER", nil),
			},
			"password": {
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				DefaultFunc: schema.EnvDefaultFunc("MKE_PASS", nil),
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"mke_clientbundle": ResourceClientBundle(),
		},
		ConfigureContextFunc: providerConfigure,
	}
}

func providerConfigure(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	username := d.Get("username").(string)
	password := d.Get("password").(string)
	endpoint := d.Get("endpoint").(string)

	if (username == "") || (password == "") || (endpoint == "") {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Unable to create MKE client",
			Detail:   "Unable to create anonymous MKE client",
		})

		return nil, diags
	}

	c, err := client.NewUnsafeSSLClient(endpoint, username, password)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Unable to create MKE client",
			Detail:   "Unable to authenticate user for authenticated MKE client",
		})

		return nil, diags
	}

	if err := c.ApiPing(ctx); err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "MKE endpoint is not healthy",
			Detail:   err.Error(),
		})
		return nil, diags
	}

	return c, diags
}
