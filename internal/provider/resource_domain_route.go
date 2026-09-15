package cpln

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"

	client "github.com/controlplane-com/terraform-provider-cpln/internal/provider/client"
	"github.com/hashicorp/terraform-plugin-framework-validators/float64validator"
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
	_ resource.ResourceWithModifyPlan  = &DomainRouteResource{}
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
	Mirror        types.List   `tfsdk:"mirror"`
	Canary        types.List   `tfsdk:"canary"`
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
	// Split the import ID into its domain, port, path and optional host segments
	parts := strings.SplitN(req.ID, ":", 4)

	// Validate that the identifier carries at least three non-empty segments
	if len(parts) < 3 || parts[0] == "" || parts[1] == "" || parts[2] == "" {
		// Report error when import identifier format is unexpected
		resp.Diagnostics.AddError(
			"Unexpected Import Identifier",
			fmt.Sprintf(
				"Expected import identifier with format: "+
					"'domain_link:domain_port:[PREFIX|REGEX]' or "+
					"'domain_link:domain_port:[PREFIX|REGEX]:[HOST_PREFIX|HOST_REGEX]'. Got: %q", req.ID,
			),
		)

		// Abort import operation on error
		return
	}

	// Extract domainLink, domainPortStr, and pathOrRegex from parts
	domainLink, domainPortStr, pathOrRegex := parts[0], parts[1], parts[2]

	// Read the trailing segment as the host the route matches on, leaving it unset when the identifier omits it
	var host *string

	if len(parts) == 4 {
		host = &parts[3]
	}

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

	// Normalize domainLink to a full self-link if the user provided a domain name instead of a link
	if !strings.HasPrefix(domainLink, "/org/") && !strings.HasPrefix(domainLink, "//") {
		// Construct the full domain link from the provided domain name
		if drr.client != nil {
			domainLink = GetSelfLink(drr.client.Org, "domain", domainLink)
		}
	}

	// Resolve which of the port's routes the identifier addresses
	route, err := drr.findImportRoute(domainLink, domainPort, pathOrRegex, host)

	// A regex may itself contain a colon, so read the last two segments as a single path when nothing matched
	if err == nil && route == nil && host != nil {
		route, err = drr.findImportRoute(domainLink, domainPort, pathOrRegex+":"+*host, nil)
	}

	// Report the failure that resolving the route ran into
	if err != nil {
		resp.Diagnostics.AddError("Unable To Import Domain Route", err.Error())

		// Abort import operation on error
		return
	}

	// Report when the port carries no route the identifier addresses
	if route == nil {
		resp.Diagnostics.AddError(
			"Domain Route Not Found",
			fmt.Sprintf(
				"Domain '%s' has no route matching '%s' at port %d.",
				GetNameFromSelfLink(domainLink), pathOrRegex, domainPort,
			),
		)

		// Abort import operation on error
		return
	}

	// Set the generated ID attribute in the Terraform state
	resp.Diagnostics.Append(
		resp.State.SetAttribute(ctx, path.Root("id"), types.StringValue(buildDomainRouteID(domainLink, int32(domainPort), *route)))...,
	)

	// Set the domain_link attribute in the Terraform state
	resp.Diagnostics.Append(
		resp.State.SetAttribute(ctx, path.Root("domain_link"), types.StringValue(domainLink))...,
	)

	// Set the domain_port attribute in the Terraform state
	resp.Diagnostics.Append(
		resp.State.SetAttribute(ctx, path.Root("domain_port"), types.Int32Value(int32(domainPort)))...,
	)

	// Set the path the route matches on
	resp.Diagnostics.Append(
		resp.State.SetAttribute(ctx, path.Root("prefix"), types.StringPointerValue(route.Prefix))...,
	)

	resp.Diagnostics.Append(
		resp.State.SetAttribute(ctx, path.Root("regex"), types.StringPointerValue(route.Regex))...,
	)

	// Set the host the route matches on so the read that follows addresses this exact route
	resp.Diagnostics.Append(
		resp.State.SetAttribute(ctx, path.Root("host_prefix"), types.StringPointerValue(route.HostPrefix))...,
	)

	resp.Diagnostics.Append(
		resp.State.SetAttribute(ctx, path.Root("host_regex"), types.StringPointerValue(route.HostRegex))...,
	)
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
				Description: "The unique identifier for this Domain Route, which pairs the domain and port with the path and the host the route matches on.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					DomainRouteIdPlanModifier(),
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
			"mirror": schema.ListNestedBlock{
				Description: "Mirror the traffic to the specified workload(s). Only works for workloads running in the same location as the primary workload(s).",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"workload_link": schema.StringAttribute{
							Description: "The workload to mirror traffic to.",
							Required:    true,
						},
						"port": schema.Int32Attribute{
							Description: "The port on the mirrored workload to send traffic to. If not provided, traffic will be mirrored to the first discovered port on the mirrored workload.",
							Optional:    true,
						},
						"percent": schema.Float64Attribute{
							Description: "The percentage of traffic to mirror to the specified workload.",
							Required:    true,
							Validators: []validator.Float64{
								float64validator.Between(0, 100),
							},
						},
					},
				},
			},
			"canary": schema.ListNestedBlock{
				Description: "Routes a weighted percentage of traffic to one or more additional workloads. The combined weight of all canaries on a route must not exceed 100; the remaining weight goes to the primary workload. Only supported on http and http2 ports.",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"workload_link": schema.StringAttribute{
							Description: "The canary workload to route a weighted percentage of traffic to.",
							Required:    true,
						},
						"port": schema.Int32Attribute{
							Description: "The port to send canary traffic to. If not provided, the first configured port on the workload is used.",
							Optional:    true,
						},
						"weight": schema.Int32Attribute{
							Description: "The percentage of traffic to send to this canary workload. A weight of 0 disables the canary so it can be toggled on and off without removing it.",
							Required:    true,
							Validators: []validator.Int32{
								int32validator.Between(0, 100),
							},
						},
					},
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

