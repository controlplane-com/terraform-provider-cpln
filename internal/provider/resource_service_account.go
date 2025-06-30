package cpln

import (
	"context"

	client "github.com/controlplane-com/terraform-provider-cpln/internal/provider/client"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure resource implements required interfaces.
var (
	_ resource.Resource                = &ServiceAccountResource{}
	_ resource.ResourceWithImportState = &ServiceAccountResource{}
)

/*** Resource Model ***/

// ServiceAccountResourceModel holds the Terraform state for the resource.
type ServiceAccountResourceModel struct {
	EntityBaseModel
	Origin types.String `tfsdk:"origin"`
}

/*** Resource Configuration ***/

// ServiceAccountResource is the resource implementation.
type ServiceAccountResource struct {
	EntityBase
	Operations EntityOperations[ServiceAccountResourceModel, client.ServiceAccount]
}

// NewServiceAccountResource returns a new instance of the resource implementation.
func NewServiceAccountResource() resource.Resource {
	return &ServiceAccountResource{}
}

// Configure configures the resource before use.
func (sar *ServiceAccountResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	sar.EntityBaseConfigure(ctx, req.ProviderData, &resp.Diagnostics)
	sar.Operations = NewEntityOperations(sar.client, &ServiceAccountResourceOperator{})
}

// ImportState sets up the import operation to map the imported ID to the "id" attribute in the state.
func (sar *ServiceAccountResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// Metadata provides the resource type name.
func (sar *ServiceAccountResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = "cpln_service_account"
}

// Schema defines the schema for the resource.
func (sar *ServiceAccountResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: MergeAttributes(sar.EntityBaseAttributes("service account"), map[string]schema.Attribute{
			"origin": schema.StringAttribute{
				Description: "Origin of the Policy. Either `builtin` or `default`.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
		}),
	}
}

// Create creates the resource.
func (sar *ServiceAccountResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	CreateGeneric(ctx, req, resp, sar.Operations)
}

// Read fetches the current state of the resource.
func (sar *ServiceAccountResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	ReadGeneric(ctx, req, resp, sar.Operations)
}

// Update modifies the resource.
func (sar *ServiceAccountResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	UpdateGeneric(ctx, req, resp, sar.Operations)
}

// Delete removes the resource.
func (sar *ServiceAccountResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	DeleteGeneric(ctx, req, resp, sar.Operations)
}

/*** Resource Operator ***/

// ServiceAccountResourceOperator is the operator for managing the state.
type ServiceAccountResourceOperator struct {
	EntityOperator[ServiceAccountResourceModel]
}

// NewAPIRequest creates a request payload from a state model.
func (saro *ServiceAccountResourceOperator) NewAPIRequest(isUpdate bool) client.ServiceAccount {
	// Initialize a new request payload
	requestPayload := client.ServiceAccount{}

	// Populate Base fields from state
	saro.Plan.Fill(&requestPayload.Base, isUpdate)

	// Return constructed request payload
	return requestPayload
}

// MapResponseToState constructs the Terraform state model from the API response payload.
func (saro *ServiceAccountResourceOperator) MapResponseToState(apiResp *client.ServiceAccount, isCreate bool) ServiceAccountResourceModel {
	// Initialize empty state model
	state := ServiceAccountResourceModel{}

	// Populate common fields from base resource data
	state.From(apiResp.Base)

	// Set specific attributes
	state.Origin = types.StringPointerValue(apiResp.Origin)

	// Return completed state model
	return state
}

// InvokeCreate invokes the Create API to create a new resource.
func (saro *ServiceAccountResourceOperator) InvokeCreate(req client.ServiceAccount) (*client.ServiceAccount, int, error) {
	return saro.Client.CreateServiceAccount(req)
}

// InvokeRead invokes the Get API to retrieve an existing resource by name.
func (saro *ServiceAccountResourceOperator) InvokeRead(name string) (*client.ServiceAccount, int, error) {
	return saro.Client.GetServiceAccount(name)
}

// InvokeUpdate invokes the Update API to update an existing resource.
func (saro *ServiceAccountResourceOperator) InvokeUpdate(req client.ServiceAccount) (*client.ServiceAccount, int, error) {
	return saro.Client.UpdateServiceAccount(req)
}

// InvokeDelete invokes the Delete API to delete a resource by name.
func (saro *ServiceAccountResourceOperator) InvokeDelete(name string) error {
	return saro.Client.DeleteServiceAccount(name)
}
