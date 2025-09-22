package cpln

import (
	"context"
	"fmt"
	"regexp"
	"strings"
	"unicode"

	client "github.com/controlplane-com/terraform-provider-cpln/internal/provider/client"
	models "github.com/controlplane-com/terraform-provider-cpln/internal/provider/models/identity"
	"github.com/controlplane-com/terraform-provider-cpln/internal/provider/validators"
	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/objectvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int32default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure resource implements required interfaces.
var (
	_ resource.Resource                = &IdentityResource{}
	_ resource.ResourceWithImportState = &IdentityResource{}
)

/*** Resource Model ***/

// IdentityResourceModel holds the Terraform state for the resource.
type IdentityResourceModel struct {
	EntityBaseModel
	Gvc                   types.String `tfsdk:"gvc"`
	Status                types.Map    `tfsdk:"status"`
	AwsAccessPolicy       types.List   `tfsdk:"aws_access_policy"`
	GcpAccessPolicy       types.List   `tfsdk:"gcp_access_policy"`
	AzureAccessPolicy     types.List   `tfsdk:"azure_access_policy"`
	NgsAccessPolicy       types.List   `tfsdk:"ngs_access_policy"`
	NetworkResource       types.Set    `tfsdk:"network_resource"`
	NativeNetworkResource types.Set    `tfsdk:"native_network_resource"`
}

/*** Resource Configuration ***/

// IdentityResource is the resource implementation.
type IdentityResource struct {
	EntityBase
	Operations EntityOperations[IdentityResourceModel, client.Identity]
}

// NewIdentityResource returns a new instance of the resource implementation.
func NewIdentityResource() resource.Resource {
	return &IdentityResource{}
}

// Configure configures the resource before use.
func (ir *IdentityResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	ir.EntityBaseConfigure(ctx, req.ProviderData, &resp.Diagnostics)
	ir.Operations = NewEntityOperations(ir.client, &IdentityResourceOperator{})
}

// ImportState sets up the import operation to map the imported ID to the "id" attribute in the state.
func (ir *IdentityResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Split the import ID
	parts := strings.SplitN(req.ID, ":", 2)

	// Validate that ID has exactly three non-empty segments
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		// Report error when import identifier format is unexpected
		resp.Diagnostics.AddError(
			"Unexpected Import Identifier",
			fmt.Sprintf(
				"Expected import identifier with format: "+
					"'gvc:identity_name'. Got: %q", req.ID,
			),
		)

		// Abort import operation on error
		return
	}

	// Extract gvc and identityName from parts
	gvc, identityName := parts[0], parts[1]

	// Set the generated ID attribute in the Terraform state
	resp.Diagnostics.Append(
		resp.State.SetAttribute(ctx, path.Root("id"), types.StringValue(identityName))...,
	)

	// Set the GVC attribute in the Terraform state
	resp.Diagnostics.Append(
		resp.State.SetAttribute(ctx, path.Root("gvc"), types.StringValue(gvc))...,
	)
}

// Metadata provides the resource type name.
func (ir *IdentityResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = "cpln_identity"
}

