package cpln

import (
	"context"
	"fmt"
	"strings"

	client "github.com/controlplane-com/terraform-provider-cpln/internal/provider/client"
	"github.com/controlplane-com/terraform-provider-cpln/internal/provider/validators"
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
	_ resource.Resource                = &ServiceAccountKeyResource{}
	_ resource.ResourceWithImportState = &ServiceAccountKeyResource{}
)

/*** Resource Model ***/

// ServiceAccountKeyResourceModel holds the Terraform state for the resource.
type ServiceAccountKeyResourceModel struct {
	ServiceAccountName types.String `tfsdk:"service_account_name"`
	Description        types.String `tfsdk:"description"`
	Name               types.String `tfsdk:"name"`
	Created            types.String `tfsdk:"created"`
	Key                types.String `tfsdk:"key"`
}

/*** Resource Configuration ***/

// ServiceAccountKeyResource is the resource implementation.
type ServiceAccountKeyResource struct {
	EntityBase
}

// NewServiceAccountKeyResource returns a new instance of the resource implementation.
func NewServiceAccountKeyResource() resource.Resource {
	return &ServiceAccountKeyResource{}
}

// Configure configures the resource before use.
func (sakr *ServiceAccountKeyResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	sakr.EntityBaseConfigure(ctx, req.ProviderData, &resp.Diagnostics)
}

// ImportState sets up the import operation.
func (sakr *ServiceAccountKeyResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Split the import ID
	parts := strings.SplitN(req.ID, ":", 2)

	// Validate that ID has exactly three non-empty segments
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		// Report error when import identifier format is unexpected
		resp.Diagnostics.AddError(
			"Unexpected Import Identifier",
			fmt.Sprintf(
				"Expected import identifier with format: "+
					"'service_account_name:key_name'. Got: %q", req.ID,
			),
		)

		// Abort import operation on error
		return
	}

	// Extract serviceAccountName and keyName from parts
	serviceAccountName, keyName := parts[0], parts[1]

	// Set the generated ID attribute in the Terraform state
	resp.Diagnostics.Append(
		resp.State.SetAttribute(ctx, path.Root("name"), types.StringValue(keyName))...,
	)

	// Set the service_account_name attribute in the Terraform state
	resp.Diagnostics.Append(
		resp.State.SetAttribute(ctx, path.Root("service_account_name"), types.StringValue(serviceAccountName))...,
	)
}

// Metadata provides the resource type name.
func (sakr *ServiceAccountKeyResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = "cpln_service_account_key"
}

