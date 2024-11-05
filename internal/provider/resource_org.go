package cpln

import (
	"context"
	"fmt"

	client "github.com/controlplane-com/terraform-provider-cpln/internal/provider/client"

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
			Type:        schema.TypeString,
			Description: "The ID, in GUID format, of the org.",
			Computed:    true,
		},
		"name": {
			Type:        schema.TypeString,
			Description: "The name of the org.",
			Computed:    true,
		},
		"description": {
			Type:             schema.TypeString,
			Description:      "The description of org.",
			Optional:         true,
			ValidateFunc:     DescriptionValidator,
			DiffSuppressFunc: DiffSuppressDescription,
		},
		"tags": {
			Type:        schema.TypeMap,
			Description: "Key-value map of the org's tags.",
			Optional:    true,
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
			ValidateFunc: TagValidator,
		},
		"self_link": {
			Type:        schema.TypeString,
			Description: "Full link to this resource. Can be referenced by other resources.",
			Computed:    true,
		},
		"status": {
			Type:        schema.TypeList,
			Description: "Status of the org.",
			Computed:    true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"account_link": {
						Type:        schema.TypeString,
						Description: "The link of the account the org belongs to.",
						Optional:    true,
					},
					"active": {
						Type:        schema.TypeBool,
						Description: "Indicates whether the org is active or not.",
						Optional:    true,
					},
				},
			},
		},
		"account_id": {
			Type:        schema.TypeString,
			Description: "The associated account ID that will be used when creating the org. Only used on org creation. The account ID can be obtained from the `Org Management & Billing` page.",
			Optional:    true,
		},
		"invitees": {
			Type:        schema.TypeSet,
			Description: "When an org is created, the list of email addresses which will receive an invitation to join the org and be assigned to the `superusers` group. The user account used when creating the org will be included in this list.",
			Optional:    true,
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},
		"session_timeout_seconds": {
			Type:        schema.TypeInt,
			Description: "The idle time (in seconds) in which the console UI will automatically sign-out the user. Default: 900 (15 minutes)",
			Optional:    true,
			Default:     900,
		},
		"auth_config": {
			Type:        schema.TypeList,
			Description: "The configuration settings and parameters related to authentication within the org.",
			Optional:    true,
			MaxItems:    1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"domain_auto_members": {
						Type:        schema.TypeSet,
						Description: "List of domains which will auto-provision users when authenticating using SAML.",
						Required:    true,
						Elem: &schema.Schema{
							Type: schema.TypeString,
						},
					},
					"saml_only": {
						Type:        schema.TypeBool,
						Description: "Enforce SAML only authentication.",
						Optional:    true,
						Default:     false,
					},
				},
			},
		},
		"observability": {
			Type:        schema.TypeList,
			Description: "The retention period (in days) for logs, metrics, and traces. Charges apply for storage beyond the 30 day default.",
			Required:    true,
			MaxItems:    1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"logs_retention_days": {
						Type:         schema.TypeInt,
						Description:  "Log retention days. Default: 30",
						Default:      30,
						Optional:     true,
						ValidateFunc: ObservabilityValidator,
					},
					"metrics_retention_days": {
						Type:         schema.TypeInt,
						Description:  "Metrics retention days. Default: 30",
						Default:      30,
						Optional:     true,
						ValidateFunc: ObservabilityValidator,
					},
					"traces_retention_days": {
						Type:         schema.TypeInt,
						Description:  "Traces retention days. Default: 30",
						Default:      30,
						Optional:     true,
						ValidateFunc: ObservabilityValidator,
					},
				},
			},
		},
		"security": {
			Type:        schema.TypeList,
			Description: "",
			Optional:    true,
			MaxItems:    1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"threat_detection": {
						Type:        schema.TypeList,
						Description: "",
						Optional:    true,
						MaxItems:    1,
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"enabled": {
									Type:        schema.TypeBool,
									Description: "Indicates whether threat detection should be forwarded or not.",
									Required:    true,
								},
								"minimum_severity": {
									Type:        schema.TypeString,
									Description: "Any threats with this severity and more severe will be sent. Others will be ignored. Valid values: `warning`, `error`, or `critical`.",
									Optional:    true,
								},
								"syslog": {
									Type:        schema.TypeList,
									Description: "Configuration for syslog forwarding.",
									Optional:    true,
									MaxItems:    1,
									Elem: &schema.Resource{
										Schema: map[string]*schema.Schema{
											"transport": {
												Type:        schema.TypeString,
												Description: "The transport-layer protocol to send the syslog messages over. If TCP is chosen, messages will be sent with TLS. Default: `tcp`.",
												Optional:    true,
												Default:     "tcp",
											},
											"host": {
												Type:        schema.TypeString,
												Description: "The hostname to send syslog messages to.",
												Required:    true,
											},
											"port": {
												Type:        schema.TypeInt,
												Description: "The port to send syslog messages to.",
												Required:    true,
											},
										},
									},
								},
							},
						},
					},
					"_sentinel": {
						Type:     schema.TypeBool,
						Optional: true,
						Default:  true,
					},
				},
			},
		},
	}
}