// Schema defines the schema for the resource.
func (ir *IdentityResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: MergeAttributes(ir.EntityBaseAttributes("identity"), map[string]schema.Attribute{
			"gvc": schema.StringAttribute{
				Description: "The GVC to which this identity belongs.",
				Required:    true,
				Validators: []validator.String{
					validators.NameValidator{},
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"status": schema.MapAttribute{
				Description: "Key-value map of identity status. Available fields: `objectName`.",
				ElementType: types.StringType,
				Computed:    true,
			},
		}),
		Blocks: map[string]schema.Block{
			"aws_access_policy": schema.ListNestedBlock{
				Description: "A set of access policy rules that defines the actions and resources that an identity can access within an AWS environment.",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"cloud_account_link": schema.StringAttribute{
							Description: "Full link to referenced cloud account.",
							Required:    true,
							Validators: []validator.String{
								validators.LinkValidator{},
							},
						},
						"policy_refs": schema.SetAttribute{
							Description: "List of policies.",
							ElementType: types.StringType,
							Optional:    true,
						},
						"role_name": schema.StringAttribute{
							Description: "Role name.",
							Optional:    true,
							Validators: []validator.String{
								stringvalidator.RegexMatches(
									regexp.MustCompile(`^([a-zA-Z0-9/+=,.@_-])+$`),
									"must contain only letters, numbers, and the symbols / += , . @ _ -",
								),
								stringvalidator.LengthAtMost(64),
							},
						},
					},
					Blocks: map[string]schema.Block{
						"trust_policy": schema.SetNestedBlock{
							Description: "The trust policy for the role.",
							NestedObject: schema.NestedBlockObject{
								Attributes: map[string]schema.Attribute{
									"version": schema.StringAttribute{
										Description: "Version of the policy.",
										Optional:    true,
										Computed:    true,
										Default:     stringdefault.StaticString("2012-10-17"),
									},
									"statement": schema.SetAttribute{
										Description: "List of statements.",
										Optional:    true,
										ElementType: types.MapType{
											ElemType: types.StringType,
										},
									},
								},
							},
						},
					},
				},
				Validators: []validator.List{
					listvalidator.SizeAtMost(1),
				},
			},
			"gcp_access_policy": schema.ListNestedBlock{
				Description: "The GCP access policy can either contain an existing service_account or multiple bindings.",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"cloud_account_link": schema.StringAttribute{
							Description: "Full link to referenced cloud account.",
							Required:    true,
							Validators: []validator.String{
								validators.LinkValidator{},
							},
						},
						"scopes": schema.StringAttribute{
							Description: "Comma delimited list of GCP scope URLs.",
							Optional:    true,
							Computed:    true,
							Default:     stringdefault.StaticString("https://www.googleapis.com/auth/cloud-platform"),
						},
						"service_account": schema.StringAttribute{
							Description: "Name of existing GCP service account.",
							Optional:    true,
							Validators: []validator.String{
								stringvalidator.RegexMatches(
									regexp.MustCompile(`^.*\.gserviceaccount\.com$`),
									"must be a valid GCP service account email (i.e. ending in .gserviceaccount.com)",
								),
							},
						},
					},
					Blocks: map[string]schema.Block{
						"binding": schema.SetNestedBlock{
							Description: "The association or connection between a particular identity, such as a user or a group, and a set of permissions or roles within the system.",
							NestedObject: schema.NestedBlockObject{
								Attributes: map[string]schema.Attribute{
									"resource": schema.StringAttribute{
										Description: "Name of resource for binding.",
										Optional:    true,
									},
									"roles": schema.SetAttribute{
										Description: "List of allowed roles.",
										ElementType: types.StringType,
										Optional:    true,
										Validators: []validator.Set{
											setvalidator.SizeAtLeast(1),
										},
									},
								},
							},
						},
					},
				},
				Validators: []validator.List{
					listvalidator.SizeAtMost(1),
				},
			},
			"azure_access_policy": schema.ListNestedBlock{
				Description: "A set of access policy rules that defines the actions and resources that an identity can access within an Azure environment.",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"cloud_account_link": schema.StringAttribute{
							Description: "Full link to referenced cloud account.",
							Required:    true,
							Validators: []validator.String{
								validators.LinkValidator{},
							},
						},
					},
					Blocks: map[string]schema.Block{
						"role_assignment": schema.SetNestedBlock{
							Description: "The process of assigning specific roles or permissions to an entity, such as a user or a service principal, within the system.",
							NestedObject: schema.NestedBlockObject{
								Attributes: map[string]schema.Attribute{
									"scope": schema.StringAttribute{
										Description: "Scope of roles.",
										Optional:    true,
									},
									"roles": schema.SetAttribute{
										Description: "List of assigned roles.",
										ElementType: types.StringType,
										Optional:    true,
										Validators: []validator.Set{
											setvalidator.SizeAtLeast(1),
										},
									},
								},
							},
						},
					},
				},
				Validators: []validator.List{
					listvalidator.SizeAtMost(1),
				},
			},
			"ngs_access_policy": schema.ListNestedBlock{
				Description: "A set of access policy rules that defines the actions and resources that an identity can access within an NGA environment.",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"cloud_account_link": schema.StringAttribute{
							Description: "Full link to referenced cloud account.",
							Required:    true,
							Validators: []validator.String{
								validators.LinkValidator{},
							},
						},
						"subs": schema.Int32Attribute{
							Description: "Max number of subscriptions per connection. Default: -1",
							Optional:    true,
							Computed:    true,
							Default:     int32default.StaticInt32(-1),
						},
						"data": schema.Int32Attribute{
							Description: "Max number of bytes a connection can send. Default: -1",
							Optional:    true,
							Computed:    true,
							Default:     int32default.StaticInt32(-1),
						},
						"payload": schema.Int32Attribute{
							Description: "Max message payload. Default: -1",
							Optional:    true,
							Computed:    true,
							Default:     int32default.StaticInt32(-1),
						},
					},
					Blocks: map[string]schema.Block{
						"pub": ir.NgsPermissions("Pub Permission."),
						"sub": ir.NgsPermissions("Sub Permission."),
						"resp": schema.ListNestedBlock{
							Description: "Reponses.",
							NestedObject: schema.NestedBlockObject{
								Attributes: map[string]schema.Attribute{
									"max": schema.Int32Attribute{
										Description: "Number of responses allowed on the replyTo subject, -1 means no limit. Default: -1",
										Optional:    true,
										Computed:    true,
										Default:     int32default.StaticInt32(1),
									},
									"ttl": schema.StringAttribute{
										Description: "Deadline to send replies on the replyTo subject [#ms(millis) | #s(econds) | m(inutes) | h(ours)]. -1 means no restriction.",
										Optional:    true,
									},
								},
							},
							Validators: []validator.List{
								listvalidator.SizeAtMost(1),
							},
						},
					},
				},
				Validators: []validator.List{
					listvalidator.SizeAtMost(1),
				},
			},
			"network_resource": schema.SetNestedBlock{
				Description: "A network resource can be configured with: - A fully qualified domain name (FQDN) and ports. - An FQDN, resolver IP, and ports. - IP's and ports.",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							Description: "Name of the Network Resource.",
							Required:    true,
						},
						"agent_link": schema.StringAttribute{
							Description: "Full link to referenced Agent.",
							Optional:    true,
							Validators: []validator.String{
								validators.LinkValidator{},
							},
						},
						"ips": schema.SetAttribute{
							Description: "List of IP addresses.",
							ElementType: types.StringType,
							Optional:    true,
							Validators: []validator.Set{
								setvalidator.ConflictsWith(
									path.MatchRelative().AtParent().AtName("fqdn"),
								),
							},
						},
						"fqdn": schema.StringAttribute{
							Description: "Fully qualified domain name.",
							Optional:    true,
							Validators: []validator.String{
								stringvalidator.ConflictsWith(
									path.MatchRelative().AtParent().AtName("ips"),
								),
							},
						},
						"resolver_ip": schema.StringAttribute{
							Description: "Resolver IP.",
							Optional:    true,
						},
						"ports": schema.SetAttribute{
							Description: "Ports to expose.",
							ElementType: types.Int32Type,
							Required:    true,
						},
					},
				},
			},
			"native_network_resource": schema.SetNestedBlock{
				Description: "~> **NOTE** The configuration of a native network resource requires the assistance of Control Plane support.",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							Description: "Name of the Native Network Resource.",
							Required:    true,
						},
						"fqdn": schema.StringAttribute{
							Description: "Fully qualified domain name.",
							Required:    true,
						},
						"ports": schema.SetAttribute{
							Description: "Ports to expose. At least one port is required.",
							ElementType: types.Int32Type,
							Required:    true,
							Validators: []validator.Set{
								setvalidator.SizeAtLeast(1),
							},
						},
					},
					Blocks: map[string]schema.Block{
						"aws_private_link": schema.ListNestedBlock{
							Description: "A feature provided by AWS that enables private connectivity between private VPCs and compute running at Control Plane without traversing the public internet.",
							NestedObject: schema.NestedBlockObject{
								Attributes: map[string]schema.Attribute{
									"endpoint_service_name": schema.StringAttribute{
										Description: "Endpoint service name.",
										Required:    true,
									},
								},
								Validators: []validator.Object{
									objectvalidator.ConflictsWith(
										path.MatchRelative().AtParent().AtParent().AtName("gcp_service_connect"),
									),
								},
							},
							Validators: []validator.List{
								listvalidator.SizeAtMost(1),
							},
						},
						"gcp_service_connect": schema.ListNestedBlock{
							Description: "Capability provided by GCP that allows private communication between private VPC networks and compute running at Control Plane.",
							NestedObject: schema.NestedBlockObject{
								Attributes: map[string]schema.Attribute{
									"target_service": schema.StringAttribute{
										Description: "Target service name.",
										Required:    true,
									},
								},
								Validators: []validator.Object{
									objectvalidator.ConflictsWith(
										path.MatchRelative().AtParent().AtParent().AtName("aws_private_link"),
									),
								},
							},
							Validators: []validator.List{
								listvalidator.SizeAtMost(1),
							},
						},
					},
				},
			},
		},
	}
}

