package msr

import (
	"context"
	"time"

	"github.com/Mirantis/terraform-provider-mirantis/mirantis/msr/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// ResourceUsers for managing MSR users
func ResourceUsers() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceUserCreate,
		ReadContext:   resourceUserRead,
		UpdateContext: resourceUserUpdate,
		DeleteContext: resourceUserDelete,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"full_name": {
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

func resourceUserCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c, ok := m.(client.Client)
	if !ok {
		return diag.Errorf("unable to cast meta interface to MSR Client")
	}

	user := client.User{
		Name:       d.Get("name").(string),
		Password:   client.GeneratePass(),
		FullName:   d.Get("full_name").(string),
		IsActive:   true,
		IsAdmin:    false,
		IsOrg:      false,
		SearchLDAP: false,
	}
	u, err := c.CreateUser(ctx, user)
	if err != nil {
		return diag.FromErr(err)
	}
	err = d.Set("last_updated", time.Now().Format(time.RFC850))
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(u.ID)

	return diag.Diagnostics{}
}

func resourceUserRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c, ok := m.(client.Client)
	if !ok {
		return diag.Errorf("unable to cast meta interface to MSR Client")
	}

	u, err := c.ReadUser(ctx, d.Get("name").(string))
	if err != nil {
		// If the user doesn't exist we should gracefully handle it
		d.SetId("")
		return diag.FromErr(err)
	}

	d.SetId(u.ID)

	return diag.Diagnostics{}
}

func resourceUserUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c, ok := m.(client.Client)

	if !ok {
		return diag.Errorf("unable to cast meta interface to MSR Client")
	}
	if d.HasChange("msr_users") {
		user := client.User{
			Name:       d.Get("name").(string),
			ID:         d.State().ID,
			FullName:   d.Get("full_name").(string),
			IsActive:   true,
			IsAdmin:    false,
			IsOrg:      false,
			SearchLDAP: false,
		}
		_, err := c.UpdateUser(ctx, user)

		if err != nil {
			return diag.FromErr(err)
		}
		err = d.Set("last_updated", time.Now().Format(time.RFC850))
		if err != nil {
			return diag.FromErr(err)
		}
	}
	return resourceUserRead(ctx, d, m)
}

func resourceUserDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c, ok := m.(client.Client)

	if !ok {
		return diag.Errorf("unable to cast meta interface to MSR Client")
	}
	err := c.DeleteUser(ctx, d.State().ID)
	if err != nil {
		return diag.FromErr(err)
	}

	// d.SetId("") is automatically called assuming delete returns no errors, but
	// it is added here for explicitness.
	d.SetId("")

	return diag.Diagnostics{}
}
