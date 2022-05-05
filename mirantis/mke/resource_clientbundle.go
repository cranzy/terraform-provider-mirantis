package msr

import (
	"context"
	"time"

	"github.com/Mirantis/terraform-provider-mirantis/mirantis/mke/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// ResourceClientBundle for managing MKE Client Bundles
func ResourceClientBundle() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceClientBundleCreate,
		ReadContext:   resourceClientBundleRead,
		UpdateContext: resourceClientBundleUpdate,
		DeleteContext: resourceClientBundleDelete,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"last_updated": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"private_key": {
				Type:      schema.TypeString,
				Computed:  true,
				Sensitive: true,
			},
			"public_key": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"ca_cert": {
				Type:      schema.TypeString,
				Computed:  true,
				Sensitive: true,
			},
			"client_certs": {
				Type:     schema.TypeList,
			  Elem: &schema.Schema{
			    Type: schema.TypeString,
			  },
				Computed: true,
			},

			"kube": {
				Type:        schema.TypeList,
				Description: "Kubernetes components from the client bundle.",
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"config_yml": {
							Type:      schema.TypeString,
							Computed:  true,
							Sensitive: true,
						},
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
						"ca_cert": {
							Type:      schema.TypeString,
							Computed:  true,
							Sensitive: true,
						},
					},
				},
			},
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func resourceClientBundleCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	c, ok := m.(client.Client)
	if !ok {
		diags = append(diags, diag.Errorf("unable to cast meta interface to MKE Client")...)
		return diags
	}

	cb, err := c.ApiClientBundleCreate(ctx)
	if err != nil {
		diags = append(diags, diag.FromErr(err)...)
		return diags
	}

	if err := d.Set("private_key", cb.PrivateKey); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}
	if err := d.Set("public_key", cb.PublicKey); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}
	if err := d.Set("ca_cert", cb.CACert); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}
	if err := d.Set("client_certs", cb.Certs); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}

	kc := cb.Kube
	if kc == nil {
		diags = append(diags, diag.Errorf("MKE Client produced no kube configuration. Is it a kube cluster?")...)
		return diags
	} else {

		m := make(map[string]interface{})

		m["config_yml"] = kc.Config
		m["host"] = kc.Host
		m["client_key"] = kc.ClientKey
		m["client_cert"] = kc.ClientCertificate
		m["ca_cert"] = kc.CACertificate

		if err := d.Set("kube", []interface{}{m}); err != nil {
			diags = append(diags, diag.FromErr(err)...)
		}

	}

	if !diags.HasError() {
		if err := d.Set("last_updated", time.Now().Format(time.RFC850)); err != nil {
			return diag.FromErr(err)
		}

		d.SetId(d.Get("name").(string))
	}

	return diags
}

func resourceClientBundleRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	return diag.Diagnostics{}
}

func resourceClientBundleUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	return diag.Diagnostics{}
}

func resourceClientBundleDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	return diag.Diagnostics{}
}