// ModifyPlan refuses a plan that would put this resource on a route the port already carries.
func (drr *DomainRouteResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	// A destroy carries no plan, and removing a route can never collide with another one
	if req.Plan.Raw.IsNull() || drr.client == nil {
		return
	}

	var plannedState DomainRouteResourceModel

	// Retrieve the route the configuration describes
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plannedState)...)

	// Abort on errors to avoid comparing partial values
	if resp.Diagnostics.HasError() {
		return
	}

	// Leave the plan alone while the domain is still unresolved, a domain being created carries no routes yet
	if plannedState.DomainLink.IsUnknown() || plannedState.DomainPort.IsUnknown() {
		return
	}

	// Resolve the identity the configuration describes
	plannedIdentity := drr.routeIdentity(plannedState)

	// A configuration without a path addresses no route, the schema validators report that on their own
	if plannedIdentity.IdentityKey() == "" {
		return
	}

	// Resolve the identity this resource addresses today, which is empty on a create and on the create half of a replacement
	var priorIdentity client.DomainRoute

	if !req.State.Raw.IsNull() {
		var priorState DomainRouteResourceModel

		// Retrieve the route this resource addresses today
		resp.Diagnostics.Append(req.State.Get(ctx, &priorState)...)

		// Abort on errors to avoid comparing partial values
		if resp.Diagnostics.HasError() {
			return
		}

		priorIdentity = drr.routeIdentity(priorState)
	}

	// Nothing can collide while the resource keeps addressing the same route
	if priorIdentity.IdentityKey() == plannedIdentity.IdentityKey() {
		return
	}

	// Look up the route the configuration describes, which something else already owns when it exists
	existing, _, err := drr.client.GetDomainRoute(GetNameFromSelfLink(plannedState.DomainLink.ValueString()), int(plannedState.DomainPort.ValueInt32()), plannedIdentity)

	// Leave the plan alone when the route is free, and when the lookup fails so the apply reports it with more context
	if err != nil || existing == nil {
		return
	}

	// Refuse the plan before an apply destroys the route this resource addresses today
	drr.addRouteIdentityTakenError(&resp.Diagnostics, priorIdentity, plannedState)
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

	// Fetch the route the state addresses by its path and host
	responsePayload, code, err := drr.client.GetDomainRoute(GetNameFromSelfLink(domainLink), domainPort, drr.routeIdentity(plannedState))

	// Handle the case where the domain is not found (HTTP 404),
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

	// Handle the case where the domain no longer carries the route, so the next plan proposes creating it again
	if responsePayload == nil {
		resp.State.RemoveResource(ctx)
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
	var priorState DomainRouteResourceModel

	// Retrieve the planned state from the Terraform configuration
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plannedState)...)

	// Retrieve the prior state, which holds the path and host the route carries before this update
	resp.Diagnostics.Append(req.State.Get(ctx, &priorState)...)

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

	// Send the update request to the API with the modified data, addressing the route by the identity it had before the update
	responsePayload, _, err := drr.client.UpdateDomainRoute(domainName, domainPort, drr.routeIdentity(priorState), &route)

	// Point at the repair when this resource's state addresses a route its configuration does not describe
	if errors.Is(err, client.ErrDomainRouteConflict) {
		drr.addRouteIdentityTakenError(&resp.Diagnostics, drr.routeIdentity(priorState), plannedState)
		return
	}

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

	// Send a delete request to the API for the route the state addresses by its path and host
	err := drr.client.RemoveDomainRoute(domainName, int(state.DomainPort.ValueInt32()), drr.routeIdentity(state))

	// Handle errors from the API delete request
	if err != nil && !errors.Is(err, client.ErrDomainRouteNotFound) {
		// If an error occurs during the delete request, add an error to diagnostics
		resp.Diagnostics.AddError("API Error", fmt.Sprintf("Error deleting domain route: %s", err))
		return
	}

	// Treat a route the domain no longer carries as already deleted, and say so because the port may still carry a route nothing manages
	if err != nil {
		resp.Diagnostics.AddWarning(
			"Domain Route Already Gone",
			fmt.Sprintf(
				"Domain '%s' carries no route %s at port %d, so there was nothing to delete. "+
					"Check the domain for a route this configuration no longer manages.",
				domainName, drr.routeIdentity(state).Identifier(), state.DomainPort.ValueInt32(),
			),
		)
	}

	// Remove the resource from Terraform's state, indicating successful deletion
	resp.State.RemoveResource(ctx)
}

