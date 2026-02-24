package cpln

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	client "github.com/controlplane-com/terraform-provider-cpln/internal/provider/client"
	"github.com/hashicorp/terraform-plugin-framework-validators/int32validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/resourcevalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int32default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int32planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure resource implements required interfaces.
var (
	_ resource.Resource                = &DomainRouteResource{}
	_ resource.ResourceWithImportState = &DomainRouteResource{}
)

/*** Resource Model ***/

// DomainRouteResourceModel holds the Terraform state for the resource.
type DomainRouteResourceModel struct {
	ID            types.String `tfsdk:"id"`
	DomainLink    types.String `tfsdk:"domain_link"`
	DomainPort    types.Int32  `tfsdk:"domain_port"`
	Prefix        types.String `tfsdk:"prefix"`
	ReplacePrefix types.String `tfsdk:"replace_prefix"`
	Regex         types.String `tfsdk:"regex"`
	WorkloadLink  types.String `tfsdk:"workload_link"`
	Port          types.Int32  `tfsdk:"port"`
	HostPrefix    types.String `tfsdk:"host_prefix"`
	HostRegex     types.String `tfsdk:"host_regex"`
	Headers       types.List   `tfsdk:"headers"`
	Replica       types.Int32  `tfsdk:"replica"`
}

/*** Resource Configuration ***/

// DomainRouteResource is the resource implementation.
type DomainRouteResource struct {
	EntityBase
}

// NewDomainRouteResource returns a new instance of the resource implementation.
func NewDomainRouteResource() resource.Resource {
	return &DomainRouteResource{}
}

// Configure configures the resource before use.
func (drr *DomainRouteResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	drr.EntityBaseConfigure(ctx, req.ProviderData, &resp.Diagnostics)
}

// ImportState sets up the import operation to map the imported ID to the "id" attribute in the state.
func (drr *DomainRouteResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Split the import ID
	parts := strings.SplitN(req.ID, ":", 3)

	// Validate that ID has exactly three non-empty segments
	if len(parts) != 3 || parts[0] == "" || parts[1] == "" || parts[2] == "" {
		// Report error when import identifier format is unexpected
		resp.Diagnostics.AddError(
			"Unexpected Import Identifier",
			fmt.Sprintf(
				"Expected import identifier with format: "+
					"'domain_link:domain_port:[PREFIX|REGEX]'. Got: %q", req.ID,
			),
		)

		// Abort import operation on error
		return
	}

	// Extract domainLink, domainPortStr, and prefixOrRegex from parts
	domainLink, domainPortStr, prefixOrRegex := parts[0], parts[1], parts[2]

	// Convert domainPortStr to integer
	portInt, err := strconv.Atoi(domainPortStr)

	// Handle error when port conversion fails
	if err != nil {
		// Report error for invalid port value in identifier
		resp.Diagnostics.AddError(
			"Invalid Import Identifier",
			fmt.Sprintf(
				"domain_port must be an integer; got %q (error: %s)",
				domainPortStr, err.Error(),
			),
		)

		// Abort import operation on error
		return
	}

	// Cast portInt to int for state attribute
	domainPort := int(portInt)

	// Flag to indicate if route is regex-based
	var isRegex bool

	// Check if API client is available before fetching domain details
	if drr.client != nil {
		// Retrieve domain details and status from API
		dom, status, err := drr.client.GetDomain(GetNameFromSelfLink(domainLink))

		// Report error if API call to fetch domain fails
		if err != nil {
			resp.Diagnostics.AddError("Error fetching domain", err.Error())
			return
		}

		// Report error if domain is not found
		if status == 404 {
			resp.Diagnostics.AddError("Domain not found", fmt.Sprintf("Domain '%s' not found", domainLink))
			return
		}

		// iterate through ports to find matching domainPort
		for _, p := range *dom.Spec.Ports {
			// Skip ports that do not match or have no number
			if p.Number == nil || *p.Number != domainPort {
				continue
			}

			// Inspect each route to determine regex or prefix match
			for _, rt := range *p.Routes {
				// Check for regex match
				if rt.Regex != nil && *rt.Regex == prefixOrRegex {
					isRegex = true
					break
				}

				// Check for prefix match
				if rt.Prefix != nil && *rt.Prefix == prefixOrRegex {
					isRegex = false
					break
				}
			}
		}
	}

	// Build stateID combining domainLink, domainPort, and prefixOrRegex
	stateID := fmt.Sprintf("%s_%d_%s", domainLink, domainPort, prefixOrRegex)

	// Set the generated ID attribute in the Terraform state
	resp.Diagnostics.Append(
		resp.State.SetAttribute(ctx, path.Root("id"), types.StringValue(stateID))...,
	)

	// Set the domain_link attribute in the Terraform state
	resp.Diagnostics.Append(
		resp.State.SetAttribute(ctx, path.Root("domain_link"), types.StringValue(domainLink))...,
	)

	// Set the domain_port attribute in the Terraform state
	resp.Diagnostics.Append(
		resp.State.SetAttribute(ctx, path.Root("domain_port"), types.Int32Value(int32(domainPort)))...,
	)

	// Set the regex or prefix attribute based on isRegex flag
	if isRegex {
		// Assign regex attribute when route is regex-based
		resp.Diagnostics.Append(
			resp.State.SetAttribute(ctx, path.Root("regex"), types.StringValue(prefixOrRegex))...,
		)
	} else {
		// Assign prefix attribute when route is prefix-based
		resp.Diagnostics.Append(
			resp.State.SetAttribute(ctx, path.Root("prefix"), types.StringValue(prefixOrRegex))...,
		)
	}
}