// Create creates the resource.
func (ir *IdentityResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	CreateGeneric(ctx, req, resp, ir.Operations)
}

// Read fetches the current state of the resource.
func (ir *IdentityResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	ReadGeneric(ctx, req, resp, ir.Operations)
}

// Update modifies the resource.
func (ir *IdentityResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	UpdateGeneric(ctx, req, resp, ir.Operations)
}

// Delete removes the resource.
func (ir *IdentityResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	DeleteGeneric(ctx, req, resp, ir.Operations)
}

/*** Schemas ***/

// NgsPermissions creates a nested list block schema using the provided description.
func (ir *IdentityResource) NgsPermissions(description string) schema.ListNestedBlock {
	return schema.ListNestedBlock{
		Description: description,
		NestedObject: schema.NestedBlockObject{
			Attributes: map[string]schema.Attribute{
				"allow": schema.SetAttribute{
					Description: "List of allow subjects.",
					ElementType: types.StringType,
					Optional:    true,
				},
				"deny": schema.SetAttribute{
					Description: "List of deny subjects.",
					ElementType: types.StringType,
					Optional:    true,
				},
			},
		},
		Validators: []validator.List{
			listvalidator.SizeAtMost(1),
		},
	}
}

/*** Resource Operator ***/

// IdentityResourceOperator is the operator for managing the state.
type IdentityResourceOperator struct {
	EntityOperator[IdentityResourceModel]
}

// NewAPIRequest creates a request payload from a state model.
func (iro *IdentityResourceOperator) NewAPIRequest(isUpdate bool) client.Identity {
	// Initialize a new request payload
	requestPayload := client.Identity{}

	// Populate Base fields from state
	iro.Plan.Fill(&requestPayload.Base, isUpdate)

	// Set specific attributes
	if isUpdate {
		requestPayload.AwsReplace = iro.buildAws(iro.Plan.AwsAccessPolicy)
		requestPayload.GcpReplace = iro.buildGcp(iro.Plan.GcpAccessPolicy)
		requestPayload.AzureReplace = iro.buildAzure(iro.Plan.AzureAccessPolicy)
		requestPayload.NgsReplace = iro.buildNgs(iro.Plan.NgsAccessPolicy)
	} else {
		requestPayload.Aws = iro.buildAws(iro.Plan.AwsAccessPolicy)
		requestPayload.Gcp = iro.buildGcp(iro.Plan.GcpAccessPolicy)
		requestPayload.Azure = iro.buildAzure(iro.Plan.AzureAccessPolicy)
		requestPayload.Ngs = iro.buildNgs(iro.Plan.NgsAccessPolicy)
	}

	requestPayload.NetworkResources = iro.buildNetworkResources(iro.Plan.NetworkResource)
	requestPayload.NativeNetworkResources = iro.buildNativeNetworkResources(iro.Plan.NativeNetworkResource)

	// Return constructed request payload
	return requestPayload
}

// MapResponseToState constructs the Terraform state model from the API response payload.
func (iro *IdentityResourceOperator) MapResponseToState(apiResp *client.Identity, isCreate bool) IdentityResourceModel {
	// Initialize empty state model
	state := IdentityResourceModel{}

	// Populate common fields from base resource data
	state.From(apiResp.Base)

	// Set specific attributes
	state.Gvc = types.StringPointerValue(BuildString(iro.Plan.Gvc))
	state.AwsAccessPolicy = iro.flattenAws(apiResp.Aws)
	state.GcpAccessPolicy = iro.flattenGcp(apiResp.Gcp)
	state.AzureAccessPolicy = iro.flattenAzure(apiResp.Azure)
	state.NgsAccessPolicy = iro.flattenNgs(apiResp.Ngs)
	state.NetworkResource = iro.flattenNetworkResources(apiResp.NetworkResources)
	state.NativeNetworkResource = iro.flattenNativeNetworkResources(apiResp.NativeNetworkResources)
	state.Status = iro.flattenStatus(apiResp.Status)

	// Return completed state model
	return state
}

// InvokeCreate invokes the Create API to create a new resource.
func (iro *IdentityResourceOperator) InvokeCreate(req client.Identity) (*client.Identity, int, error) {
	return iro.Client.CreateIdentity(req, iro.Plan.Gvc.ValueString())
}

