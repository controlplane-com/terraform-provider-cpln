package cpln

import (
	"context"

	client "github.com/controlplane-com/terraform-provider-cpln/internal/provider/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceIpSet() *schema.Resource {

	return &schema.Resource{
		CreateContext: resourceIpSetCreate,
		ReadContext:   resourceIpSetRead,
		UpdateContext: resourceIpSetUpdate,
		DeleteContext: resourceIpSetDelete,
		Schema: map[string]*schema.Schema{
			"cpln_id": {
				Type:        schema.TypeString,
				Description: "The ID, in GUID format, of the IpSet.",
				Computed:    true,
			},
			"name": {
				Type:         schema.TypeString,
				Description:  "Name of the IpSet.",
				Required:     true,
				ForceNew:     true,
				ValidateFunc: NameValidator,
			},
			"description": {
				Type:             schema.TypeString,
				Description:      "Description of the IpSet.",
				Optional:         true,
				ValidateFunc:     DescriptionValidator,
				DiffSuppressFunc: DiffSuppressDescription,
			},
			"tags": {
				Type:         schema.TypeMap,
				Description:  "Key-value map of resource tags.",
				Optional:     true,
				Elem:         StringSchema(),
				ValidateFunc: TagValidator,
			},
			"self_link": {
				Type:        schema.TypeString,
				Description: "Full link to this resource. Can be referenced by other resources.",
				Computed:    true,
			},
			"link": {
				Type:        schema.TypeString,
				Description: "The self link of a workload.",
				Optional:    true,
			},
			"location": {
				Type:        schema.TypeList,
				Description: "",
				Optional:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Description: "The self link of a location.",
							Required:    true,
						},
						"retention_policy": {
							Type:        schema.TypeString,
							Description: "",
							Required:    true,
						},
					},
				},
			},
			"status": {
				Type:        schema.TypeList,
				Description: "Status of the IpSet.",
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"ip_address": {
							Type:        schema.TypeList,
							Description: "",
							Computed:    true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": {
										Type:        schema.TypeString,
										Description: "",
										Computed:    true,
									},
									"ip": {
										Type:        schema.TypeString,
										Description: "",
										Computed:    true,
									},
									"id": {
										Type:        schema.TypeString,
										Description: "",
										Computed:    true,
									},
									"state": {
										Type:        schema.TypeString,
										Description: "",
										Computed:    true,
									},
									"created": {
										Type:        schema.TypeString,
										Description: "",
										Computed:    true,
									},
								},
							},
						},
						"error": {
							Type:        schema.TypeString,
							Description: "",
							Computed:    true,
						},
					},
				},
			},
		},
	}
}

func resourceIpSetCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	c := m.(*client.Client)

	// Define & Build
	ipSet := client.IpSet{
		Spec: &client.IpSetSpec{},
	}

	ipSet.Name = GetString(d.Get("name"))
	ipSet.Description = GetString(d.Get("description"))
	ipSet.Tags = GetStringMap(d.Get("tags"))

	if d.Get("link") != nil {
		ipSet.Spec.Link = GetString(d.Get("link"))
	}

	if d.Get("location") != nil {
		ipSet.Spec.Locations = buildIpSetLocations(d.Get("location").([]interface{}))
	}

	// Create
	newIpSet, code, err := c.CreateIpSet(ipSet)

	if code == 409 {
		return ResourceExistsHelper()
	}

	if err != nil {
		return diag.FromErr(err)
	}

	return setIpSet(d, newIpSet)
}

func resourceIpSetRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	c := m.(*client.Client)

	ipSet, code, err := c.GetIpSet(d.Id())

	if code == 404 {
		d.SetId("")
		return nil
	}

	if err != nil {
		return diag.FromErr(err)
	}

	return setIpSet(d, ipSet)
}

func resourceIpSetUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	if d.HasChanges("description", "tags", "link", "location") {
		c := m.(*client.Client)

		// Define & Build
		ipSetToUpdate := client.IpSet{
			SpecReplace: &client.IpSetSpec{},
		}

		ipSetToUpdate.Name = GetString(d.Get("name"))
		ipSetToUpdate.Description = GetDescriptionString(d.Get("description"), *ipSetToUpdate.Name)
		ipSetToUpdate.Tags = GetTagChanges(d)
		ipSetToUpdate.SpecReplace.Link = GetString(d.Get("link"))

		if d.Get("location") != nil {
			ipSetToUpdate.SpecReplace.Locations = buildIpSetLocations(d.Get("location").([]interface{}))
		}

		// Update
		updatedIpSet, _, err := c.UpdateIpSet(ipSetToUpdate)

		if err != nil {
			return diag.FromErr(err)
		}

		return setIpSet(d, updatedIpSet)
	}

	return nil
}

func resourceIpSetDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	c := m.(*client.Client)

	err := c.DeleteIpSet(d.Id())

	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")
	return nil
}

func setIpSet(d *schema.ResourceData, ipSet *client.IpSet) diag.Diagnostics {

	if ipSet == nil {
		d.SetId("")
		return nil
	}

	d.SetId(*ipSet.Name)

	if err := SetBase(d, ipSet.Base); err != nil {
		return diag.FromErr(err)
	}

	if ipSet.Spec != nil {

		if err := d.Set("link", ipSet.Spec.Link); err != nil {
			return diag.FromErr(err)
		}

		if err := d.Set("location", flattenIpSetLocations(ipSet.Spec.Locations)); err != nil {
			return diag.FromErr(err)
		}
	}

	if err := d.Set("status", flattenIpSetStatus(ipSet.Status)); err != nil {
		return diag.FromErr(err)
	}

	if err := SetSelfLink(ipSet.Links, d); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

/*** Build ***/

func buildIpSetLocations(specs []interface{}) *[]client.IpSetLocation {

	if len(specs) == 0 {
		return nil
	}

	output := []client.IpSetLocation{}

	for _, _spec := range specs {

		spec := _spec.(map[string]interface{})

		ipSetLocation := client.IpSetLocation{
			Name:            GetString(spec["name"]),
			RetentionPolicy: GetString(spec["retention_policy"]),
		}

		output = append(output, ipSetLocation)
	}

	return &output
}

/*** Flatten ***/

func flattenIpSetLocations(locations *[]client.IpSetLocation) []interface{} {

	if locations == nil {
		return nil
	}

	specs := []interface{}{}

	for _, location := range *locations {

		spec := map[string]interface{}{
			"name":             *location.Name,
			"retention_policy": *location.RetentionPolicy,
		}

		specs = append(specs, spec)
	}

	return specs
}

func flattenIpSetStatus(status *client.IpSetStatus) []interface{} {

	if status == nil {
		return nil
	}

	spec := map[string]interface{}{}

	if status.IpAddresses != nil {
		spec["ip_address"] = flattenIpAddresses(status.IpAddresses)
	}

	if status.Error != nil {
		spec["error"] = *status.Error
	}

	return []interface{}{
		spec,
	}
}

func flattenIpAddresses(ipAddresses *[]client.IpAddress) []interface{} {

	if ipAddresses == nil {
		return nil
	}

	specs := []interface{}{}

	for _, ipAddress := range *ipAddresses {

		spec := map[string]interface{}{}

		if ipAddress.Name != nil {
			spec["name"] = *ipAddress.Name
		}

		if ipAddress.Ip != nil {
			spec["ip"] = *ipAddress.Ip
		}

		if ipAddress.Id != nil {
			spec["id"] = *ipAddress.Id
		}

		if ipAddress.State != nil {
			spec["state"] = *ipAddress.State
		}

		if ipAddress.Created != nil {
			spec["created"] = *ipAddress.Created
		}

		specs = append(specs, spec)
	}

	return specs
}
