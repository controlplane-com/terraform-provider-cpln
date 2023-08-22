package cpln

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceCloudAccount() *schema.Resource {

	return &schema.Resource{
		ReadContext: dataSourceCloudAccountRead,
		Schema: map[string]*schema.Schema{
			"aws_identifiers": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func dataSourceCloudAccountRead(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	aws_identifiers := []string{"arn:aws:iam::957753459089:user/controlplane-driver", "arn:aws:iam::957753459089:role/controlplane-driver"}
	if err := d.Set("aws_identifiers", aws_identifiers); err != nil {
		return diag.FromErr(err)
	}

	d.SetId("static-cloud-account")

	return nil
}
