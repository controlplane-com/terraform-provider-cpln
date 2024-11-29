package cpln

import (
	"context"
	"fmt"
	client "github.com/controlplane-com/terraform-provider-cpln/internal/provider/client"
	models "github.com/controlplane-com/terraform-provider-cpln/internal/provider/models/cloud_account"
	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/resourcevalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"regexp"
)

// Ensure resource defined types fully satisfy framework interfaces
var (
	_ resource.Resource                = &CloudAccountResource{}
	_ resource.ResourceWithImportState = &CloudAccountResource{}
)

/*** Resource Model ***/

// CloudAccountResourceModel holds the resource data structure for the Terraform state.
type CloudAccountResourceModel struct {
	BaseResourceModel
	Aws   types.List `tfsdk:"aws"`
	Azure types.List `tfsdk:"azure"`
	Gcp   types.List `tfsdk:"gcp"`
	Ngs   types.List `tfsdk:"ngs"`
}

/*** Resource Configuration ***/

// CloudAccountResource is the resource implementation.
type CloudAccountResource struct {
	client *client.Client
}

// NewCloudAccountResource is a helper function to simplify resource implementation.
func NewCloudAccountResource() resource.Resource {
	return &CloudAccountResource{}
}

// ImportState sets up the import operation to map the imported ID to the "id" attribute in the state.
func (r *CloudAccountResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// Metadata provides the resource type name.
func (r *CloudAccountResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = "cpln_cloud_account"
}

// Schema defines the schema for the resource.
func (r *CloudAccountResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: BaseResourceAttributes("Cloud Account"),
		Blocks: map[string]schema.Block{
			"aws": schema.ListNestedBlock{
				Description: "Contains AWS cloud account configuration.",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"role_arn": schema.StringAttribute{
							Required:    true,
							Description: "Amazon Resource Name (ARN) Role.",
							Validators: []validator.String{
								stringvalidator.RegexMatches(
									regexp.MustCompile(`^arn:(aws|aws-us-gov|aws-cn):iam::[0-9]+:role/.+`),
									"must be a valid Amazon Resource Name (ARN) for an IAM role, formatted as 'arn:partition:service:region:account-id:resource'",
								),
							},
							PlanModifiers: []planmodifier.String{
								stringplanmodifier.RequiresReplace(),
							},
						},
					},
				},
				Validators: []validator.List{
					listvalidator.SizeAtMost(1),
				},
			},
			"azure": schema.ListNestedBlock{
				Description: "Contains Azure cloud account configuration.",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"secret_link": schema.StringAttribute{
							Required:    true,
							Description: "Full link to an Azure secret. (e.g., /org/ORG_NAME/secret/AZURE_SECRET).",
							Validators: []validator.String{
								stringvalidator.RegexMatches(
									regexp.MustCompile(`(/org/[^/]+/.*)|(//.+)`),
									"must be a valid link of an Azure secret within Control Plane",
								),
							},
							PlanModifiers: []planmodifier.String{
								stringplanmodifier.RequiresReplace(),
							},
						},
					},
				},
			},
			"gcp": schema.ListNestedBlock{
				Description: "Contains GCP cloud account configuration.",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"project_id": schema.StringAttribute{
							Required:    true,
							Description: "GCP project ID. Obtained from the GCP cloud console.",
							Validators: []validator.String{
								stringvalidator.RegexMatches(
									regexp.MustCompile(`[a-z]([a-z]|-|[0-9])+`),
									"must be a valid project id",
								),
								stringvalidator.LengthBetween(6, 30),
							},
							PlanModifiers: []planmodifier.String{
								stringplanmodifier.RequiresReplace(),
							},
						},
					},
				},
				Validators: []validator.List{
					listvalidator.SizeAtMost(1),
				},
			},
			"ngs": schema.ListNestedBlock{
				Description: "Contains NGS cloud account configuration.",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"secret_link": schema.StringAttribute{
							Required:    true,
							Description: "Full link to a NATS Account Secret secret. (e.g., /org/ORG_NAME/secret/NATS_ACCOUNT_SECRET).",
							Validators: []validator.String{
								stringvalidator.RegexMatches(
									regexp.MustCompile(`(/org/[^/]+/.*)|(//.+)`),
									"must be a valid link of a NATS secret within Control Plane",
								),
							},
							PlanModifiers: []planmodifier.String{
								stringplanmodifier.RequiresReplace(),
							},
						},
					},
				},
				Validators: []validator.List{
					listvalidator.SizeAtMost(1),
				},
			},
		},
	}
}