// InvokeRead invokes the Get API to retrieve an existing resource by name.
func (iro *IdentityResourceOperator) InvokeRead(name string) (*client.Identity, int, error) {
	return iro.Client.GetIdentity(name, iro.Plan.Gvc.ValueString())
}

// InvokeUpdate invokes the Update API to update an existing resource.
func (iro *IdentityResourceOperator) InvokeUpdate(req client.Identity) (*client.Identity, int, error) {
	return iro.Client.UpdateIdentity(req, iro.Plan.Gvc.ValueString())
}

// InvokeDelete invokes the Delete API to delete a resource by name.
func (iro *IdentityResourceOperator) InvokeDelete(name string) error {
	return iro.Client.DeleteIdentity(name, iro.Plan.Gvc.ValueString())
}

// Builders //

// buildAws constructs a IdentityAws struct from the given Terraform state.
func (iro *IdentityResourceOperator) buildAws(state types.List) *client.IdentityAws {
	// Convert Terraform list into model blocks using generic helper
	blocks, ok := BuildList[models.AwsAccessPolicyModel](iro.Ctx, iro.Diags, state)

	// Return nil if conversion failed or list was empty
	if !ok {
		return nil
	}

	// Extract the very first block from the blocks slice
	block := blocks[0]

	// Construct and return the output
	return &client.IdentityAws{
		CloudAccountLink: BuildString(block.CloudAccountLink),
		PolicyRefs:       iro.BuildSetString(block.PolicyRefs),
		RoleName:         BuildString(block.RoleName),
		TrustPolicy:      iro.buildAwsTrustPolicy(block.TrustPolicy),
	}
}

// buildAwsTrustPolicy constructs a IdentityAwsTrustPolicy struct from the given Terraform state.
func (iro *IdentityResourceOperator) buildAwsTrustPolicy(state types.Set) *client.IdentityAwsTrustPolicy {
	// Convert Terraform set into model blocks using generic helper
	blocks, ok := BuildSet[models.AwsAccessPolicyTrustPolicyModel](iro.Ctx, iro.Diags, state)

	// Return nil if conversion failed or set was empty
	if !ok {
		return nil
	}

	// Extract the very first block from the blocks slice
	block := blocks[0]

	// Construct and return the output
	return &client.IdentityAwsTrustPolicy{
		Version:   BuildString(block.Version),
		Statement: iro.buildAwsTrustPolicyStatement(block.Statement),
	}
}

// buildAwsTrustPolicyStatement constructs a []map[string]interface{} struct from the given Terraform state.
func (iro *IdentityResourceOperator) buildAwsTrustPolicyStatement(state types.Set) *[]map[string]interface{} {
	// Exit early if set itself is null or unknown
	if state.IsNull() || state.IsUnknown() {
		return nil
	}

	// Prepare an intermediate slice to unmarshal Terraform values
	var intermediate []types.Map

	// Decode Terraform set elements into the intermediate slice
	iro.Diags.Append(state.ElementsAs(iro.Ctx, &intermediate, false)...)

	// Abort if any diagnostic errors occurred during decoding
	if iro.Diags.HasError() {
		return nil
	}

	// Build the output slice, preallocating for efficiency
	output := make([]map[string]interface{}, 0, len(intermediate))

	// Iterate and extract each known string value
	for _, elem := range intermediate {
		// Skip null or unknown entries
		if elem.IsNull() || elem.IsUnknown() {
			continue
		}

		// Add the element to the output slice
		output = append(output, *iro.BuildMapString(elem))
	}

	// Return a pointer to the output
	return &output
}

// buildGcp constructs a IdentityAws struct from the given Terraform state.
func (iro *IdentityResourceOperator) buildGcp(state types.List) *client.IdentityGcp {
	// Convert Terraform list into model blocks using generic helper
	blocks, ok := BuildList[models.GcpAccessPolicyModel](iro.Ctx, iro.Diags, state)

	// Return nil if conversion failed or list was empty
	if !ok {
		return nil
	}

	// Extract the very first block from the blocks slice
	block := blocks[0]

	// Construct and return the output
	return &client.IdentityGcp{
		CloudAccountLink: BuildString(block.CloudAccountLink),
		Scopes:           iro.buildGcpScopes(block.Scopes),
		ServiceAccount:   BuildString(block.ServiceAccount),
		Bindings:         iro.buildGcpBinding(block.Binding),
	}
}

// buildGcpScopes constructs a *[]string slice from the given Terraform state.
func (iro *IdentityResourceOperator) buildGcpScopes(state types.String) *[]string {
	// Build the state string into a golang string
	scopes := BuildString(state)

	// If input is nil or empty, return nil
	if scopes == nil {
		return nil
	}

	// Split by comma
	parts := strings.Split(*scopes, ",")

	// Trim spaces from each part
	for i := range parts {
		parts[i] = strings.TrimSpace(parts[i])
	}

	// Return the output
	return &parts
}

// buildGcpBinding constructs a []client.IdentityGcpBinding from the given Terraform state.
func (iro *IdentityResourceOperator) buildGcpBinding(state types.Set) *[]client.IdentityGcpBinding {
	// Convert Terraform set into model blocks using generic helper
	blocks, ok := BuildSet[models.GcpAccessPolicyBindingModel](iro.Ctx, iro.Diags, state)

	// Return nil if conversion failed or set was empty
	if !ok {
		return nil
	}

	// Prepare the output slice
	output := []client.IdentityGcpBinding{}

	// Iterate over each block and construct an output item
	for _, block := range blocks {
		// Construct the item
		item := client.IdentityGcpBinding{
			Resource: BuildString(block.Resource),
			Roles:    iro.BuildSetString(block.Roles),
		}

		// Add the item to the output slice
		output = append(output, item)
	}

	// Return a pointer to the output
	return &output
}