// Metadata provides the resource type name.
func (drr *DomainRouteResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = "cpln_domain_route"
}

// Schema defines the schema for the resource.
func (drr *DomainRouteResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The unique identifier for this Domain Route.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"domain_link": schema.StringAttribute{
				Description: "The self link of the domain to add the route to.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"domain_port": schema.Int32Attribute{
				Description: "The port the route corresponds to. Default: 443",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Int32{
					int32planmodifier.RequiresReplace(),
					int32planmodifier.UseStateForUnknown(),
				},
				Default: int32default.StaticInt32(443),
			},
			"prefix": schema.StringAttribute{
				Description: "The path will match any unmatched path prefixes for the subdomain.",
				Optional:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"replace_prefix": schema.StringAttribute{
				Description: "A path prefix can be configured to be replaced when forwarding the request to the Workload.",
				Optional:    true,
			},
			"regex": schema.StringAttribute{
				Description: "Used to match URI paths. Uses the google re2 regex syntax.",
				Optional:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"workload_link": schema.StringAttribute{
				Description: "The link of the workload to map the prefix to.",
				Required:    true,
			},
			"port": schema.Int32Attribute{
				Description: "For the linked workload, the port to route traffic to.",
				Optional:    true,
			},
			"host_prefix": schema.StringAttribute{
				Description: "This option allows forwarding traffic for different host headers to different workloads. This will only be used when the target GVC has dedicated load balancing enabled and the Domain is configured for wildcard support. Please contact us on Slack or at support@controlplane.com for additional details.",
				Optional:    true,
			},
			"host_regex": schema.StringAttribute{
				Description: "A regex to match the host header. This will only be used when the target GVC has dedicated load balancing enabled and the Domain is configure for wildcard support. Contact your account manager for details.",
				Optional:    true,
			},
			"replica": schema.Int32Attribute{
				Description: "The replica number of a stateful workload to route to. If not provided, traffic will be routed to all replicas.",
				Optional:    true,
				Validators: []validator.Int32{
					int32validator.AtLeast(0), // Ensures replica >= 0
				},
			},
		},
		Blocks: map[string]schema.Block{
			"headers": schema.ListNestedBlock{
				Description: "Modify the headers for all http requests for this route.",
				NestedObject: schema.NestedBlockObject{
					Blocks: map[string]schema.Block{
						"request": schema.ListNestedBlock{
							Description: "Manipulates HTTP headers.",
							NestedObject: schema.NestedBlockObject{
								Attributes: map[string]schema.Attribute{
									"set": schema.MapAttribute{
										Description: "Sets or overrides headers to all http requests for this route.",
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
				},
				Validators: []validator.List{
					listvalidator.SizeAtMost(1),
				},
			},
		},
	}
}

// ConfigValidators enforces mutual exclusivity between attributes.
func (drr *DomainRouteResource) ConfigValidators(ctx context.Context) []resource.ConfigValidator {
	return []resource.ConfigValidator{
		resourcevalidator.ExactlyOneOf(path.MatchRoot("prefix"), path.MatchRoot("regex")),
		resourcevalidator.Conflicting(path.MatchRoot("host_prefix"), path.MatchRoot("host_regex")),
	}
}

// Create creates the resource.
func (drr *DomainRouteResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plannedState DomainRouteResourceModel

	// Retrieve the planned state from the Terraform configuration
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plannedState)...)

	// Abort on errors to avoid partial or inconsistent state
	if resp.Diagnostics.HasError() {
		return
	}

	// Serialize domain operations to prevent read-modify-write race conditions
	domainName := GetNameFromSelfLink(plannedState.DomainLink.ValueString())
	mu := GetDomainLock(domainName)
	mu.Lock()
	defer mu.Unlock()

	// Initialize a new request payload structure and populate it with the planned state
	_, domainPort, route := drr.buildRequest(ctx, &resp.Diagnostics, plannedState)

	// Return if an error has occurred during the request payload creation
	if resp.Diagnostics.HasError() {
		return
	}

	// Send the create request to the API client
	responsePayload, _, err := drr.client.AddDomainRoute(domainName, domainPort, route)

	// Handle any other errors that occurred during the API request
	if err != nil {
		resp.Diagnostics.AddError("API Error", fmt.Sprintf("Error creating domain route: %s", err))
		return
	}

	// Map the API response to the Terraform state
	finalState := drr.buildState(ctx, &resp.Diagnostics, plannedState, plannedState.DomainLink.ValueString(), domainPort, responsePayload)

	// Return if an error has occurred during the state creation
	if resp.Diagnostics.HasError() {
		return
	}

	// Set the resource state in Terraform
	resp.Diagnostics.Append(resp.State.Set(ctx, &finalState)...)
}

// Read fetches the current state of the resource.
func (drr *DomainRouteResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var plannedState DomainRouteResourceModel

	// Retrieve the planned state from the Terraform configuration
	resp.Diagnostics.Append(req.State.Get(ctx, &plannedState)...)

	// Abort on errors to avoid partial or inconsistent state
	if resp.Diagnostics.HasError() {
		return
	}

	// Extract necessary values from the planned state
	domainLink := plannedState.DomainLink.ValueString()
	domainPort := int(plannedState.DomainPort.ValueInt32())

	// Fetch the domain route
	responsePayload, code, err := drr.client.GetDomainRoute(GetNameFromSelfLink(domainLink), domainPort, plannedState.Prefix.ValueStringPointer(), plannedState.Regex.ValueStringPointer())

	// Handle the case where the route is not found (HTTP 404),
	// indicating it has been deleted outside of Terraform. Remove it from state
	if code == 404 {
		resp.State.RemoveResource(ctx)
		return
	}

	// Handle any other errors that occur during the API call
	if err != nil {
		resp.Diagnostics.AddError("API Error", fmt.Sprintf("Error reading domain route: %s", err))
		return
	}

	// Map the API response to the Terraform state
	finalState := drr.buildState(ctx, &resp.Diagnostics, plannedState, domainLink, domainPort, responsePayload)

	// Return if an error has occurred during the state creation
	if resp.Diagnostics.HasError() {
		return
	}

	// Set the updated state in Terraform
	resp.Diagnostics.Append(resp.State.Set(ctx, &finalState)...)
}

// Update modifies the resource.
func (drr *DomainRouteResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plannedState DomainRouteResourceModel

	// Retrieve the planned state from the Terraform configuration
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plannedState)...)

	// Abort on errors to avoid partial or inconsistent state
	if resp.Diagnostics.HasError() {
		return
	}

	// Serialize domain operations to prevent read-modify-write race conditions
	domainName := GetNameFromSelfLink(plannedState.DomainLink.ValueString())
	mu := GetDomainLock(domainName)
	mu.Lock()
	defer mu.Unlock()

	// Initialize a new request payload structure and populate it with the planned state
	_, domainPort, route := drr.buildRequest(ctx, &resp.Diagnostics, plannedState)

	// Return if an error has occurred during the request payload creation
	if resp.Diagnostics.HasError() {
		return
	}

	// Send the update request to the API with the modified data
	responsePayload, _, err := drr.client.UpdateDomainRoute(domainName, domainPort, &route)

	// Handle errors from the API update request
	if err != nil {
		resp.Diagnostics.AddError("API Error", fmt.Sprintf("Error updating domain route: %s", err))
		return
	}

	// Map the API response to the Terraform finalState
	finalState := drr.buildState(ctx, &resp.Diagnostics, plannedState, plannedState.DomainLink.ValueString(), domainPort, responsePayload)

	// Return if an error has occurred during the state creation
	if resp.Diagnostics.HasError() {
		return
	}

	// Set the updated state in Terraform
	resp.Diagnostics.Append(resp.State.Set(ctx, &finalState)...)
}

// Delete removes the resource.
func (drr *DomainRouteResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state DomainRouteResourceModel

	// Retrieve the state from the Terraform configuration
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	// Abort on errors to avoid partial or inconsistent state
	if resp.Diagnostics.HasError() {
		return
	}

	// Serialize domain operations to prevent read-modify-write race conditions
	domainName := GetNameFromSelfLink(state.DomainLink.ValueString())
	mu := GetDomainLock(domainName)
	mu.Lock()
	defer mu.Unlock()

	// Send a delete request to the API using the name from the state
	err := drr.client.RemoveDomainRoute(domainName, int(state.DomainPort.ValueInt32()), state.Prefix.ValueStringPointer(), state.Regex.ValueStringPointer())

	// Handle errors from the API delete request
	if err != nil {
		// If an error occurs during the delete request, add an error to diagnostics
		resp.Diagnostics.AddError("API Error", fmt.Sprintf("Error deleting domain route: %s", err))
		return
	}

	// Remove the resource from Terraform's state, indicating successful deletion
	resp.State.RemoveResource(ctx)
}

/*** Helpers ***/

// buildRequest creates a request payload from a state model.
func (drr *DomainRouteResource) buildRequest(ctx context.Context, diags *diag.Diagnostics, plan DomainRouteResourceModel) (string, int, client.DomainRoute) {
	// Initialize a new request payload
	route := client.DomainRoute{}

	// Set specific attributes
	route.Prefix = BuildString(plan.Prefix)
	route.ReplacePrefix = BuildString(plan.ReplacePrefix)
	route.Regex = BuildString(plan.Regex)
	route.WorkloadLink = BuildString(plan.WorkloadLink)
	route.Port = BuildInt(plan.Port)
	route.HostPrefix = BuildString(plan.HostPrefix)
	route.HostRegex = BuildString(plan.HostRegex)
	route.Headers = BuildRouteHeaders(ctx, diags, plan.Headers)
	route.Replica = BuildInt(plan.Replica)

	// Return constructed request payload
	return GetNameFromSelfLink(plan.DomainLink.ValueString()), int(plan.DomainPort.ValueInt32()), route
}

// buildState creates a state model from response payload.
func (drr *DomainRouteResource) buildState(ctx context.Context, diags *diag.Diagnostics, plan DomainRouteResourceModel, domainLink string, domainPort int, route *client.DomainRoute) DomainRouteResourceModel {
	// Initialize empty state model
	state := DomainRouteResourceModel{}

	// Route prefix or regex details
	var plannedPrefixOrRegex string

	// Determine planned prefix or regex
	if !plan.Prefix.IsNull() && !plan.Prefix.IsUnknown() {
		plannedPrefixOrRegex = fmt.Sprintf("prefix: %s", plan.Prefix.ValueString())
	} else {
		plannedPrefixOrRegex = fmt.Sprintf("regex: %s", plan.Regex.ValueString())
	}

	// In case the route is nil, then it was never found
	if route == nil {
		// Add an error to the diagnostics
		diags.AddError(
			"Route Doesn't Exist",
			fmt.Sprintf(
				"The planned route doesn't exist in the domain, route details: domain link: %s, domain port: %d, %s",
				plan.DomainLink.ValueString(),
				plan.DomainPort.ValueInt32(),
				plannedPrefixOrRegex,
			),
		)

		// Return an empty state
		return state
	}

	// Declare a variable to hold either the prefix or the regex
	var prefixOrRegex string

	// Choose prefix when provided
	if route.Prefix != nil {
		// Use the prefix value when available
		prefixOrRegex = *route.Prefix
	} else {
		// No prefix provided, so regex must be present
		prefixOrRegex = *route.Regex
	}

	// Set specific attributes
	state.ID = types.StringValue(fmt.Sprintf("%s_%d_%s", domainLink, domainPort, prefixOrRegex))
	state.DomainLink = types.StringValue(domainLink)
	state.DomainPort = types.Int32Value(int32(domainPort))
	state.Prefix = types.StringPointerValue(route.Prefix)
	state.ReplacePrefix = types.StringPointerValue(route.ReplacePrefix)
	state.Regex = types.StringPointerValue(route.Regex)
	state.WorkloadLink = types.StringPointerValue(route.WorkloadLink)
	state.Port = FlattenInt(route.Port)
	state.HostPrefix = types.StringPointerValue(route.HostPrefix)
	state.HostRegex = types.StringPointerValue(route.HostRegex)
	state.Headers = FlattenRouteHeaders(ctx, diags, route.Headers)
	state.Replica = FlattenInt(route.Replica)

	// Return completed state model
	return state
}
