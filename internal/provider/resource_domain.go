package cpln

import (
	"context"
	"time"

	client "terraform-provider-cpln/internal/provider/client"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceDomain() *schema.Resource {

	return &schema.Resource{
		CreateContext: resourceDomainCreate,
		ReadContext:   resourceDomainRead,
		UpdateContext: resourceDomainUpdate,
		DeleteContext: resourceDomainDelete,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				// ValidateFunc: NameValidator,
			},
			"description": {
				Type:             schema.TypeString,
				Optional:         true,
				ValidateFunc:     DescriptionDomainValidator,
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
		},
		Importer: &schema.ResourceImporter{},
	}
}

func resourceDomainCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	domain := client.Domain{}
	domain.Name = GetString(d.Get("name"))
	domain.Description = GetString(d.Get("description"))
	domain.Tags = GetStringMap(d.Get("tags"))

	c := m.(*client.Client)
	count := 0

	for {

		newDomain, code, err := c.CreateDomain(domain)

		if code == 409 {
			return ResourceExistsHelper()
		}

		if code != 400 {

			if err != nil {
				return diag.FromErr(err)
			}

			return setDomain(d, newDomain)
		}

		if count++; count > 16 {
			// Exit loop after timeout

			var diags diag.Diagnostics

			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  err.Error(),
			})

			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "Unable to verify domain ownership",
				Detail:   "Please review and run terraform apply again",
			})

			return diags
		}

		time.Sleep(15 * time.Second)
	}
}

func resourceDomainRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	c := m.(*client.Client)
	domain, code, err := c.GetDomain(d.Id())

	if code == 404 {
		d.SetId("")
		return nil
	}

	if err != nil {
		return diag.FromErr(err)
	}

	return setDomain(d, domain)
}

func resourceDomainUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	if d.HasChanges("description", "tags") {

		domainToUpdate := client.Domain{}
		domainToUpdate.Name = GetString(d.Get("name"))

		if d.HasChange("description") {
			domainToUpdate.Description = GetDescriptionString(d.Get("description"), *domainToUpdate.Name)
		}

		if d.HasChange("tags") {
			domainToUpdate.Tags = GetTagChanges(d)
		}

		c := m.(*client.Client)
		updatedDomain, _, err := c.UpdateDomain(domainToUpdate)
		if err != nil {
			return diag.FromErr(err)
		}

		return setDomain(d, updatedDomain)
	}

	return nil
}

func setDomain(d *schema.ResourceData, domain *client.Domain) diag.Diagnostics {

	if domain == nil {
		d.SetId("")
		return nil
	}

	d.SetId(*domain.Name)

	if err := d.Set("name", domain.Name); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("description", domain.Description); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("tags", domain.Tags); err != nil {
		return diag.FromErr(err)
	}

	if err := SetSelfLink(domain.Links, d); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceDomainDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	c := m.(*client.Client)
	err := c.DeleteDomain(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")

	return nil
}
