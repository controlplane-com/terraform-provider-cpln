package cpln

import (
	"context"

	client "github.com/controlplane-com/terraform-provider-cpln/internal/provider/client"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceLocation() *schema.Resource {

	return &schema.Resource{
		ReadContext: dataSourceLocationRead,
		Schema:      client.LocationSchema(),
	}
}

func dataSourceLocationRead(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	c := m.(*client.Client)

	locationName := d.Get("name").(string)

	location, _, err := c.GetLocation(locationName)

	if err != nil {
		return diag.FromErr(err)
	}

	return setLocation(d, location)
}

func setLocation(d *schema.ResourceData, location *client.Location) diag.Diagnostics {

	if location == nil {
		d.SetId("")
		return nil
	}

	d.SetId(*location.Name)

	if err := SetBase(d, location.Base); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("cloud_provider", location.Provider); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("region", location.Region); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("enabled", location.Spec.Enabled); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("geo", flattenLocationGeo(location.Status.Geo)); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("ip_ranges", flattenIpRanges(location.Status.IpRanges)); err != nil {
		return diag.FromErr(err)
	}

	if err := SetSelfLink(location.Links, d); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func flattenLocationGeo(geo *client.LocationGeo) []interface{} {
	if geo == nil {
		return nil
	}

	spec := make(map[string]interface{})

	if geo.Lat != nil {
		spec["lat"] = *geo.Lat
	}

	if geo.Lon != nil {
		spec["lon"] = *geo.Lon
	}

	if geo.Country != nil {
		spec["country"] = *geo.Country
	}

	if geo.State != nil {
		spec["state"] = *geo.State
	}

	if geo.City != nil {
		spec["city"] = *geo.City
	}

	if geo.Continent != nil {
		spec["continent"] = *geo.Continent
	}

	return []interface{}{
		spec,
	}
}

func flattenIpRanges(ipRanges *[]string) []interface{} {

	if len(*ipRanges) > 0 {

		l := make([]interface{}, len(*ipRanges))

		for i, ip := range *ipRanges {
			l[i] = ip
		}

		return l
	}

	return make([]interface{}, 0)
}
