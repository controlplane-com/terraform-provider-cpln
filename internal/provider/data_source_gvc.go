package cpln

import (
	"context"

	client "github.com/controlplane-com/terraform-provider-cpln/internal/provider/client"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceGvc() *schema.Resource {

	return &schema.Resource{
		ReadContext: dataSourceGvcRead,
		Schema:      GvcSchema(),
	}
}

func dataSourceGvcRead(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	c := m.(*client.Client)

	gvcName := d.Get("name").(string)

	org, _, err := c.GetOrg()
	if err != nil {
		return diag.FromErr(err)
	}

	gvc, _, err := c.GetGvc(gvcName)

	if err != nil {
		return diag.FromErr(err)
	}

	return setGvc(d, gvc, *org.Name)
}