// Schema defines the schema for the resource.
func (sakr *ServiceAccountKeyResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"service_account_name": schema.StringAttribute{
				Description: "The name of an existing Service Account this key will belong to.",
				Required:    true,
				Validators: []validator.String{
					validators.NameValidator{},
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"description": schema.StringAttribute{
				Description: "Description of the Service Account Key.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Description: "The generated name of the key.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"created": schema.StringAttribute{
				Description: "The timestamp, in UTC, when the key was created.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"key": schema.StringAttribute{
				Description: "The generated key.",
				Computed:    true,
				Sensitive:   true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

// Create creates the resource.
func (sakr *ServiceAccountKeyResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plannedState ServiceAccountKeyResourceModel

	// Retrieve the planned state from the Terraform configuration
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plannedState)...)

	// Abort on errors to avoid partial or inconsistent state
	if resp.Diagnostics.HasError() {
		return
	}

	// Initialize a new request payload structure and populate it with the planned state
	serviceAccountName, description := sakr.buildRequest(plannedState)

	// Return if an error has occurred during the request payload creation
	if resp.Diagnostics.HasError() {
		return
	}

	// Send the create request to the API client
	responsePayload, err := sakr.client.AddServiceAccountKey(serviceAccountName, description)

	// Handle any other errors that occurred during the API request
	if err != nil {
		resp.Diagnostics.AddError("API Error", fmt.Sprintf("Error creating service account key: %s", err))
		return
	}

	// Map the API response to the Terraform state
	finalState := sakr.buildState(serviceAccountName, responsePayload)

	// Return if an error has occurred during the state creation
	if resp.Diagnostics.HasError() {
		return
	}

	// Set the resource state in Terraform
	resp.Diagnostics.Append(resp.State.Set(ctx, &finalState)...)
}

// Read fetches the current state of the resource.
func (sakr *ServiceAccountKeyResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var plannedState ServiceAccountKeyResourceModel

	// Retrieve the planned state from the Terraform configuration
	resp.Diagnostics.Append(req.State.Get(ctx, &plannedState)...)

	// Abort on errors to avoid partial or inconsistent state
	if resp.Diagnostics.HasError() {
		return
	}

	// Extract necessary values from the planned state
	serviceAccountName := plannedState.ServiceAccountName.ValueString()
	keyName := plannedState.Name.ValueString()

	// Fetch the domain route
	responsePayload, code, err := sakr.client.GetServiceAccount(serviceAccountName)

	// Handle the case where the route is not found (HTTP 404),
	// indicating it has been deleted outside of Terraform. Remove it from state
	if code == 404 {
		resp.State.RemoveResource(ctx)
		return
	}

	// Handle any other errors that occur during the API call
	if err != nil {
		resp.Diagnostics.AddError("API Error", fmt.Sprintf("Error reading service account key: %s", err))
		return
	}

	// Attempt to find the key within keys if the keys are not nil
	if responsePayload.Keys != nil {
		// Iterate over the keys of the service account
		for _, serviceAccountKey := range *responsePayload.Keys {
			if serviceAccountKey.Name == keyName {
				// Map the API response to the Terraform state
				finalState := sakr.buildState(serviceAccountName, &serviceAccountKey)

				// Return if an error has occurred during the state creation
				if resp.Diagnostics.HasError() {
					return
				}

				// Set the updated state in Terraform
				resp.Diagnostics.Append(resp.State.Set(ctx, &finalState)...)
				return
			}
		}
	}

	// If we got here then the key was never found, in this case let's delete this Terraform resource
	resp.State.RemoveResource(ctx)
}

// Update modifies the resource.
func (sakr *ServiceAccountKeyResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plannedState ServiceAccountKeyResourceModel

	// Retrieve the planned state from the Terraform configuration
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plannedState)...)

	// Abort on errors to avoid partial or inconsistent state
	if resp.Diagnostics.HasError() {
		return
	}

	// Set the state in Terraform
	resp.Diagnostics.Append(resp.State.Set(ctx, &plannedState)...)
}

// Delete removes the resource.
func (sakr *ServiceAccountKeyResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state ServiceAccountKeyResourceModel

	// Retrieve the state from the Terraform configuration
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	// Abort on errors to avoid partial or inconsistent state
	if resp.Diagnostics.HasError() {
		return
	}

	// Extract necessary values from the planned state
	serviceAccountName := state.ServiceAccountName.ValueString()
	keyName := state.Name.ValueString()

	// Send a delete request to the API using the name from the state
	err := sakr.client.RemoveServiceAccountKey(serviceAccountName, keyName)

	// Handle errors from the API delete request
	if err != nil {
		// If an error occurs during the delete request, add an error to diagnostics
		resp.Diagnostics.AddError("API Error", fmt.Sprintf("Error deleting service account key: %s", err))
		return
	}

	// Remove the resource from Terraform's state, indicating successful deletion
	resp.State.RemoveResource(ctx)
}

/*** Helpers ***/

// buildRequest creates a request payload from a state model.
func (sakr *ServiceAccountKeyResource) buildRequest(state ServiceAccountKeyResourceModel) (string, string) {
	return *BuildString(state.ServiceAccountName), *BuildString(state.Description)
}

// buildState creates a state model from response payload.
func (sakr *ServiceAccountKeyResource) buildState(serviceAccountName string, serviceAccountKey *client.ServiceAccountKey) ServiceAccountKeyResourceModel {
	// Initialize empty state model
	state := ServiceAccountKeyResourceModel{}

	// Set specific attributes
	state.ServiceAccountName = types.StringValue(serviceAccountName)
	state.Description = types.StringPointerValue(serviceAccountKey.Description)
	state.Name = types.StringValue(serviceAccountKey.Name)
	state.Created = types.StringPointerValue(serviceAccountKey.Created)
	state.Key = types.StringValue(serviceAccountKey.Key)

	// Return completed state model
	return state
}
