package cpln

import (
	"context"
	"regexp"

	client "github.com/controlplane-com/terraform-provider-cpln/internal/provider/client"
	models "github.com/controlplane-com/terraform-provider-cpln/internal/provider/models/cloud_account"
	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/resourcevalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure resource implements required interfaces.
var (
	_ resource.Resource                = &CloudAccountResource{}
	_ resource.ResourceWithImportState = &CloudAccountResource{}
)

/*** Resource Model ***/

// CloudAccountResourceModel holds the Terraform state for the resource.
type CloudAccountResourceModel struct {
	EntityBaseModel
	Aws                   types.List   `tfsdk:"aws"`
	Azure                 types.List   `tfsdk:"azure"`
	Gcp                   types.List   `tfsdk:"gcp"`
	Ngs                   types.List   `tfsdk:"ngs"`
	GcpServiceAccountName types.String `tfsdk:"gcp_service_account_name"`
	GcpRoles              types.Set    `tfsdk:"gcp_roles"`
}

/*** Resource Configuration ***/

// CloudAccountResource is the resource implementation.
type CloudAccountResource struct {
	EntityBase
	Operations EntityOperations[CloudAccountResourceModel, client.CloudAccount]
}

// NewCloudAccountResource returns a new instance of the resource implementation.
func NewCloudAccountResource() resource.Resource {
	return &CloudAccountResource{}
}

// Configure configures the resource before use.
func (car *CloudAccountResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	car.EntityBaseConfigure(ctx, req.ProviderData, &resp.Diagnostics)
	car.Operations = NewEntityOperations(car.client, &CloudAccountResourceOperator{})
}

// ImportState sets up the import operation to map the imported ID to the "id" attribute in the state.
func (car *CloudAccountResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// Metadata provides the resource type name.
func (car *CloudAccountResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = "cpln_cloud_account"
}

// Schema defines the schema for the resource.
func (car *CloudAccountResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: MergeAttributes(car.EntityBaseAttributes("Cloud Account"), map[string]schema.Attribute{
			"gcp_service_account_name": schema.StringAttribute{
				Description: "GCP service account name used during the configuration of the cloud account at GCP.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"gcp_roles": schema.SetAttribute{
				Description: "GCP roles used during the configuration of the cloud account at GCP.",
				Computed:    true,
				ElementType: types.StringType,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
			},
		}),
		Blocks: map[string]schema.Block{
			"aws": schema.ListNestedBlock{
				Description: "Contains AWS cloud account configuration.",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"role_arn": schema.StringAttribute{
							Description: "Amazon Resource Name (ARN) Role.",
							Optional:    true,
							Validators: []validator.String{
								stringvalidator.RegexMatches(
									regexp.MustCompile(`^arn:(aws|aws-us-gov|aws-cn):iam::[0-9]+:role/.+`),
									"must be a valid Amazon Resource Name (ARN) for an IAM role, formatted as 'arn:partition:service:region:account-id:resource'",
								),
							},
							PlanModifiers: []planmodifier.String{
								car.RequiresReplaceOnChangeOrRemoval(),
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
							Description: "Full link to an Azure secret. (e.g., /org/ORG_NAME/secret/AZURE_SECRET).",
							Optional:    true,
							Validators: []validator.String{
								stringvalidator.RegexMatches(
									regexp.MustCompile(`(/org/[^/]+/.*)|(//.+)`),
									"must be a valid link of an Azure secret within Control Plane",
								),
							},
							PlanModifiers: []planmodifier.String{
								car.RequiresReplaceOnChangeOrRemoval(),
							},
						},
					},
				},
				Validators: []validator.List{
					listvalidator.SizeAtMost(1),
				},
			},
			"gcp": schema.ListNestedBlock{
				Description: "Contains GCP cloud account configuration.",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"project_id": schema.StringAttribute{
							Description: "GCP project ID. Obtained from the GCP cloud console.",
							Optional:    true,
							Validators: []validator.String{
								stringvalidator.RegexMatches(
									regexp.MustCompile(`[a-z]([a-z]|-|[0-9])+`),
									"must be a valid project id",
								),
								stringvalidator.LengthBetween(6, 30),
							},
							PlanModifiers: []planmodifier.String{
								car.RequiresReplaceOnChangeOrRemoval(),
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
							Description: "Full link to a NATS Account Secret secret. (e.g., /org/ORG_NAME/secret/NATS_ACCOUNT_SECRET).",
							Optional:    true,
							Validators: []validator.String{
								stringvalidator.RegexMatches(
									regexp.MustCompile(`(/org/[^/]+/.*)|(//.+)`),
									"must be a valid link of a NATS secret within Control Plane",
								),
							},
							PlanModifiers: []planmodifier.String{
								car.RequiresReplaceOnChangeOrRemoval(),
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
func (car *CloudAccountResource) ConfigValidators(ctx context.Context) []resource.ConfigValidator {
	expressions := []path.Expression{
		path.MatchRoot("aws"),
		path.MatchRoot("azure"),
		path.MatchRoot("gcp"),
		path.MatchRoot("ngs"),
	}

	return []resource.ConfigValidator{
		resourcevalidator.ExactlyOneOf(expressions...),
	}
}

// Create creates the resource.
func (car *CloudAccountResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	CreateGeneric(ctx, req, resp, car.Operations)
}

// Read fetches the current state of the resource.
func (car *CloudAccountResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	ReadGeneric(ctx, req, resp, car.Operations)
}

// Update modifies the resource.
func (car *CloudAccountResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	UpdateGeneric(ctx, req, resp, car.Operations)
}

// Delete removes the resource.
func (car *CloudAccountResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	DeleteGeneric(ctx, req, resp, car.Operations)
}

/*** Plan Modifiers ***/

// RequiresReplaceOnChangeOrRemoval returns a plan modifier that recreates the resource when an existing string attribute is modified or removed.
func (car *CloudAccountResource) RequiresReplaceOnChangeOrRemoval() planmodifier.String {
	return stringplanmodifier.RequiresReplaceIf(
		func(ctx context.Context,
			req planmodifier.StringRequest,
			resp *stringplanmodifier.RequiresReplaceIfFuncResponse,
		) {
			// If the prior state had a non-empty string, require replace
			if len(req.StateValue.ValueString()) != 0 {
				resp.RequiresReplace = true
			}
		},
		"Recreate resource when modifying or removing an existing value",
		"If this attribute previously had a non-empty value, any change (including removal) will cause the resource to be replaced.",
	)
}

/*** Resource Operator ***/

// CloudAccountResourceOperator is the operator for managing the state.
type CloudAccountResourceOperator struct {
	EntityOperator[CloudAccountResourceModel]
}

// NewAPIRequest creates a request payload from a state model.
func (caro *CloudAccountResourceOperator) NewAPIRequest(isUpdate bool) client.CloudAccount {
	// Initialize a new request payload
	requestPayload := client.CloudAccount{
		Provider: caro.getCloudAccountProviderName(caro.Plan),
	}

	// Populate Base fields from state
	caro.Plan.Fill(&requestPayload.Base, isUpdate)

	// Set specific attributes
	if requestPayload.Provider != nil {
		// Build specified cloud provider data
		switch *requestPayload.Provider {
		case "aws":
			if data := caro.buildAws(caro.Plan.Aws); data != nil {
				requestPayload.Data = &client.CloudAccountConfig{
					RoleArn: data,
				}
			}
		case "azure":
			if data := caro.buildAzure(caro.Plan.Azure); data != nil {
				requestPayload.Data = &client.CloudAccountConfig{
					SecretLink: data,
				}
			}
		case "gcp":
			if data := caro.buildGcp(caro.Plan.Gcp); data != nil {
				requestPayload.Data = &client.CloudAccountConfig{
					ProjectId: data,
				}
			}
		case "ngs":
			if data := caro.buildNgs(caro.Plan.Ngs); data != nil {
				requestPayload.Data = &client.CloudAccountConfig{
					SecretLink: data,
				}
			}
		}
	}

	// Return constructed request payload
	return requestPayload
}

// MapResponseToState creates a state model from response payload.
func (caro *CloudAccountResourceOperator) MapResponseToState(cloudAccount *client.CloudAccount, isCreate bool) CloudAccountResourceModel {
	// Initialize empty state model
	state := CloudAccountResourceModel{}

	// Populate common fields from base resource data
	state.From(cloudAccount.Base)

	// Set specific attributes
	state.Aws = caro.flattenAws(cloudAccount)
	state.Azure = caro.flattenAzure(cloudAccount)
	state.Gcp = caro.flattenGcp(cloudAccount)
	state.Ngs = caro.flattenNgs(cloudAccount)
	state.GcpServiceAccountName = types.StringValue("cpln-" + caro.Client.Org + "@cpln-prod01.iam.gserviceaccount.com")
	state.GcpRoles = FlattenSetString(&GcpRoles)

	// Return completed state model
	return state
}

// InvokeCreate invokes the Create API to create a new resource.
func (caro *CloudAccountResourceOperator) InvokeCreate(req client.CloudAccount) (*client.CloudAccount, int, error) {
	return caro.Client.CreateCloudAccount(req)
}

// InvokeRead invokes the Get API to retrieve an existing resource by name.
func (caro *CloudAccountResourceOperator) InvokeRead(name string) (*client.CloudAccount, int, error) {
	return caro.Client.GetCloudAccount(name)
}

// InvokeUpdate invokes the Update API to update an existing resource.
func (caro *CloudAccountResourceOperator) InvokeUpdate(req client.CloudAccount) (*client.CloudAccount, int, error) {
	return caro.Client.UpdateCloudAccount(req)
}

// InvokeDelete invokes the Delete API to delete a resource by name.
func (caro *CloudAccountResourceOperator) InvokeDelete(name string) error {
	return caro.Client.DeleteCloudAccount(name)
}

// getCloudAccountProviderName determines the provider name based on the non-null and known state of cloud provider attributes.
func (caro *CloudAccountResourceOperator) getCloudAccountProviderName(state CloudAccountResourceModel) *string {
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

// buildAws maps the Terraform state for the AWS block to the target CloudAccountConfig struct.
func (caro *CloudAccountResourceOperator) buildAws(state types.List) *string {
	// Convert Terraform list into model blocks using generic helper
	blocks, ok := BuildList[models.AwsModel](caro.Ctx, caro.Diags, state)

	// Return nil if conversion failed or list was empty
	if !ok {
		return nil
	}

	// Assign the properties from the single block to the target configuration object
	return blocks[0].RoleArn.ValueStringPointer()
}

// buildAzure maps the Terraform state for the Azure block to the target CloudAccountConfig struct.
func (caro *CloudAccountResourceOperator) buildAzure(state types.List) *string {
	// Convert Terraform list into model blocks using generic helper
	blocks, ok := BuildList[models.AzureModel](caro.Ctx, caro.Diags, state)

	// Return nil if conversion failed or list was empty
	if !ok {
		return nil
	}

	// Assign the properties from the single block to the target configuration object
	return blocks[0].SecretLink.ValueStringPointer()
}

// buildGcp maps the Terraform state for the GCP block to the target CloudAccountConfig struct.
func (caro *CloudAccountResourceOperator) buildGcp(state types.List) *string {
	// Convert Terraform list into model blocks using generic helper
	blocks, ok := BuildList[models.GcpModel](caro.Ctx, caro.Diags, state)

	// Return nil if conversion failed or list was empty
	if !ok {
		return nil
	}

	// Assign the properties from the single block to the target configuration object
	return blocks[0].ProjectId.ValueStringPointer()
}

// buildNgs maps the Terraform state for the NGS block to the target CloudAccountConfig struct.
func (caro *CloudAccountResourceOperator) buildNgs(state types.List) *string {
	// Convert Terraform list into model blocks using generic helper
	blocks, ok := BuildList[models.NgsModel](caro.Ctx, caro.Diags, state)

	// Return nil if conversion failed or list was empty
	if !ok {
		return nil
	}

	// Assign the properties from the single block to the target configuration object
	return blocks[0].SecretLink.ValueStringPointer()
}

// Flatteners //

// flattenAws maps the CloudAccountConfig struct to a Terraform state list.
func (caro *CloudAccountResourceOperator) flattenAws(cloudAccount *client.CloudAccount) types.List {
	// Initialize a default block
	block := models.AwsModel{
		RoleArn: types.StringNull(),
	}

	// Return a nil list if the provider isn't the same
	if cloudAccount == nil || *cloudAccount.Provider != "aws" {
		return types.ListNull(block.AttributeTypes())
	}

	// Populate the properties in the block with the data from the input data
	if cloudAccount.Data != nil {
		block.RoleArn = types.StringPointerValue(cloudAccount.Data.RoleArn)
	}

	// Return the successfully created types.List
	return FlattenList(caro.Ctx, caro.Diags, []models.AwsModel{block})
}

// flattenAzure maps the CloudAccountConfig struct to a Terraform state list.
func (caro *CloudAccountResourceOperator) flattenAzure(cloudAccount *client.CloudAccount) types.List {
	// Initialize a default block
	block := models.AzureModel{
		SecretLink: types.StringNull(),
	}

	// Return a nil list if the provider isn't the same
	if cloudAccount == nil || *cloudAccount.Provider != "azure" {
		return types.ListNull(block.AttributeTypes())
	}

	// Populate the properties in the block with the data from the input data
	if cloudAccount.Data != nil {
		block.SecretLink = types.StringPointerValue(cloudAccount.Data.SecretLink)
	}

	// Return the successfully created types.List
	return FlattenList(caro.Ctx, caro.Diags, []models.AzureModel{block})
}

// flattenGcp maps the CloudAccountConfig struct to a Terraform state list.
func (caro *CloudAccountResourceOperator) flattenGcp(cloudAccount *client.CloudAccount) types.List {
	// Initialize a default block
	block := models.GcpModel{
		ProjectId: types.StringNull(),
	}

	// Return a nil list if the provider isn't the same
	if cloudAccount == nil || *cloudAccount.Provider != "gcp" {
		return types.ListNull(block.AttributeTypes())
	}

	// Populate the properties in the block with the data from the input data
	if cloudAccount.Data != nil {
		block.ProjectId = types.StringPointerValue(cloudAccount.Data.ProjectId)
	}

	// Return the successfully created types.List
	return FlattenList(caro.Ctx, caro.Diags, []models.GcpModel{block})
}

// flattenNgs maps the CloudAccountConfig struct to a Terraform state list.
func (caro *CloudAccountResourceOperator) flattenNgs(cloudAccount *client.CloudAccount) types.List {
	// Initialize a default block
	block := models.NgsModel{
		SecretLink: types.StringNull(),
	}

	// Return a nil list if the provider isn't the same
	if cloudAccount == nil || *cloudAccount.Provider != "ngs" {
		return types.ListNull(block.AttributeTypes())
	}

	// Populate the properties in the block with the data from the input data
	if cloudAccount.Data != nil {
		block.SecretLink = types.StringPointerValue(cloudAccount.Data.SecretLink)
	}

	// Return the successfully created types.List
	return FlattenList(caro.Ctx, caro.Diags, []models.NgsModel{block})
}
