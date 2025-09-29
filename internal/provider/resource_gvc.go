package cpln

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	client "github.com/controlplane-com/terraform-provider-cpln/internal/provider/client"
	commonmodels "github.com/controlplane-com/terraform-provider-cpln/internal/provider/models/common"
	models "github.com/controlplane-com/terraform-provider-cpln/internal/provider/models/gvc"
	"github.com/controlplane-com/terraform-provider-cpln/internal/provider/validators"
	"github.com/hashicorp/terraform-plugin-framework-validators/int32validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int32default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure resource implements required interfaces.
var (
	_ resource.Resource                = &GvcResource{}
	_ resource.ResourceWithImportState = &GvcResource{}
)

/*** Resource Model ***/

// GvcResourceModel holds the Terraform state for the resource.
type GvcResourceModel struct {
	EntityBaseModel
	Alias                types.String `tfsdk:"alias"`
	Locations            types.Set    `tfsdk:"locations"`
	PullSecrets          types.Set    `tfsdk:"pull_secrets"`
	Domain               types.String `tfsdk:"domain"`
	EndpointNamingFormat types.String `tfsdk:"endpoint_naming_format"`
	Env                  types.Map    `tfsdk:"env"`
	LightstepTracing     types.List   `tfsdk:"lightstep_tracing"`
	OtelTracing          types.List   `tfsdk:"otel_tracing"`
	ControlPlaneTracing  types.List   `tfsdk:"controlplane_tracing"`
	Sidecar              types.List   `tfsdk:"sidecar"`
	LoadBalancer         types.List   `tfsdk:"load_balancer"`
	Keda                 types.List   `tfsdk:"keda"`
}

/*** Resource Configuration ***/

// GvcResource is the resource implementation.
type GvcResource struct {
	EntityBase
	Operations EntityOperations[GvcResourceModel, client.Gvc]
}

// NewGvcResource returns a new instance of the resource implementation.
func NewGvcResource() resource.Resource {
	resource := GvcResource{}
	resource.EntityBase.RequiresReplace = resource.RequiresReplace
	return &resource
}

// Configure configures the resource before use.
func (gr *GvcResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	gr.EntityBaseConfigure(ctx, req.ProviderData, &resp.Diagnostics)
	gr.Operations = NewEntityOperations(gr.client, &GvcResourceOperator{})
}

