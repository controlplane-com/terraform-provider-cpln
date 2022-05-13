package cpln

import (
	"context"
	"strconv"
	client "terraform-provider-cpln/internal/provider/client"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceGvcs() *schema.Resource {

	return &schema.Resource{
		ReadContext: dataSourceGvcsRead,
		Schema: map[string]*schema.Schema{
			"kind": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"items": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{

						"kind": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"version": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"tags": {
							Type:     schema.TypeMap,
							Computed: true,
						},

						"created": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"lastmodified": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"links": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{

									"ref": {
										Type:     schema.TypeString,
										Computed: true,
									},

									"href": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},

						"alias": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},

			"links": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{

						"ref": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"href": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},

			"query": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"kind": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"fetch": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceGvcsRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	c := m.(*client.Client)

	var diags diag.Diagnostics

	gvcs, err := c.GetGvcs()

	if err != nil {
		return diag.FromErr(err)
	}

	gvcItems := flattenGvcData(&gvcs.Items)

	if err := d.Set("items", gvcItems); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(strconv.FormatInt(time.Now().Unix(), 10))

	return diags
}

func flattenGvcData(gvcItems *[]client.Gvc) []interface{} {

	if gvcItems != nil {

		gvcs := make([]interface{}, len(*gvcItems))

		for i, gvcItem := range *gvcItems {

			gvc := make(map[string]interface{})

			gvc["kind"] = gvcItem.Kind
			gvc["id"] = gvcItem.ID
			gvc["name"] = gvcItem.Name

			gvcs[i] = gvc
		}

		return gvcs
	}

	return make([]interface{}, 0)
}
