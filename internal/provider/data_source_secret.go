package cpln

import (
	"context"

	client "github.com/controlplane-com/terraform-provider-cpln/internal/provider/client"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceSecret() *schema.Resource {

	return &schema.Resource{
		ReadContext: dataSourceSecretRead,
		Schema:      secretSchema(false),
	}
}

func dataSourceSecretRead(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	c := m.(*client.Client)

	name := d.Get("name").(string)

	secret, _, err := c.GetSecret(name)

	if err != nil {
		return diag.FromErr(err)
	}

	return setSecret(d, secret)
}
