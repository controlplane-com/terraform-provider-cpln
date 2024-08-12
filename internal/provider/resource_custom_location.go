package cpln

import (
	"context"

	client "github.com/controlplane-com/terraform-provider-cpln/internal/provider/client"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceCustomLocation() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceCustomLocationCreate,
		ReadContext:   resourceCustomLocationRead,
		UpdateContext: resourceCustomLocationUpdate,
		DeleteContext: resourceCustomLocationDelete,
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
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Optional:     true,
				ValidateFunc: TagValidator,
			},
			"cloud_provider": {
				Type:        schema.TypeString,
				Description: "Cloud Provider of the location.",
				Required:    true,
			},
			"enabled": {
				Type:        schema.TypeBool,
				Description: "Indication if location is enabled.",
				Required:    true,
			},
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

func resourceCustomLocationCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	c := m.(*client.Client)

	// Define & Build
	locationName := d.Get("name").(string)
	isLocationEnabled := d.Get("enabled").(bool)
	description := GetString(d.Get("description").(string))
	tags := GetStringMap(d.Get("tags"))
	provider := d.Get("cloud_provider").(string)

	if provider != "byok" {
		return diag.Errorf("provider must be byok for custom locations")
	}

	location := &client.Location{}
	location.Name = &locationName
	location.Spec = &client.LocationSpec{}
	location.Spec.Enabled = GetBool(isLocationEnabled)
	location.Provider = &provider
	location.Tags = tags
	location.Description = description
	location.Status = nil
	createdLocation, _, err := c.CreateCustomLocation(*location)

	if err != nil {
		return diag.FromErr(err)
	}

	return setCustomLocationResource(d, createdLocation)
}

func resourceCustomLocationRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	c := m.(*client.Client)

	location, code, err := c.GetLocation(d.Id())

	if code == 404 {
		d.SetId("")
		return nil
	}

	if err != nil {
		return diag.FromErr(err)
	}

	return setCustomLocationResource(d, location)
}

func resourceCustomLocationUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	if d.HasChanges("tags", "enabled", "description") {
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

		if d.HasChange("description") {
			locationToUpdate.Description = GetString(d.Get("description").(string))
		}

		// Update
		updatedLocation, _, err := c.UpdateLocation(locationToUpdate)

		if err != nil {
			return diag.FromErr(err)
		}

		return setCustomLocationResource(d, updatedLocation)
	}

	return nil
}

func resourceCustomLocationDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	c := m.(*client.Client)

	locationName := d.Get("name").(string)

	// Get location
	location, _, err := c.GetLocation(locationName)

	if err != nil {
		return diag.FromErr(err)
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

func setCustomLocationResource(d *schema.ResourceData, location *client.Location) diag.Diagnostics {
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

	if err := d.Set("enabled", location.Spec.Enabled); err != nil {
		return diag.FromErr(err)
	}

	if err := SetSelfLink(location.Links, d); err != nil {
		return diag.FromErr(err)
	}

	return nil
}
