package cpln

import (
	"context"
	"os"

	client "github.com/controlplane-com/terraform-provider-cpln/internal/provider/client"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ provider.Provider = &CplnProvider{}
)

// CplnProvider is the provider implementation.
type CplnProvider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
	client  *client.Client
}

// CplnProviderModel maps provider schema data to a Go type.
type CplnProviderModel struct {
	Org          types.String `tfsdk:"org"`
	Endpoint     types.String `tfsdk:"endpoint"`
	Profile      types.String `tfsdk:"profile"`
	Token        types.String `tfsdk:"token"`
	RefreshToken types.String `tfsdk:"refresh_token"`
}

// New is a helper function to simplify provider server and testing implementation.
func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &CplnProvider{
			version: version,
		}
	}
}

// Metadata returns the provider type name.
func (p *CplnProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "cpln"
	resp.Version = p.version
}

// Schema defines the provider-level schema for configuration data.
func (p *CplnProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"org": schema.StringAttribute{
				Optional:    true,
				Description: "The Control Plane org that this provider will perform actions against. Can be specified with the CPLN_ORG environment variable.",
			},
			"endpoint": schema.StringAttribute{
				Optional:    true,
				Description: "The Control Plane Data Service API endpoint. Default is: https://api.cpln.io. Can be specified with the CPLN_ENDPOINT environment variable.",
			},
			"profile": schema.StringAttribute{
				Optional:    true,
				Description: "The user/service account profile that this provider will use to authenticate to the data service. Can be specified with the CPLN_PROFILE environment variable.",
			},
			"token": schema.StringAttribute{
				Optional:    true,
				Sensitive:   true,
				Description: "A generated token that can be used to authenticate to the data service API. Can be specified with the CPLN_TOKEN environment variable.",
			},
			"refresh_token": schema.StringAttribute{
				Optional:    true,
				Sensitive:   true,
				Description: "A generated token that can be used to authenticate to the data service API. Can be specified with the CPLN_REFRESH_TOKEN environment variable. Used when the provider is required to create an org or update the auth_config property. Refer to the section above on how to obtain the refresh token.",
			},
		},
	}
}

// Configure prepares a Control Plane API client for data sources and resources.
func (p *CplnProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	// Retrieve provider data from configuration
	var config CplnProviderModel

	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Give config attributes environment variables as default values if unknown or null
	if config.Org.IsNull() || config.Org.IsUnknown() {
		config.Org = types.StringValue(os.Getenv("CPLN_ORG"))
	}

	if config.Endpoint.IsNull() || config.Endpoint.IsUnknown() {
		if endpoint := os.Getenv("CPLN_ENDPOINT"); endpoint != "" {
			config.Endpoint = types.StringValue(endpoint)
		}
	}

	if config.Profile.IsNull() || config.Profile.IsUnknown() {
		config.Profile = types.StringValue(os.Getenv("CPLN_PROFILE"))
	}

	if config.Token.IsNull() || config.Token.IsUnknown() {
		config.Token = types.StringValue(os.Getenv("CPLN_TOKEN"))
	}

	if config.RefreshToken.IsNull() || config.RefreshToken.IsUnknown() {
		config.RefreshToken = types.StringValue(os.Getenv("CPLN_REFRESH_TOKEN"))
	}

	// Create a new cpln client using the configuration values
	c, err := client.NewClient(
		config.Org.ValueStringPointer(),
		config.Endpoint.ValueStringPointer(),
		config.Profile.ValueStringPointer(),
		config.Token.ValueStringPointer(),
		config.RefreshToken.ValueStringPointer(),
		p.version,
	)

	// Handle client initialization error
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to Create Cpln API Client",
			"An error occurred while attempting to create the Cpln API client. "+
				"Please verify your configuration or try again. "+
				"If the issue persists, consider reaching out to the provider's support team.\n\n"+
				"Detailed Error: "+err.Error(),
		)

		return
	}

	// Set provider client
	p.client = c

	// Make the cpln client available during DataSource and Resource type Configure methods
	resp.DataSourceData = c
	resp.ResourceData = c
}

// DataSources defines the data sources implemented in the provider.
func (p *CplnProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewCloudAccountDataSource,
		NewGvcDataSource,
		NewHelmTemplateDataSource,
		NewImageDataSource,
		NewImagesDataSource,
		NewLocationDataSource,
		NewLocationsDataSource,
		NewOrgDataSource,
		NewSecretDataSource,
		NewWorkloadDataSource,
	}
}

// Resources defines the resources implemented in the provider.
func (p *CplnProvider) Resources(_ context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewAgentResource,
		NewAuditContextResource,
		NewCatalogTemplateResource,
		NewCloudAccountResource,
		NewCustomLocationResource,
		NewDomainRouteResource,
		NewDomainResource,
		NewGroupResource,
		NewGvcResource,
		NewHelmReleaseResource,
		NewIdentityResource,
		NewIpSetResource,
		NewLocationResource,
		NewMk8sResource,
		NewMk8sKubeconfigResource,
		NewOrgLoggingResource,
		NewOrgTracingResource,
		NewOrgResource,
		NewPolicyResource,
		NewSecretResource,
		NewServiceAccountKeyResource,
		NewServiceAccountResource,
		NewVolumeSetResource,
		NewWorkloadResource,
	}
}