// ImportState sets up the import operation to map the imported ID to the "id" attribute in the state.
func (gr *GvcResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// Metadata provides the resource type name.
func (gr *GvcResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = "cpln_gvc"
}

// Schema defines the schema for the resource.
func (gr *GvcResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: MergeAttributes(gr.EntityBaseAttributes("Global Virtual Cloud"), map[string]schema.Attribute{
			"alias": schema.StringAttribute{
				Description: "The alias name of the GVC.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"locations": schema.SetAttribute{
				MarkdownDescription: "A list of [locations](https://docs.controlplane.com/reference/location#current) making up the Global Virtual Cloud.",
				ElementType:         types.StringType,
				Optional:            true,
			},
			"pull_secrets": schema.SetAttribute{
				MarkdownDescription: "A list of [pull secret](https://docs.controlplane.com/reference/gvc#pull-secrets) names used to authenticate to any private image repository referenced by Workloads within the GVC.",
				ElementType:         types.StringType,
				Optional:            true,
			},
			"domain": schema.StringAttribute{
				Description:        "Custom domain name used by associated workloads.",
				DeprecationMessage: "Selecting a domain on a GVC will be deprecated in the future. Use the 'cpln_domain resource' instead.",
				Optional:           true,
			},
			"endpoint_naming_format": schema.StringAttribute{
				Description: "Customizes the subdomain format for the canonical workload endpoint. `default` leaves it as '${workloadName}-${gvcName}.cpln.app'. `org` follows the scheme '${workloadName}-${gvcName}.${org}.cpln.app'.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					gr.RequiresReplace(),
					stringplanmodifier.UseStateForUnknown(),
				},
				Validators: []validator.String{
					stringvalidator.OneOf("default", "org"),
				},
			},
			"env": schema.MapAttribute{
				Description: "Key-value array of resource environment variables.",
				ElementType: types.StringType,
				Optional:    true,
			},
		}),
		Blocks: map[string]schema.Block{
			"lightstep_tracing":    gr.LightstepTracingSchema(),
			"otel_tracing":         gr.OtelTracingSchema(),
			"controlplane_tracing": gr.ControlPlaneTracingSchema(),
			"sidecar": schema.ListNestedBlock{
				Description: "",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"envoy": schema.StringAttribute{
							Description: "",
							Required:    true,
						},
					},
				},
				Validators: []validator.List{
					listvalidator.SizeAtMost(1),
				},
			},
			"load_balancer": schema.ListNestedBlock{
				Description: "Dedicated load balancer configuration.",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"dedicated": schema.BoolAttribute{
							Description: "Creates a dedicated load balancer in each location and enables additional Domain features: custom ports, protocols and wildcard hostnames. Charges apply for each location.",
							Optional:    true,
						},
						"trusted_proxies": schema.Int32Attribute{
							Description: "Controls the address used for request logging and for setting the X-Envoy-External-Address header. If set to 1, then the last address in an existing X-Forwarded-For header will be used in place of the source client IP address. If set to 2, then the second to last address in an existing X-Forwarded-For header will be used in place of the source client IP address. If the XFF header does not have at least two addresses or does not exist then the source client IP address will be used instead.",
							Optional:    true,
							Computed:    true,
							Default:     int32default.StaticInt32(0),
							Validators: []validator.Int32{
								int32validator.AtLeast(0),
								int32validator.AtMost(2),
							},
						},
						"ipset": schema.StringAttribute{
							Description: "The link or the name of the IP Set that will be used for this load balancer.",
							Optional:    true,
						},
					},
					Blocks: map[string]schema.Block{
						"multi_zone": schema.ListNestedBlock{
							Description: "",
							NestedObject: schema.NestedBlockObject{
								Attributes: map[string]schema.Attribute{
									"enabled": schema.BoolAttribute{
										Description: "",
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
						"redirect": schema.ListNestedBlock{
							Description: "Specify the url to be redirected to for different http status codes.",
							NestedObject: schema.NestedBlockObject{
								Blocks: map[string]schema.Block{
									"class": schema.ListNestedBlock{
										Description: "Specify the redirect url for all status codes in a class.",
										NestedObject: schema.NestedBlockObject{
											Attributes: map[string]schema.Attribute{
												"status_5xx": schema.StringAttribute{
													Description: "Specify the redirect url for any 500 level status code.",
													Optional:    true,
												},
												"status_401": schema.StringAttribute{
													Description: "An optional url redirect for 401 responses. Supports envoy format strings to include request information. E.g. https://your-oauth-server/oauth2/authorize?return_to=%REQ(:path)%&client_id=your-client-id",
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
					},
				},
				Validators: []validator.List{
					listvalidator.SizeAtMost(1),
				},
			},
			"keda": schema.ListNestedBlock{
				Description: "KEDA configuration for the GVC.",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"enabled": schema.BoolAttribute{
							Description: "Enable KEDA for this GVC. KEDA is a Kubernetes-based event-driven autoscaler that allows you to scale workloads based on external events. When enabled, a keda operator will be deployed in the GVC and workloads in the GVC can use KEDA to scale based on external metrics.",
							Optional:    true,
							Computed:    true,
							Default:     booldefault.StaticBool(false),
						},
						"identity_link": schema.StringAttribute{
							Description: "A link to an Identity resource that will be used for KEDA. This will allow the keda operator to access cloud and network resources.",
							Optional:    true,
							Validators: []validator.String{
								validators.LinkValidator{},
							},
						},
						"secrets": schema.SetAttribute{
							Description: "A list of secrets to be used as TriggerAuthentication objects. The TriggerAuthentication object will be named after the secret and can be used by triggers on workloads in this GVC.",
							ElementType: types.StringType,
							Optional:    true,
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

// Create creates the resource.
func (gr *GvcResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	CreateGeneric(ctx, req, resp, gr.Operations)
}

// Read fetches the current state of the resource.
func (gr *GvcResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	ReadGeneric(ctx, req, resp, gr.Operations)
}

// Update modifies the resource.
func (gr *GvcResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	UpdateGeneric(ctx, req, resp, gr.Operations)
}

// Delete removes the resource.
func (gr *GvcResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	DeleteGeneric(ctx, req, resp, gr.Operations)
}

/*** Plan Modifiers ***/

// RequiresReplace returns a plan modifier that blocks GVC replacement unless the configuration contains the tag `cpln/terraformReplace = "true"`.
func (gr *GvcResource) RequiresReplace() planmodifier.String {
	return stringplanmodifier.RequiresReplaceIf(
		func(ctx context.Context, req planmodifier.StringRequest, resp *stringplanmodifier.RequiresReplaceIfFuncResponse) {
			// Build human friendly strings for old and new values
			oldStr := StringifyStringValue(req.StateValue)
			newStr := StringifyStringValue(req.PlanValue)

			// Capture the attribute path for the error message
			attrPath := req.Path.String()

			// Skip if it is an endpoint naming format change from null to default
			if attrPath == "endpoint_naming_format" && strings.EqualFold(oldStr, "<null>") && strings.EqualFold(newStr, "default") {
				return
			}

			// Whole resource removed from config
			if req.Config.Raw.IsNull() {
				return
			}

			// Whole resource plan is null when running destroy
			if req.Plan.Raw.IsNull() {
				return
			}

			// Skip when either old or new value is unknown to avoid blocking on computed/plan-time unknowns
			if req.StateValue.IsUnknown() || req.PlanValue.IsUnknown() {
				return
			}

			// Skip logic if the tracked attribute did not change between state and plan
			if req.StateValue.Equal(req.PlanValue) {
				return
			}

			// Declare a holder for the tags attribute
			var tags types.Map

			// Read the tags attribute from configuration (not state) so we see user intent
			resp.Diagnostics.Append(req.Config.GetAttribute(ctx, path.Root("tags"), &tags)...)

			// Abort early if fetching the tags produced diagnostics
			if resp.Diagnostics.HasError() {
				return
			}

			// Declare a native Go map to make key lookups ergonomic
			var m map[string]string

			// Convert the Terraform types.Map into a Go map
			resp.Diagnostics.Append(tags.ElementsAs(ctx, &m, false)...)

			// Abort early if conversion produced diagnostics
			if resp.Diagnostics.HasError() {
				return
			}

			// Allow replacement when the user explicitly opted in with the tag
			if strings.EqualFold(m["cpln/terraformReplace"], "true") {
				resp.RequiresReplace = true
				return
			}

			// Deny replacement and emit a blocking error explaining how to allow the replacement
			resp.Diagnostics.AddError(
				"GVC replacement blocked to avoid downtime and cascading deletes",
				fmt.Sprintf(
					"Replacing a GVC deletes it and then recreates it, causing downtime. Deleting a GVC also deletes all workloads, identities, and volumesets associated with it. If those were created by Terraform, they will be recreated after the new GVC is created.\n\nChange detected:\n- %s:\n    old: %s\n    new: %s\n\nTo proceed, set the tag \"cpln/terraformReplace\" = \"true\" in the resource configuration",
					attrPath, oldStr, newStr,
				),
			)
		},
		"Require an explicit opt-in tag to allow GVC replacement.",
		"Blocks GVC replacement by default to prevent downtime and unintended cascading deletes. Replacing a GVC is a delete-and-recreate operation that removes associated workloads, identities, and volumesets. Terraform-managed dependents will be recreated after the new GVC exists. To opt in, add the tag `cpln/terraformReplace = \"true\"`.",
	)
}

/*** Resource Operator ***/

// GvcResourceOperator is the operator for managing the state.
type GvcResourceOperator struct {
	EntityOperator[GvcResourceModel]
}

// NewAPIRequest creates a request payload from a state model.
func (gro *GvcResourceOperator) NewAPIRequest(isUpdate bool) client.Gvc {
	// Initialize a new request payload
	requestPayload := client.Gvc{}

	// Initialize the GVC spec struct
	var spec *client.GvcSpec = &client.GvcSpec{}

	// Populate Base fields from state
	gro.Plan.Fill(&requestPayload.Base, isUpdate)

	// Map planned state attributes to the API struct
	if isUpdate {
		requestPayload.SpecReplace = spec
	} else {
		requestPayload.Spec = spec
	}

	// Set specific attributes
	spec.StaticPlacement = gro.buildStaticPlacement(gro.Plan.Locations)
	spec.PullSecretLinks = gro.buildPullSecrets(gro.Plan.PullSecrets)
	spec.Domain = BuildString(gro.Plan.Domain)
	spec.EndpointNamingFormat = BuildString(gro.Plan.EndpointNamingFormat)
	spec.Tracing = gro.BuildTracing(gro.Plan.LightstepTracing, gro.Plan.OtelTracing, gro.Plan.ControlPlaneTracing)
	spec.Sidecar = gro.buildSidecar(gro.Plan.Sidecar)
	spec.Env = gro.buildEnv(gro.Plan.Env)
	spec.LoadBalancer = gro.buildLoadBalancer(gro.Plan.LoadBalancer)
	spec.Keda = gro.buildKeda(gro.Plan.Keda)

	// Return constructed request payload
	return requestPayload
}

// MapResponseToState creates a state model from response payload.
func (gro *GvcResourceOperator) MapResponseToState(gvc *client.Gvc, isCreate bool) GvcResourceModel {
	// Initialize a new request payload
	state := GvcResourceModel{}

	// Populate common fields from base resource data
	state.From(gvc.Base)

	// Set attributes that are not related to spec
	state.Alias = types.StringPointerValue(gvc.Alias)

	// Just in case GVC spec is nil
	if gvc.Spec != nil {
		// Extract tracing configurations from spec
		lightstepTracing, otelTracing, cplnTracing := gro.FlattenTracing(gvc.Spec.Tracing)

		// Set specific attributes
		state.Locations = gro.flattenStaticPlacement(gvc.Spec.StaticPlacement)
		state.PullSecrets = gro.flattenPullSecrets(gvc.Spec.PullSecretLinks)
		state.Domain = types.StringPointerValue(gvc.Spec.Domain)
		state.EndpointNamingFormat = types.StringPointerValue(gvc.Spec.EndpointNamingFormat)
		state.Env = gro.flattenEnv(gvc.Spec.Env)
		state.LightstepTracing = lightstepTracing
		state.OtelTracing = otelTracing
		state.ControlPlaneTracing = cplnTracing
		state.Sidecar = gro.flattenSidecar(gvc.Spec.Sidecar)
		state.LoadBalancer = gro.flattenLoadBalancer(gro.Plan.LoadBalancer, gvc.Spec.LoadBalancer)
		state.Keda = gro.flattenKeda(gvc.Spec.Keda)
	} else {
		state.Locations = types.SetNull(types.StringType)
		state.PullSecrets = types.SetNull(types.StringType)
		state.Domain = types.StringNull()
		state.EndpointNamingFormat = types.StringNull()
		state.Env = types.MapNull(types.StringType)
		state.LightstepTracing = types.ListNull(commonmodels.LightstepTracingModel{}.AttributeTypes())
		state.OtelTracing = types.ListNull(commonmodels.OtelTracingModel{}.AttributeTypes())
		state.ControlPlaneTracing = types.ListNull(commonmodels.ControlPlaneTracingModel{}.AttributeTypes())
		state.Sidecar = types.ListNull(models.SidecarModel{}.AttributeTypes())
		state.LoadBalancer = types.ListNull(models.LoadBalancerModel{}.AttributeTypes())
		state.Keda = types.ListNull(models.KedaModel{}.AttributeTypes())
	}

	// Return completed state model
	return state
}

// InvokeCreate invokes the Create API to create a new resource.
func (aro *GvcResourceOperator) InvokeCreate(req client.Gvc) (*client.Gvc, int, error) {
	return aro.Client.CreateGvc(req)
}

// InvokeRead invokes the Get API to retrieve an existing resource by name.
func (aro *GvcResourceOperator) InvokeRead(name string) (*client.Gvc, int, error) {
	return aro.Client.GetGvc(name)
}

// InvokeUpdate invokes the Update API to update an existing resource.
func (aro *GvcResourceOperator) InvokeUpdate(req client.Gvc) (*client.Gvc, int, error) {
	return aro.Client.UpdateGvc(req)
}

// InvokeDelete invokes the Delete API to delete a resource by name.
func (aro *GvcResourceOperator) InvokeDelete(name string) error {
	return aro.Client.DeleteGvc(name)
}

// Builders //

// buildStaticPlacement constructs a client.StaticPlacement from Terraform state.
func (gro *GvcResourceOperator) buildStaticPlacement(state types.Set) *client.GvcStaticPlacement {
	// If the state is unknown or null, there is no block to process, so exit early
	if state.IsNull() || state.IsUnknown() {
		return nil
	}

	// Construct a slice of strings to hold the location links
	locationLinks := []string{}

	// Build the set of location names from the Terraform state
	locationNames := BuildSetString(gro.Ctx, gro.Diags, state)

	// If the location names are not nil, iterate through them and create links
	if locationNames != nil {
		for _, locationName := range *locationNames {
			locationLinks = append(locationLinks, fmt.Sprintf("/org/%s/location/%s", gro.Client.Org, locationName))
		}
	}

	// Construct and return the output
	return &client.GvcStaticPlacement{
		LocationLinks: &locationLinks,
	}
}

// buildPullSecrets constructs a []string from Terraform state.
func (gro *GvcResourceOperator) buildPullSecrets(state types.Set) *[]string {
	// If the state is unknown or null, there is no block to process, so exit early
	if state.IsNull() || state.IsUnknown() {
		return nil
	}

	// Construct a slice of strings to hold the pull secret links
	pullSecretLinks := []string{}

	// Build the set of pull secret names from the Terraform state
	pullSecretNames := BuildSetString(gro.Ctx, gro.Diags, state)

	// If the pull secret names are not nil, iterate through them and create links
	if pullSecretNames != nil {
		for _, pullSecretName := range *pullSecretNames {
			pullSecretLinks = append(pullSecretLinks, fmt.Sprintf("/org/%s/secret/%s", gro.Client.Org, pullSecretName))
		}
	}

	// Return the output
	return &pullSecretLinks
}

// buildSidecar constructs a GvcSidecar struct from the given Terraform state.
func (gro *GvcResourceOperator) buildSidecar(state types.List) *client.GvcSidecar {
	// Convert Terraform list into model blocks using generic helper
	blocks, ok := BuildList[models.SidecarModel](gro.Ctx, gro.Diags, state)

	// Return nil if conversion failed or list was empty
	if !ok {
		return nil
	}

	// Extract the very first block from the blocks slice
	block := blocks[0]

	// Build the envoy string
	envoyString := BuildString(block.Envoy)

	// Return empty GvcSidecar if envoy string is nil
	if envoyString == nil {
		return &client.GvcSidecar{}
	}

	// Attempt to unmarshal `envoy`
	var envoy interface{}
	err := json.Unmarshal([]byte(*envoyString), &envoy)
	if err != nil {
		gro.Diags.AddError("Unable to Unmarshall Sidecar Envoy", fmt.Sprintf("Error occurred during unmarshaling 'envoy' value. Error: %s", err.Error()))
	}

	// Construct and return the result
	return &client.GvcSidecar{
		Envoy: &envoy,
	}
}

// buildEnv constructs a []client.NameValue from Terraform state.
func (gro *GvcResourceOperator) buildEnv(state types.Map) *[]client.WorkloadContainerNameValue {
	// Convert Terraform state map to a Go map[string]interface{}
	envMap := BuildMapString(gro.Ctx, gro.Diags, state)

	// Return nil if the converted map is nil
	if envMap == nil {
		return nil
	}

	// Initialize output slice for NameValue entries
	output := []client.WorkloadContainerNameValue{}

	// Loop through each entry in the state-derived map
	for key, value := range *envMap {
		// Create a NameValue with key pointer and value pointer
		item := client.WorkloadContainerNameValue{
			Name:  &key,
			Value: StringPointerFromInterface(value),
		}

		// Add the item to the output slice
		output = append(output, item)
	}

	// Return a pointer to the assembled slice of NameValue entries
	return &output
}

// buildLoadBalancer constructs a LoadBalancer struct from the given Terraform state.
func (gro *GvcResourceOperator) buildLoadBalancer(state types.List) *client.GvcLoadBalancer {
	// Convert Terraform list into model blocks using generic helper
	blocks, ok := BuildList[models.LoadBalancerModel](gro.Ctx, gro.Diags, state)

	// Return nil if conversion failed or list was empty
	if !ok {
		return nil
	}

	// Extract the very first block from the blocks slice
	block := blocks[0]

	// Construct and return the output
	return &client.GvcLoadBalancer{
		Dedicated:      BuildBool(block.Dedicated),
		MultiZone:      gro.buildLoadBalancerMultiZone(block.MultiZone),
		TrustedProxies: BuildInt(block.TrustedProxies),
		Redirect:       gro.buildLoadBalancerRedirect(block.Redirect),
		IpSet:          gro.BuildLoadBalancerIpSet(block.IpSet, gro.Client.Org),
	}
}

// buildLoadBalancerMultiZone constructs a LoadBalancerMultiZone struct from the given Terraform state.
func (gro *GvcResourceOperator) buildLoadBalancerMultiZone(state types.List) *client.GvcLoadBalancerMultiZone {
	// Convert Terraform list into model blocks using generic helper
	blocks, ok := BuildList[models.LoadBalancerMultiZoneModel](gro.Ctx, gro.Diags, state)

	// Return nil if conversion failed or list was empty
	if !ok {
		return nil
	}

	// Extract the very first block from the blocks slice
	block := blocks[0]

	// Construct and return the result
	return &client.GvcLoadBalancerMultiZone{
		Enabled: BuildBool(block.Enabled),
	}
}

// buildLoadBalancerRedirect constructs a Redirect struct from the given Terraform state.
func (gro *GvcResourceOperator) buildLoadBalancerRedirect(state types.List) *client.GvcLoadBalancerRedirect {
	// Convert Terraform list into model blocks using generic helper
	blocks, ok := BuildList[models.LoadBalancerRedirectModel](gro.Ctx, gro.Diags, state)

	// Return nil if conversion failed or list was empty
	if !ok {
		return nil
	}

	// Extract the very first block from the blocks slice
	block := blocks[0]

	// Construct and return the result
	return &client.GvcLoadBalancerRedirect{
		Class: gro.buildLoadBalancerRedirectClass(block.Class),
	}
}

// buildLoadBalancerRedirectClass constructs a RedirectClass struct from the given Terraform state.
func (gro *GvcResourceOperator) buildLoadBalancerRedirectClass(state types.List) *client.GvcLoadBalancerRedirectClass {
	// Convert Terraform list into model blocks using generic helper
	blocks, ok := BuildList[models.LoadBalancerRedirectClassModel](gro.Ctx, gro.Diags, state)

	// Return nil if conversion failed or list was empty
	if !ok {
		return nil
	}

	// Extract the very first block from the blocks slice
	block := blocks[0]

	// Construct and return the result
	return &client.GvcLoadBalancerRedirectClass{
		Status5XX: BuildString(block.Status5xx),
		Status401: BuildString(block.Status401),
	}
}

// buildKeda constructs a GvcKeda from the given Terraform state.
func (gro *GvcResourceOperator) buildKeda(state types.List) *client.GvcKeda {
	// Convert Terraform list into model blocks using generic helper
	blocks, ok := BuildList[models.KedaModel](gro.Ctx, gro.Diags, state)

	// Return nil if conversion failed or list was empty
	if !ok {
		return nil
	}

	// Take the first (and only) block
	block := blocks[0]

	// Construct and return the output
	return &client.GvcKeda{
		Enabled:      BuildBool(block.Enabled),
		IdentityLink: BuildString(block.IdentityLink),
		Secrets:      gro.BuildSetString(block.Secrets),
	}
}

// Flatteners //

// flattenStaticPlacement transforms client.StaticPlacement into a Terraform types.List.
func (gro *GvcResourceOperator) flattenStaticPlacement(input *client.GvcStaticPlacement) types.Set {
	// Check if the input is nil
	if input == nil {
		return types.SetNull(types.StringType)
	}

	// Construct a slice of strings to hold the location names
	locationNames := []string{}

	// Check if the LocationLinks field is not nil and iterate through it
	if input.LocationLinks != nil {
		// Iterate through the location links and extract the names
		for _, locationLink := range *input.LocationLinks {
			locationNames = append(locationNames, strings.TrimPrefix(locationLink, fmt.Sprintf("/org/%s/location/", gro.Client.Org)))
		}
	}

	// Flatten the location names into a Terraform types.Set
	return FlattenSetString(&locationNames)
}

// flattenPullSecrets transforms []string into a Terraform types.Set.
func (gro *GvcResourceOperator) flattenPullSecrets(input *[]string) types.Set {
	// Check if the input is nil
	if input == nil {
		return types.SetNull(types.StringType)
	}

	// Construct a slice of strings to hold the pull secret names
	pullSecretNames := []string{}

	// Iterate through the pull secret links and extract the names
	for _, pullSecretLink := range *input {
		pullSecretNames = append(pullSecretNames, strings.TrimPrefix(pullSecretLink, fmt.Sprintf("/org/%s/secret/", gro.Client.Org)))
	}

	// Flatten the pull secret names into a Terraform types.Set
	return FlattenSetString(&pullSecretNames)
}

// flattenEnv transforms []client.NameValue into a Terraform types.List.
func (gro *GvcResourceOperator) flattenEnv(input *[]client.WorkloadContainerNameValue) types.Map {
	// Check if the input is nil
	if input == nil {
		return types.MapNull(types.StringType)
	}

	// Prepare a native map for conversion
	envMap := map[string]interface{}{}

	// Populate native map from NameValue entries
	for _, item := range *input {
		// Skip this record just in case the name was nil
		if item.Name == nil {
			continue
		}

		// Dereference the record name
		key := *item.Name

		// Initialize with a nil value
		envMap[key] = nil

		// If the value is not nil, update the map key
		if item.Value != nil {
			envMap[key] = *item.Value
		}
	}

	// Convert native map to Terraform types.Map
	return FlattenMapString(&envMap)
}

// flattenSidecar transforms client.GvcSidecar into a Terraform types.List.
func (gro *GvcResourceOperator) flattenSidecar(input *client.GvcSidecar) types.List {
	// Get attribute types
	elementType := models.SidecarModel{}.AttributeTypes()

	// Check if the input is nil
	if input == nil {
		// Return a null list
		return types.ListNull(elementType)
	}

	// Build a single block
	block := models.SidecarModel{
		Envoy: types.StringNull(),
	}

	// Marshal Envoy JSON when Envoy field is provided
	if input.Envoy != nil {
		jsonOut, err := json.Marshal(*input.Envoy)
		if err != nil {
			// Record a diagnostic error if marshaling fails
			gro.Diags.AddError("Unable to unmarshal Sidecar Envoy", fmt.Sprintf("Error during Envoy JSON marshal: %s", err.Error()))
		}

		// Build the planned sidecar
		plannedSidecar, ok := BuildList[models.SidecarModel](gro.Ctx, gro.Diags, gro.Plan.Sidecar)

		// Set Envoy attribute to the marshaled JSON string
		if ok && len(plannedSidecar) > 0 {
			// Set Envoy attribute to the preserved JSON string
			block.Envoy = PreserveJSONFormatting(string(jsonOut), plannedSidecar[0].Envoy)
		} else {
			// Set Envoy attribute to the marshaled JSON string
			block.Envoy = types.StringValue(string(jsonOut))
		}

	}

	// Return the successfully created types.List
	return FlattenList(gro.Ctx, gro.Diags, []models.SidecarModel{block})
}

// flattenLoadBalancer transforms client.LoadBalancer into a Terraform types.List.
func (gro *GvcResourceOperator) flattenLoadBalancer(state types.List, input *client.GvcLoadBalancer) types.List {
	// Get attribute types
	elementType := models.LoadBalancerModel{}.AttributeTypes()

	// Check if the input is nil
	if input == nil {
		// Return a null list
		return types.ListNull(elementType)
	}

	// Initialize the ipSetState to null
	ipSetState := types.StringNull()

	// Build the list of LoadBalancerModel blocks from the state
	stateBlocks, ok := BuildList[models.LoadBalancerModel](gro.Ctx, gro.Diags, state)

	// If the state is not nil and contains blocks, extract the IP Set from the first block
	if ok {
		ipSetState = stateBlocks[0].IpSet
	}

	// Build a single block
	block := models.LoadBalancerModel{
		Dedicated:      types.BoolPointerValue(input.Dedicated),
		MultiZone:      gro.flattenLoadBalancerMultiZone(input.MultiZone),
		TrustedProxies: FlattenInt(input.TrustedProxies),
		Redirect:       gro.flattenLoadBalancerRedirect(input.Redirect),
		IpSet:          gro.FlattenLoadBalancerIpSet(ipSetState, input.IpSet, gro.Client.Org),
	}

	// Return the successfully created types.List
	return FlattenList(gro.Ctx, gro.Diags, []models.LoadBalancerModel{block})
}

// flattenLoadBalancerMultiZone transforms client.LoadBalancerMultiZone into a Terraform types.List.
func (gro *GvcResourceOperator) flattenLoadBalancerMultiZone(input *client.GvcLoadBalancerMultiZone) types.List {
	// Get attribute types
	elementType := models.LoadBalancerMultiZoneModel{}.AttributeTypes()

	// Check if the input is nil
	if input == nil {
		// Return a null list
		return types.ListNull(elementType)
	}

	// Build a single block
	block := models.LoadBalancerMultiZoneModel{
		Enabled: types.BoolPointerValue(input.Enabled),
	}

	// Return the successfully created types.List
	return FlattenList(gro.Ctx, gro.Diags, []models.LoadBalancerMultiZoneModel{block})
}

// flattenLoadBalancerRedirect transforms client.Redirect into a Terraform types.List.
func (gro *GvcResourceOperator) flattenLoadBalancerRedirect(input *client.GvcLoadBalancerRedirect) types.List {
	// Get attribute types
	elementType := models.LoadBalancerRedirectModel{}.AttributeTypes()

	// Check if the input is nil
	if input == nil {
		// Return a null list
		return types.ListNull(elementType)
	}

	// Build a single block
	block := models.LoadBalancerRedirectModel{
		Class: gro.flattenLoadBalancerRedirectClass(input.Class),
	}

	// Return the successfully created types.List
	return FlattenList(gro.Ctx, gro.Diags, []models.LoadBalancerRedirectModel{block})
}

// flattenLoadBalancerRedirectClass transforms client.RedirectClass into a Terraform types.List.
func (gro *GvcResourceOperator) flattenLoadBalancerRedirectClass(input *client.GvcLoadBalancerRedirectClass) types.List {
	// Get attribute types
	elementType := models.LoadBalancerRedirectClassModel{}.AttributeTypes()

	// Check if the input is nil
	if input == nil {
		// Return a null list
		return types.ListNull(elementType)
	}

	// Build a single block
	block := models.LoadBalancerRedirectClassModel{
		Status5xx: types.StringPointerValue(input.Status5XX),
		Status401: types.StringPointerValue(input.Status401),
	}

	// Return the successfully created types.List
	return FlattenList(gro.Ctx, gro.Diags, []models.LoadBalancerRedirectClassModel{block})
}

// flattenKeda transforms *client.GvcKeda into a Terraform types.List.
func (gro *GvcResourceOperator) flattenKeda(input *client.GvcKeda) types.List {
	// Get attribute types
	elementType := models.KedaModel{}.AttributeTypes()

	// Check if the input is nil
	if input == nil {
		// Return a null list
		return types.ListNull(elementType)
	}

	// Build a single block
	block := models.KedaModel{
		Enabled:      types.BoolPointerValue(input.Enabled),
		IdentityLink: types.StringPointerValue(input.IdentityLink),
		Secrets:      FlattenSetString(input.Secrets),
	}

	// Return the successfully created types.List
	return FlattenList(gro.Ctx, gro.Diags, []models.KedaModel{block})
}
