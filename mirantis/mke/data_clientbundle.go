package msr

import (
	"context"

	"github.com/Mirantis/terraform-provider-mirantis/mirantis/mke/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// dataSourceClientBundle for retrieving MSR accounts in bulk
func dataSourceClientBundle() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceClientBundleRead,
		Schema: map[string]*schema.Schema{

			"kube": {
				Type:        schema.TypeList,
				Description: "Kubernetes components from the client bundle.",
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"host": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"client_key": {
							Type:      schema.TypeString,
							Computed:  true,
							Sensitive: true,
						},
						"client_cert": {
							Type:      schema.TypeString,
							Computed:  true,
							Sensitive: true,
						},
						"ca_certificate": {
							Type:      schema.TypeString,
							Computed:  true,
							Sensitive: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceClientBundleRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	c, ok := m.(client.Client)
	if !ok {
		diags = append(diags, diag.Errorf("unable to cast meta interface to MKE Client")...)
		return diags
	}

	clientBundle, err := c.ApiClientBundle(ctx)
	if err != nil {
		diags = append(diags, diag.FromErr(err)...)
		return diags
	}

	kubeConfig := clientBundle.Kube
	if kubeConfig == nil {
		diags = append(diags, diag.Errorf("MKE Client produced no kube configuration. Is it a kube cluster?")...)
		return diags
	} else {

		m := make(map[string]interface{})

		m["host"] = kubeConfig.Host
		m["client_key"] = kubeConfig.ClientKey
		m["client_cert"] = kubeConfig.ClientCertificate
		m["ca_certificate"] = kubeConfig.CACertificate

		if err := d.Set("kube", []interface{}{m}); err != nil {
			diags = append(diags, diag.FromErr(err)...)
		}

	}

	if !diags.HasError() {
		d.SetId("cbid")
	}

	return diags
}
