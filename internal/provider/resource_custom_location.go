package cpln

import (
	"context"

	client "github.com/controlplane-com/terraform-provider-cpln/internal/provider/client"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure resource implements required interfaces.
var (
	_ resource.Resource                = &CustomLocationResource{}
	_ resource.ResourceWithImportState = &CustomLocationResource{}
)

/*** Resource Model ***/

// CustomLocationResourceModel holds the Terraform state for the resource.
type CustomLocationResourceModel struct {
	EntityBaseModel
	Origin        types.String `tfsdk:"origin"`
	CloudProvider types.String `tfsdk:"cloud_provider"`
	Region        types.String `tfsdk:"region"`
	Enabled       types.Bool   `tfsdk:"enabled"`
}

/*** Resource Configuration ***/

// CustomLocationResource is the resource implementation.
type CustomLocationResource struct {
	EntityBase
	Operations EntityOperations[CustomLocationResourceModel, client.Location]
}

// NewCustomLocationResource returns a new instance of the resource implementation.
func NewCustomLocationResource() resource.Resource {
	return &CustomLocationResource{}
}

// Configure configures the resource before use.
func (clr *CustomLocationResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	clr.EntityBaseConfigure(ctx, req.ProviderData, &resp.Diagnostics)
	clr.Operations = NewEntityOperations(clr.client, &CustomLocationResourceOperator{})
}

// ImportState sets up the import operation to map the imported ID to the "id" attribute in the state.
func (clr *CustomLocationResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// Metadata provides the resource type name.
func (clr *CustomLocationResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = "cpln_custom_location"
}

// Schema defines the schema for the resource.
func (clr *CustomLocationResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: MergeAttributes(clr.EntityBaseAttributes("Custom Location"), map[string]schema.Attribute{
			"origin": schema.StringAttribute{
				Description: "",
				Computed:    true,
			},
			"cloud_provider": schema.StringAttribute{
				Description: "Cloud Provider of the custom location.",
				Required:    true,
				Validators:  []validator.String{stringvalidator.OneOf(CustomLocationCloudProviders...)},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
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
				Description: "Indication if the custom location is enabled.",
				Required:    true,
			},
		}),
	}
}

// Create creates the resource.
func (clr *CustomLocationResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	CreateGeneric(ctx, req, resp, clr.Operations)
}

// Read fetches the current state of the resource.
func (clr *CustomLocationResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	ReadGeneric(ctx, req, resp, clr.Operations)
}

// Update modifies the resource.
func (clr *CustomLocationResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	UpdateGeneric(ctx, req, resp, clr.Operations)
}

// Delete removes the resource.
func (clr *CustomLocationResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	DeleteGeneric(ctx, req, resp, clr.Operations)
}

/*** Resource Operator ***/

// CustomLocationResourceOperator is the operator for managing the state.
type CustomLocationResourceOperator struct {
	EntityOperator[CustomLocationResourceModel]
}

// NewAPIRequest creates a request payload from a state model.
func (clro *CustomLocationResourceOperator) NewAPIRequest(isUpdate bool) client.Location {
	// Initialize a new request payload
	requestPayload := client.Location{}

	// Populate Base fields from state
	clro.Plan.Fill(&requestPayload.Base, isUpdate)

	// Set specific attributes
	requestPayload.Provider = BuildString(clro.Plan.CloudProvider)
	requestPayload.Spec = &client.LocationSpec{
		Enabled: clro.Plan.Enabled.ValueBoolPointer(),
	}

	// Return the request payload object
	return requestPayload
}

// MapResponseToState creates a state model from response payload.
func (clro *CustomLocationResourceOperator) MapResponseToState(location *client.Location, isCreate bool) CustomLocationResourceModel {
	// Initialize empty state model
	state := CustomLocationResourceModel{}

	// Populate common fields from base resource data
	state.From(location.Base)

	// Set specific attributes
	state.Origin = types.StringPointerValue(location.Origin)
	state.CloudProvider = types.StringPointerValue(location.Provider)
	state.Region = types.StringPointerValue(location.Region)
	state.Enabled = types.BoolPointerValue(location.Spec.Enabled) // Spec is always defined, no need to check for nil

	// Return the built state
	return state
}

// InvokeCreate invokes the Create API to create a new resource.
func (clro *CustomLocationResourceOperator) InvokeCreate(req client.Location) (*client.Location, int, error) {
	return clro.Client.CreateCustomLocation(req)
}

// InvokeRead invokes the Get API to retrieve an existing resource by name.
func (clro *CustomLocationResourceOperator) InvokeRead(name string) (*client.Location, int, error) {
	return clro.Client.GetLocation(name)
}

// InvokeUpdate invokes the Update API to update an existing resource.
func (clro *CustomLocationResourceOperator) InvokeUpdate(req client.Location) (*client.Location, int, error) {
	return clro.Client.UpdateLocation(req)
}

// InvokeDelete invokes the Delete API to delete a resource by name.
func (clro *CustomLocationResourceOperator) InvokeDelete(name string) error {
	return clro.Client.DeleteCustomLocation(name)
}