/*** Plan Modifiers ***/

// domainRouteIdPlanModifier keeps the planned id in step with the attributes that address the route.
type domainRouteIdPlanModifier struct{}

// DomainRouteIdPlanModifier returns a plan modifier that derives the id from the attributes that address the route.
func DomainRouteIdPlanModifier() planmodifier.String {
	return domainRouteIdPlanModifier{}
}

// Description returns a plain text description of the plan modifier's behavior.
func (m domainRouteIdPlanModifier) Description(ctx context.Context) string {
	return "Derives the id from the domain link, the domain port, the prefix or regex, and the host prefix or host regex."
}

// MarkdownDescription returns a markdown description of the plan modifier's behavior.
func (m domainRouteIdPlanModifier) MarkdownDescription(ctx context.Context) string {
	return m.Description(ctx)
}

// PlanModifyString sets the planned id to the value the route's own attributes produce.
func (m domainRouteIdPlanModifier) PlanModifyString(ctx context.Context, req planmodifier.StringRequest, resp *planmodifier.StringResponse) {
	// Declare the attributes that together address a single route within a domain port
	var domainLink, prefix, regex, hostPrefix, hostRegex types.String
	var domainPort types.Int32

	// Read each of them from the plan
	resp.Diagnostics.Append(req.Plan.GetAttribute(ctx, path.Root("domain_link"), &domainLink)...)
	resp.Diagnostics.Append(req.Plan.GetAttribute(ctx, path.Root("domain_port"), &domainPort)...)
	resp.Diagnostics.Append(req.Plan.GetAttribute(ctx, path.Root("prefix"), &prefix)...)
	resp.Diagnostics.Append(req.Plan.GetAttribute(ctx, path.Root("regex"), &regex)...)
	resp.Diagnostics.Append(req.Plan.GetAttribute(ctx, path.Root("host_prefix"), &hostPrefix)...)
	resp.Diagnostics.Append(req.Plan.GetAttribute(ctx, path.Root("host_regex"), &hostRegex)...)

	// Abort on errors to avoid deriving an identifier from partial values
	if resp.Diagnostics.HasError() {
		return
	}

	// Leave the id unknown while any attribute it is derived from is still unresolved
	if domainLink.IsUnknown() || domainPort.IsUnknown() || prefix.IsUnknown() || regex.IsUnknown() || hostPrefix.IsUnknown() || hostRegex.IsUnknown() {
		resp.PlanValue = types.StringUnknown()
		return
	}

	// Build the route identity the configuration describes
	identity := client.DomainRoute{
		Prefix:     BuildString(prefix),
		Regex:      BuildString(regex),
		HostPrefix: BuildString(hostPrefix),
		HostRegex:  BuildString(hostRegex),
	}

	// Store the identifier that identity produces
	resp.PlanValue = types.StringValue(buildDomainRouteID(domainLink.ValueString(), domainPort.ValueInt32(), identity))
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
	route.Mirror = BuildRouteMirror(ctx, diags, plan.Mirror)
	route.Canaries = BuildRouteCanary(ctx, diags, plan.Canary)

	// Return constructed request payload
	return GetNameFromSelfLink(plan.DomainLink.ValueString()), int(plan.DomainPort.ValueInt32()), route
}

