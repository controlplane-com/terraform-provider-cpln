package cpln

import (
	"context"
	"encoding/json"
	"fmt"
	client "github.com/controlplane-com/terraform-provider-cpln/internal/provider/client"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure resource defined types fully satisfy framework interfaces
var (
	_ resource.Resource                = &AgentResource{}
	_ resource.ResourceWithImportState = &AgentResource{}
)

/*** Resource Model ***/

// AgentResourceModel holds the resource data structure for the Terraform state.
type AgentResourceModel struct {
	BaseResourceModel
	UserData types.String `tfsdk:"user_data"`
}

/*** Resource Configuration ***/

// AgentResource is the resource implementation.
type AgentResource struct {
	client *client.Client
}

// NewAgentResource is a helper function to simplify resource implementation.
func NewAgentResource() resource.Resource {
	return &AgentResource{}
}

// ImportState sets up the import operation to map the imported ID to the "id" attribute in the state.
func (r *AgentResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// Metadata provides the resource type name.
func (r *AgentResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = "cpln_agent"
}

// Schema defines the schema for the resource.
func (r *AgentResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: MergeAttributes(BaseResourceAttributes("Agent"), map[string]schema.Attribute{
			"user_data": schema.StringAttribute{
				Computed:    true,
				Sensitive:   true,
				Description: "The JSON output needed when creating an agent.",
			},
		}),
	}
}

// Configure assigns the provider's client to the resource for API interactions.
func (r *AgentResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
func (r *AgentResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plannedState AgentResourceModel

	// Retrieve the planned state from the Terraform configuration
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plannedState)...)

	// Check if any errors have been added to the diagnostics in previous operations
	// If there are errors, stop further execution to prevent inconsistent or partial changes
	if resp.Diagnostics.HasError() {
		return
	}

	// Initialize a new request payload structure and populate it with the input state
	requestPayload := client.Agent{}
	UpdateBaseClientFromState(&requestPayload.Base, plannedState.BaseResourceModel)

	// Send the create request to the API client
	responsePayload, code, err := r.client.CreateAgent(requestPayload)

	// Handle cases where the resource already exists, returning an error to inform the user
	if code == 409 {
		resp.Diagnostics.AddError("Resource already exists", "The requestPayload resource already exists. You can import the existing resource using the Terraform import command.")
		return
	}

	// Handle any other errors that occurred during the API request
	if err != nil {
		resp.Diagnostics.AddError("API Error", fmt.Sprintf("Error creating requestPayload: %s", err))
		return
	}

	// Marshal the BootstrapConfig from the response to JSON for the user_data attribute
	userData, err := json.Marshal(responsePayload.Status.BootstrapConfig)
	if err != nil {
		resp.Diagnostics.AddError("JSON Marshal Error", fmt.Sprintf("Error marshalling user_data: %s", err))
		return
	}

	// Map the API response to the Terraform state
	finalState := AgentResourceModel{
		UserData: types.StringValue(string(userData)),
	}

	// Map shared attributes from the API response's Base object to the Terraform state model
	UpdateStateFromBaseClient(&finalState.BaseResourceModel, responsePayload.Base)

	// Set the resource state in Terraform
	resp.Diagnostics.Append(resp.State.Set(ctx, &finalState)...)
}

// Read fetches the current state of the resource.
func (r *AgentResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var plannedState AgentResourceModel

	// Retrieve the planned state from the Terraform configuration
	resp.Diagnostics.Append(req.State.Get(ctx, &plannedState)...)

	// Check if any errors have been added to the diagnostics in previous operations
	// If there are errors, stop further execution to prevent inconsistent or partial changes
	if resp.Diagnostics.HasError() {
		return
	}

	// Call the API to get the latest state for the resource
	responsePayload, code, err := r.client.GetAgent(plannedState.ID.ValueString())

	// Handle the case where the resource is not found (HTTP 404),
	// indicating it has been deleted outside of Terraform. Remove it from state
	if code == 404 {
		resp.State.RemoveResource(ctx)
		return
	}

	// Handle any other errors that occur during the API call
	if err != nil {
		resp.Diagnostics.AddError("API Error", fmt.Sprintf("Error reading agent: %s", err))
		return
	}

	// Map the API response to the Terraform state
	finalState := AgentResourceModel{}

	// Map shared attributes from the API response's Base object to the Terraform state model
	UpdateStateFromBaseClient(&finalState.BaseResourceModel, responsePayload.Base)

	// Set the updated state in Terraform
	resp.Diagnostics.Append(resp.State.Set(ctx, &finalState)...)
}

// Update modifies the resource.
func (r *AgentResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plannedState AgentResourceModel

	// Retrieve the planned state from the Terraform configuration
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plannedState)...)

	// Check if any errors have been added to the diagnostics in previous operations
	// If there are errors, stop further execution to prevent inconsistent or partial changes
	if resp.Diagnostics.HasError() {
		return
	}

	// Initialize a new request payload structure and populate it with the input state
	requestPayload := client.Agent{}
	UpdateReplaceBaseClientFromState(&requestPayload.Base, plannedState.BaseResourceModel)

	// Send the update request to the API with the modified data
	responsePayload, _, err := r.client.UpdateAgent(requestPayload)

	// Handle errors from the API update request
	if err != nil {
		resp.Diagnostics.AddError("API Error", fmt.Sprintf("Error updating agent: %s", err))
		return
	}

	// Map the API response to the Terraform state
	state := AgentResourceModel{
		UserData: types.StringNull(),
	}

	// Map shared attributes from the API response's Base object to the Terraform state model
	UpdateStateFromBaseClient(&state.BaseResourceModel, responsePayload.Base)

	// Set the updated state in Terraform, appending any diagnostics from this operation
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// Delete removes the resource.
func (r *AgentResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state AgentResourceModel

	// Retrieve the state from the Terraform configuration
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	// Check if any errors have been added to the diagnostics in previous operations
	// If there are errors, stop further execution to prevent inconsistent or partial changes
	if resp.Diagnostics.HasError() {
		return
	}

	// Send a delete request to the API using the name from the state
	err := r.client.DeleteAgent(state.Name.ValueString())

	// Handle errors from the API delete request
	if err != nil {
		// If an error occurs during the delete request, add an error to diagnostics
		resp.Diagnostics.AddError("API Error", fmt.Sprintf("Error deleting agent: %s", err))
		return
	}

	// Remove the resource from Terraform's state, indicating successful deletion
	resp.State.RemoveResource(ctx)
}
