package cpln

import (
	"context"

	client "terraform-provider-cpln/internal/provider/client"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceCloudAccount() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceCloudAccountRead,
		Schema:      CloudAccountSchema(),
	}
}

func dataSourceCloudAccountRead(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*client.Client)

	org, _, err := c.GetOrg()
	if err != nil {
		return diag.FromErr(err)
	}

	cloudAccount, _, err := c.GetCloudAccount(d.Get("name").(string))

	if err != nil {
		return diag.FromErr(err)
	}

	return setCloudAccount(d, cloudAccount, *org.Name)
}
