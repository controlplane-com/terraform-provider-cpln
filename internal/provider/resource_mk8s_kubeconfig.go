package cpln

import (
	"context"
	"fmt"

	"github.com/controlplane-com/terraform-provider-cpln/internal/provider/validators"
	"github.com/hashicorp/terraform-plugin-framework-validators/resourcevalidator"
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
	_ resource.Resource = &Mk8sKubeconfigResource{}
)

/*** Resource Model ***/

// Mk8sKubeconfigResourceModel holds the Terraform state for the resource.
type Mk8sKubeconfigResourceModel struct {
	ID             types.String `tfsdk:"id"`
	Name           types.String `tfsdk:"name"`
	Profile        types.String `tfsdk:"profile"`
	ServiceAccount types.String `tfsdk:"service_account"`
	Kubeconfig     types.String `tfsdk:"kubeconfig"`
}

/*** Resource Configuration ***/

// Mk8sKubeconfigResource is the resource implementation.
type Mk8sKubeconfigResource struct {
	EntityBase
}

// NewMk8sKubeconfigResource returns a new instance of the resource implementation.
func NewMk8sKubeconfigResource() resource.Resource {
	return &Mk8sKubeconfigResource{}
}

// Configure configures the resource before use.
func (mkr *Mk8sKubeconfigResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	mkr.EntityBaseConfigure(ctx, req.ProviderData, &resp.Diagnostics)
}

// Metadata provides the resource type name.
func (mkr *Mk8sKubeconfigResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = "cpln_mk8s_kubeconfig"
}

// Schema defines the schema for the resource.
func (mkr *Mk8sKubeconfigResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The unique identifier for this MK8s.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Description: "Name of the MK8s to create the Kubeconfig for.",
				Required:    true,
				Validators: []validator.String{
					validators.NameValidator{},
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"profile": schema.StringAttribute{
				Description: "Profile name to extract the token from.",
				Optional:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"service_account": schema.StringAttribute{
				Description: "A service account to add a key to.",
				Optional:    true,
				Validators: []validator.String{
					validators.NameValidator{},
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"kubeconfig": schema.StringAttribute{
				Description: "The Kubeconfig of your MK8s cluster in YAML format.",
				Computed:    true,
				Sensitive:   true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

// ConfigValidators enforces mutual exclusivity between attributes.
func (mkr *Mk8sKubeconfigResource) ConfigValidators(ctx context.Context) []resource.ConfigValidator {
	return []resource.ConfigValidator{
		resourcevalidator.ExactlyOneOf(path.MatchRoot("profile"), path.MatchRoot("service_account")),
	}
}

// Create creates the resource.
func (mkr *Mk8sKubeconfigResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plannedState Mk8sKubeconfigResourceModel

	// Retrieve the planned state from the Terraform configuration
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plannedState)...)

	// Abort on errors to avoid partial or inconsistent state
	if resp.Diagnostics.HasError() {
		return
	}

	// Initialize a new request payload structure and populate it with the planned state
	mk8sName, profileName, serviceAccountName := mkr.buildRequest(plannedState)

	// Return if an error has occurred during the request payload creation
	if resp.Diagnostics.HasError() {
		return
	}

	// Create a new MK8s Kubeconfig using the API client
	kubeconfig, err := mkr.client.CreateMk8sKubeconfig(mk8sName, profileName, serviceAccountName)

	// Handle any other errors that occurred during the API request
	if err != nil {
		resp.Diagnostics.AddError("API Error", fmt.Sprintf("Error creating mk8s kubeconfig: %s", err))
		return
	}

	// Map the API response to the Terraform state
	finalState := mkr.buildState(mk8sName, profileName, serviceAccountName, *kubeconfig)

	// Return if an error has occurred during the state creation
	if resp.Diagnostics.HasError() {
		return
	}

	// Set the resource state in Terraform
	resp.Diagnostics.Append(resp.State.Set(ctx, &finalState)...)
}

// Read fetches the current state of the resource.
func (mkr *Mk8sKubeconfigResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var plannedState Mk8sKubeconfigResourceModel

	// Retrieve the planned state from the Terraform configuration
	resp.Diagnostics.Append(req.State.Get(ctx, &plannedState)...)

	// Abort on errors to avoid partial or inconsistent state
	if resp.Diagnostics.HasError() {
		return
	}

	// Set the updated state in Terraform
	resp.Diagnostics.Append(resp.State.Set(ctx, &plannedState)...)
}

// Update modifies the resource.
func (mkr *Mk8sKubeconfigResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plannedState Mk8sKubeconfigResourceModel

	// Retrieve the planned state from the Terraform configuration
	resp.Diagnostics.Append(req.State.Get(ctx, &plannedState)...)

	// Abort on errors to avoid partial or inconsistent state
	if resp.Diagnostics.HasError() {
		return
	}

	// Set the updated state in Terraform
	resp.Diagnostics.Append(resp.State.Set(ctx, &plannedState)...)
}

// Delete removes the resource.
func (mkr *Mk8sKubeconfigResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Remove the resource from Terraform's state, indicating successful deletion
	resp.State.RemoveResource(ctx)
}

/*** Helpers ***/

// buildRequest creates a request payload from a state model.
func (mkr *Mk8sKubeconfigResource) buildRequest(state Mk8sKubeconfigResourceModel) (string, *string, *string) {
	return state.Name.ValueString(), state.Profile.ValueStringPointer(), state.ServiceAccount.ValueStringPointer()
}

// buildState creates a state model from response payload.
func (mkr *Mk8sKubeconfigResource) buildState(mk8sName string, profile *string, serviceAccount *string, kubeconfig string) Mk8sKubeconfigResourceModel {
	// Initialize empty state model
	state := Mk8sKubeconfigResourceModel{}

	// Set specific attributes
	state.ID = types.StringValue(mkr.getMk8sKubeconfigUnqiueId(mk8sName, profile, serviceAccount))
	state.Name = types.StringValue(mk8sName)
	state.Profile = types.StringPointerValue(profile)
	state.ServiceAccount = types.StringPointerValue(serviceAccount)
	state.Kubeconfig = types.StringValue(kubeconfig)

	// Return completed state model
	return state
}

// getMk8sKubeconfigUnqiueId generates a unique identifier for the Mk8sKubeconfig resource.
func (mkr *Mk8sKubeconfigResource) getMk8sKubeconfigUnqiueId(mk8sName string, profile *string, serviceAccount *string) string {
	// Create an identity based on the profile name
	if profile != nil && len(*profile) != 0 {
		return fmt.Sprintf("%s:profile:%s", mk8sName, *profile)
	}

	// If profile is empty, then service account shouldn't, let's create an identity based on it
	return fmt.Sprintf("%s:service_account:%s", mk8sName, *serviceAccount)
}