// ConfigValidators enforces mutual exclusivity between attributes.
func (r *CloudAccountResource) ConfigValidators(ctx context.Context) []resource.ConfigValidator {
	expressions := []path.Expression{
		path.MatchRoot("aws"),
		path.MatchRoot("azure"),
		path.MatchRoot("gcp"),
		path.MatchRoot("ngs"),
	}

	return []resource.ConfigValidator{
		resourcevalidator.Conflicting(expressions...),
		resourcevalidator.AtLeastOneOf(expressions...),
	}
}

// Configure assigns the provider's client to the resource for API interactions.
func (r *CloudAccountResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
func (r *CloudAccountResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plannedState CloudAccountResourceModel

	// Retrieve the planned state from the Terraform configuration
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plannedState)...)

	// Check if any errors have been added to the diagnostics in previous operations
	// If there are errors, stop further execution to prevent inconsistent or partial changes
	if resp.Diagnostics.HasError() {
		return
	}

	// Initialize a new request payload structure and populate it with the input state
	requestPayload := r.getCloudAccountRequest(ctx, &resp.Diagnostics, plannedState, false)

	// Return if an error has occurred during the request payload creation
	if resp.Diagnostics.HasError() {
		return
	}

	// Send the create request to the API client
	responsePayload, code, err := r.client.CreateCloudAccount(requestPayload)

	// Handle cases where the resource already exists, returning an error to inform the user
	if code == 409 {
		resp.Diagnostics.AddError("Resource already exists", "The resource already exists. You can import the existing resource using the Terraform import command.")
		return
	}

	// Handle any other errors that occurred during the API request
	if err != nil {
		resp.Diagnostics.AddError("API Error", fmt.Sprintf("Error creating cloud account: %s", err))
		return
	}

	// Map the API response to the Terraform state
	finalState := r.getCloudAccountState(ctx, &resp.Diagnostics, responsePayload)

	// Return if an error has occurred during the state creation
	if resp.Diagnostics.HasError() {
		return
	}

	// Set the resource state in Terraform
	resp.Diagnostics.Append(resp.State.Set(ctx, &finalState)...)
}

// Read fetches the current state of the resource.
func (r *CloudAccountResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var plannedState CloudAccountResourceModel

	// Retrieve the planned state from the Terraform configuration
	resp.Diagnostics.Append(req.State.Get(ctx, &plannedState)...)

	// Check if any errors have been added to the diagnostics in previous operations
	// If there are errors, stop further execution to prevent inconsistent or partial changes
	if resp.Diagnostics.HasError() {
		return
	}

	// Call the API to get the latest state for the resource
	responsePayload, code, err := r.client.GetCloudAccount(plannedState.ID.ValueString())

	// Handle the case where the resource is not found (HTTP 404),
	// indicating it has been deleted outside of Terraform. Remove it from state
	if code == 404 {
		resp.State.RemoveResource(ctx)
		return
	}

	// Handle any other errors that occur during the API call
	if err != nil {
		resp.Diagnostics.AddError("API Error", fmt.Sprintf("Error reading cloud account: %s", err))
		return
	}

	// Map the API response to the Terraform state
	finalState := r.getCloudAccountState(ctx, &resp.Diagnostics, responsePayload)

	// Return if an error has occurred during the state creation
	if resp.Diagnostics.HasError() {
		return
	}

	// Set the updated state in Terraform
	resp.Diagnostics.Append(resp.State.Set(ctx, &finalState)...)
}

