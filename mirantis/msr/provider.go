package msr

import (
	"context"

	"github.com/Mirantis/terraform-provider-mirantis/mirantis/msr/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"host": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("MSR_HOST_URL", nil),
			},
			"username": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("MSR_ADMIN_USER", nil),
			},
			"password": {
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				DefaultFunc: schema.EnvDefaultFunc("MSR_ADMIN_PASS", nil),
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"msr_user": ResourceUser(),
			"msr_org":  ResourceOrg(),
			"msr_team": ResourceTeam(),
		},
		ConfigureContextFunc: providerConfigure,
	}
}

func providerConfigure(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	username := d.Get("username").(string)
	password := d.Get("password").(string)

	var host string

	hVal, ok := d.GetOk("host")
	if ok {
		tempHost := hVal.(string)
		host = tempHost
	}

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	if (username != "") && (password != "") {
		c, err := client.NewClient(host, username, password)
		if err != nil {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "Unable to create MSR client",
				Detail:   "Unable to authenticate user for authenticated MSR client",
			})

			return nil, diags
		}

		healthy, err := c.IsHealthy(ctx)
		if !healthy {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "MSR endpoint is not healthy",
				Detail:   err.Error(),
			})
			return nil, diags
		}
		return c, diags
	}

	diags = append(diags, diag.Diagnostic{
		Severity: diag.Error,
		Summary:  "Unable to create MSR client",
		Detail:   "Unable to create anonymous MSR client",
	})

	return nil, diags
}
