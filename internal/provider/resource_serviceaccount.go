package cpln

import (
	"context"

	client "github.com/controlplane-com/terraform-provider-cpln/internal/provider/client"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceServiceAccount() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceServiceAccountCreate,
		ReadContext:   resourceServiceAccountRead,
		UpdateContext: resourceServiceAccountUpdate,
		DeleteContext: resourceServiceAccountDelete,
		Schema: map[string]*schema.Schema{
			"cpln_id": {
				Type:        schema.TypeString,
				Description: "The ID, in GUID format, of the Service Account.",
				Computed:    true,
			},
			"name": {
				Type:         schema.TypeString,
				Description:  "Name of the Service Account.",
				ForceNew:     true,
				Required:     true,
				ValidateFunc: NameValidator,
			},
			"description": {
				Type:             schema.TypeString,
				Description:      "Description of the Service Account.",
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
			"origin": {
				Type:        schema.TypeString,
				Description: "Origin of the Policy. Either `builtin` or `default`.",
				Computed:    true,
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

func resourceServiceAccountCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	sa := client.ServiceAccount{}
	sa.Name = GetString(d.Get("name"))
	sa.Description = GetString(d.Get("description"))
	sa.Tags = GetStringMap(d.Get("tags"))
	sa.Origin = GetString(d.Get("origin"))

	c := m.(*client.Client)
	newSa, code, err := c.CreateServiceAccount(sa)

	if code == 409 {
		return ResourceExistsHelper()
	}

	if err != nil {
		return diag.FromErr(err)
	}

	return setServiceAccount(d, newSa)
}

func resourceServiceAccountRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	c := m.(*client.Client)
	sa, code, err := c.GetServiceAccount(d.Id())

	if code == 404 {
		d.SetId("")
		return nil
	}

	if err != nil {
		return diag.FromErr(err)
	}

	return setServiceAccount(d, sa)
}

func setServiceAccount(d *schema.ResourceData, sa *client.ServiceAccount) diag.Diagnostics {

	if sa == nil {
		d.SetId("")
		return nil
	}

	d.SetId(*sa.Name)

	if err := SetBase(d, sa.Base); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("origin", sa.Origin); err != nil {
		return diag.FromErr(err)
	}

	if err := SetSelfLink(sa.Links, d); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceServiceAccountUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	if d.HasChanges("description", "tags") {

		saToUpdate := client.ServiceAccount{}
		saToUpdate.Name = GetString(d.Get("name"))
		saToUpdate.Description = GetDescriptionString(d.Get("description"), *saToUpdate.Name)
		saToUpdate.Tags = GetTagChanges(d)

		c := m.(*client.Client)
		updatedSa, _, err := c.UpdateServiceAccount(saToUpdate)
		if err != nil {
			return diag.FromErr(err)
		}

		return setServiceAccount(d, updatedSa)
	}

	return nil
}

func resourceServiceAccountDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	c := m.(*client.Client)
	err := c.DeleteServiceAccount(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")

	return nil
}
