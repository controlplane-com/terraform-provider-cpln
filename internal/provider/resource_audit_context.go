package cpln

import (
	"context"
	"fmt"
	client "github.com/controlplane-com/terraform-provider-cpln/internal/provider/client"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

// Ensure resource defined types fully satisfy framework interfaces.
var (
	_ resource.Resource                = &AuditContextResource{}
	_ resource.ResourceWithImportState = &AuditContextResource{}
)

/*** Resource Model ***/

// AuditContextResourceModel holds the resource data structure for the Terraform state.
type AuditContextResourceModel struct {
	BaseResourceModel
}

/*** Resource Configuration ***/

// AuditContextResource is the resource implementation.
type AuditContextResource struct {
	client *client.Client
}

// NewAuditContextResource is a helper function to simplify resource implementation.
func NewAuditContextResource() resource.Resource {
	return &AuditContextResource{}
}

// ImportState sets up the import operation to map the imported ID to the "id" attribute in the state.
func (r *AuditContextResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// Metadata provides the resource type name.
func (r *AuditContextResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = "cpln_audit_context"
}

// Schema defines the schema for the resource.
func (r *AuditContextResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: BaseResourceAttributes("Audit Context"),
	}
}

// Configure assigns the provider's client to the resource for API interactions.
func (r *AuditContextResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Ensure ProviderData is present; return early if the provider has not been properly configured
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
	r.client = c
}

// Create sets up the resource's Create operation.
func (r *AuditContextResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plannedState AuditContextResourceModel

	// Retrieve the planned state from the Terraform configuration
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plannedState)...)

	// Check if any errors have been added to the diagnostics in previous operations
	// If there are errors, stop further execution to prevent inconsistent or partial changes
	if resp.Diagnostics.HasError() {
		return
	}

	// Initialize a new request payload structure and populate it with the input state
	requestPayload := client.AuditContext{}
	UpdateBaseClientFromState(&requestPayload.Base, plannedState.BaseResourceModel)

	// Send the create request to the API client
	responsePayload, code, err := r.client.CreateAuditContext(requestPayload)

	// Handle cases where the resource already exists, returning an error to inform the user
	if code == 409 {
		resp.Diagnostics.AddError("Resource already exists", "The audit context resource already exists. You can import the existing resource using the Terraform import command.")
		return
	}

	// Handle any other errors that occurred during the API request
	if err != nil {
		resp.Diagnostics.AddError("API Error", fmt.Sprintf("Error creating audit context: %s", err))
		return
	}

	// Map the API response to the Terraform state
	finalState := AuditContextResourceModel{}

	// Map shared attributes from the API response's Base object to the Terraform state model
	UpdateStateFromBaseClient(&finalState.BaseResourceModel, responsePayload.Base)

	// Set the resource state in Terraform
	resp.Diagnostics.Append(resp.State.Set(ctx, &finalState)...)
}

// Read fetches the current state of the resource.
func (r *AuditContextResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var plannedState AuditContextResourceModel

	// Retrieve the planned state from the Terraform configuration
	resp.Diagnostics.Append(req.State.Get(ctx, &plannedState)...)

	// Check if any errors have been added to the diagnostics in previous operations
	// If there are errors, stop further execution to prevent inconsistent or partial changes
	if resp.Diagnostics.HasError() {
		return
	}

	// Call the API to get the latest state for the resource
	responsePayload, code, err := r.client.GetAuditContext(plannedState.ID.ValueString())

	// Handle the case where the resource is not found (HTTP 404),
	// indicating it has been deleted outside of Terraform. Remove it from state
	if code == 404 {
		resp.State.RemoveResource(ctx)
		return
	}

	// Handle any other errors that occur during the API call
	if err != nil {
		resp.Diagnostics.AddError("API Error", fmt.Sprintf("Error reading audit context: %s", err))
		return
	}

	// Map the API response to the Terraform state
	state := AuditContextResourceModel{}

	// Map shared attributes from the API response's Base object to the Terraform state model
	UpdateStateFromBaseClient(&state.BaseResourceModel, responsePayload.Base)

	// Set the updated state in Terraform
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// Update modifies the resource.
func (r *AuditContextResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plannedState AuditContextResourceModel

	// Retrieve the planned state from the Terraform configuration
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plannedState)...)

	// Check if any errors have been added to the diagnostics in previous operations
	// If there are errors, stop further execution to prevent inconsistent or partial changes
	if resp.Diagnostics.HasError() {
		return
	}

	// Initialize a new request payload structure and populate it with the input state
	requestPayload := client.AuditContext{}
	UpdateReplaceBaseClientFromState(&requestPayload.Base, plannedState.BaseResourceModel)

	// Send the update request to the API with the modified data
	responsePayload, _, err := r.client.UpdateAuditContext(requestPayload)

	// Handle errors from the API update request
	if err != nil {
		resp.Diagnostics.AddError("API Error", fmt.Sprintf("Error updating audit context: %s", err))
		return
	}

	// Map the API response to the Terraform state
	finalState := AuditContextResourceModel{}

	// Map shared attributes from the API response's Base object to the Terraform state model
	UpdateStateFromBaseClient(&finalState.BaseResourceModel, responsePayload.Base)

	// Set the updated state in Terraform, appending any diagnostics from this operation
	resp.Diagnostics.Append(resp.State.Set(ctx, &finalState)...)
}

// Delete removes the resource.
func (r *AuditContextResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Remove the resource from Terraform's state, indicating successful deletion
	resp.State.RemoveResource(ctx)
}
