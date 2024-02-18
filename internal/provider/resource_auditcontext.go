package cpln

import (
	"context"

	client "github.com/controlplane-com/terraform-provider-cpln/internal/provider/client"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAuditContext() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAuditContextCreate,
		ReadContext:   resourceAuditContextRead,
		UpdateContext: resourceAuditContextUpdate,
		DeleteContext: resourceAuditContextDelete,
		Schema: map[string]*schema.Schema{
			"cpln_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:         schema.TypeString,
				ForceNew:     true,
				Required:     true,
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
		},
		Importer: &schema.ResourceImporter{},
	}
}

func resourceAuditContextCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	auditCtx := client.AuditContext{}
	auditCtx.Name = GetString(d.Get("name"))
	auditCtx.Description = GetString(d.Get("description"))
	auditCtx.Tags = GetStringMap(d.Get("tags"))

	c := m.(*client.Client)
	newAuditCtx, code, err := c.CreateAuditContext(auditCtx)

	if code == 409 {
		return ResourceExistsHelper()
	}

	if err != nil {
		return diag.FromErr(err)
	}

	return setAuditContext(d, c.Org, newAuditCtx)
}

func resourceAuditContextRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	c := m.(*client.Client)
	auditCtx, code, err := c.GetAuditContext(d.Id())

	if code == 404 {
		d.SetId("")
		return nil
	}

	if err != nil {
		return diag.FromErr(err)
	}

	return setAuditContext(d, c.Org, auditCtx)
}

func resourceAuditContextUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	if d.HasChanges("description", "tags") {

		c := m.(*client.Client)

		auditCtxToUpdate := client.AuditContext{}
		auditCtxToUpdate.Name = GetString(d.Get("name"))

		if d.HasChange("description") {
			auditCtxToUpdate.Description = GetDescriptionString(d.Get("description"), *auditCtxToUpdate.Name)
		}

		if d.HasChange("tags") {
			auditCtxToUpdate.Tags = GetTagChanges(d)
		}

		updatedAuditCtx, _, err := c.UpdateAuditContext(auditCtxToUpdate)
		if err != nil {
			return diag.FromErr(err)
		}

		return setAuditContext(d, c.Org, updatedAuditCtx)
	}

	return nil
}

func resourceAuditContextDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	d.SetId("")
	return nil
}

func setAuditContext(d *schema.ResourceData, org string, auditCtx *client.AuditContext) diag.Diagnostics {

	if auditCtx == nil {
		d.SetId("")
		return nil
	}

	d.SetId(*auditCtx.Name)

	if err := SetBase(d, auditCtx.Base); err != nil {
		return diag.FromErr(err)
	}

	if err := SetSelfLink(auditCtx.Links, d); err != nil {
		return diag.FromErr(err)
	}

	return nil
}
