package msr

import (
	"context"
	"time"

	"github.com/Mirantis/terraform-provider-mirantis/mirantis/msr/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// ResourceTeam for managing MSR team
func ResourceTeam() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceTeamCreate,
		ReadContext:   resourceTeamRead,
		UpdateContext: resourceTeamUpdate,
		DeleteContext: resourceTeamDelete,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"org_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"last_updated": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func resourceTeamCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c, ok := m.(client.Client)
	if !ok {
		return diag.Errorf("unable to cast meta interface to MSR Client")
	}

	team := client.Team{
		Name:        d.Get("name").(string),
		Description: d.Get("description").(string),
	}
	t, err := c.CreateTeam(ctx, d.Get("org_id").(string), team)
	if err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("last_updated", time.Now().Format(time.RFC850)); err != nil {
		return diag.FromErr(err)
	}
	d.SetId(t.ID)

	return diag.Diagnostics{}
}

func resourceTeamRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c, ok := m.(client.Client)
	if !ok {
		return diag.Errorf("unable to cast meta interface to MSR Client")
	}

	t, err := c.ReadTeam(ctx, d.Get("org_id").(string), d.State().ID)
	if err != nil {
		// If the user doesn't exist we should gracefully handle it
		d.SetId("")
		return diag.FromErr(err)
	}

	d.SetId(t.ID)

	return diag.Diagnostics{}
}

func resourceTeamUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c, ok := m.(client.Client)

	if !ok {
		return diag.Errorf("unable to cast meta interface to MSR Client")
	}
	if d.HasChange("msr_team") {
		team := client.Team{
			ID:          d.State().ID,
			Description: d.Get("description").(string),
		}
		if _, err := c.UpdateTeam(ctx, d.Get("org_id").(string), team); err != nil {
			return diag.FromErr(err)
		}
		if err := d.Set("last_updated", time.Now().Format(time.RFC850)); err != nil {
			return diag.FromErr(err)
		}
	}
	return resourceTeamRead(ctx, d, m)
}

func resourceTeamDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c, ok := m.(client.Client)

	if !ok {
		return diag.Errorf("unable to cast meta interface to MSR Client")
	}
	if err := c.DeleteTeam(ctx, d.Get("org_id").(string), d.State().ID); err != nil {
		return diag.FromErr(err)
	}

	// d.SetId("") is automatically called assuming delete returns no errors, but
	// it is added here for explicitness.
	d.SetId("")

	return diag.Diagnostics{}
}
