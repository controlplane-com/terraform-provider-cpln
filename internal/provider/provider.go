package cpln

import (
	"context"

	client "github.com/controlplane-com/terraform-provider-cpln/internal/provider/client"

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
			"refresh_token": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("CPLN_REFRESH_TOKEN", ""),
			},
		},

		ResourcesMap: map[string]*schema.Resource{
			"cpln_agent":               resourceAgent(),
			"cpln_audit_context":       resourceAuditContext(),
			"cpln_cloud_account":       resourceCloudAccount(),
			"cpln_custom_location":     resourceCustomLocation(),
			"cpln_domain":              resourceDomain(),
			"cpln_domain_route":        resourceDomainRoute(),
			"cpln_group":               resourceGroup(),
			"cpln_gvc":                 resourceGvc(),
			"cpln_identity":            resourceIdentity(),
			"cpln_ipset":               resourceIpSet(),
			"cpln_location":            resourceLocation(),
			"cpln_mk8s":                resourceMk8s(),
			"cpln_org_logging":         resourceOrgLogging(),
			"cpln_org_tracing":         resourceOrgTracing(),
			"cpln_org":                 resourceOrg(),
			"cpln_policy":              resourcePolicy(),
			"cpln_secret":              resourceSecret(),
			"cpln_service_account":     resourceServiceAccount(),
			"cpln_service_account_key": resourceServiceAccountKey(),
			"cpln_volume_set":          resourceVolumeSet(),
			"cpln_workload":            resourceWorkload(),
		},

		DataSourcesMap: map[string]*schema.Resource{
			// "cpln_gvcs": dataSourceGvcs(),
			"cpln_cloud_account": dataSourceCloudAccount(),
			"cpln_gvc":           dataSourceGvc(),
			"cpln_image":         dataSourceImage(),
			"cpln_images":        dataSourceImages(),
			"cpln_location":      dataSourceLocation(),
			"cpln_locations":     dataSourceLocations(),
			"cpln_org":           dataSourceOrg(),
		},

		ConfigureContextFunc: providerConfigure,
	}
}

func providerConfigure(_ context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {

	org := d.Get("org").(string)
	host := d.Get("endpoint").(string)
	profile := d.Get("profile").(string)
	token := d.Get("token").(string)
	refreshToken := d.Get("refresh_token").(string)

	var diags diag.Diagnostics

	httpClient, err := client.NewClient(&org, &host, &profile, &token, &refreshToken)

	if err != nil {
		return nil, diag.FromErr(err)
	}

	return httpClient, diags
}
