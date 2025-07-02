package cpln

import (
	"context"

	client "github.com/controlplane-com/terraform-provider-cpln/internal/provider/client"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

// Ensure resource implements required interfaces.
var (
	_ resource.Resource                = &AuditContextResource{}
	_ resource.ResourceWithImportState = &AuditContextResource{}
)

/*** Resource Model ***/

// AuditContextResourceModel holds the Terraform state for the resource.
type AuditContextResourceModel struct {
	EntityBaseModel
}

/*** Resource Configuration ***/

// AuditContextResource is the resource implementation.
type AuditContextResource struct {
	EntityBase
	Operations EntityOperations[AuditContextResourceModel, client.AuditContext]
}

// NewAuditContextResource returns a new instance of the resource implementation.
func NewAuditContextResource() resource.Resource {
	return &AuditContextResource{}
}

// Configure configures the resource before use.
func (acr *AuditContextResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	acr.EntityBaseConfigure(ctx, req.ProviderData, &resp.Diagnostics)
	acr.Operations = NewEntityOperations(acr.client, &AuditContextResourceOperator{})
}

// ImportState sets up the import operation to map the imported ID to the "id" attribute in the state.
func (acr *AuditContextResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// Metadata provides the resource type name.
func (acr *AuditContextResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = "cpln_audit_context"
}

// Schema defines the schema for the resource.
func (acr *AuditContextResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: acr.EntityBaseAttributes("Audit Context"),
	}
}

// Create creates the resource.
func (acr *AuditContextResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	CreateGeneric(ctx, req, resp, acr.Operations)
}

// Read fetches the current state of the resource.
func (acr *AuditContextResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	ReadGeneric(ctx, req, resp, acr.Operations)
}

// Update modifies the resource.
func (acr *AuditContextResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	UpdateGeneric(ctx, req, resp, acr.Operations)
}

// Delete removes the resource.
func (acr *AuditContextResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Audit contexts are immutable once created; on Delete just remove from state
	resp.State.RemoveResource(ctx)
}

/*** Resource Operator ***/

// AuditContextResourceOperator is the operator for managing the state.
type AuditContextResourceOperator struct {
	EntityOperator[AuditContextResourceModel]
}

// NewAPIRequest creates a request payload from a state model.
func (acro *AuditContextResourceOperator) NewAPIRequest(isUpdate bool) client.AuditContext {
	// Initialize a new request payload
	requestPayload := client.AuditContext{}

	// Populate Base fields from state
	acro.Plan.Fill(&requestPayload.Base, isUpdate)

	// Return constructed request payload
	return requestPayload
}

// MapResponseToState constructs the Terraform state model from the API response payload.
func (acro *AuditContextResourceOperator) MapResponseToState(auditctx *client.AuditContext, isCreate bool) AuditContextResourceModel {
	// Initialize empty state model
	state := AuditContextResourceModel{}

	// Populate common fields from base resource data
	state.From(auditctx.Base)

	// Return completed state model
	return state
}

// InvokeCreate invokes the Create API to create a new resource.
func (acro *AuditContextResourceOperator) InvokeCreate(req client.AuditContext) (*client.AuditContext, int, error) {
	return acro.Client.CreateAuditContext(req)
}

// InvokeRead invokes the Get API to retrieve an existing resource by name.
func (acro *AuditContextResourceOperator) InvokeRead(name string) (*client.AuditContext, int, error) {
	return acro.Client.GetAuditContext(name)
}

// InvokeUpdate invokes the Update API to update an existing resource.
func (acro *AuditContextResourceOperator) InvokeUpdate(req client.AuditContext) (*client.AuditContext, int, error) {
	return acro.Client.UpdateAuditContext(req)
}

// InvokeDelete invokes the Delete API to delete a resource by name.
func (acro *AuditContextResourceOperator) InvokeDelete(name string) error {
	return nil
}
