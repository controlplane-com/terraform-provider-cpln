package cpln

import (
	"context"

	client "github.com/controlplane-com/terraform-provider-cpln/internal/provider/client"
	models "github.com/controlplane-com/terraform-provider-cpln/internal/provider/models/location"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure data source implements required interfaces.
var (
	_ datasource.DataSource              = &LocationsDataSource{}
	_ datasource.DataSourceWithConfigure = &LocationsDataSource{}
)

/*** Data Source Model ***/

// LocationsDataSourceModel holds the Terraform state for the data source.
type LocationsDataSourceModel struct {
	Locations types.List `tfsdk:"locations"`
}

/*** Data Source Configuration ***/

// LocationsDataSource is the data source implementation.
type LocationsDataSource struct {
	EntityBase
}

// NewLocationsDataSource returns a new instance of the data source implementation.
func NewLocationsDataSource() datasource.DataSource {
	return &LocationsDataSource{}
}

// Metadata provides the data source type name.
func (d *LocationsDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = "cpln_locations"
}

// Configure configures the data source before use.
func (d *LocationsDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	d.EntityBaseConfigure(ctx, req.ProviderData, &resp.Diagnostics)
}

// Schema defines the schema for the data source.
func (d *LocationsDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"locations": schema.ListNestedAttribute{
				Description: "List of all images of the org.",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
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
				},
			},
		},
	}
}

// Read fetches the current state of the resource.
func (d *LocationsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	// Declare variable to hold existing state
	var state LocationsDataSourceModel

	// Populate state from request and capture diagnostics
	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)

	// Abort if diagnostics errors occurred
	if resp.Diagnostics.HasError() {
		// Exit early on error
		return
	}

	// Create a new operator instance
	operator := LocationsDataSourceOperator{
		Ctx:    ctx,
		Diags:  &resp.Diagnostics,
		Client: d.client,
		Plan:   state,
	}

	// Invoke API to read resource details
	apiResp, err := operator.InvokeRead()

	// Handle API invocation errors
	if err != nil {
		// Report API error
		resp.Diagnostics.AddError("API error", err.Error())

		// Exit on API error
		return
	}

	// Build new state from API response
	newState := operator.MapResponseToState(apiResp)

	// Abort if diagnostics errors occurred
	if resp.Diagnostics.HasError() {
		return
	}

	// Persist updated state into Terraform
	resp.Diagnostics.Append(resp.State.Set(ctx, &newState)...)
}

/*** Data Source Operator ***/

// LocationsDataSourceOperator is the operator for managing the state.
type LocationsDataSourceOperator struct {
	Ctx    context.Context
	Diags  *diag.Diagnostics
	Client *client.Client
	Plan   LocationsDataSourceModel
}

// MapResponseToState creates a state model from response payload.
func (leo *LocationsDataSourceOperator) MapResponseToState(locations *client.Locations) LocationsDataSourceModel {
	// Initialize a new request payload
	state := LocationsDataSourceModel{}

	// Set specific attributes
	state.Locations = leo.flattenLocations(locations)

	// Return completed state model
	return state
}

// InvokeRead invokes the Get API to retrieve an existing resource by name.
func (leo *LocationsDataSourceOperator) InvokeRead() (*client.Locations, error) {
	return leo.Client.GetLocations()
}

// Flatteners //

// flattenLocations transforms *client.Locations into a Terraform types.List.
func (leo *LocationsDataSourceOperator) flattenLocations(input *client.Locations) types.List {
	// Get attribute types
	elementType := models.LocationModel{}.AttributeTypes()

	// Check if the input is nil
	if input == nil {
		// Return a null list
		return types.ListNull(elementType)
	}

	// Define the blocks slice
	var blocks []models.LocationModel

	// Iterate over the slice and construct the blocks
	for _, item := range input.Items {
		// Construct a block
		block := models.LocationModel{
			CplnId:        types.StringPointerValue(item.ID),
			Name:          types.StringPointerValue(item.Name),
			Description:   types.StringPointerValue(item.Description),
			Tags:          FlattenTags(item.Tags),
			SelfLink:      FlattenSelfLink(item.Links),
			CloudProvider: types.StringPointerValue(item.Provider),
			Region:        types.StringPointerValue(item.Region),
		}

		// Handle the case where spec could be nil
		if item.Spec != nil {
			block.Enabled = types.BoolPointerValue(item.Spec.Enabled)
		} else {
			block.Enabled = types.BoolNull()
		}

		// Handle the case where status could be nil
		if item.Status != nil {
			block.Geo = leo.flattenGeo(item.Status.Geo)
			block.IpRanges = FlattenSetString(item.Status.IpRanges)
		} else {
			block.Geo = types.ListNull(models.GeoModel{}.AttributeTypes())
			block.IpRanges = types.SetNull(types.StringType)
		}

		// Append the constructed block to the blocks slice
		blocks = append(blocks, block)
	}

	// Return the successfully created types.List
	return FlattenList(leo.Ctx, leo.Diags, blocks)
}

// flattenGeo transforms *client.LocationGeo into a Terraform types.List.
func (leo *LocationsDataSourceOperator) flattenGeo(input *client.LocationGeo) types.List {
	// Get attribute types
	elementType := models.GeoModel{}.AttributeTypes()

	// Check if the input is nil
	if input == nil {
		// Return a null list
		return types.ListNull(elementType)
	}

	// Build a single block
	block := models.GeoModel{
		Lat:       types.Float32PointerValue(input.Lat),
		Lon:       types.Float32PointerValue(input.Lon),
		Country:   types.StringPointerValue(input.Country),
		State:     types.StringPointerValue(input.State),
		City:      types.StringPointerValue(input.City),
		Continent: types.StringPointerValue(input.Continent),
	}

	// Return the successfully created types.List
	return FlattenList(leo.Ctx, leo.Diags, []models.GeoModel{block})
}
