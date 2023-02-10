package cpln

import (
	"context"

	client "terraform-provider-cpln/internal/provider/client"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Provider -
func Provider() *schema.Provider {

	return &schema.Provider{

		Schema: map[string]*schema.Schema{
			"org": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("CPLN_ORG", ""),
			},
			"endpoint": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("CPLN_ENDPOINT", "https://api.cpln.io"),
			},
			"profile": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("CPLN_PROFILE", ""),
			},
			"token": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("CPLN_TOKEN", ""),
			},
		},

		ResourcesMap: map[string]*schema.Resource{
			"cpln_agent":               resourceAgent(),
			"cpln_audit_context":       resourceAuditContext(),
			"cpln_cloud_account":       resourceCloudAccount(),
			"cpln_domain":              resourceDomain(),
			"cpln_group":               resourceGroup(),
			"cpln_gvc":                 resourceGvc(),
			"cpln_identity":            resourceIdentity(),
			"cpln_org_logging":         resourceOrgLogging(),
			"cpln_org_tracing":         resourceOrgTracing(),
			"cpln_policy":              resourcePolicy(),
			"cpln_secret":              resourceSecret(),
			"cpln_service_account":     resourceServiceAccount(),
			"cpln_service_account_key": resourceServiceAccountKey(),
			"cpln_workload":            resourceWorkload(),
		},

		DataSourcesMap: map[string]*schema.Resource{
			// "cpln_gvcs": dataSourceGvcs(),
			"cpln_gvc":       dataSourceGvc(),
			"cpln_location":  dataSourceLocation(),
			"cpln_locations": dataSourceLocations(),
			"cpln_org":       dataSourceOrg(),
		},

		ConfigureContextFunc: providerConfigure,
	}
}

func providerConfigure(_ context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {

	org := d.Get("org").(string)
	host := d.Get("endpoint").(string)
	profile := d.Get("profile").(string)
	token := d.Get("token").(string)

	var diags diag.Diagnostics

	httpClient, err := client.NewClient(&org, &host, &profile, &token)

	if err != nil {
		return nil, diag.FromErr(err)
	}

	return httpClient, diags
}
