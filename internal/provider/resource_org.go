package cpln

import (
	"context"
	client "terraform-provider-cpln/internal/provider/client"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

/*** Main ***/
func resourceOrg() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceOrgCreate,
		ReadContext:   resourceOrgRead,
		UpdateContext: resourceOrgUpdate,
		DeleteContext: resourceOrgDelete,
		Schema:        orgSchema(),
		Importer: &schema.ResourceImporter{
			StateContext: importStateOrg,
		},
	}
}

func orgSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"cpln_id": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"name": {
			Type:     schema.TypeString,
			Computed: true,
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
		"status": {
			Type:     schema.TypeList,
			Computed: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"account_link": {
						Type:     schema.TypeString,
						Optional: true,
					},
					"active": {
						Type:     schema.TypeBool,
						Optional: true,
					},
				},
			},
		},
		"create_org": {
			Type:     schema.TypeBool,
			Optional: true,
			Default:  false,
		},
		"account_id": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"invitees": {
			Type:     schema.TypeSet,
			Optional: true,
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},
		"session_timeout_seconds": {
			Type:     schema.TypeInt,
			Optional: true,
			Default:  900,
		},
		"auth_config": {
			Type:     schema.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"domain_auto_members": {
						Type:     schema.TypeSet,
						Required: true,
						Elem: &schema.Schema{
							Type: schema.TypeString,
						},
					},
					"saml_only": {
						Type:     schema.TypeBool,
						Optional: true,
						Default:  false,
					},
				},
			},
		},
		"observability": {
			Type:     schema.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"logs_retention_days": {
						Type:         schema.TypeInt,
						Default:      30,
						Optional:     true,
						ValidateFunc: ObservabilityValidator,
					},
					"metrics_retention_days": {
						Type:         schema.TypeInt,
						Default:      30,
						Optional:     true,
						ValidateFunc: ObservabilityValidator,
					},
					"traces_retention_days": {
						Type:         schema.TypeInt,
						Default:      30,
						Optional:     true,
						ValidateFunc: ObservabilityValidator,
					},
				},
			},
		},
	}
}

func importStateOrg(ctx context.Context, d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {

	// TODO: Need to review and implement

	// c := m.(*client.Client)

	// Get then set account id
	// account, _, err := c.GetOrgAccount(d.Id())

	// if err != nil {
	// 	return nil, fmt.Errorf("import org %s failed. Error: %s", d.Id(), err)
	// }

	// d.Set("account_id", account.ID)

	// // Set invitees to empty
	// d.Set("invitees", []interface{}{})

	return []*schema.ResourceData{d}, nil
}

func resourceOrgCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	c := m.(*client.Client)
	org := client.Org{}

	// org.Name = GetString(d.Get("name"))
	org.Name = &c.Org
	org.Description = GetString(d.Get("description"))
	org.Tags = GetStringMap(d.Get("tags"))

	// createOrg := GetBool(d.Get("create_org"))

	// if *createOrg {

	// accountId := d.Get("account_id").(string)
	// invitees := []string{}

	// for _, value := range d.Get("invitees").(*schema.Set).List() {
	// 	invitees = append(invitees, value.(string))
	// }

	// if accountId == "" {
	// 	return diag.FromErr(fmt.Errorf("account id must not be empty"))
	// }

	// if len(invitees) == 0 {
	// 	return diag.FromErr(fmt.Errorf("invitees must not be empty"))
	// }

	// createOrgRequest := client.CreateOrgRequest{
	// 	Org:      &org,
	// 	Invitees: &invitees,
	// }

	// // Make the request to create the org
	// createdOrg, code, err := c.CreateOrg(accountId, createOrgRequest)

	// if code == 409 {
	// 	return ResourceExistsHelper()
	// }

	// if err != nil {
	// 	return diag.FromErr(fmt.Errorf("org %s cannot be created. Error: %s", *org.Name, err))
	// }

	// }

	org.SpecReplace = &client.OrgSpec{
		AuthConfig:            buildAuthConfig(d.Get("auth_config").([]interface{})),
		Observability:         buildObservability(d.Get("observability").([]interface{})),
		SessionTimeoutSeconds: GetInt(d.Get("session_timeout_seconds").(int)),
	}

	// if d.Get("session_timeout_seconds") != nil {
	// 	createdOrg.SpecReplace.SessionTimeoutSeconds = GetInt(d.Get("session_timeout_seconds").(int))
	// }

	// Make the request to update the org
	updatedOrg, _, err := c.UpdateOrg(org)
	if err != nil {
		return diag.FromErr(err)
	}

	// // Set invitees
	// flattenedInvitees := []interface{}{}

	// for _, value := range invitees {
	// 	flattenedInvitees = append(flattenedInvitees, value)
	// }

	// if err := d.Set("invitees", flattenedInvitees); err != nil {
	// 	return diag.FromErr(err)
	// }

	return setOrg(d, updatedOrg)
}

func resourceOrgRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	c := m.(*client.Client)
	org, _, err := c.GetSpecificOrg(d.Id())

	if err != nil {
		return diag.FromErr(err)
	}

	return setOrg(d, org)
}

func resourceOrgUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	if d.HasChanges("description", "tags", "session_timeout_seconds", "auth_config", "observability") {

		c := m.(*client.Client)

		orgToUpdate := client.Org{
			SpecReplace: &client.OrgSpec{},
		}

		orgToUpdate.Name = GetString(d.Get("name"))
		orgToUpdate.Description = GetDescriptionString(d.Get("description"), *orgToUpdate.Name)
		orgToUpdate.Tags = GetTagChanges(d)

		if d.Get("session_timeout_seconds") != nil {
			orgToUpdate.SpecReplace.SessionTimeoutSeconds = GetInt(d.Get("session_timeout_seconds").(int))
		}

		orgToUpdate.SpecReplace.AuthConfig = buildAuthConfig(d.Get("auth_config").([]interface{}))
		orgToUpdate.SpecReplace.Observability = buildObservability(d.Get("observability").([]interface{}))

		// Make the request to update the org
		updatedOrg, _, err := c.UpdateOrg(orgToUpdate)
		if err != nil {
			return diag.FromErr(err)
		}

		return setOrg(d, updatedOrg)
	}

	return nil
}

func resourceOrgDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	// c := m.(*client.Client)

	// org := client.Org{
	// 	SpecReplace: &client.OrgSpec{},
	// }

	// org.Name = GetString(d.Get("name"))

	// _, _, err := c.UpdateOrg(org)
	// if err != nil {
	// 	return diag.FromErr(err)
	// }

	d.SetId("")

	return nil
}

func setOrg(d *schema.ResourceData, org *client.Org) diag.Diagnostics {

	if org == nil {
		d.SetId("")
		return nil
	}

	d.SetId(*org.Name)

	if err := SetBase(d, org.Base); err != nil {
		return diag.FromErr(err)
	}

	if err := SetSelfLink(org.Links, d); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("status", flattenOrgStatus(org.Status)); err != nil {
		return diag.FromErr(err)
	}

	if org.Spec != nil {
		if err := d.Set("session_timeout_seconds", org.Spec.SessionTimeoutSeconds); err != nil {
			return diag.FromErr(err)
		}

		if err := d.Set("auth_config", flattenAuthConfig(org.Spec.AuthConfig)); err != nil {
			return diag.FromErr(err)
		}

		if err := d.Set("observability", flattenObservability(org.Spec.Observability)); err != nil {
			return diag.FromErr(err)
		}
	}

	return nil
}

/*** Build ***/
func buildAuthConfig(specs []interface{}) *client.AuthConfig {
	if len(specs) == 0 || specs[0] == nil {
		return nil
	}

	spec := specs[0].(map[string]interface{})
	output := client.AuthConfig{
		SamlOnly: GetBool(spec["saml_only"].(bool)),
	}

	domainAutoMembers := []string{}
	for _, value := range spec["domain_auto_members"].(*schema.Set).List() {
		domainAutoMembers = append(domainAutoMembers, value.(string))
	}

	output.DomainAutoMembers = &domainAutoMembers

	return &output
}

func buildObservability(specs []interface{}) *client.Observability {
	if len(specs) == 0 || specs[0] == nil {
		return nil
	}

	spec := specs[0].(map[string]interface{})
	output := client.Observability{
		LogsRetentionDays:    GetInt(spec["logs_retention_days"].(int)),
		MetricsRetentionDays: GetInt(spec["metrics_retention_days"].(int)),
		TracesRetentionDays:  GetInt(spec["traces_retention_days"].(int)),
	}

	return &output
}

/*** Flatten ***/
func flattenOrgStatus(status *client.OrgStatus) []interface{} {
	if status == nil {
		return nil
	}

	output := map[string]interface{}{}

	if status.AccountLink != nil {
		output["account_link"] = *status.AccountLink
	}

	if status.Active != nil {
		output["active"] = *status.Active
	}

	return []interface{}{
		output,
	}
}

func flattenAuthConfig(spec *client.AuthConfig) []interface{} {
	if spec == nil {
		return nil
	}

	output := map[string]interface{}{
		"saml_only": *spec.SamlOnly,
	}

	if len(*spec.DomainAutoMembers) > 0 {
		output["domain_auto_members"] = []interface{}{}
		for _, value := range *spec.DomainAutoMembers {
			output["domain_auto_members"] = append(output["domain_auto_members"].([]interface{}), value)
		}
	}

	return []interface{}{
		output,
	}
}

func flattenObservability(spec *client.Observability) []interface{} {
	if spec == nil {
		return nil
	}

	output := map[string]interface{}{
		"logs_retention_days":    *spec.LogsRetentionDays,
		"metrics_retention_days": *spec.MetricsRetentionDays,
		"traces_retention_days":  *spec.TracesRetentionDays,
	}

	return []interface{}{
		output,
	}
}