// buildState creates a state model from response payload.
func (drr *DomainRouteResource) buildState(ctx context.Context, diags *diag.Diagnostics, plan DomainRouteResourceModel, domainLink string, domainPort int, route *client.DomainRoute) DomainRouteResourceModel {
	// Initialize empty state model
	state := DomainRouteResourceModel{}

	// In case the route is nil, then it was never found
	if route == nil {
		// Add an error to the diagnostics
		diags.AddError(
			"Route Doesn't Exist",
			fmt.Sprintf(
				"The planned route doesn't exist in the domain, route details: domain link: %s, domain port: %d, route %s",
				plan.DomainLink.ValueString(),
				plan.DomainPort.ValueInt32(),
				drr.routeIdentity(plan).Identifier(),
			),
		)

		// Return an empty state
		return state
	}

	// Set specific attributes
	state.ID = types.StringValue(buildDomainRouteID(domainLink, int32(domainPort), *route))
	state.DomainLink = types.StringValue(domainLink)
	state.DomainPort = types.Int32Value(int32(domainPort))
	state.Prefix = types.StringPointerValue(route.Prefix)
	state.ReplacePrefix = types.StringPointerValue(route.ReplacePrefix)
	state.Regex = types.StringPointerValue(route.Regex)
	state.WorkloadLink = FlattenLinkString(plan.WorkloadLink, route.WorkloadLink, drr.client.Org)
	state.Port = FlattenInt(route.Port)
	state.HostPrefix = types.StringPointerValue(route.HostPrefix)
	state.HostRegex = types.StringPointerValue(route.HostRegex)
	state.Headers = FlattenRouteHeaders(ctx, diags, route.Headers)
	state.Replica = FlattenInt(route.Replica)
	state.Mirror = FlattenRouteMirror(ctx, diags, plan.Mirror, route.Mirror, drr.client.Org)
	state.Canary = FlattenRouteCanary(ctx, diags, plan.Canary, route.Canaries, drr.client.Org)

	// Return completed state model
	return state
}

// routeIdentity returns the path and host values that address a single route within a domain port.
func (drr *DomainRouteResource) routeIdentity(state DomainRouteResourceModel) client.DomainRoute {
	// Return only the attributes the API treats as the route's identity
	return client.DomainRoute{
		Prefix:     BuildString(state.Prefix),
		Regex:      BuildString(state.Regex),
		HostPrefix: BuildString(state.HostPrefix),
		HostRegex:  BuildString(state.HostRegex),
	}
}

// findImportRoute returns the single route on a domain port that the given path addresses, narrowed to a host when one is given.
func (drr *DomainRouteResource) findImportRoute(domainLink string, domainPort int, pathOrRegex string, host *string) (*client.DomainRoute, error) {
	// The route can only be resolved against the domain that carries it
	if drr.client == nil {
		return nil, errors.New("the provider is not configured, unable to look up the domain")
	}

	// Retrieve the domain the route belongs to
	domainName := GetNameFromSelfLink(domainLink)
	domain, code, err := drr.client.GetDomain(domainName)

	// Report that the domain is missing rather than the raw transport error
	if code == 404 {
		return nil, fmt.Errorf("domain '%s' not found", domainName)
	}

	// Report any other error the lookup ran into
	if err != nil {
		return nil, fmt.Errorf("error fetching domain '%s': %w", domainName, err)
	}

	// Report when the domain carries no ports at all
	if domain.Spec == nil || domain.Spec.Ports == nil {
		return nil, fmt.Errorf("domain '%s' has no ports configured", domainName)
	}

	// Collect every route on the port whose path matches, narrowed to one host when the identifier carries a host segment
	matches := []client.DomainRoute{}

	for _, port := range *domain.Spec.Ports {
		// Skip ports that carry another number or no routes at all
		if port.Number == nil || *port.Number != domainPort || port.Routes == nil {
			continue
		}

		for _, route := range *port.Routes {
			// Skip routes whose path differs from the one being imported
			if (route.Prefix == nil || *route.Prefix != pathOrRegex) && (route.Regex == nil || *route.Regex != pathOrRegex) {
				continue
			}

			// Skip routes on another host once the identifier carries a host segment, where an empty segment means the route matches every host
			if host != nil && routeHost(route) != *host {
				continue
			}

			matches = append(matches, route)
		}
	}

	// Report that nothing matched so the caller can retry with another reading of the identifier
	if len(matches) == 0 {
		return nil, nil
	}

	// Refuse to guess when the path alone addresses several routes that differ only by host
	if len(matches) > 1 {
		// Collect the hosts the matched routes carry so the error names every candidate
		hosts := []string{}

		for _, match := range matches {
			hosts = append(hosts, routeHost(match))
		}

		return nil, fmt.Errorf(
			"%d routes at port %d of domain '%s' match '%s', one per host (%q). Append the host to the import identifier, for example '%s:%d:%s:%s'. A trailing colon with no host imports the route that matches every host",
			len(matches), domainPort, domainName, pathOrRegex, hosts, domainLink, domainPort, pathOrRegex, hosts[0],
		)
	}

	// Return the only route the identifier addresses
	return &matches[0], nil
}

