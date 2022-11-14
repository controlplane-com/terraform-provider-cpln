package cpln

import (
	"context"
	"strconv"
	client "terraform-provider-cpln/internal/provider/client"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceLocations() *schema.Resource {

	return &schema.Resource{
		ReadContext: dataSourceLocationsRead,
		Schema: map[string]*schema.Schema{
			"locations": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: client.LocationsSchema(),
				},
			},
		},
	}
}

func dataSourceLocationsRead(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	c := m.(*client.Client)

	locations, err := c.GetLocations()

	if err != nil {
		return diag.FromErr(err)
	}

	locationItems := flattenLocationData(&locations.Items)

	if err := d.Set("locations", locationItems); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(strconv.FormatInt(time.Now().Unix(), 10))

	return nil
}

func flattenLocationData(locationItems *[]client.Location) []interface{} {

	if locationItems != nil {

		locations := make([]interface{}, len(*locationItems))

		for i, locationItem := range *locationItems {

			location := make(map[string]interface{})

			location["cpln_id"] = locationItem.ID
			location["name"] = locationItem.ID
			location["description"] = locationItem.Name
			location["tags"] = GetTags(locationItem.Tags)
			location["cloud_provider"] = locationItem.Provider
			location["region"] = locationItem.Region
			location["enabled"] = locationItem.Spec.Enabled
			location["ip_ranges"] = flattenIpRanges(locationItem.Status.IpRanges)
			location["self_link"] = GetSelfLink(locationItem.Links)

			locations[i] = location
		}

		return locations
	}

	return make([]interface{}, 0)
}
