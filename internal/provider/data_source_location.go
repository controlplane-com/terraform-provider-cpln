package cpln

import (
	"context"

	client "github.com/controlplane-com/terraform-provider-cpln/internal/provider/client"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure data source implements required interfaces.
var (
	_ datasource.DataSource              = &LocationDataSource{}
	_ datasource.DataSourceWithConfigure = &LocationDataSource{}
)

/*** Data Source Configuration ***/

// LocationDataSource is the data source implementation.
type LocationDataSource struct {
	EntityBase
	Operations EntityOperations[LocationResourceModel, client.Location]
}

// NewLocationDataSource returns a new instance of the data source implementation.
func NewLocationDataSource() datasource.DataSource {
	return &LocationDataSource{}
}

// Metadata provides the data source type name.
func (d *LocationDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = "cpln_location"
}

// Configure configures the data source before use.
func (d *LocationDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	d.EntityBaseConfigure(ctx, req.ProviderData, &resp.Diagnostics)
	d.Operations = NewEntityOperations(d.client, &LocationResourceOperator{})
}

// Schema defines the schema for the data source.
func (d *LocationDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The unique identifier for this location.",
				Computed:    true,
			},
			"cpln_id": schema.StringAttribute{
				Description: "The ID, in GUID format, of the location.",
				Computed:    true,
			},
			"name": schema.StringAttribute{
				Description: "Name of the location.",
				Required:    true,
			},
			"description": schema.StringAttribute{
				Description: "Description of the location.",
				Computed:    true,
			},
			"tags": schema.MapAttribute{
				Description: "Key-value map of resource tags.",
				ElementType: types.StringType,
				Optional:    true,
				Computed:    true,
			},
			"self_link": schema.StringAttribute{
				Description: "Full link to this resource. Can be referenced by other resources.",
				Computed:    true,
			},
			"origin": schema.StringAttribute{
				Description: "",
				Computed:    true,
			},
			"cloud_provider": schema.StringAttribute{
				Description: "Cloud Provider of the location.",
				Computed:    true,
			},
			"region": schema.StringAttribute{
				Description: "Region of the location.",
				Computed:    true,
			},
			"enabled": schema.BoolAttribute{
				Description: "Indication if location is enabled.",
				Computed:    true,
			},
			"geo": schema.ListNestedAttribute{
				Description: "",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"lat": schema.Float32Attribute{
							Description: "Latitude of the location.",
							Computed:    true,
						},
						"lon": schema.Float32Attribute{
							Description: "Longitude of the location.",
							Computed:    true,
						},
						"country": schema.StringAttribute{
							Description: "Country of the location.",
							Computed:    true,
						},
						"state": schema.StringAttribute{
							Description: "State of the location.",
							Computed:    true,
						},
						"city": schema.StringAttribute{
							Description: "City of the location.",
							Computed:    true,
						},
						"continent": schema.StringAttribute{
							Description: "Continent of the location.",
							Computed:    true,
						},
					},
				},
			},
			"ip_ranges": schema.SetAttribute{
				Description: "A list of IP ranges of the location.",
				ElementType: types.StringType,
				Computed:    true,
			},
		},
	}
}

// Read fetches the current state of the resource.
func (d *LocationDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	// Declare variable to hold existing state
	var state LocationResourceModel

	// Populate state from request and capture diagnostics
	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)

	// Abort if diagnostics errors occurred
	if resp.Diagnostics.HasError() {
		// Exit early on error
		return
	}

	// Create a new operator instance
	operator := d.Operations.NewOperator(ctx, &resp.Diagnostics, state)

	// Invoke API to read resource details
	apiResp, code, err := operator.InvokeRead(state.Name.ValueString())

	// Remove resource from state if not found
	if code == 404 {
		// Drop resource from Terraform state
		resp.State.RemoveResource(ctx)
		return
	}

	// Handle API invocation errors
	if err != nil {
		// Report API error
		resp.Diagnostics.AddError("API error", err.Error())

		// Exit on API error
		return
	}

	// Build new state from API response
	newState := operator.MapResponseToState(apiResp, true)

	// Abort if diagnostics errors occurred
	if resp.Diagnostics.HasError() {
		return
	}

	// Persist updated state into Terraform
	resp.Diagnostics.Append(resp.State.Set(ctx, &newState)...)
}