// Update modifies the resource.
func (r *CloudAccountResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plannedState CloudAccountResourceModel

	// Retrieve the planned state from the Terraform configuration
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plannedState)...)

	// Check if any errors have been added to the diagnostics in previous operations
	// If there are errors, stop further execution to prevent inconsistent or partial changes
	if resp.Diagnostics.HasError() {
		return
	}

	// Initialize a new request payload structure and populate it with the input state
	requestPayload := r.getCloudAccountRequest(ctx, &resp.Diagnostics, plannedState, true)

	// Return if an error has occurred during the request payload creation
	if resp.Diagnostics.HasError() {
		return
	}

	// Send the update request to the API with the modified data
	responsePayload, _, err := r.client.UpdateCloudAccount(requestPayload)

	// Handle errors from the API update request
	if err != nil {
		resp.Diagnostics.AddError("API Error", fmt.Sprintf("Error updating cloud account: %s", err))
		return
	}

	// Map the API response to the Terraform state
	finalState := r.getCloudAccountState(ctx, &resp.Diagnostics, responsePayload)

	// Return if an error has occurred during the state creation
	if resp.Diagnostics.HasError() {
		return
	}

	// Set the updated state in Terraform, appending any diagnostics from this operation
	resp.Diagnostics.Append(resp.State.Set(ctx, &finalState)...)
}

// Delete removes the resource.
func (r *CloudAccountResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state CloudAccountResourceModel

	// Retrieve the state from the Terraform configuration
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	// Check if any errors have been added to the diagnostics in previous operations
	// If there are errors, stop further execution to prevent inconsistent or partial changes
	if resp.Diagnostics.HasError() {
		return
	}

	// Send a delete request to the API using the name from the state
	err := r.client.DeleteCloudAccount(state.Name.ValueString())

	// Handle errors from the API delete request
	if err != nil {
		// If an error occurs during the delete request, add an error to diagnostics
		resp.Diagnostics.AddError("API Error", fmt.Sprintf("Error deleting cloud account: %s", err))
		return
	}

	// Remove the resource from Terraform state, indicating successful deletion
	resp.State.RemoveResource(ctx)
}

/*** Helpers ***/

// getCloudAccountRequest creates a request payload from a state model.
func (r *CloudAccountResource) getCloudAccountRequest(ctx context.Context, diags *diag.Diagnostics, state CloudAccountResourceModel, isUpdate bool) client.CloudAccount {
	// Define a new request payload object
	requestPayload := client.CloudAccount{
		Provider: r.getProviderName(state),
	}

	// Map planned state attributes to the API request's Base object
	if isUpdate {
		UpdateReplaceBaseClientFromState(&requestPayload.Base, state.BaseResourceModel)
	} else {
		UpdateBaseClientFromState(&requestPayload.Base, state.BaseResourceModel)
	}

	// Set specific attributes
	if requestPayload.Provider != nil {
		requestPayload.Data = &client.CloudAccountConfig{}

		// Build specified cloud provider
		switch *requestPayload.Provider {
		case "aws":
			requestPayload.Data.RoleArn = r.buildAws(ctx, diags, state.Aws)
		case "azure":
			requestPayload.Data.SecretLink = r.buildAzure(ctx, diags, state.Azure)
		case "gcp":
			requestPayload.Data.ProjectId = r.buildGcp(ctx, diags, state.Gcp)
		case "ngs":
			requestPayload.Data.SecretLink = r.buildNgs(ctx, diags, state.Ngs)
		}
	}

	// Return the request payload object
	return requestPayload
}

// getCloudAccountState creates a state model from response payload.
func (r *CloudAccountResource) getCloudAccountState(ctx context.Context, diags *diag.Diagnostics, cloudAccount *client.CloudAccount) CloudAccountResourceModel {
	// Define the state
	state := CloudAccountResourceModel{}

	// Map shared attributes from the API response's Base object to the Terraform state model
	UpdateStateFromBaseClient(&state.BaseResourceModel, cloudAccount.Base)

	// Update specific attributes
	state.Aws = r.flattenAws(ctx, diags, cloudAccount)
	state.Azure = r.flattenAzure(ctx, diags, cloudAccount)
	state.Gcp = r.flattenGcp(ctx, diags, cloudAccount)
	state.Ngs = r.flattenNgs(ctx, diags, cloudAccount)

	// Return the built state
	return state
}