// buildAzure constructs a IdentityAzure from the given Terraform state.
func (iro *IdentityResourceOperator) buildAzure(state types.List) *client.IdentityAzure {
	// Convert Terraform list into model blocks using generic helper
	blocks, ok := BuildList[models.AzureAccessPolicyModel](iro.Ctx, iro.Diags, state)

	// Return nil if conversion failed or list was empty
	if !ok {
		return nil
	}

	// Take the first (and only) block
	block := blocks[0]

	// Construct and return the output
	return &client.IdentityAzure{
		CloudAccountLink: BuildString(block.CloudAccountLink),
		RoleAssignments:  iro.buildAzureRoleAssignments(block.RoleAssignment),
	}
}

// buildAzureRoleAssignments constructs a []client.IdentityAzureRoleAssignment from the given Terraform state.
func (iro *IdentityResourceOperator) buildAzureRoleAssignments(state types.Set) *[]client.IdentityAzureRoleAssignment {
	// Convert Terraform set into model blocks using generic helper
	blocks, ok := BuildSet[models.AzureAccessPolicyRoleAssignmentModel](iro.Ctx, iro.Diags, state)

	// Return nil if conversion failed or set was empty
	if !ok {
		return nil
	}

	// Prepare the output slice
	output := []client.IdentityAzureRoleAssignment{}

	// Iterate over each block and construct an output item
	for _, block := range blocks {
		// Construct the item
		item := client.IdentityAzureRoleAssignment{
			Scope: BuildString(block.Scope),
			Roles: iro.BuildSetString(block.Roles),
		}

		// Add the item to the output slice
		output = append(output, item)
	}

	// Return a pointer to the output
	return &output
}

// buildNgs constructs a IdentityNgs from the given Terraform state.
func (iro *IdentityResourceOperator) buildNgs(state types.List) *client.IdentityNgs {
	// Convert Terraform list into model blocks using generic helper
	blocks, ok := BuildList[models.NgsAccessPolicyModel](iro.Ctx, iro.Diags, state)

	// Return nil if conversion failed or list was empty
	if !ok {
		return nil
	}

	// Take the first (and only) block
	block := blocks[0]

	// Construct and return the output
	return &client.IdentityNgs{
		CloudAccountLink: BuildString(block.CloudAccountLink),
		Subs:             BuildInt(block.Subs),
		Data:             BuildInt(block.Data),
		Payload:          BuildInt(block.Payload),
		Pub:              iro.buildNgsPerm(block.Pub),
		Sub:              iro.buildNgsPerm(block.Sub),
		Resp:             iro.buildNgsResp(block.Resp),
	}
}

// buildNgsPerm constructs a IdentityNgs from the given Terraform state.
func (iro *IdentityResourceOperator) buildNgsPerm(state types.List) *client.IdentityNgsPerm {
	// Convert Terraform list into model blocks using generic helper
	blocks, ok := BuildList[models.NgsAccessPolicyPermissionModel](iro.Ctx, iro.Diags, state)

	// Return nil if conversion failed or list was empty
	if !ok {
		return nil
	}

	// Take the first (and only) block
	block := blocks[0]

	// Construct and return the output
	return &client.IdentityNgsPerm{
		Allow: iro.BuildSetString(block.Allow),
		Deny:  iro.BuildSetString(block.Deny),
	}
}

// buildNgsResp constructs a IdentityNgsResp from the given Terraform state.
func (iro *IdentityResourceOperator) buildNgsResp(state types.List) *client.IdentityNgsResp {
	// Convert Terraform list into model blocks using generic helper
	blocks, ok := BuildList[models.NgsAccessPolicyResponsesModel](iro.Ctx, iro.Diags, state)

	// Return nil if conversion failed or list was empty
	if !ok {
		return nil
	}

	// Take the first (and only) block
	block := blocks[0]

	// Construct and return the output
	return &client.IdentityNgsResp{
		Max: BuildInt(block.Max),
		TTL: BuildString(block.TTL),
	}
}

// buildNetworkResources constructs a []client.IdentityNetworkResource from the given Terraform state.
func (iro *IdentityResourceOperator) buildNetworkResources(state types.Set) *[]client.IdentityNetworkResource {
	// Convert Terraform set into model blocks using generic helper
	blocks, ok := BuildSet[models.NetworkResourceModel](iro.Ctx, iro.Diags, state)

	// Return nil if conversion failed or set was empty
	if !ok {
		return nil
	}

	// Prepare the output slice
	output := []client.IdentityNetworkResource{}

	// Iterate over each block and construct an output item
	for _, block := range blocks {
		// Construct the item
		item := client.IdentityNetworkResource{
			Name:       BuildString(block.Name),
			AgentLink:  BuildString(block.AgentLink),
			IPs:        iro.BuildSetString(block.IPs),
			FQDN:       BuildString(block.FQDN),
			ResolverIP: BuildString(block.ResolverIP),
			Ports:      iro.BuildSetInt(block.Ports),
		}

		// Add the item to the output slice
		output = append(output, item)
	}

	// Return a pointer to the output
	return &output
}

// buildNativeNetworkResources constructs a []client.IdentityNativeNetworkResource from the given Terraform state.
func (iro *IdentityResourceOperator) buildNativeNetworkResources(state types.Set) *[]client.IdentityNativeNetworkResource {
	// Convert Terraform set into model blocks using generic helper
	blocks, ok := BuildSet[models.NativeNetworkResourceModel](iro.Ctx, iro.Diags, state)

	// Return nil if conversion failed or set was empty
	if !ok {
		return nil
	}

	// Prepare the output slice
	output := []client.IdentityNativeNetworkResource{}

	// Iterate over each block and construct an output item
	for _, block := range blocks {
		// Construct the item
		item := client.IdentityNativeNetworkResource{
			Name:              BuildString(block.Name),
			FQDN:              BuildString(block.FQDN),
			Ports:             iro.BuildSetInt(block.Ports),
			AWSPrivateLink:    iro.buildAwsPrivateLink(block.AwsPrivateLink),
			GCPServiceConnect: iro.buildGcpServiceConnect(block.GcpServiceConnect),
		}

		// Add the item to the output slice
		output = append(output, item)
	}

	// Return a pointer to the output
	return &output
}