func importStateOrg(ctx context.Context, d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	return []*schema.ResourceData{d}, nil
}

func resourceOrgCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	c := m.(*client.Client)

	currentOrg, _, err := c.GetOrg()

	if err != nil {

		accountId := d.Get("account_id").(string)
		invitees := []string{}

		for _, value := range d.Get("invitees").(*schema.Set).List() {
			invitees = append(invitees, value.(string))
		}

		if accountId != "" && len(invitees) > 0 {

			org := client.Org{}

			org.Name = &c.Org
			org.Description = GetString(d.Get("description"))
			org.Tags = GetStringMap(d.Get("tags"))

			createOrgRequest := client.CreateOrgRequest{
				Org:      &org,
				Invitees: &invitees,
			}

			responseCode := 0

			// Make the request to create the org
			currentOrg, responseCode, err = c.CreateOrg(accountId, createOrgRequest)

			if err != nil {
				if responseCode == 409 {
					currentOrg = &client.Org{}
					currentOrg.Name = &c.Org
				} else {
					return diag.FromErr(fmt.Errorf("org %s cannot be created. Error: %s", *org.Name, err))
				}
			}

		} else {
			return diag.FromErr(err)
		}
	}

	currentOrg.Description = GetString(d.Get("description"))
	currentOrg.Tags = GetStringMap(d.Get("tags"))
	currentOrg.Spec = nil
	currentOrg.SpecReplace = &client.OrgSpec{
		AuthConfig:            buildAuthConfig(d.Get("auth_config").([]interface{})),
		Observability:         buildObservability(d.Get("observability").([]interface{})),
		SessionTimeoutSeconds: GetInt(d.Get("session_timeout_seconds").(int)),
	}

	if d.Get("security") != nil {
		currentOrg.SpecReplace.Security = buildOrgSecurity(d.Get("security").([]interface{}))
	}

	// Make the request to update the org
	updatedOrg, _, err := c.UpdateOrg(*currentOrg)
	if err != nil {
		return diag.FromErr(err)
	}

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

	if d.HasChanges("description", "tags", "session_timeout_seconds", "auth_config", "observability", "security") {

		c := m.(*client.Client)

		orgToUpdate := client.Org{
			SpecReplace: &client.OrgSpec{},
		}

		orgToUpdate.Name = GetString(d.Get("name"))
		orgToUpdate.Description = GetDescriptionString(d.Get("description"), *orgToUpdate.Name)
		orgToUpdate.Tags = GetTagChanges(d)
		orgToUpdate.SpecReplace.SessionTimeoutSeconds = GetInt(d.Get("session_timeout_seconds").(int))
		orgToUpdate.SpecReplace.AuthConfig = buildAuthConfig(d.Get("auth_config").([]interface{}))
		orgToUpdate.SpecReplace.Observability = buildObservability(d.Get("observability").([]interface{}))

		if d.Get("security") != nil {
			orgToUpdate.SpecReplace.Security = buildOrgSecurity(d.Get("security").([]interface{}))
		}

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

	c := m.(*client.Client)

	org := client.Org{
		Base: client.Base{
			Name:        GetString(d.Get("name")),
			Description: GetString(d.Get("name")),
			TagsReplace: GetStringMap(nil),
		},
		SpecReplace: &client.OrgSpec{
			SessionTimeoutSeconds: GetInt(900),
			Observability: &client.Observability{
				LogsRetentionDays:    GetInt(30),
				MetricsRetentionDays: GetInt(30),
				TracesRetentionDays:  GetInt(30),
			},
			AuthConfig: nil,
			Security:   nil,
		},
	}

	_, _, err := c.UpdateOrg(org)
	if err != nil {
		return diag.FromErr(err)
	}

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

		if err := d.Set("security", flattenOrgSecurity(org.Spec.Security)); err != nil {
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

	spec := specs[0].(map[string]interface{})

	return &client.Observability{
		LogsRetentionDays:    GetInt(spec["logs_retention_days"].(int)),
		MetricsRetentionDays: GetInt(spec["metrics_retention_days"].(int)),
		TracesRetentionDays:  GetInt(spec["traces_retention_days"].(int)),
	}
}

func buildOrgSecurity(specs []interface{}) *client.OrgSecurity {

	if len(specs) == 0 || specs[0] == nil {
		return nil
	}

	spec := specs[0].(map[string]interface{})
	output := client.OrgSecurity{}

	if spec["threat_detection"] != nil {
		output.ThreatDetection = buildOrgThreatDetection(spec["threat_detection"].([]interface{}))
	}

	return &output
}

func buildOrgThreatDetection(specs []interface{}) *client.OrgThreatDetection {

	if len(specs) == 0 || specs[0] == nil {
		return nil
	}

	spec := specs[0].(map[string]interface{})
	output := client.OrgThreatDetection{}

	if spec["enabled"] != nil {
		output.Enabled = GetBool(spec["enabled"])
	}

	if spec["minimum_severity"] != nil {
		output.MinimumSeverity = GetString(spec["minimum_severity"])
	}

	if spec["syslog"] != nil {
		output.Syslog = buildOrgThreatDetectionSyslog(spec["syslog"].([]interface{}))
	}

	return &output
}

func buildOrgThreatDetectionSyslog(specs []interface{}) *client.OrgThreatDetectionSyslog {

	if len(specs) == 0 || specs[0] == nil {
		return nil
	}

	spec := specs[0].(map[string]interface{})
	output := client.OrgThreatDetectionSyslog{}

	if spec["transport"] != nil {
		output.Transport = GetString(spec["transport"])
	}

	if spec["host"] != nil {
		output.Host = GetString(spec["host"])
	}

	if spec["port"] != nil {
		output.Port = GetInt(spec["port"])
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

func flattenOrgSecurity(spec *client.OrgSecurity) []interface{} {

	if spec == nil {
		return nil
	}

	output := map[string]interface{}{
		"_sentinel": true,
	}

	if spec.ThreatDetection != nil {
		output["threat_detection"] = flattenOrgThreatDetection(spec.ThreatDetection)
	}

	return []interface{}{
		output,
	}
}

func flattenOrgThreatDetection(spec *client.OrgThreatDetection) []interface{} {

	if spec == nil {
		return nil
	}

	output := map[string]interface{}{}

	if spec.Enabled != nil {
		output["enabled"] = *spec.Enabled
	}

	if spec.MinimumSeverity != nil {
		output["minimum_severity"] = *spec.MinimumSeverity
	}

	if spec.Syslog != nil {
		output["syslog"] = flattenOrgThreatDetectionSyslog(spec.Syslog)
	}

	return []interface{}{
		output,
	}
}

func flattenOrgThreatDetectionSyslog(spec *client.OrgThreatDetectionSyslog) []interface{} {

	if spec == nil {
		return nil
	}

	output := map[string]interface{}{
		"port": *spec.Port,
	}

	if spec.Transport != nil {
		output["transport"] = *spec.Transport
	}

	if spec.Host != nil {
		output["host"] = *spec.Host
	}

	return []interface{}{
		output,
	}
}