// getProviderName determines the provider name based on the non-null and known state of cloud provider attributes.
func (r *CloudAccountResource) getProviderName(state CloudAccountResourceModel) *string {
	// Define a mapping of provider names to their corresponding state attributes for easier iteration.
	providers := map[string]types.List{
		"aws":   state.Aws,
		"azure": state.Azure,
		"gcp":   state.Gcp,
		"ngs":   state.Ngs,
	}

	// Iterate over the provider mapping to find the first non-null and known provider.
	for name, attr := range providers {
		if !attr.IsNull() && !attr.IsUnknown() {
			// Return the name of the provider as a pointer.
			return &name
		}
	}

	// If no provider is found, return nil.
	return nil
}

// Builders //

// buildAws maps the Terraform state for the AWS block to the target CloudAccountConfig object.
func (r *CloudAccountResource) buildAws(ctx context.Context, diags *diag.Diagnostics, state types.List) *string {
	// If the state is unknown or null, there is no block to process, so exit early
	if state.IsNull() || state.IsUnknown() {
		return nil
	}

	// Declare a slice to hold the unmarshalled block from the Terraform state
	var blocks []models.Aws

	// Convert the Terraform state list into a Go slice of block models
	diags.Append(state.ElementsAs(ctx, &blocks, false)...)

	// If there are any diagnostics errors during the conversion, stop further processing
	if diags.HasError() {
		return nil
	}

	// Assign the properties from the single block to the target configuration object
	return blocks[0].RoleArn.ValueStringPointer()
}

// buildAzure maps the Terraform state for the Azure block to the target CloudAccountConfig object.
func (r *CloudAccountResource) buildAzure(ctx context.Context, diags *diag.Diagnostics, state types.List) *string {
	// If the state is unknown or null, there is no block to process, so exit early
	if state.IsNull() || state.IsUnknown() {
		return nil
	}

	// Declare a slice to hold the unmarshalled block from the Terraform state
	var blocks []models.Azure

	// Convert the Terraform state list into a Go slice of block models
	diags.Append(state.ElementsAs(ctx, &blocks, false)...)

	// If there are any diagnostics errors during the conversion, stop further processing
	if diags.HasError() {
		return nil
	}

	// Assign the properties from the single block to the target configuration object
	return blocks[0].SecretLink.ValueStringPointer()
}

// buildGcp maps the Terraform state for the GCP block to the target CloudAccountConfig object.
func (r *CloudAccountResource) buildGcp(ctx context.Context, diags *diag.Diagnostics, state types.List) *string {
	// If the state is unknown or null, there is no block to process, so exit early
	if state.IsNull() || state.IsUnknown() {
		return nil
	}

	// Declare a slice to hold the unmarshalled block from the Terraform state
	var blocks []models.Gcp

	// Convert the Terraform state list into a Go slice of block models
	diags.Append(state.ElementsAs(ctx, &blocks, false)...)

	// If there are any diagnostics errors during the conversion, stop further processing
	if diags.HasError() {
		return nil
	}

	// Assign the properties from the single block to the target configuration object
	return blocks[0].ProjectId.ValueStringPointer()
}

// buildNgs maps the Terraform state for the NGS block to the target CloudAccountConfig object.
func (r *CloudAccountResource) buildNgs(ctx context.Context, diags *diag.Diagnostics, state types.List) *string {
	// If the state is unknown or null, there is no block to process, so exit early
	if state.IsNull() || state.IsUnknown() {
		return nil
	}

	// Declare a slice to hold the unmarshalled block from the Terraform state
	var blocks []models.Ngs

	// Convert the Terraform state list into a Go slice of block models
	diags.Append(state.ElementsAs(ctx, &blocks, false)...)

	// If there are any diagnostics errors during the conversion, stop further processing
	if diags.HasError() {
		return nil
	}

	// Assign the properties from the single block to the target configuration object
	return blocks[0].SecretLink.ValueStringPointer()
}

// Flattens //

