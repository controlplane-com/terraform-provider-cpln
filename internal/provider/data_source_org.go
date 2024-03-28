package cpln

import (
	"context"

	client "github.com/controlplane-com/terraform-provider-cpln/internal/provider/client"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceOrg() *schema.Resource {

	return &schema.Resource{
		ReadContext: dataSourceOrgRead,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Description: "The name of the organization.",
				Computed:    true,
			},
			"cpln_id": {
				Type:        schema.TypeString,
				Description: "The ID, in GUID format, of the organization.",
				Computed:    true,
			},
		},
	}
}

func dataSourceOrgRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	c := m.(*client.Client)

	var diags diag.Diagnostics

	org, _, err := c.GetOrg()
	if err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("name", org.Name); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("cpln_id", org.ID); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(*org.ID)

	return diags
}
