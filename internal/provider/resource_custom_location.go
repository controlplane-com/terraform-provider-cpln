package cpln

import (
	"context"
	"fmt"

	client "github.com/controlplane-com/terraform-provider-cpln/internal/provider/client"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var allowedProviders = []string{"byok"}

func resourceCustomLocation() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceCustomLocationCreate,
		ReadContext:   resourceCustomLocationRead,
		UpdateContext: resourceCustomLocationUpdate,
		DeleteContext: resourceCustomLocationDelete,
		Schema: map[string]*schema.Schema{
			"cpln_id": {
				Type:        schema.TypeString,
				Description: "The ID, in GUID format, of the custom location.",
				Computed:    true,
			},
			"name": {
				Type:         schema.TypeString,
				Description:  "Name of the custom location.",
				ForceNew:     true,
				Required:     true,
				ValidateFunc: NameValidator,
			},
			"description": {
				Type:             schema.TypeString,
				Description:      "Description of the custom location.",
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
				Description: "Cloud Provider of the custom location.",
				Required:    true,
				ForceNew:    true,
				ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
					v := val.(string)
					for _, provider := range allowedProviders {
						if v == provider {
							return
						}
					}
					errs = append(errs, fmt.Errorf("%q must be one of the allowed providers: %v", key, allowedProviders))
					return
				},
			},
			"enabled": {
				Type:        schema.TypeBool,
				Description: "Indication if the custom location is enabled.",
				Required:    true,
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

	location := &client.Location{}
	location.Name = &locationName
	location.Spec = &client.LocationSpec{}
	location.Spec.Enabled = GetBool(isLocationEnabled)
	location.Provider = &provider
	location.Tags = tags
	location.Description = description
	createdLocation, code, err := c.CreateCustomLocation(*location)

	if code == 409 {
		return ResourceExistsHelper()
	}

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
			locationToUpdate.Description = GetDescriptionString(d.Get("description"), *locationToUpdate.Name)
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
	err := c.DeleteCustomLocation(locationName)

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