// flattenAws maps the CloudAccountConfig object to a Terraform state list.
func (r *CloudAccountResource) flattenAws(ctx context.Context, diags *diag.Diagnostics, cloudAccount *client.CloudAccount) types.List {
	// Initialize a default block
	block := models.Aws{
		RoleArn: types.StringNull(),
	}

	// Check if the input data is valid
	if cloudAccount == nil || cloudAccount.Provider == nil || *cloudAccount.Provider != "aws" || cloudAccount.Data == nil || cloudAccount.Data.RoleArn == nil {
		return types.ListNull(block.AttributeTypes())
	}

	// Populate the properties in the block with the data from the input data
	block.RoleArn = types.StringPointerValue(cloudAccount.Data.RoleArn)

	// Convert the populated block into a Terraform-compatible list of objects
	list, d := types.ListValueFrom(ctx, block.AttributeTypes(), []models.Aws{block})

	// Append any diagnostics from the conversion process
	// If there are errors, return a null list with the appropriate type
	if d.HasError() {
		diags.Append(d...)
		return types.ListNull(block.AttributeTypes())
	}

	// Return the successfully created Terraform-compatible list
	return list
}

// flattenAzure maps the CloudAccountConfig object to a Terraform state list.
func (r *CloudAccountResource) flattenAzure(ctx context.Context, diags *diag.Diagnostics, cloudAccount *client.CloudAccount) types.List {
	// Initialize a default block
	block := models.Azure{
		SecretLink: types.StringNull(),
	}

	// Check if the input data is valid
	if cloudAccount == nil || cloudAccount.Provider == nil || *cloudAccount.Provider != "azure" || cloudAccount.Data == nil || cloudAccount.Data.SecretLink == nil {
		return types.ListNull(block.AttributeTypes())
	}

	// Populate the properties in the block with the data from the input data
	block.SecretLink = types.StringPointerValue(cloudAccount.Data.SecretLink)

	// Convert the populated block into a Terraform-compatible list of objects
	list, d := types.ListValueFrom(ctx, block.AttributeTypes(), []models.Azure{block})

	// Append any diagnostics from the conversion process
	// If there are errors, return a null list with the appropriate type
	if d.HasError() {
		diags.Append(d...)
		return types.ListNull(block.AttributeTypes())
	}

	// Return the successfully created Terraform-compatible list
	return list
}

// flattenGcp maps the CloudAccountConfig object to a Terraform state list.
func (r *CloudAccountResource) flattenGcp(ctx context.Context, diags *diag.Diagnostics, cloudAccount *client.CloudAccount) types.List {
	// Initialize a default block
	block := models.Gcp{
		ProjectId: types.StringNull(),
	}

	// Check if the input data is valid
	if cloudAccount == nil || cloudAccount.Provider == nil || *cloudAccount.Provider != "gcp" || cloudAccount.Data == nil || cloudAccount.Data.ProjectId == nil {
		return types.ListNull(block.AttributeTypes())
	}

	// Populate the properties in the block with the data from the input data
	block.ProjectId = types.StringPointerValue(cloudAccount.Data.ProjectId)

	// Convert the populated block into a Terraform-compatible list of objects
	list, d := types.ListValueFrom(ctx, block.AttributeTypes(), []models.Gcp{block})

	// Append any diagnostics from the conversion process
	// If there are errors, return a null list with the appropriate type
	if d.HasError() {
		diags.Append(d...)
		return types.ListNull(block.AttributeTypes())
	}

	// Return the successfully created Terraform-compatible list
	return list
}

// flattenNgs maps the CloudAccountConfig object to a Terraform state list.
func (r *CloudAccountResource) flattenNgs(ctx context.Context, diags *diag.Diagnostics, cloudAccount *client.CloudAccount) types.List {
	// Initialize a default block
	block := models.Ngs{
		SecretLink: types.StringNull(),
	}

	// Check if the input data is valid
	if cloudAccount == nil || cloudAccount.Provider == nil || *cloudAccount.Provider != "ngs" || cloudAccount.Data == nil || cloudAccount.Data.SecretLink == nil {
		return types.ListNull(block.AttributeTypes())
	}

	// Populate the properties in the block with the data from the input data
	block.SecretLink = types.StringPointerValue(cloudAccount.Data.SecretLink)

	// Convert the populated block into a Terraform-compatible list of objects
	list, d := types.ListValueFrom(ctx, block.AttributeTypes(), []models.Ngs{block})

	// Append any diagnostics from the conversion process
	// If there are errors, return a null list with the appropriate type
	if d.HasError() {
		diags.Append(d...)
		return types.ListNull(block.AttributeTypes())
	}

	// Return the successfully created Terraform-compatible list
	return list
}
