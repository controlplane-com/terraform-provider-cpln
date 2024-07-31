package cpln

import (
	"context"

	client "github.com/controlplane-com/terraform-provider-cpln/internal/provider/client"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceLocation() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceLocationCreate,
		ReadContext:   resourceLocationRead,
		UpdateContext: resourceLocationUpdate,
		DeleteContext: resourceLocationDelete,
		Schema: map[string]*schema.Schema{
			"cpln_id": {
				Type:        schema.TypeString,
				Description: "The ID, in GUID format, of the location.",
				Computed:    true,
			},
			"name": {
				Type:         schema.TypeString,
				Description:  "Name of the Location.",
				ForceNew:     true,
				Required:     true,
				ValidateFunc: NameValidator,
			},
			"description": {
				Type:             schema.TypeString,
				Description:      "Description of the location.",
				Optional:         true,
				ValidateFunc:     DescriptionValidator,
				DiffSuppressFunc: DiffSuppressDescription,
			},
			"tags": {
				Type:        schema.TypeMap,
				Description: "Key-value map of resource tags.",
				Optional:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				ValidateFunc: TagValidator,
			},
			"cloud_provider": {
				Type:        schema.TypeString,
				Description: "Cloud Provider of the location.",
				Optional:    true,
			},
			"region": {
				Type:        schema.TypeString,
				Description: "Region of the location.",
				Computed:    true,
			},
			"enabled": {
				Type:        schema.TypeBool,
				Description: "Indication if location is enabled.",
				Required:    true,
			},
			"geo": client.GeoSchema(),
			"ip_ranges": {
				Type:        schema.TypeSet,
				Description: "A list of IP ranges of the location.",
				Computed:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"self_link": {
				Type:        schema.TypeString,
				Description: "Full link to this resource. Can be referenced by other resources.",
				Computed:    true,
			},
		},
		Importer: &schema.ResourceImporter{},
	}
}

func resourceLocationCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	c := m.(*client.Client)

	// Define & Build
	name := d.Get("name").(string)
	description := GetString(d.Get("description"))
	provider := GetString(d.Get(("cloud_provider")))
	isEnabled := d.Get("enabled").(bool)

	if provider != nil {
		// Define the BYOK location object
		byokLocation := &client.Location{}
		byokLocation.Name = &name
		byokLocation.Description = description
		byokLocation.Provider = provider
		byokLocation.Spec = &client.LocationSpec{
			Enabled: &isEnabled,
		}

		// Create
		newLocation, code, err := c.CreateLocation(*byokLocation)

		if code == 409 {
			return ResourceExistsHelper()
		}

		if err != nil {
			return diag.FromErr(err)
		}

		// Update state
		return setLocationResource(d, newLocation)
	}

	// Get location
	location, _, err := c.GetLocation(name)

	if err != nil {
		return diag.FromErr(err)
	}

	// Handle nil just to be on the safe side
	if location.Spec == nil {
		location.Spec = &client.LocationSpec{}
	}

	// Set user's enabled value
	location.Spec.Enabled = GetBool(isEnabled)

	// Set tags
	location.Tags = GetStringMap(d.Get("tags"))

	// Remove status
	location.Status = nil

	// Send an update request
	updatedLocation, _, err := c.UpdateLocation(*location)

	if err != nil {
		return diag.FromErr(err)
	}

	return setLocationResource(d, updatedLocation)
}

func resourceLocationRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	c := m.(*client.Client)

	location, code, err := c.GetLocation(d.Id())

	if code == 404 {
		d.SetId("")
		return nil
	}

	if err != nil {
		return diag.FromErr(err)
	}

	return setLocationResource(d, location)
}

func resourceLocationUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	if d.HasChanges("tags", "enabled", "cloud_provider") {
		c := m.(*client.Client)

		// Define & Build
		locationToUpdate := client.Location{
			Spec: &client.LocationSpec{
				Enabled: GetBool(d.Get("enabled").(bool)),
			},
		}

		locationToUpdate.Name = GetString(d.Get("name").(string))

		if d.HasChange("tags") {
			locationToUpdate.Tags = GetTagChanges(d)
		}

		if d.HasChange("cloud_provider") {
			locationToUpdate.Provider = GetString(d.Get("cloud_provider"))
		}

		// Update
		updatedLocation, _, err := c.UpdateLocation(locationToUpdate)

		if err != nil {
			return diag.FromErr(err)
		}

		return setLocationResource(d, updatedLocation)
	}

	return nil
}

func resourceLocationDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	c := m.(*client.Client)

	name := d.Get("name").(string)

	// Get location
	location, _, err := c.GetLocation(name)

	if err != nil {
		return diag.FromErr(err)
	}

	if *location.Provider == "byok" {
		// Delete BYOK location
		err := c.DeleteLocation(d.Id())

		if err != nil {
			return diag.FromErr(err)
		}

		d.SetId("")
		return nil
	}

	// Handle nil just to be on the safe side
	if location.Spec == nil {
		location.Spec = &client.LocationSpec{}
	}

	// Set enabled to its default value
	location.Spec.Enabled = GetBool(true)

	// Remove status
	location.Status = nil

	// Send an update request
	_, _, err = c.UpdateLocation(*location)

	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")
	return nil
}

func setLocationResource(d *schema.ResourceData, location *client.Location) diag.Diagnostics {
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

	if location.Status != nil {
		if location.Status.Geo != nil {
			if err := d.Set("geo", flattenLocationGeo(location.Status.Geo)); err != nil {
				return diag.FromErr(err)
			}
		}

		if location.Status.IpRanges != nil {
			if err := d.Set("ip_ranges", flattenIpRanges(location.Status.IpRanges)); err != nil {
				return diag.FromErr(err)
			}
		}
	}

	if err := SetSelfLink(location.Links, d); err != nil {
		return diag.FromErr(err)
	}

	return nil
}