// routeHost returns the host a route matches on, which is empty when the route matches every host.
func routeHost(route client.DomainRoute) string {
	// Return whichever host matcher the route carries
	switch {
	case route.HostPrefix != nil:
		return *route.HostPrefix
	case route.HostRegex != nil:
		return *route.HostRegex
	}

	return ""
}

// addRouteIdentityTakenError reports that the port already carries the route the configuration describes, and names the repair.
func (drr *DomainRouteResource) addRouteIdentityTakenError(diags *diag.Diagnostics, priorIdentity client.DomainRoute, plannedState DomainRouteResourceModel) {
	// Describe the route the configuration addresses
	domainLink := plannedState.DomainLink.ValueString()
	domainPort := plannedState.DomainPort.ValueInt32()
	planned := drr.routeIdentity(plannedState)
	importIdentifier := buildImportIdentifier(domainLink, domainPort, planned)

	// A resource that addresses no route yet is being created over a route that already exists
	if priorIdentity.IdentityKey() == "" {
		diags.AddError(
			"Domain Route Identity Already Taken",
			fmt.Sprintf(
				"Domain '%s' already carries a route %s at port %d.\n\n"+
					"Import that route rather than creating a second one:\n\n"+
					"  terraform import cpln_domain_route.NAME '%s'\n\n"+
					"If this configuration means to add another route, give it a path or a host no other route on the port uses.",
				GetNameFromSelfLink(domainLink), planned.Identifier(), domainPort, importIdentifier,
			),
		)

		return
	}

	// Report both readings of the collision, since either the state or the configuration is the one to correct
	diags.AddError(
		"Domain Route Identity Already Taken",
		fmt.Sprintf(
			"Domain '%s' already carries a route %s at port %d, while this resource addresses the route %s.\n\n"+
				"Either another resource already manages the route this configuration describes, in which case drop this resource from state and import it at its own identity:\n\n"+
				"  terraform state rm cpln_domain_route.NAME\n"+
				"  terraform import cpln_domain_route.NAME '%s'\n\n"+
				"or the configuration moves this route onto a host another route already uses, in which case choose a host no other route on the port uses.",
			GetNameFromSelfLink(domainLink), planned.Identifier(), domainPort, priorIdentity.Identifier(), importIdentifier,
		),
	)
}

// buildImportIdentifier returns the identifier terraform import uses to address a single route within a domain port.
func buildImportIdentifier(domainLink string, domainPort int32, route client.DomainRoute) string {
	// Resolve the path the route matches on
	var pathOrRegex string

	switch {
	case route.Prefix != nil:
		pathOrRegex = *route.Prefix
	case route.Regex != nil:
		pathOrRegex = *route.Regex
	}

	// The host segment is always written, left empty for a route that matches every host
	return fmt.Sprintf("%s:%d:%s:%s", domainLink, domainPort, pathOrRegex, routeHost(route))
}

// buildDomainRouteID returns the identifier that addresses a single route within a domain port.
func buildDomainRouteID(domainLink string, domainPort int32, route client.DomainRoute) string {
	// Resolve the path the route matches on
	var pathOrRegex string

	switch {
	case route.Prefix != nil:
		pathOrRegex = *route.Prefix
	case route.Regex != nil:
		pathOrRegex = *route.Regex
	}

	// A route that matches every host keeps the identifier earlier provider versions issued
	host := routeHost(route)

	if host == "" {
		return fmt.Sprintf("%s_%d_%s", domainLink, domainPort, pathOrRegex)
	}

	// A route narrowed to one host carries that host so routes sharing a path stay distinct
	return fmt.Sprintf("%s_%d_%s_%s", domainLink, domainPort, pathOrRegex, host)
}