// buildAwsPrivateLink constructs a IdentityAwsPrivateLink from the given Terraform state.
func (iro *IdentityResourceOperator) buildAwsPrivateLink(state types.List) *client.IdentityAwsPrivateLink {
	// Convert Terraform list into model blocks using generic helper
	blocks, ok := BuildList[models.NativeNetworkResourceAwsPrivateLinkModel](iro.Ctx, iro.Diags, state)

	// Return nil if conversion failed or list was empty
	if !ok {
		return nil
	}

	// Take the first (and only) block
	block := blocks[0]

	// Construct and return the output
	return &client.IdentityAwsPrivateLink{
		EndpointServiceName: BuildString(block.EndpointServiceName),
	}
}

// buildGcpServiceConnect constructs a IdentityGcpServiceConnect from the given Terraform state.
func (iro *IdentityResourceOperator) buildGcpServiceConnect(state types.List) *client.IdentityGcpServiceConnect {
	// Convert Terraform list into model blocks using generic helper
	blocks, ok := BuildList[models.NativeNetworkResourceGcpServiceConnectModel](iro.Ctx, iro.Diags, state)

	// Return nil if conversion failed or list was empty
	if !ok {
		return nil
	}

	// Take the first (and only) block
	block := blocks[0]

	// Construct and return the output
	return &client.IdentityGcpServiceConnect{
		TargetService: BuildString(block.TargetService),
	}
}

// Flatteners //

// flattenAws transforms *client.IdentityAws into a Terraform types.List.
func (iro *IdentityResourceOperator) flattenAws(input *client.IdentityAws) types.List {
	// Get attribute types
	elementType := models.AwsAccessPolicyModel{}.AttributeTypes()

	// Check if the input is nil
	if input == nil {
		// Return a null list
		return types.ListNull(elementType)
	}

	// Build a single block
	block := models.AwsAccessPolicyModel{
		CloudAccountLink: types.StringPointerValue(input.CloudAccountLink),
		PolicyRefs:       FlattenSetString(input.PolicyRefs),
		RoleName:         types.StringPointerValue(input.RoleName),
		TrustPolicy:      iro.flattenAwsTrustPolicy(input.TrustPolicy),
	}

	// Return the successfully created types.List
	return FlattenList(iro.Ctx, iro.Diags, []models.AwsAccessPolicyModel{block})
}

// flattenAwsTrustPolicy transforms *client.IdentityAwsTrustPolicy into a Terraform types.Set.
func (iro *IdentityResourceOperator) flattenAwsTrustPolicy(input *client.IdentityAwsTrustPolicy) types.Set {
	// Get attribute types
	elementType := models.AwsAccessPolicyTrustPolicyModel{}.AttributeTypes()

	// Check if the input is nil
	if input == nil {
		// Return a null set
		return types.SetNull(elementType)
	}

	// Build a single block
	block := models.AwsAccessPolicyTrustPolicyModel{
		Version:   types.StringPointerValue(input.Version),
		Statement: iro.flattenAwsTrustPolicyStatement(input.Statement),
	}

	// Return the successfully created types.Set
	return FlattenSet(iro.Ctx, iro.Diags, []models.AwsAccessPolicyTrustPolicyModel{block})
}

// flattenAwsTrustPolicyStatement transforms *[]map[string]interface{} into a Terraform types.List.
func (iro *IdentityResourceOperator) flattenAwsTrustPolicyStatement(input *[]map[string]interface{}) types.Set {
	// Get attribute types
	elementType := types.MapType{ElemType: types.StringType}

	// Check if the input is nil
	if input == nil {
		// Return a null list
		return types.SetNull(elementType)
	}

	// Define the blocks slice
	var blocks []types.Map

	// Iterate over the slice and construct the blocks
	for _, item := range *input {
		// Construct a block
		block := FlattenMapString(&item)

		// Append the constructed block to the blocks slice
		blocks = append(blocks, block)
	}

	// Convert the slice of blocks into a Terraform list while collecting diagnostics
	l, d := types.SetValueFrom(iro.Ctx, elementType, blocks)

	// Merge any diagnostics from the conversion into the main diagnostics
	iro.Diags.Append(d...)

	// If the conversion produced errors, return a null list
	if d.HasError() {
		return types.SetNull(elementType)
	}

	// Return the successfully created types.List
	return l
}

// flattenGcp transforms *client.IdentityGcp into a Terraform types.List.
func (iro *IdentityResourceOperator) flattenGcp(input *client.IdentityGcp) types.List {
	// Get attribute types
	elementType := models.GcpAccessPolicyModel{}.AttributeTypes()

	// Check if the input is nil
	if input == nil {
		// Return a null list
		return types.ListNull(elementType)
	}

	// Declare a variable to hold planned scopes of the GCP block
	var plannedScopes *string

	// Build the planned GCP
	plannedGcp, ok := BuildList[models.GcpAccessPolicyModel](iro.Ctx, iro.Diags, iro.Plan.GcpAccessPolicy)

	// Extract the planned scopes from the planned GCP block
	if ok && len(plannedGcp) != 0 {
		plannedScopes = BuildString(plannedGcp[0].Scopes)
	}

	// Build a single block
	block := models.GcpAccessPolicyModel{
		CloudAccountLink: types.StringPointerValue(input.CloudAccountLink),
		Scopes:           iro.flattenGcpScopes(plannedScopes, input.Scopes),
		ServiceAccount:   types.StringPointerValue(input.ServiceAccount),
		Binding:          iro.flattenGcpBinding(input.Bindings),
	}

	// Return the successfully created types.List
	return FlattenList(iro.Ctx, iro.Diags, []models.GcpAccessPolicyModel{block})
}

