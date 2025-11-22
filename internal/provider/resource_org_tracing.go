package cpln

import (
	"context"
	"fmt"

	client "github.com/controlplane-com/terraform-provider-cpln/internal/provider/client"
	commonmodels "github.com/controlplane-com/terraform-provider-cpln/internal/provider/models/common"
	"github.com/hashicorp/terraform-plugin-framework-validators/resourcevalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure resource implements required interfaces.
var (
	_ resource.Resource                = &OrgTracingResource{}
	_ resource.ResourceWithImportState = &OrgTracingResource{}
)

/*** Resource Model ***/

// OrgTracingResourceModel holds the Terraform state for the resource.
type OrgTracingResourceModel struct {
	ID                  types.String `tfsdk:"id"`
	CplnID              types.String `tfsdk:"cpln_id"`
	Name                types.String `tfsdk:"name"`
	Description         types.String `tfsdk:"description"`
	Tags                types.Map    `tfsdk:"tags"`
	SelfLink            types.String `tfsdk:"self_link"`
	LightstepTracing    types.List   `tfsdk:"lightstep_tracing"`
	OtelTracing         types.List   `tfsdk:"otel_tracing"`
	ControlPlaneTracing types.List   `tfsdk:"controlplane_tracing"`
}

/*** Resource Configuration ***/

// OrgTracingResource is the resource implementation.
type OrgTracingResource struct {
	EntityBase
}

// NewOrgTracingResource returns a new instance of the resource implementation.
func NewOrgTracingResource() resource.Resource {
	return &OrgTracingResource{}
}

// Configure configures the resource before use.
func (otr *OrgTracingResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	otr.EntityBaseConfigure(ctx, req.ProviderData, &resp.Diagnostics)
}

// ImportState sets up the import operation to map the imported ID to the "id" attribute in the state.
func (otr *OrgTracingResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// Metadata provides the resource type name.
func (otr *OrgTracingResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = "cpln_org_tracing"
}

// Schema defines the schema for the resource.
func (otr *OrgTracingResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The unique identifier for this Org Tracing.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"cpln_id": schema.StringAttribute{
				Description: "The ID, in GUID format, of the Org.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Description: "Name of the Org.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"description": schema.StringAttribute{
				Description: "Description of the Org.",
				Computed:    true,
			},
			"tags": schema.MapAttribute{
				Description: "Key-value map of resource tags.",
				ElementType: types.StringType,
				Computed:    true,
			},
			"self_link": schema.StringAttribute{
				Description: "Full link to this resource. Can be referenced by other resources.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
		},
		Blocks: map[string]schema.Block{
			"lightstep_tracing":    otr.LightstepTracingSchema(),
			"otel_tracing":         otr.OtelTracingSchema(),
			"controlplane_tracing": otr.ControlPlaneTracingSchema(),
		},
	}
}

// ConfigValidators enforces mutual exclusivity between attributes.
func (otr *OrgTracingResource) ConfigValidators(ctx context.Context) []resource.ConfigValidator {
	expressions := []path.Expression{
		path.MatchRoot("lightstep_tracing"),
		path.MatchRoot("otel_tracing"),
		path.MatchRoot("controlplane_tracing"),
	}

	return []resource.ConfigValidator{
		resourcevalidator.ExactlyOneOf(expressions...),
	}
}

// Create creates the resource.
func (otr *OrgTracingResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Acquire lock to ensure only one operation modifies the resource at a time
	orgOperationLock.Lock()
	defer orgOperationLock.Unlock()

	// Declare variable to hold the planned state from Terraform configuration
	var plannedState OrgTracingResourceModel

	// Retrieve the planned state from the Terraform configuration
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plannedState)...)

	// Abort on errors to avoid partial or inconsistent state
	if resp.Diagnostics.HasError() {
		return
	}

	// Initialize a new request payload structure and populate it with the planned state
	tracing := otr.buildRequest(ctx, &resp.Diagnostics, plannedState)

	// Return if an error has occurred during the request payload creation
	if resp.Diagnostics.HasError() {
		return
	}

	// Send the create request to the API client
	responsePayload, _, err := otr.client.UpdateOrgTracing(tracing)

	// Handle any other errors that occurred during the API request
	if err != nil {
		resp.Diagnostics.AddError("API Error", fmt.Sprintf("Error creating org tracing: %s", err))
		return
	}

	// Map the API response to the Terraform state
	finalState := otr.buildState(ctx, &resp.Diagnostics, responsePayload)

	// Return if an error has occurred during the state creation
	if resp.Diagnostics.HasError() {
		return
	}

	// Set the resource state in Terraform
	resp.Diagnostics.Append(resp.State.Set(ctx, &finalState)...)
}

