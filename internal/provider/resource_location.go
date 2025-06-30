package cpln

import (
	"context"
	"fmt"
	"slices"

	client "github.com/controlplane-com/terraform-provider-cpln/internal/provider/client"
	models "github.com/controlplane-com/terraform-provider-cpln/internal/provider/models/location"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure resource implements required interfaces.
var (
	_ resource.Resource                = &LocationResource{}
	_ resource.ResourceWithImportState = &LocationResource{}
)

/*** Resource Model ***/

// LocationResourceModel holds the Terraform state for the resource.
type LocationResourceModel struct {
	EntityBaseModel
	CloudProvider types.String `tfsdk:"cloud_provider"`
	Region        types.String `tfsdk:"region"`
	Enabled       types.Bool   `tfsdk:"enabled"`
	Geo           types.List   `tfsdk:"geo"`
	IpRanges      types.Set    `tfsdk:"ip_ranges"`
}

/*** Resource Configuration ***/

// LocationResource is the resource implementation.
type LocationResource struct {
	EntityBase
	Operations EntityOperations[LocationResourceModel, client.Location]
}

// NewLocationResource returns a new instance of the resource implementation.
func NewLocationResource() resource.Resource {
	return &LocationResource{
		EntityBase: EntityBase{
			IsDescriptionComputed: true,
		},
	}
}

// Configure configures the resource before use.
func (lr *LocationResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	lr.EntityBaseConfigure(ctx, req.ProviderData, &resp.Diagnostics)
	lr.Operations = NewEntityOperations(lr.client, &LocationResourceOperator{})
}

// ImportState sets up the import operation to map the imported ID to the "id" attribute in the state.
func (lr *LocationResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// Metadata provides the resource type name.
func (lr *LocationResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = "cpln_location"
}

// Schema defines the schema for the resource.
func (lr *LocationResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: MergeAttributes(lr.EntityBaseAttributes("location"), map[string]schema.Attribute{
			"cloud_provider": schema.StringAttribute{
				Description: "Cloud Provider of the location.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"region": schema.StringAttribute{
				Description: "Region of the location.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"enabled": schema.BoolAttribute{
				Description: "Indication if location is enabled.",
				Required:    true,
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
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
			},
		}),
	}
}

// Create creates the resource.
func (lr *LocationResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	CreateGeneric(ctx, req, resp, lr.Operations)
}

// Read fetches the current state of the resource.
func (lr *LocationResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	ReadGeneric(ctx, req, resp, lr.Operations)
}

// Update modifies the resource.
func (lr *LocationResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	UpdateGeneric(ctx, req, resp, lr.Operations)
}

// Delete removes the resource.
func (lr *LocationResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	DeleteGeneric(ctx, req, resp, lr.Operations)
}

/*** Resource Operator ***/

// LocationResourceOperator is the operator for managing the state.
type LocationResourceOperator struct {
	EntityOperator[LocationResourceModel]
}

// NewAPIRequest creates a request payload from a state model.
func (leo *LocationResourceOperator) NewAPIRequest(isUpdate bool) client.Location {
	// Initialize a new request payload
	requestPayload := client.Location{
		Spec: &client.LocationSpec{},
	}

	// Populate Base fields from state
	leo.Plan.Fill(&requestPayload.Base, isUpdate)

	// Set specific attributes
	requestPayload.Spec.Enabled = BuildBool(leo.Plan.Enabled)

	// Return constructed request payload
	return requestPayload
}

// MapResponseToState constructs the Terraform state model from the API response payload.
func (leo *LocationResourceOperator) MapResponseToState(apiResp *client.Location, isCreate bool) LocationResourceModel {
	// Initialize empty state model
	state := LocationResourceModel{}

	// Populate common fields from base resource data
	state.From(apiResp.Base)

	// Set specific attributes
	state.CloudProvider = types.StringPointerValue(apiResp.Provider)
	state.Region = types.StringPointerValue(apiResp.Region)
	state.Enabled = types.BoolPointerValue(apiResp.Spec.Enabled)

	// Just in case Status is nil
	if apiResp.Status != nil {
		state.Geo = leo.flattenGeo(apiResp.Status.Geo)
		state.IpRanges = FlattenSetString(apiResp.Status.IpRanges)
	} else {
		state.Geo = types.ListNull(models.GeoModel{}.AttributeTypes())
		state.IpRanges = types.SetNull(types.StringType)
	}

	// Return completed state model
	return state
}

// InvokeCreate invokes the Create API to create a new resource.
func (leo *LocationResourceOperator) InvokeCreate(req client.Location) (*client.Location, int, error) {
	// Fetch the location by name
	l, code, err := leo.Client.GetLocation(*req.Name)

	// If the location doesn't exist, then maybe the user attempted to create it, let them know that
	if code == 404 {
		return nil, 0, fmt.Errorf("location '/org/%s/location/%s' does not exist. Did you want to create a BYOK location? Please refer to the 'cpln_custom_location' resource. You can find more info here: https://registry.terraform.io/providers/controlplane-com/cpln/latest/docs/resources/custom_location", leo.Client.Org, *req.Name)
	}

	// If this is one of the custom locations, tell users to use the custom location resource
	if req.Provider != nil && IsCustomLocation(*req.Provider) {
		return nil, 0, fmt.Errorf("you are trying to create an existing '%s' location, please use the 'cpln_custom_location' resource and import using the 'terraform import' command. You can find more info here: https://registry.terraform.io/providers/controlplane-com/cpln/latest/docs/resources/custom_location", *req.Provider)
	}

	// If the location already exists, return it
	return l, code, err
}

// InvokeRead invokes the Get API to retrieve an existing resource by name.
func (leo *LocationResourceOperator) InvokeRead(name string) (*client.Location, int, error) {
	// Fetch the location by name
	l, code, err := leo.Client.GetLocation(name)

	// If this is one of the custom locations, tell users to use the custom location resource
	if l.Provider != nil && IsCustomLocation(*l.Provider) {
		return l, code, fmt.Errorf("the location '/org/%s/location/%s' is a '%s' location, please use the 'cpln_custom_location' resource and import using the 'terraform import' command. You can find more info here: https://registry.terraform.io/providers/controlplane-com/cpln/latest/docs/resources/custom_location#import-syntax", leo.Client.Org, name, *l.Provider)
	}

	// Return the result
	return l, code, err
}

// InvokeUpdate invokes the Update API to update an existing resource.
func (leo *LocationResourceOperator) InvokeUpdate(req client.Location) (*client.Location, int, error) {
	return leo.Client.UpdateLocation(req)
}

// InvokeDelete invokes the Delete API to delete a resource by name.
func (leo *LocationResourceOperator) InvokeDelete(name string) error {
	// Retrieve the location resource using client
	location, _, err := leo.Client.GetLocation(name)

	// Return error if retrieval fails
	if err != nil {
		return err
	}

	// Initialize spec if nil to avoid nil pointer
	if location.Spec == nil {
		// Create a new LocationSpec for the location
		location.Spec = &client.LocationSpec{}
	}

	// Set the enabled flag to default true
	location.Spec.Enabled = BuildBool(types.BoolValue(true))

	// Clear the status to indicate deletion
	location.Status = nil

	// Set the TagsReplace just so we don't get a panic
	location.TagsReplace = location.Tags
	location.Tags = nil

	// Update the location resource on the server
	_, _, err = leo.Client.UpdateLocationToDefault(*location)

	// Return error if update fails
	if err != nil {
		return err
	}

	// Return nil on success
	return nil
}

// Flatteners //

// flattenGeo transforms *client.LocationGeo into a Terraform types.List.
func (leo *LocationResourceOperator) flattenGeo(input *client.LocationGeo) types.List {
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

/*** Helpers ***/

// IsCustomLocation checks if the given provider is a custom location provider.
func IsCustomLocation(provider string) bool {
	return slices.Contains(AllowedCustomLocationProviders, provider)
}