// FlattenGcpScopes converts a *[]string into a Terraform types.String, preserving comma+whitespace separators from a prior plan string when available; new items use "," with no spaces.
func (iro *IdentityResourceOperator) flattenGcpScopes(state *string, input *[]string) types.String {
	// Return null when the slice pointer is nil (attribute absent)
	if input == nil {
		return types.StringNull()
	}

	// Return empty string when the slice is present but has no elements
	if len(*input) == 0 {
		return types.StringValue("")
	}

	// Prepare a container for exact separators ("," plus any following whitespace)
	var seps []string

	// Extract separators only if we have a prior state string
	if state != nil {
		// Scan the state string for commas
		for i := 0; i < len(*state); i++ {
			// Check for a comma boundary
			if (*state)[i] == ',' {
				// Advance past any whitespace immediately after the comma
				j := i + 1
				for j < len(*state) && unicode.IsSpace(rune((*state)[j])) {
					j++
				}

				// Capture the exact comma+whitespace sequence for reuse
				seps = append(seps, (*state)[i:j])
			}
		}
	}

	// Use a strings.Builder for efficient concatenation
	var b strings.Builder

	// Write the first element as-is
	b.WriteString((*input)[0])

	// Append remaining elements, preserving separators when available
	for i := 1; i < len(*input); i++ {
		// Reuse the saved separator if one exists at this position
		if i-1 < len(seps) {
			b.WriteString(seps[i-1])
		} else {
			// Fallback to a plain comma with no spaces for new positions
			b.WriteString(",")
		}

		// Write the current item
		b.WriteString((*input)[i])
	}

	// Return the final Terraform string value
	return types.StringValue(b.String())
}

// flattenGcpBinding transforms *[]client.IdentityGcpBinding into a Terraform types.Set.
func (iro *IdentityResourceOperator) flattenGcpBinding(input *[]client.IdentityGcpBinding) types.Set {
	// Get attribute types
	elementType := models.GcpAccessPolicyBindingModel{}.AttributeTypes()

	// Check if the input is nil
	if input == nil {
		// Return a null set
		return types.SetNull(elementType)
	}

	// Define the blocks slice
	var blocks []models.GcpAccessPolicyBindingModel

	// Iterate over the slice and construct the blocks
	for _, item := range *input {
		// Construct a block
		block := models.GcpAccessPolicyBindingModel{
			Resource: types.StringPointerValue(item.Resource),
			Roles:    FlattenSetString(item.Roles),
		}

		// Append the constructed block to the blocks slice
		blocks = append(blocks, block)
	}

	// Return the successfully created types.Set
	return FlattenSet(iro.Ctx, iro.Diags, blocks)
}

// flattenAzure transforms *client.IdentityAzure into a Terraform types.List.
func (iro *IdentityResourceOperator) flattenAzure(input *client.IdentityAzure) types.List {
	// Get attribute types
	elementType := models.AzureAccessPolicyModel{}.AttributeTypes()

	// Check if the input is nil
	if input == nil {
		// Return a null list
		return types.ListNull(elementType)
	}

	// Build a single block
	block := models.AzureAccessPolicyModel{
		CloudAccountLink: types.StringPointerValue(input.CloudAccountLink),
		RoleAssignment:   iro.flattenAzureRoleAssignment(input.RoleAssignments),
	}

	// Return the successfully created types.List
	return FlattenList(iro.Ctx, iro.Diags, []models.AzureAccessPolicyModel{block})
}

// flattenAzureRoleAssignment transforms *[]client.IdentityAzureRoleAssignment into a Terraform types.Set.
func (iro *IdentityResourceOperator) flattenAzureRoleAssignment(input *[]client.IdentityAzureRoleAssignment) types.Set {
	// Get attribute types
	elementType := models.AzureAccessPolicyRoleAssignmentModel{}.AttributeTypes()

	// Check if the input is nil
	if input == nil {
		// Return a null set
		return types.SetNull(elementType)
	}

	// Define the blocks slice
	var blocks []models.AzureAccessPolicyRoleAssignmentModel

	// Iterate over the slice and construct the blocks
	for _, item := range *input {
		// Construct a block
		block := models.AzureAccessPolicyRoleAssignmentModel{
			Scope: types.StringPointerValue(item.Scope),
			Roles: FlattenSetString(item.Roles),
		}

		// Append the constructed block to the blocks slice
		blocks = append(blocks, block)
	}

	// Return the successfully created types.Set
	return FlattenSet(iro.Ctx, iro.Diags, blocks)
}

// flattenNgs transforms *client.IdentityNgs into a Terraform types.List.
func (iro *IdentityResourceOperator) flattenNgs(input *client.IdentityNgs) types.List {
	// Get attribute types
	elementType := models.NgsAccessPolicyModel{}.AttributeTypes()

	// Check if the input is nil
	if input == nil {
		// Return a null list
		return types.ListNull(elementType)
	}

	// Build a single block
	block := models.NgsAccessPolicyModel{
		CloudAccountLink: types.StringPointerValue(input.CloudAccountLink),
		Subs:             FlattenInt(input.Subs),
		Data:             FlattenInt(input.Data),
		Payload:          FlattenInt(input.Payload),
		Pub:              iro.flattenNgsPerm(input.Pub),
		Sub:              iro.flattenNgsPerm(input.Sub),
		Resp:             iro.flattenNgsResp(input.Resp),
	}

	// Return the successfully created types.List
	return FlattenList(iro.Ctx, iro.Diags, []models.NgsAccessPolicyModel{block})
}

