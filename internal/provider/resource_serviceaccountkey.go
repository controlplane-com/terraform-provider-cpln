package cpln

import (
	"context"
	"fmt"
	"strings"
	client "terraform-provider-cpln/internal/provider/client"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceServiceAccountKey() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceServiceAccountKeyCreate,
		ReadContext:   resourceServiceAccountKeyRead,
		// UpdateContext: resourceServiceAccountKeyUpdate,
		DeleteContext: resourceServiceAccountKeyDelete,
		Schema: map[string]*schema.Schema{
			"service_account_name": {
				Type:         schema.TypeString,
				ForceNew:     true,
				Required:     true,
				ValidateFunc: NameValidator,
			},
			"description": {
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
				// ValidateFunc: DescriptionValidator,
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"created": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"key": {
				Type:      schema.TypeString,
				Sensitive: true,
				Computed:  true,
			},
		},
		Importer: &schema.ResourceImporter{
			StateContext: importStateServiceAccountKey,
		},
	}
}

func importStateServiceAccountKey(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {

	parts := strings.SplitN(d.Id(), ":", 2)

	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		return nil, fmt.Errorf("unexpected format of ID (%s), expected ID syntax: 'service_account_name:key_name'. Example: 'terraform import cpln_service_account_key.RESOURCE_NAME SERVICE_ACCOUNT_NAME:KEY_NAME'", d.Id())
	}

	d.Set("service_account_name", parts[0])
	d.SetId(parts[1])

	return []*schema.ResourceData{d}, nil
}

func resourceServiceAccountKeyCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	serviceAccountName := d.Get("service_account_name").(string)
	keyDescription := d.Get("description").(string)

	c := m.(*client.Client)
	key, err := c.AddServiceAccountKey(serviceAccountName, keyDescription)
	if err != nil {
		return diag.FromErr(err)
	}

	return setServiceAccountKey(d, key)
}

func resourceServiceAccountKeyRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	keyName := d.Id()
	serviceAccountName := d.Get("service_account_name").(string)

	c := m.(*client.Client)
	sa, code, err := c.GetServiceAccount(serviceAccountName)

	if code == 404 {
		d.SetId("")
		return nil
	}

	if err != nil {
		return diag.FromErr(err)
	}

	if *sa.Keys != nil {
		for _, k := range *sa.Keys {
			if k.Name == keyName {
				return setServiceAccountKey(d, &k)
			}
		}
	}

	return resourceServiceAccountKeyDelete(ctx, d, m)
}

func setServiceAccountKey(d *schema.ResourceData, saKey *client.ServiceAccountKey) diag.Diagnostics {

	if saKey == nil {
		d.SetId("")
		return nil
	}

	d.SetId(saKey.Name)

	if err := d.Set("name", saKey.Name); err != nil {
		return diag.FromErr(err)
	}

	if saKey.Key != "" {
		if err := d.Set("key", saKey.Key); err != nil {
			return diag.FromErr(err)
		}
	}

	if err := d.Set("description", saKey.Description); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

// func resourceServiceAccountKeyUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

// 	var diags diag.Diagnostics

// 	diags = append(diags, diag.Diagnostic{
// 		Severity: diag.Error,
// 		Summary:  "Unable to update Service Account Key",
// 		Detail:   "Unable to update Service Account Key. Key can only be created or deleted",
// 	})

// 	return diags
// }

func resourceServiceAccountKeyDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	serviceAccountName := d.Get("service_account_name").(string)

	c := m.(*client.Client)
	err := c.RemoveServiceAccountKey(serviceAccountName, d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")

	return nil
}
