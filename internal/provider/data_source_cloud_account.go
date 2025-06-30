package cpln

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"

	client "github.com/controlplane-com/terraform-provider-cpln/internal/provider/client"
)

// Ensure data source implements required interfaces.
var (
	_ datasource.DataSource              = &CloudAccountDataSource{}
	_ datasource.DataSourceWithConfigure = &CloudAccountDataSource{}
)

// CloudAccountDataSourceModel holds the Terraform state for the data source.
type CloudAccountDataSourceModel struct {
	Id             types.String `tfsdk:"id"`
	AwsIdentifiers types.Set    `tfsdk:"aws_identifiers"`
}

// CloudAccountDataSource is the data source implementation.
type CloudAccountDataSource struct {
	client *client.Client
}

// NewCloudAccountDataSource returns a new instance of the data source implementation.
func NewCloudAccountDataSource() datasource.DataSource {
	return &CloudAccountDataSource{}
}

// Metadata provides the data source type name.
func (d *CloudAccountDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = "cpln_cloud_account"
}

// Configure configures the data source before use.
func (d *CloudAccountDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	// Assert that ProviderData is of the expected type *client.Client
	c, ok := req.ProviderData.(*client.Client)
	if !ok {
		// If the type assertion fails, add an error diagnostic indicating the issue and specifying the expected type
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *client.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}

	// Assign the retrieved client to the resource's client field, enabling API calls
	d.client = c
}

// Schema defines the schema for the data source.
func (d *CloudAccountDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The ID of the cloud account to retrieve.",
				Computed:    true,
			},
			"aws_identifiers": schema.SetAttribute{
				Description: "Unique identifiers associated with resources and services within an Amazon Web Services (AWS) environment.",
				Computed:    true,
				ElementType: types.StringType,
			},
		},
	}
}

// Read retrieves the AWS identifiers for the given cloud account ID and sets the data source state.
func (d *CloudAccountDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	// Declare variable to hold existing state
	var state CloudAccountDataSourceModel

	// Populate state from request and capture diagnostics
	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)

	// Abort if diagnostics errors occurred
	if resp.Diagnostics.HasError() {
		// Exit early on error
		return
	}

	// Build new state
	newState := CloudAccountDataSourceModel{
		Id:             types.StringValue("static-cloud-account"),
		AwsIdentifiers: FlattenSetString(&CloudAccountIdentifiers),
	}

	// Abort if diagnostics errors occurred
	if resp.Diagnostics.HasError() {
		return
	}

	// Persist updated state into Terraform
	resp.Diagnostics.Append(resp.State.Set(ctx, &newState)...)
}
