package cpln

import (
	"context"
	"encoding/json"
	"fmt"

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
	_ resource.Resource                = &AgentResource{}
	_ resource.ResourceWithImportState = &AgentResource{}
)

/*** Resource Model ***/

// AgentResourceModel holds the Terraform state for the resource.
type AgentResourceModel struct {
	EntityBaseModel
	UserData types.String `tfsdk:"user_data"`
}

/*** Resource Configuration ***/

// AgentResource is the resource implementation.
type AgentResource struct {
	EntityBase
	Operations EntityOperations[AgentResourceModel, client.Agent]
}

// NewAgentResource returns a new instance of the resource implementation.
func NewAgentResource() resource.Resource {
	return &AgentResource{}
}

// Configure configures the resource before use.
func (ar *AgentResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	ar.EntityBaseConfigure(ctx, req.ProviderData, &resp.Diagnostics)
	ar.Operations = NewEntityOperations(ar.client, &AgentResourceOperator{})
}

// ImportState sets up the import operation to map the imported ID to the "id" attribute in the state.
func (ar *AgentResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// Metadata provides the resource type name.
func (ar *AgentResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = "cpln_agent"
}

// Schema defines the schema for the resource.
func (ar *AgentResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: MergeAttributes(ar.EntityBaseAttributes("Agent"), map[string]schema.Attribute{
			"user_data": schema.StringAttribute{
				Description: "The JSON output needed when creating an agent.",
				Computed:    true,
				Sensitive:   true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
		}),
	}
}

// Create creates the resource.
func (ar *AgentResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	CreateGeneric(ctx, req, resp, ar.Operations)
}

// Read fetches the current state of the resource.
func (ar *AgentResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	ReadGeneric(ctx, req, resp, ar.Operations)
}

// Update modifies the resource.
func (ar *AgentResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	UpdateGeneric(ctx, req, resp, ar.Operations)
}

// Delete removes the resource.
func (ar *AgentResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	DeleteGeneric(ctx, req, resp, ar.Operations)
}

/*** Resource Operator ***/

// AgentResourceOperator is the operator for managing the state.
type AgentResourceOperator struct {
	EntityOperator[AgentResourceModel]
}

// NewAPIRequest creates a request payload from a state model.
func (aro *AgentResourceOperator) NewAPIRequest(isUpdate bool) client.Agent {
	// Initialize a new request payload
	requestPayload := client.Agent{}

	// Populate Base fields from state
	aro.Plan.Fill(&requestPayload.Base, isUpdate)

	// Return constructed request payload
	return requestPayload
}

// MapResponseToState constructs the Terraform state model from the API response payload.
func (aro *AgentResourceOperator) MapResponseToState(agent *client.Agent, isCreate bool) AgentResourceModel {
	// Initialize empty state model
	state := AgentResourceModel{}

	// On create operation, include bootstrap config as user_data
	if isCreate && agent.Status != nil && agent.Status.BootstrapConfig != nil {
		// Convert BootstrapConfig to JSON string
		userData, err := json.Marshal(agent.Status.BootstrapConfig)
		if err != nil {
			// Report JSON marshalling error
			aro.Diags.AddError("JSON Marshal Error", fmt.Sprintf("Error marshalling user_data: %s", err))
			return state
		}

		// Set user_data attribute in state
		state.UserData = types.StringValue(string(userData))
	} else {
		state.UserData = aro.Plan.UserData // :D
	}

	// Populate common fields from base resource data
	state.From(agent.Base)

	// Return completed state model
	return state
}

// InvokeCreate invokes the Create API to create a new resource.
func (aro *AgentResourceOperator) InvokeCreate(req client.Agent) (*client.Agent, int, error) {
	return aro.Client.CreateAgent(req)
}

// InvokeRead invokes the Get API to retrieve an existing resource by name.
func (aro *AgentResourceOperator) InvokeRead(name string) (*client.Agent, int, error) {
	return aro.Client.GetAgent(name)
}

// InvokeUpdate invokes the Update API to update an existing resource.
func (aro *AgentResourceOperator) InvokeUpdate(req client.Agent) (*client.Agent, int, error) {
	return aro.Client.UpdateAgent(req)
}

// InvokeDelete invokes the Delete API to delete a resource by name.
func (aro *AgentResourceOperator) InvokeDelete(name string) error {
	return aro.Client.DeleteAgent(name)
}
