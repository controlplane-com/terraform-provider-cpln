package cpln

import (
	"context"

	client "terraform-provider-cpln/internal/provider/client"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceGvc() *schema.Resource {

	return &schema.Resource{
		ReadContext: dataSourceGvcRead,
		Schema: map[string]*schema.Schema{
			"cpln_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: NameValidator,
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"domain": {
				Type:     schema.TypeString,
				Computed: true,
			},
			// "alias": {
			// 	Type:     schema.TypeString,
			// 	Computed: true,
			// },
			"pull_secrets": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"locations": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"tags": {
				Type:     schema.TypeMap,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"self_link": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"lightstep_tracing": client.LightstepSchema(),
		},
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