// Read fetches the current state of the resource.
func (otr *OrgTracingResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var plannedState OrgTracingResourceModel

	// Retrieve the planned state from the Terraform configuration
	resp.Diagnostics.Append(req.State.Get(ctx, &plannedState)...)

	// Abort on errors to avoid partial or inconsistent state
	if resp.Diagnostics.HasError() {
		return
	}

	// Fetch the org
	responsePayload, code, err := otr.client.GetOrg()

	// Handle the case where the org is not found (HTTP 404),
	// indicating it has been deleted outside of Terraform. Remove it from state
	if code == 404 {
		resp.State.RemoveResource(ctx)
		return
	}

	// Handle any other errors that occur during the API call
	if err != nil {
		resp.Diagnostics.AddError("API Error", fmt.Sprintf("Error reading org tracing: %s", err))
		return
	}

	// Map the API response to the Terraform state
	finalState := otr.buildState(ctx, &resp.Diagnostics, responsePayload)

	// Return if an error has occurred during the state creation
	if resp.Diagnostics.HasError() {
		return
	}

	// Set the updated state in Terraform
	resp.Diagnostics.Append(resp.State.Set(ctx, &finalState)...)
}

// Update modifies the resource.
func (otr *OrgTracingResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Acquire lock to ensure only one operation modifies the resource at a time
	orgOperationLock.Lock()
	defer orgOperationLock.Unlock()

	var plannedState OrgTracingResourceModel

	// Retrieve the planned state from the Terraform configuration
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plannedState)...)

	// Abort on errors to avoid partial or inconsistent state
	if resp.Diagnostics.HasError() {
		return
	}

	// Initialize a new request payload structure and populate it with the planned state
	tracing := otr.buildRequest(ctx, &resp.Diagnostics, plannedState)

	// Return if an error has occurred during the request payload creation
	if resp.Diagnostics.HasError() {
		return
	}

	// Send the update request to the API with the modified data
	responsePayload, _, err := otr.client.UpdateOrgTracing(tracing)

	// Handle errors from the API update request
	if err != nil {
		resp.Diagnostics.AddError("API Error", fmt.Sprintf("Error updating org tracing: %s", err))
		return
	}

	// Map the API response to the Terraform finalState
	finalState := otr.buildState(ctx, &resp.Diagnostics, responsePayload)

	// Return if an error has occurred during the state creation
	if resp.Diagnostics.HasError() {
		return
	}

	// Set the updated state in Terraform
	resp.Diagnostics.Append(resp.State.Set(ctx, &finalState)...)
}

// Delete removes the resource.
func (otr *OrgTracingResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Acquire lock to ensure only one operation modifies the resource at a time
	orgOperationLock.Lock()
	defer orgOperationLock.Unlock()

	var state OrgTracingResourceModel

	// Retrieve the state from the Terraform configuration
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	// Abort on errors to avoid partial or inconsistent state
	if resp.Diagnostics.HasError() {
		return
	}

	// Send a delete request to the API using the name from the state
	_, _, err := otr.client.UpdateOrgTracing(nil)

	// Handle errors from the API delete request
	if err != nil {
		// If an error occurs during the delete request, add an error to diagnostics
		resp.Diagnostics.AddError("API Error", fmt.Sprintf("Error deleting org tracing: %s", err))
		return
	}

	// Remove the resource from Terraform's state, indicating successful deletion
	resp.State.RemoveResource(ctx)
}

/*** Operations ***/

// buildRequest creates a request payload from a state model.
func (otr *OrgTracingResource) buildRequest(ctx context.Context, diags *diag.Diagnostics, state OrgTracingResourceModel) *client.Tracing {
	return BuildTracing(ctx, diags, state.LightstepTracing, state.OtelTracing, state.ControlPlaneTracing)
}

// buildState creates a state model from response payload.
func (otr *OrgTracingResource) buildState(ctx context.Context, diags *diag.Diagnostics, apiResp *client.Org) OrgTracingResourceModel {
	// Initialize empty state model
	state := OrgTracingResourceModel{}

	// Set specific attributes
	state.ID = types.StringPointerValue(apiResp.Name)
	state.CplnID = types.StringPointerValue(apiResp.ID)
	state.Name = types.StringPointerValue(apiResp.Name)
	state.Description = types.StringPointerValue(apiResp.Description)
	state.Tags = FlattenTags(apiResp.Tags)
	state.SelfLink = FlattenSelfLink(apiResp.Links)

	// Only process tracing if Spec is non-nil
	if apiResp.Spec != nil {
		// Extract tracing configurations from spec
		lightstepTracing, otelTracing, cplnTracing := FlattenTracing(ctx, diags, apiResp.Spec.Tracing)

		// Set specific attributes
		state.LightstepTracing = lightstepTracing
		state.OtelTracing = otelTracing
		state.ControlPlaneTracing = cplnTracing
	} else {
		state.LightstepTracing = types.ListNull(commonmodels.LightstepTracingModel{}.AttributeTypes())
		state.OtelTracing = types.ListNull(commonmodels.OtelTracingModel{}.AttributeTypes())
		state.ControlPlaneTracing = types.ListNull(commonmodels.ControlPlaneTracingModel{}.AttributeTypes())
	}

	// Return completed state model
	return state
}