// flattenNgsPerm transforms *client.IdentityNgsPerm into a Terraform types.List.
func (iro *IdentityResourceOperator) flattenNgsPerm(input *client.IdentityNgsPerm) types.List {
	// Get attribute types
	elementType := models.NgsAccessPolicyPermissionModel{}.AttributeTypes()

	// Check if the input is nil
	if input == nil {
		// Return a null list
		return types.ListNull(elementType)
	}

	// Build a single block
	block := models.NgsAccessPolicyPermissionModel{
		Allow: FlattenSetString(input.Allow),
		Deny:  FlattenSetString(input.Deny),
	}

	// Return the successfully created types.List
	return FlattenList(iro.Ctx, iro.Diags, []models.NgsAccessPolicyPermissionModel{block})
}

// flattenNgsResp transforms *client.IdentityNgsResp into a Terraform types.List.
func (iro *IdentityResourceOperator) flattenNgsResp(input *client.IdentityNgsResp) types.List {
	// Get attribute types
	elementType := models.NgsAccessPolicyResponsesModel{}.AttributeTypes()

	// Check if the input is nil
	if input == nil {
		// Return a null list
		return types.ListNull(elementType)
	}

	// Build a single block
	block := models.NgsAccessPolicyResponsesModel{
		Max: FlattenInt(input.Max),
		TTL: types.StringPointerValue(input.TTL),
	}

	// Return the successfully created types.List
	return FlattenList(iro.Ctx, iro.Diags, []models.NgsAccessPolicyResponsesModel{block})
}

// flattenNetworkResources transforms *[]client.IdentityNetworkResource into a Terraform types.Set.
func (iro *IdentityResourceOperator) flattenNetworkResources(input *[]client.IdentityNetworkResource) types.Set {
	// Get attribute types
	elementType := models.NetworkResourceModel{}.AttributeTypes()

	// Check if the input is nil
	if input == nil {
		// Return a null set
		return types.SetNull(elementType)
	}

	// Define the blocks slice
	var blocks []models.NetworkResourceModel

	// Iterate over the slice and construct the blocks
	for _, item := range *input {
		// Construct a block
		block := models.NetworkResourceModel{
			Name:       types.StringPointerValue(item.Name),
			AgentLink:  types.StringPointerValue(item.AgentLink),
			IPs:        FlattenSetString(item.IPs),
			FQDN:       types.StringPointerValue(item.FQDN),
			ResolverIP: types.StringPointerValue(item.ResolverIP),
			Ports:      FlattenSetInt(item.Ports),
		}

		// Append the constructed block to the blocks slice
		blocks = append(blocks, block)
	}

	// Return the successfully created types.Set
	return FlattenSet(iro.Ctx, iro.Diags, blocks)
}

// flattenNativeNetworkResources transforms *[]client.IdentityNativeNetworkResource into a Terraform types.Set.
func (iro *IdentityResourceOperator) flattenNativeNetworkResources(input *[]client.IdentityNativeNetworkResource) types.Set {
	// Get attribute types
	elementType := models.NativeNetworkResourceModel{}.AttributeTypes()

	// Check if the input is nil
	if input == nil {
		// Return a null set
		return types.SetNull(elementType)
	}

	// Define the blocks slice
	var blocks []models.NativeNetworkResourceModel

	// Iterate over the slice and construct the blocks
	for _, item := range *input {
		// Construct a block
		block := models.NativeNetworkResourceModel{
			Name:              types.StringPointerValue(item.Name),
			FQDN:              types.StringPointerValue(item.FQDN),
			Ports:             FlattenSetInt(item.Ports),
			AwsPrivateLink:    iro.flattenAwsPrivateLink(item.AWSPrivateLink),
			GcpServiceConnect: iro.flattenGcpServiceConnect(item.GCPServiceConnect),
		}

		// Append the constructed block to the blocks slice
		blocks = append(blocks, block)
	}

	// Return the successfully created types.Set
	return FlattenSet(iro.Ctx, iro.Diags, blocks)
}

// flattenAwsPrivateLink transforms *client.IdentityAwsPrivateLink into a Terraform types.List.
func (iro *IdentityResourceOperator) flattenAwsPrivateLink(input *client.IdentityAwsPrivateLink) types.List {
	// Get attribute types
	elementType := models.NativeNetworkResourceAwsPrivateLinkModel{}.AttributeTypes()

	// Check if the input is nil
	if input == nil {
		// Return a null list
		return types.ListNull(elementType)
	}

	// Build a single block
	block := models.NativeNetworkResourceAwsPrivateLinkModel{
		EndpointServiceName: types.StringPointerValue(input.EndpointServiceName),
	}

	// Return the successfully created types.List
	return FlattenList(iro.Ctx, iro.Diags, []models.NativeNetworkResourceAwsPrivateLinkModel{block})
}

// flattenGcpServiceConnect transforms *client.IdentityGcpServiceConnect into a Terraform types.List.
func (iro *IdentityResourceOperator) flattenGcpServiceConnect(input *client.IdentityGcpServiceConnect) types.List {
	// Get attribute types
	elementType := models.NativeNetworkResourceGcpServiceConnectModel{}.AttributeTypes()

	// Check if the input is nil
	if input == nil {
		// Return a null list
		return types.ListNull(elementType)
	}

	// Build a single block
	block := models.NativeNetworkResourceGcpServiceConnectModel{
		TargetService: types.StringPointerValue(input.TargetService),
	}

	// Return the successfully created types.List
	return FlattenList(iro.Ctx, iro.Diags, []models.NativeNetworkResourceGcpServiceConnectModel{block})
}

// flattenStatus transforms *client.IdentityGcpServiceConnect into a Terraform types.Map.
func (iro *IdentityResourceOperator) flattenStatus(input *client.IdentityStatus) types.Map {
	// Check if the input is nil
	if input == nil {
		// Return a null list
		return types.MapNull(types.StringType)
	}

	// Build a single block
	statusMap := map[string]interface{}{
		"objectName": input.ObjectName,
	}

	// Return the successfully created types.List
	return FlattenMapString(&statusMap)
}
