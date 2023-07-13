package cpln

import (
	"context"
	"fmt"
	"strings"
	client "terraform-provider-cpln/internal/provider/client"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

/*** Main ***/
func resourceSpicedb() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceSpicedbCreate,
		ReadContext:   resourceSpicedbRead,
		UpdateContext: resourceSpicedbUpdate,
		DeleteContext: resourceSpicedbDelete,
		Schema: map[string]*schema.Schema{
			"cpln_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: NameValidator,
			},
			"description": {
				Type:             schema.TypeString,
				Optional:         true,
				ValidateFunc:     DescriptionValidator,
				DiffSuppressFunc: DiffSuppressDescription,
			},
			"tags": {
				Type:     schema.TypeMap,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				ValidateFunc: TagValidator,
			},
			"self_link": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"alias": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"external_endpoint": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"version": {
				Type:     schema.TypeString,
				Required: true,
			},
			"locations": {
				Type:     schema.TypeSet,
				Required: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
		Importer: &schema.ResourceImporter{},
	}
}

func setSpicedb(d *schema.ResourceData, spicedb *client.Spicedb, org string) diag.Diagnostics {

	if spicedb == nil {
		d.SetId("")
		return nil
	}

	d.SetId(*spicedb.Name)

	if err := SetBase(d, spicedb.Base); err != nil {
		return diag.FromErr(err)
	}

	if err := SetSelfLink(spicedb.Links, d); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("status", flattenSpicedbStatus(spicedb.Status)); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("version", spicedb.Spec.Version); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("locations", flattenSpicedbLocations(spicedb.Spec.Locations, org)); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceSpicedbCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	c := m.(*client.Client)
	spicedb := client.Spicedb{}

	spicedb.Name = GetString(d.Get("name"))
	spicedb.Description = GetDescriptionString(d.Get("description"), *spicedb.Name)
	spicedb.Tags = GetStringMap(d.Get("tags"))

	spicedb.Spec = &client.ClusterSpec{
		Version:   GetString(d.Get("version")),
		Locations: buildSpicedbLocations(d.Get("locations"), c.Org),
	}

	newSpicedb, code, err := c.CreateSpicedb(spicedb)

	if code == 409 {
		return ResourceExistsHelper()
	}

	if err != nil {
		return diag.FromErr(err)
	}

	return setSpicedb(d, newSpicedb, c.Org)
}

func resourceSpicedbRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*client.Client)
	spicedb, code, err := c.GetSpicedb(d.Id())

	if code == 404 {
		d.SetId("")
		return nil
	}

	if err != nil {
		return diag.FromErr(err)
	}

	return setSpicedb(d, spicedb, c.Org)
}

func resourceSpicedbUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	if d.HasChanges("description", "tags", "version", "locations") {

		c := m.(*client.Client)
		spicedbToUpdate := client.Spicedb{}
		spicedbToUpdate.Name = GetString(d.Get("name"))
		spicedbToUpdate.Spec = &client.ClusterSpec{
			Version:   GetString(d.Get("version")),
			Locations: buildSpicedbLocations(d.Get("locations"), c.Org),
		}

		if d.HasChange("description") {
			spicedbToUpdate.Description = GetDescriptionString(d.Get("description"), *spicedbToUpdate.Name)
		}

		if d.HasChange("tags") {
			spicedbToUpdate.Tags = GetTagChanges(d)
		}

		// Perform update
		updatedSpicedb, _, err := c.UpdateSpicedb(spicedbToUpdate)

		if err != nil {
			return diag.FromErr(err)
		}

		return setSpicedb(d, updatedSpicedb, c.Org)
	}

	return nil
}

func resourceSpicedbDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*client.Client)
	err := c.DeleteSpicedb(d.Id())

	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")

	return nil
}

/*** Build Functions ***/
func buildSpicedbLocations(spec interface{}, org string) *[]string {

	locations := []string{}

	for _, location := range spec.(*schema.Set).List() {
		locations = append(locations, fmt.Sprintf("/org/%s/location/%s", org, location))
	}

	return &locations
}

/*** Flatten Functions ***/
func flattenSpicedbStatus(status *client.ClusterStatus) []interface{} {

	spec := map[string]interface{}{}

	if status.ExternalEndpoint != nil {
		spec["external_endpoint"] = *status.ExternalEndpoint
	}

	return []interface{}{
		spec,
	}
}

func flattenSpicedbLocations(locations *[]string, org string) interface{} {

	spec := make([]interface{}, len(*locations))

	for i, location := range *locations {
		location = strings.TrimPrefix(location, fmt.Sprintf("/org/%s/location/", org))
		spec[i] = location
	}

	return spec
}
