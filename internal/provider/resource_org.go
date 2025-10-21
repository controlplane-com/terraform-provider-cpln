package cpln

import (
	"context"
	"fmt"

	client "github.com/controlplane-com/terraform-provider-cpln/internal/provider/client"
	models "github.com/controlplane-com/terraform-provider-cpln/internal/provider/models/org"
	"github.com/hashicorp/terraform-plugin-framework-validators/int32validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int32default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure resource implements required interfaces.
var (
	_ resource.Resource                = &OrgResource{}
	_ resource.ResourceWithImportState = &OrgResource{}
)

/*** Resource Model ***/

// OrgResourceModel holds the Terraform state for the resource.
type OrgResourceModel struct {
	EntityBaseModel
	AccountId             types.String `tfsdk:"account_id"`
	Invitees              types.Set    `tfsdk:"invitees"`
	SessionTimeoutSeconds types.Int32  `tfsdk:"session_timeout_seconds"`
	AuthConfig            types.List   `tfsdk:"auth_config"`
	Observability         types.List   `tfsdk:"observability"`
	Security              types.List   `tfsdk:"security"`
	Status                types.List   `tfsdk:"status"`
}

/*** Resource Configuration ***/

// OrgResource is the resource implementation.
type OrgResource struct {
	EntityBase
	Operations EntityOperations[OrgResourceModel, client.Org]
}

// NewOrgResource returns a new instance of the resource implementation.
func NewOrgResource() resource.Resource {
	// Initialize a new OrgResource struct
	resource := OrgResource{}

	// Mark the Name field as computed
	resource.IsNameComputed = true

	// Return a pointer to the new resource instance
	return &resource
}

// Configure configures the resource before use.
func (or *OrgResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	or.EntityBaseConfigure(ctx, req.ProviderData, &resp.Diagnostics)
	or.Operations = NewEntityOperations(or.client, &OrgResourceOperator{})
}

// ImportState sets up the import operation to map the imported ID to the "id" attribute in the state.
func (or *OrgResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// Metadata provides the resource type name.
func (or *OrgResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = "cpln_org"
}

// Schema defines the schema for the resource.
func (or *OrgResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: MergeAttributes(or.EntityBaseAttributes("Organization"), map[string]schema.Attribute{
			"account_id": schema.StringAttribute{
				Description: "The associated account ID that will be used when creating the org. Only used on org creation. The account ID can be obtained from the `Org Management & Billing` page.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"invitees": schema.SetAttribute{
				Description: "When an org is created, the list of email addresses which will receive an invitation to join the org and be assigned to the `superusers` group. The user account used when creating the org will be included in this list.",
				ElementType: types.StringType,
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
			},
			"session_timeout_seconds": schema.Int32Attribute{
				Description: "The idle time (in seconds) in which the console UI will automatically sign-out the user. Default: 900 (15 minutes)",
				Optional:    true,
				Computed:    true,
				Default:     int32default.StaticInt32(900),
			},
			"status": schema.ListNestedAttribute{
				Description: "Status of the org.",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"account_link": schema.StringAttribute{
							Description: "The link of the account the org belongs to.",
							Computed:    true,
						},
						"active": schema.BoolAttribute{
							Description: "Indicates whether the org is active or not.",
							Computed:    true,
						},
						"endpoint_prefix": schema.StringAttribute{
							Description: "",
							Computed:    true,
						},
					},
				},
			},
		}),
		Blocks: map[string]schema.Block{
			"auth_config": schema.ListNestedBlock{
				Description: "The configuration settings and parameters related to authentication within the org.",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"domain_auto_members": schema.SetAttribute{
							Description: "List of domains which will auto-provision users when authenticating using SAML.",
							ElementType: types.StringType,
							Required:    true,
						},
						"saml_only": schema.BoolAttribute{
							Description: "Enforce SAML only authentication.",
							Optional:    true,
							Computed:    true,
							Default:     booldefault.StaticBool(false),
						},
					},
				},
				Validators: []validator.List{
					listvalidator.SizeAtMost(1),
				},
			},
			"observability": schema.ListNestedBlock{
				Description: "The retention period (in days) for logs, metrics, and traces. Charges apply for storage beyond the 30 day default.",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"logs_retention_days": schema.Int32Attribute{
							Description: "Log retention days. Default: 30",
							Optional:    true,
							Computed:    true,
							Default:     int32default.StaticInt32(30),
							Validators: []validator.Int32{
								int32validator.AtLeast(0),
							},
						},
						"metrics_retention_days": schema.Int32Attribute{
							Description: "Metrics retention days. Default: 30",
							Optional:    true,
							Computed:    true,
							Default:     int32default.StaticInt32(30),
							Validators: []validator.Int32{
								int32validator.AtLeast(0),
							},
						},
						"traces_retention_days": schema.Int32Attribute{
							Description: "Traces retention days. Default: 30",
							Optional:    true,
							Computed:    true,
							Default:     int32default.StaticInt32(30),
							Validators: []validator.Int32{
								int32validator.AtLeast(0),
							},
						},
						"default_alert_emails": schema.SetAttribute{
							Description: "These emails are configured as alert recipients in Grafana when the 'grafana-default-email' contact delivery type is 'Email'.",
							ElementType: types.StringType,
							Optional:    true,
							Computed:    true,
							Default:     setdefault.StaticValue(types.SetValueMust(types.StringType, []attr.Value{})),
						},
					},
				},
				Validators: []validator.List{
					listvalidator.SizeAtMost(1),
				},
			},
			"security": schema.ListNestedBlock{
				Description: "",
				NestedObject: schema.NestedBlockObject{
					Blocks: map[string]schema.Block{
						"threat_detection": schema.ListNestedBlock{
							Description: "",
							NestedObject: schema.NestedBlockObject{
								Attributes: map[string]schema.Attribute{
									"enabled": schema.BoolAttribute{
										Description: "Indicates whether threat detection should be forwarded or not.",
										Required:    true,
									},
									"minimum_severity": schema.StringAttribute{
										Description: "Any threats with this severity and more severe will be sent. Others will be ignored. Valid values: `warning`, `error`, or `critical`.",
										Optional:    true,
									},
								},
								Blocks: map[string]schema.Block{
									"syslog": schema.ListNestedBlock{
										Description: "Configuration for syslog forwarding.",
										NestedObject: schema.NestedBlockObject{
											Attributes: map[string]schema.Attribute{
												"transport": schema.StringAttribute{
													Description: "The transport-layer protocol to send the syslog messages over. If TCP is chosen, messages will be sent with TLS. Default: `tcp`.",
													Optional:    true,
													Computed:    true,
													Default:     stringdefault.StaticString("tcp"),
												},
												"host": schema.StringAttribute{
													Description: "The hostname to send syslog messages to.",
													Required:    true,
												},
												"port": schema.Int32Attribute{
													Description: "The port to send syslog messages to.",
													Required:    true,
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
					},
				},
				Validators: []validator.List{
					listvalidator.SizeAtMost(1),
				},
			},
		},
	}
}

// ModifyPlan modifies the plan for the resource.
func (or *OrgResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	// If this is a destroy plan, leave everything null and return immediately
	if req.Plan.Raw.IsNull() {
		return
	}

	// Declare variable to store desired resource plan
	var plan OrgResourceModel

	// Populate plan variable from request and capture diagnostics
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)

	// Abort if any diagnostics errors occurred
	if resp.Diagnostics.HasError() {
		// Exit early on error
		return
	}

	// Modify autoscaling in options if specified
	if plan.Description.IsNull() || plan.Description.IsUnknown() {
		plan.Description = types.StringValue(or.client.Org)
	}

	// Persist new plan into Terraform
	resp.Diagnostics.Append(resp.Plan.Set(ctx, &plan)...)
}

// Create creates the resource.
func (or *OrgResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	CreateGeneric(ctx, req, resp, or.Operations)
}

// Read fetches the current state of the resource.
func (or *OrgResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	ReadGeneric(ctx, req, resp, or.Operations)
}

// Update modifies the resource.
func (or *OrgResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	UpdateGeneric(ctx, req, resp, or.Operations)
}

// Delete removes the resource.
func (or *OrgResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	DeleteGeneric(ctx, req, resp, or.Operations)
}

/*** Resource Operator ***/

// OrgResourceOperator is the operator for managing the state.
type OrgResourceOperator struct {
	EntityOperator[OrgResourceModel]
}

// NewAPIRequest creates a request payload from a state model.
func (oro *OrgResourceOperator) NewAPIRequest(isUpdate bool) client.Org {
	// Initialize a new request payload
	requestPayload := client.Org{}

	// Populate Base fields from state
	oro.Plan.Fill(&requestPayload.Base, isUpdate)

	// Initialize the Org spec struct
	var spec *client.OrgSpec = &client.OrgSpec{}

	// Map planned state attributes to the API struct
	if isUpdate {
		requestPayload.SpecReplace = spec
	} else {
		requestPayload.Spec = spec
	}

	// Set specific attributes
	spec.SessionTimeoutSeconds = BuildInt(oro.Plan.SessionTimeoutSeconds)
	spec.AuthConfig = oro.buildAuthConfig(oro.Plan.AuthConfig)
	spec.Observability = oro.buildObservability(oro.Plan.Observability)
	spec.Security = oro.buildSecurity(oro.Plan.Security)

	// Return constructed request payload
	return requestPayload
}

// MapResponseToState constructs the Terraform state model from the API response payload.
func (oro *OrgResourceOperator) MapResponseToState(apiResp *client.Org, isCreate bool) OrgResourceModel {
	// Initialize empty state model
	state := OrgResourceModel{}

	// Populate common fields from base resource data
	state.From(apiResp.Base)

	// On create operation, include the account id and invitees in teh state
	if isCreate {
		state.AccountId = types.StringPointerValue(BuildString(oro.Plan.AccountId))
		state.Invitees = FlattenSetString(oro.BuildSetString(oro.Plan.Invitees))
	}

	// Set specific attributes
	state.SessionTimeoutSeconds = FlattenInt(apiResp.Spec.SessionTimeoutSeconds)
	state.AuthConfig = oro.flattenAuthConfig(apiResp.Spec.AuthConfig)
	state.Observability = oro.flattenObservability(apiResp.Spec.Observability)
	state.Security = oro.flattenSecurity(apiResp.Spec.Security)
	state.Status = oro.flattenStatus(apiResp.Status)

	// Return completed state model
	return state
}

// InvokeCreate invokes the Create API to create a new resource.
func (oro *OrgResourceOperator) InvokeCreate(req client.Org) (*client.Org, int, error) {
	// Attempt to fetch the org
	currentOrg, code, err := oro.Client.GetOrg()

	// In case there is any error, attempt to create the org
	if err != nil {
		// Convert the planned state AccountId to a string pointer
		accountId := BuildString(oro.Plan.AccountId)

		// Build a set of invitee strings from the planned state
		invitees := oro.BuildSetString(oro.Plan.Invitees)

		// Ensure accountId and invitees are present before proceeding into the creation
		if accountId != nil && *accountId != "" && invitees != nil && len(*invitees) > 0 {
			// Initialize a new Org struct for creation
			org := client.Org{}

			// Set the org name to the client’s configured Org name
			org.Name = req.Name

			// Set the org description from the planned state
			org.Description = req.Description

			// Build and set tags from the planned state
			org.Tags = req.Tags

			// Prepare the CreateOrgRequest payload with org and invitees
			createOrgRequest := client.CreateOrgRequest{
				Org:      &org,
				Invitees: invitees,
			}

			// Initialize a variable to capture the response code
			responseCode := 0

			// Make the request to create the org
			currentOrg, responseCode, err = oro.Client.CreateOrg(*accountId, createOrgRequest)

			// Handle any errors from the create request
			if err != nil {
				// If a 409 Conflict occurs, the org already exists; set currentOrg accordingly
				if responseCode == 409 {
					currentOrg = &client.Org{}
					currentOrg.Name = req.Name
				} else {
					// Return on other errors preventing org creation
					return nil, responseCode, fmt.Errorf("org %s cannot be created. Error: %s", *org.Name, err)
				}
			}
		} else {
			// Return early if required accountId or invitees are missing
			return nil, code, err
		}
	}

	// Copy spec from the request object to the currentOrg object
	currentOrg.Description = req.Description
	currentOrg.Tags = req.Tags
	currentOrg.Spec = nil
	currentOrg.SpecReplace = req.Spec

	// Update the org
	updateOrg, code, err := oro.Client.UpdateOrg(*currentOrg)

	// Handle any errors from the update operation
	if err != nil {
		return nil, code, err
	}

	// Return the fetched or newly created org, along with code and error
	return updateOrg, code, err
}

// InvokeRead invokes the Get API to retrieve an existing resource by name.
func (oro *OrgResourceOperator) InvokeRead(name string) (*client.Org, int, error) {
	return oro.Client.GetOrg()
}

// InvokeUpdate invokes the Update API to update an existing resource.
func (oro *OrgResourceOperator) InvokeUpdate(req client.Org) (*client.Org, int, error) {
	return oro.Client.UpdateOrg(req)
}

// InvokeDelete invokes the Delete API to delete a resource by name.
func (oro *OrgResourceOperator) InvokeDelete(name string) error {
	// Initialize an Org struct with base and spec replacement fields
	org := client.Org{
		Base: client.Base{
			Name:        BuildString(oro.Plan.Name),
			Description: BuildString(oro.Plan.Name),
			TagsReplace: &map[string]any{},
		},
		SpecReplace: &client.OrgSpec{
			SessionTimeoutSeconds: IntPointer(900),
			Observability: &client.Observability{
				LogsRetentionDays:    IntPointer(30),
				MetricsRetentionDays: IntPointer(30),
				TracesRetentionDays:  IntPointer(30),
			},
			AuthConfig: nil,
			Security:   nil,
		},
	}

	// Call UpdateOrg API to apply changes (deletion is represented by updating to default state)
	_, _, err := oro.Client.UpdateOrg(org)

	// If an error occurred during the API call, return it
	if err != nil {
		return err
	}

	// Return nil when deletion (update) is successful
	return nil
}

// Builders //

// buildAuthConfig constructs a AuthConfig from the given Terraform state.
func (oro *OrgResourceOperator) buildAuthConfig(state types.List) *client.AuthConfig {
	// Convert Terraform list into model blocks using generic helper
	blocks, ok := BuildList[models.AuthConfigModel](oro.Ctx, oro.Diags, state)

	// Return nil if conversion failed or list was empty
	if !ok {
		return nil
	}

	// Take the first (and only) block
	block := blocks[0]

	// Construct and return the output
	return &client.AuthConfig{
		DomainAutoMembers: oro.BuildSetString(block.DomainAutoMembers),
		SamlOnly:          BuildBool(block.SamlOnly),
	}
}

// buildObservability constructs a Observability from the given Terraform state.
func (oro *OrgResourceOperator) buildObservability(state types.List) *client.Observability {
	// Convert Terraform list into model blocks using generic helper
	blocks, ok := BuildList[models.ObservabilityModel](oro.Ctx, oro.Diags, state)

	// Return nil if conversion failed or list was empty
	if !ok {
		return nil
	}

	// Take the first (and only) block
	block := blocks[0]

	// Construct and return the output
	return &client.Observability{
		LogsRetentionDays:    BuildInt(block.LogsRetentionDays),
		MetricsRetentionDays: BuildInt(block.MetricsRetentionDays),
		TracesRetentionDays:  BuildInt(block.TracesRetentionDays),
		DefaultAlertEmails:   oro.BuildSetString(block.DefaultAlertEmails),
	}
}

// buildSecurity constructs a OrgSecurity from the given Terraform state.
func (oro *OrgResourceOperator) buildSecurity(state types.List) *client.OrgSecurity {
	// Convert Terraform list into model blocks using generic helper
	blocks, ok := BuildList[models.SecurityModel](oro.Ctx, oro.Diags, state)

	// Return nil if conversion failed or list was empty
	if !ok {
		return nil
	}

	// Take the first (and only) block
	block := blocks[0]

	// Construct and return the output
	return &client.OrgSecurity{
		ThreatDetection: oro.buildSecurityThreatDetection(block.ThreatDetection),
	}
}

// buildSecurityThreatDetection constructs a OrgThreatDetection from the given Terraform state.
func (oro *OrgResourceOperator) buildSecurityThreatDetection(state types.List) *client.OrgThreatDetection {
	// Convert Terraform list into model blocks using generic helper
	blocks, ok := BuildList[models.SecurityThreatDetectionModel](oro.Ctx, oro.Diags, state)

	// Return nil if conversion failed or list was empty
	if !ok {
		return nil
	}

	// Take the first (and only) block
	block := blocks[0]

	// Construct and return the output
	return &client.OrgThreatDetection{
		Enabled:         BuildBool(block.Enabled),
		MinimumSeverity: BuildString(block.MinimumSeverity),
		Syslog:          oro.buildSecurityThreatDetectionSyslog(block.Syslog),
	}
}

// buildSecurityThreatDetectionSyslog constructs a OrgThreatDetectionSyslog from the given Terraform state.
func (oro *OrgResourceOperator) buildSecurityThreatDetectionSyslog(state types.List) *client.OrgThreatDetectionSyslog {
	// Convert Terraform list into model blocks using generic helper
	blocks, ok := BuildList[models.SecurityThreatDetectionSyslogModel](oro.Ctx, oro.Diags, state)

	// Return nil if conversion failed or list was empty
	if !ok {
		return nil
	}

	// Take the first (and only) block
	block := blocks[0]

	// Construct and return the output
	return &client.OrgThreatDetectionSyslog{
		Transport: BuildString(block.Transport),
		Host:      BuildString(block.Host),
		Port:      BuildInt(block.Port),
	}
}

// Flatteners //

// flattenAuthConfig transforms *client.AuthConfig into a types.List.
func (oro *OrgResourceOperator) flattenAuthConfig(input *client.AuthConfig) types.List {
	// Get attribute types
	elementType := models.AuthConfigModel{}.AttributeTypes()

	// Check if the input is nil
	if input == nil {
		// Return a null list
		return types.ListNull(elementType)
	}

	// Build a single block
	block := models.AuthConfigModel{
		DomainAutoMembers: FlattenSetString(input.DomainAutoMembers),
		SamlOnly:          types.BoolPointerValue(input.SamlOnly),
	}

	// Return the successfully created types.List
	return FlattenList(oro.Ctx, oro.Diags, []models.AuthConfigModel{block})
}

// flattenObservability transforms *client.Observability into a types.List.
func (oro *OrgResourceOperator) flattenObservability(input *client.Observability) types.List {
	// Get attribute types
	elementType := models.ObservabilityModel{}.AttributeTypes()

	// Check if the input is nil
	if input == nil {
		// Return a null list
		return types.ListNull(elementType)
	}

	// Build a single block
	block := models.ObservabilityModel{
		LogsRetentionDays:    FlattenInt(input.LogsRetentionDays),
		MetricsRetentionDays: FlattenInt(input.MetricsRetentionDays),
		TracesRetentionDays:  FlattenInt(input.TracesRetentionDays),
		DefaultAlertEmails:   FlattenSetString(input.DefaultAlertEmails),
	}

	// Return the successfully created types.List
	return FlattenList(oro.Ctx, oro.Diags, []models.ObservabilityModel{block})
}

// flattenSecurity transforms *client.OrgSecurity into a types.List.
func (oro *OrgResourceOperator) flattenSecurity(input *client.OrgSecurity) types.List {
	// Get attribute types
	elementType := models.SecurityModel{}.AttributeTypes()

	// Check if the input is nil
	if input == nil {
		// Return a null list
		return types.ListNull(elementType)
	}

	// Build a single block
	block := models.SecurityModel{
		ThreatDetection: oro.flattenSecurityThreatDetection(input.ThreatDetection),
	}

	// Return the successfully created types.List
	return FlattenList(oro.Ctx, oro.Diags, []models.SecurityModel{block})
}

// flattenSecurityThreatDetection transforms *client.OrgThreatDetection into a types.List.
func (oro *OrgResourceOperator) flattenSecurityThreatDetection(input *client.OrgThreatDetection) types.List {
	// Get attribute types
	elementType := models.SecurityThreatDetectionModel{}.AttributeTypes()

	// Check if the input is nil
	if input == nil {
		// Return a null list
		return types.ListNull(elementType)
	}

	// Build a single block
	block := models.SecurityThreatDetectionModel{
		Enabled:         types.BoolPointerValue(input.Enabled),
		MinimumSeverity: types.StringPointerValue(input.MinimumSeverity),
		Syslog:          oro.flattenSecurityThreatDetectionSyslog(input.Syslog),
	}

	// Return the successfully created types.List
	return FlattenList(oro.Ctx, oro.Diags, []models.SecurityThreatDetectionModel{block})
}

// flattenSecurityThreatDetectionSyslog transforms *client.OrgThreatDetectionSyslog into a types.List.
func (oro *OrgResourceOperator) flattenSecurityThreatDetectionSyslog(input *client.OrgThreatDetectionSyslog) types.List {
	// Get attribute types
	elementType := models.SecurityThreatDetectionSyslogModel{}.AttributeTypes()

	// Check if the input is nil
	if input == nil {
		// Return a null list
		return types.ListNull(elementType)
	}

	// Build a single block
	block := models.SecurityThreatDetectionSyslogModel{
		Transport: types.StringPointerValue(input.Transport),
		Host:      types.StringPointerValue(input.Host),
		Port:      FlattenInt(input.Port),
	}

	// Return the successfully created types.List
	return FlattenList(oro.Ctx, oro.Diags, []models.SecurityThreatDetectionSyslogModel{block})
}

// flattenStatus transforms *client.OrgStatus into a Terraform types.List.
func (oro *OrgResourceOperator) flattenStatus(input *client.OrgStatus) types.List {
	// Get attribute types
	elementType := models.StatusModel{}.AttributeTypes()

	// Check if the input is nil
	if input == nil {
		// Return a null list
		return types.ListNull(elementType)
	}

	// Build a single block
	block := models.StatusModel{
		AccountLink:    types.StringPointerValue(input.AccountLink),
		Active:         types.BoolPointerValue(input.Active),
		EndpointPrefix: types.StringPointerValue(input.EndpointPrefix),
	}

	// Return the successfully created types.List
	return FlattenList(oro.Ctx, oro.Diags, []models.StatusModel{block})
}
