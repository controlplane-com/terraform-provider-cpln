package cpln

import (
	"context"
	"fmt"
	"strings"

	client "github.com/controlplane-com/terraform-provider-cpln/internal/provider/client"
	models "github.com/controlplane-com/terraform-provider-cpln/internal/provider/models/common"
	"github.com/controlplane-com/terraform-provider-cpln/internal/provider/modifiers"
	"github.com/controlplane-com/terraform-provider-cpln/internal/provider/validators"
	"github.com/hashicorp/terraform-plugin-framework-validators/float64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/int32validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

/*** Generic Interfaces ***/

// HasEntityID defines an interface for types that provide an entity ID.
type HasEntityID interface {
	// GetID returns the entity ID as a types.String.
	GetID() types.String
}

/*** Entity Base Model ***/

// EntityBaseModel holds the shared attributes for Terraform entities.
type EntityBaseModel struct {
	ID          types.String `tfsdk:"id"`
	CplnID      types.String `tfsdk:"cpln_id"`
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`
	Tags        types.Map    `tfsdk:"tags"`
	SelfLink    types.String `tfsdk:"self_link"`
}

// Fill updates the api client from EntityBaseModel.
func (b *EntityBaseModel) Fill(api *client.Base, isUpdate bool) {
	api.Name = BuildString(b.Name)
	api.Description = BuildString(b.Description)

	if isUpdate {
		api.TagsReplace = BuildTags(b.Tags)
	} else {
		api.Tags = BuildTags(b.Tags)
	}
}

// From updates the state attributes from the api client.
func (b *EntityBaseModel) From(api client.Base) {
	b.ID = types.StringPointerValue(api.Name)
	b.CplnID = types.StringPointerValue(api.ID)
	b.Name = types.StringPointerValue(api.Name)
	b.Description = types.StringPointerValue(api.Description)
	b.Tags = FlattenTags(api.Tags)
	b.SelfLink = FlattenSelfLink(api.Links)
}

// GetID returns the ID field from the entity base model.
func (b EntityBaseModel) GetID() types.String {
	// Return the stored ID value
	return b.ID
}

/*** Entity Base ***/

// EntityBase is the base entity (resource/data-source) implementation.
type EntityBase struct {
	client                *client.Client
	IsNameComputed        bool
	IsDescriptionComputed bool
	RequiresReplace       func() planmodifier.String
}

// Configure assigns the provider's client to the entity for API interactions.
func (r *EntityBase) EntityBaseConfigure(ctx context.Context, providerData any, diags *diag.Diagnostics) {
	// Ensure ProviderData is present; return early if the provider has not been properly configured
	if providerData == nil {
		return
	}

	// Assert that ProviderData is of the expected type *client.Client
	c, ok := providerData.(*client.Client)
	if !ok {
		// If the type assertion fails, add an error diagnostic indicating the issue and specifying the expected type
		diags.AddError(
			"Unexpected Configure Type",
			fmt.Sprintf("Expected *client.Client, got: %T. Please report this issue to the provider developers.", providerData),
		)
		return
	}

	// Assign the retrieved client to the entity's client field, enabling API calls
	r.client = c
}

// EntityBaseAttributes returns a map of attributes for a given entity name.
func (r *EntityBase) EntityBaseAttributes(entityName string) map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"id": schema.StringAttribute{
			Description: fmt.Sprintf("The unique identifier for this %s.", entityName),
			Computed:    true,
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
			},
		},
		"cpln_id": schema.StringAttribute{
			Description: fmt.Sprintf("The ID, in GUID format, of the %s.", entityName),
			Computed:    true,
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
			},
		},
		"name":        r.NameSchema(entityName),
		"description": r.DescriptionSchema(entityName),
		"tags": schema.MapAttribute{
			Description: "Key-value map of resource tags.",
			ElementType: types.StringType,
			Optional:    true,
			Computed:    true,
			Validators: []validator.Map{
				validators.TagValidator{},
			},
			PlanModifiers: []planmodifier.Map{
				modifiers.TagPlanModifier{},
			},
		},
		"self_link": schema.StringAttribute{
			Description: "Full link to this resource. Can be referenced by other resources.",
			Computed:    true,
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
			},
		},
	}
}

// Schemas //

// NameSchema returns a StringAttribute schema for the given entity's name, computing or requiring it based on IsNameComputed.
func (r *EntityBase) NameSchema(entityName string) schema.StringAttribute {
	// Build the common attribute description string just once
	description := fmt.Sprintf("Name of the %s.", entityName)

	// Check if the entity name should be marked as computed
	if r.IsNameComputed {
		// Return a computed StringAttribute with only a description
		return schema.StringAttribute{
			Description: description,
			Computed:    true,
		}
	}

	// Initialize a slice to hold string validators for the name
	nameValidators := []validator.String{}

	// Only validate the name if this entity is not a domain
	if entityName != "Domain" {
		// Append the NameValidator to enforce naming rules
		nameValidators = append(nameValidators, validators.NameValidator{})
	}

	// Define a plan modifier for resource replacement
	var requiresReplace planmodifier.String

	// Use the RequiresReplace function to get the plan modifier if found
	if r.RequiresReplace != nil {
		requiresReplace = r.RequiresReplace()
	} else {
		// Default to a no-op plan modifier if not defined
		requiresReplace = stringplanmodifier.RequiresReplace()
	}

	// Return a required StringAttribute with description, validators, and a replacement modifier
	return schema.StringAttribute{
		Description: description,
		Required:    true,
		Validators:  nameValidators,
		PlanModifiers: []planmodifier.String{
			requiresReplace,
			stringplanmodifier.UseStateForUnknown(),
		},
	}
}

// DescriptionSchema returns a StringAttribute schema for the entity's description.
func (r *EntityBase) DescriptionSchema(entityName string) schema.StringAttribute {
	// Build the common attribute description string just once
	description := fmt.Sprintf("Description of the %s.", entityName)

	// Check if the entity description should be marked as computed
	if r.IsDescriptionComputed {
		return schema.StringAttribute{
			Description: description,
			Computed:    true,
		}
	}

	// Build and return the StringAttribute with validators and modifiers
	return schema.StringAttribute{
		Description: description,
		Optional:    true,
		Computed:    true,
		Validators: []validator.String{
			validators.DescriptionValidator{},
		},
		PlanModifiers: []planmodifier.String{
			modifiers.DescriptionPlanModifier{},
		},
	}
}

// QuerySchema returns the nested block schema for query configuration.
func (r *EntityBase) QuerySchema() schema.NestedBlockObject {
	return schema.NestedBlockObject{
		Attributes: map[string]schema.Attribute{
			"fetch": schema.StringAttribute{
				Description: "Type of fetch. Specify either: `links` or `items`. Default: `items`.",
				Optional:    true,
				Computed:    true,
				Validators: []validator.String{
					stringvalidator.OneOf("items", "links"),
				},
				Default: stringdefault.StaticString("items"),
			},
		},
		Blocks: map[string]schema.Block{
			"spec": schema.ListNestedBlock{
				Description: "",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"match": schema.StringAttribute{
							Description: "Type of match. Available values: `all`, `any`, `none`. Default: `all`.",
							Optional:    true,
							Computed:    true,
							Validators: []validator.String{
								stringvalidator.OneOf("all", "any", "none"),
							},
							Default: stringdefault.StaticString("all"),
						},
					},
					Blocks: map[string]schema.Block{
						"terms": schema.ListNestedBlock{
							Description: "Terms can only contain one of the following attributes: `property`, `rel`, `tag`.",
							NestedObject: schema.NestedBlockObject{
								Attributes: map[string]schema.Attribute{
									"op": schema.StringAttribute{
										Description: "Type of query operation. Available values: `=`, `>`, `>=`, `<`, `<=`, `!=`, `exists`, `!exists`. Default: `=`.",
										Optional:    true,
										Computed:    true,
										Validators: []validator.String{
											stringvalidator.OneOf("=", ">", ">=", "<", "<=", "!=", "~", "exists", "!exists"),
										},
										Default: stringdefault.StaticString("="),
									},
									"property": schema.StringAttribute{
										Description: "Property to use for query evaluation.",
										Optional:    true,
									},
									"rel": schema.StringAttribute{
										Description: "Relation to use for query evaluation.",
										Optional:    true,
									},
									"tag": schema.StringAttribute{
										Description: "Tag key to use for query evaluation.",
										Optional:    true,
									},
									"value": schema.StringAttribute{
										Description: "Testing value for query evaluation.",
										Optional:    true,
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
		},
	}
}

// LightstepTracingSchema returns the nested block schema for query configuration.
func (r *EntityBase) LightstepTracingSchema() schema.ListNestedBlock {
	return schema.ListNestedBlock{
		Description: "",
		NestedObject: schema.NestedBlockObject{
			Attributes: map[string]schema.Attribute{
				"sampling": schema.Float64Attribute{
					Description: "Determines what percentage of requests should be traced.",
					Required:    true,
					Validators: []validator.Float64{
						float64validator.Between(0.0, 100.0),
					},
				},
				"endpoint": schema.StringAttribute{
					Description: "Tracing Endpoint Workload. Either the canonical endpoint or internal endpoint.",
					Required:    true,
				},
				"credentials": schema.StringAttribute{
					Description: "Full link to referenced Opaque Secret.",
					Optional:    true,
				},
				"custom_tags": r.CustomTagsTracingSchema(),
			},
		},
		Validators: []validator.List{
			listvalidator.SizeAtMost(1),
		},
	}
}

// OtelTracingSchema returns the nested block schema for query configuration.
func (r *EntityBase) OtelTracingSchema() schema.ListNestedBlock {
	return schema.ListNestedBlock{
		Description: "",
		NestedObject: schema.NestedBlockObject{
			Attributes: map[string]schema.Attribute{
				"sampling": schema.Float64Attribute{
					Description: "Determines what percentage of requests should be traced.",
					Required:    true,
					Validators: []validator.Float64{
						float64validator.Between(0.0, 100.0),
					},
				},
				"endpoint": schema.StringAttribute{
					Description: "Tracing Endpoint Workload. Either the canonical endpoint or internal endpoint.",
					Required:    true,
				},
				"custom_tags": r.CustomTagsTracingSchema(),
			},
		},
		Validators: []validator.List{
			listvalidator.SizeAtMost(1),
		},
	}
}

// ControlPlaneTracingSchema returns the nested block schema for query configuration.
func (r *EntityBase) ControlPlaneTracingSchema() schema.ListNestedBlock {
	return schema.ListNestedBlock{
		Description: "",
		NestedObject: schema.NestedBlockObject{
			Attributes: map[string]schema.Attribute{
				"sampling": schema.Float64Attribute{
					Description: "Determines what percentage of requests should be traced.",
					Required:    true,
					Validators: []validator.Float64{
						float64validator.Between(0.0, 100.0),
					},
				},
				"custom_tags": r.CustomTagsTracingSchema(),
			},
		},
		Validators: []validator.List{
			listvalidator.SizeAtMost(1),
		},
	}
}

// CustomTagsTracingSchema returns the nested block schema for query configuration.
func (r *EntityBase) CustomTagsTracingSchema() schema.MapAttribute {
	return schema.MapAttribute{
		Description: "Key-value map of custom tags.",
		ElementType: types.StringType,
		Optional:    true,
	}
}

// Validators //

// GetPortValidators returns an Int32 validator that enforces a valid TCP port range (80–65535)
func (r *EntityBase) GetPortValidators() []validator.Int32 {
	return []validator.Int32{
		int32validator.Between(80, 65535),
		int32validator.NoneOf(8012, 8022, 9090, 9091, 15000, 15006, 15020, 15021, 15090, 41000),
	}
}

/*** Entity Operator Interface ***/

// EntityOperatorInterface is a generic interface for entity operations.
type EntityOperatorInterface[Plan any, APIObject any] interface {
	Init(ctx context.Context, diags *diag.Diagnostics, client *client.Client, plan Plan)
	NewAPIRequest(isUpdate bool) APIObject
	MapResponseToState(src *APIObject, isCreate bool) Plan
	InvokeCreate(req APIObject) (*APIObject, int, error)
	InvokeRead(name string) (*APIObject, int, error)
	InvokeUpdate(req APIObject) (*APIObject, int, error)
	InvokeDelete(name string) error
}

/*** Entity Operator ***/

// EntityOperator is a generic interface for entity operations.
type EntityOperator[Plan any] struct {
	Ctx    context.Context
	Diags  *diag.Diagnostics
	Client *client.Client
	Plan   Plan
}

// Init initializes the entity operator with context, diagnostics, plan, and client.
func (eo *EntityOperator[Plan]) Init(
	ctx context.Context,
	diags *diag.Diagnostics,
	client *client.Client,
	plan Plan,
) {
	eo.Ctx = ctx
	eo.Diags = diags
	eo.Client = client
	eo.Plan = plan
}

// Builders //

// BuildSetString constructs a string slice from the given Terraform state.
func (eo *EntityOperator[Plan]) BuildSetString(state types.Set) *[]string {
	return BuildSetString(eo.Ctx, eo.Diags, state)
}

// BuildSetInt constructs a string slice from the given Terraform state.
func (eo *EntityOperator[Plan]) BuildSetInt(state types.Set) *[]int {
	return BuildSetInt(eo.Ctx, eo.Diags, state)
}

// BuildStringSet constructs a string slice from the given Terraform state.
func (eo *EntityOperator[Plan]) BuildMapString(state types.Map) *map[string]interface{} {
	return BuildMapString(eo.Ctx, eo.Diags, state)
}

// BuildQuery constructs a Query struct from the given Terraform state.
func (eo *EntityOperator[Plan]) BuildQuery(state types.List) *client.Query {
	return BuildQuery(eo.Ctx, eo.Diags, state)
}

// BuildTracing constructs a client.Tracing from Terraform state.
func (eo *EntityOperator[Plan]) BuildTracing(lightstepTracingState types.List, otelTracingState types.List, cplnTracingState types.List) *client.Tracing {
	return BuildTracing(eo.Ctx, eo.Diags, lightstepTracingState, otelTracingState, cplnTracingState)
}

// BuildLoadBalancerIpSet constructs and formats a LoadBalancer IP set path from the Terraform state and organization.
func (eo *EntityOperator[Plan]) BuildLoadBalancerIpSet(state types.String, org string) *string {
	// Build the ipset string
	ipset := BuildString(state)

	// Return nil if the ipset is nil
	if ipset == nil {
		return nil
	}

	// Format the ipset string and return it
	return eo.FormatIpSetPath(*ipset, org)
}

// formatIpSetPath formats the IP Set path based on the provided IP Set value.
func (eo *EntityOperator[Plan]) FormatIpSetPath(ipSetParam string, org string) *string {
	// Assume this is an IP Set name until proven otherwise
	ipsetName := ipSetParam

	// If the IP Set is a full path, return it as is
	if strings.HasPrefix(ipSetParam, "/org/") || strings.HasPrefix(ipSetParam, "//ipset/") {
		return &ipSetParam
	}

	// Construct the full path and return a pointer to it
	result := fmt.Sprintf("/org/%s/ipset/%s", org, ipsetName)
	return &result
}

// Flatteners //

// FlattenQuery transforms client.Query into a Terraform types.List.
func (eo *EntityOperator[Plan]) FlattenQuery(input *client.Query) types.List {
	return FlattenQuery(eo.Ctx, eo.Diags, input)
}

// FlattenTracing transforms a client.Tracing object into separate Terraform types.List values for each provider.
func (eo *EntityOperator[Plan]) FlattenTracing(input *client.Tracing) (types.List, types.List, types.List) {
	return FlattenTracing(eo.Ctx, eo.Diags, input)
}

// FlattenLoadBalancerIpSet normalizes the LoadBalancer IP set string based on existing state and input.
func (eo *EntityOperator[Plan]) FlattenLoadBalancerIpSet(state types.String, input *string, org string) types.String {
	// Return null if there is no input value
	if input == nil {
		return types.StringNull()
	}

	// Dereference the input string once for readability
	inputVal := *input

	// Extract the current IP set from state, if any
	currentIpSet := BuildString(state)

	// If state has no existing IP set, just use the raw input value
	if currentIpSet == nil {
		return types.StringValue(inputVal)
	}

	// Prepare to compute the normalized output
	var value string

	// Choose behavior based on the prefix of the existing IP set
	switch {
	case strings.HasPrefix(*currentIpSet, "/org/"):
		// Preserve the entire input when the IP set is already namespaced under /org/
		value = inputVal

	case strings.HasPrefix(*currentIpSet, "//ipset/"):
		// For double-slash ipset references, keep only the suffix after "/ipset/"
		parts := strings.SplitN(inputVal, "/ipset/", 2)
		if len(parts) == 2 {
			value = "//ipset/" + parts[1]
		} else {
			// Fallback to full input if split did not yield two parts
			value = inputVal
		}

	default:
		// Strip the organization-specific ipset prefix from the input
		prefix := fmt.Sprintf("/org/%s/ipset/", org)
		value = strings.TrimPrefix(inputVal, prefix)
	}

	// Return the normalized IP set as a Terraform String value
	return types.StringValue(value)
}

/*** Entity Operations ***/

// EntityOperations bundles the provider-specific callbacks.
type EntityOperations[Plan any, APIObject any] struct {
	IdFromPlan  func(plan Plan) string
	NewOperator func(ctx context.Context, diags *diag.Diagnostics, plan Plan) EntityOperatorInterface[Plan, APIObject]
}

// NewEntityOperations initializes a new EntityOperations instance.
func NewEntityOperations[
	Plan HasEntityID, APIObject any, Operator EntityOperatorInterface[Plan, APIObject],
](
	client *client.Client,
	prototype Operator,
) EntityOperations[Plan, APIObject] {
	return EntityOperations[Plan, APIObject]{
		IdFromPlan: func(p Plan) string { return p.GetID().ValueString() },
		NewOperator: func(ctx context.Context, diags *diag.Diagnostics, plan Plan) EntityOperatorInterface[Plan, APIObject] {
			prototype.Init(ctx, diags, client, plan)
			return prototype
		},
	}
}

/*** Generic Functions ***/

// CreateGeneric executes the create operation for a given resource using provided operations.
func CreateGeneric[Plan any, APIObject any](
	ctx context.Context,
	req resource.CreateRequest,
	resp *resource.CreateResponse,
	ops EntityOperations[Plan, APIObject],
) {
	// Declare variable to store desired resource plan
	var plan Plan

	// Populate plan variable from request and capture diagnostics
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)

	// Abort if any diagnostics errors occurred
	if resp.Diagnostics.HasError() {
		// Exit early on error
		return
	}

	// Create a new operator instance
	operator := ops.NewOperator(ctx, &resp.Diagnostics, plan)

	// Create a new API request using the operator
	apiReq := operator.NewAPIRequest(false)

	// Invoke API to create resource and capture response, status code, and error
	apiResp, code, err := operator.InvokeCreate(apiReq)

	// Handle conflict when resource already exists
	if code == 409 {
		// Report resource conflict with guidance
		resp.Diagnostics.AddError("Resource already exists", "Use `terraform import` to bring it under management.")

		// Exit on conflict
		return
	}

	// Handle API invocation errors
	if err != nil {
		// Report API error
		resp.Diagnostics.AddError("API error", err.Error())

		// Abort on API error
		return
	}

	// Build new state from API response
	state := operator.MapResponseToState(apiResp, true)

	// Abort if diagnostics errors occurred
	if resp.Diagnostics.HasError() {
		return
	}

	// Persist new state into Terraform
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// ReadGeneric handles reading the entity state and removes entity on 404.
func ReadGeneric[Plan any, APIObject any](
	ctx context.Context,
	req resource.ReadRequest,
	resp *resource.ReadResponse,
	ops EntityOperations[Plan, APIObject],
) {
	// Declare variable to hold existing state
	var state Plan

	// Populate state from request and capture diagnostics
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	// Abort if diagnostics errors occurred
	if resp.Diagnostics.HasError() {
		// Exit early on error
		return
	}

	// Create a new operator instance
	operator := ops.NewOperator(ctx, &resp.Diagnostics, state)

	// Extract resource ID from current state
	id := ops.IdFromPlan(state)

	// Invoke API to read resource details
	apiResp, code, err := operator.InvokeRead(id)

	// Remove resource from state if not found
	if code == 404 {
		// Drop resource from Terraform state
		resp.State.RemoveResource(ctx)
		return
	}

	// Handle API invocation errors
	if err != nil {
		// Report API error
		resp.Diagnostics.AddError("API error", err.Error())

		// Exit on API error
		return
	}

	// Build new state from API response
	newState := operator.MapResponseToState(apiResp, true)

	// Abort if diagnostics errors occurred
	if resp.Diagnostics.HasError() {
		return
	}

	// Persist updated state into Terraform
	resp.Diagnostics.Append(resp.State.Set(ctx, &newState)...)
}

// UpdateGeneric handles updating the resource based on the planned changes.
func UpdateGeneric[Plan any, APIObject any](
	ctx context.Context,
	req resource.UpdateRequest,
	resp *resource.UpdateResponse,
	ops EntityOperations[Plan, APIObject],
) {
	// Declare variable to store planned changes
	var plan Plan

	// Populate plan variable from request and capture diagnostics
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)

	// Abort if diagnostics errors occurred
	if resp.Diagnostics.HasError() {
		return
	}

	// Create a new operator instance
	operator := ops.NewOperator(ctx, &resp.Diagnostics, plan)

	// Create a new API request using the operator
	apiReq := operator.NewAPIRequest(true)

	// Invoke API to apply update and capture response
	apiResp, _, err := operator.InvokeUpdate(apiReq)

	// Handle API invocation errors
	if err != nil {
		// Report API error
		resp.Diagnostics.AddError("API error", err.Error())

		// Exit on API error
		return
	}

	// Build new state from API response
	newState := operator.MapResponseToState(apiResp, true)

	// Abort if diagnostics errors occurred
	if resp.Diagnostics.HasError() {
		return
	}

	// Persist updated state into Terraform
	resp.Diagnostics.Append(resp.State.Set(ctx, &newState)...)
}

// DeleteGeneric handles deleting the resource and removing it from the state.
func DeleteGeneric[Plan any, APIObject any](
	ctx context.Context,
	req resource.DeleteRequest,
	resp *resource.DeleteResponse,
	ops EntityOperations[Plan, APIObject],
) {
	// Declare variable to hold existing state
	var state Plan

	// Populate state from request and capture diagnostics
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	// Abort if diagnostics errors occurred
	if resp.Diagnostics.HasError() {
		return
	}

	// Create a new operator instance
	operator := ops.NewOperator(ctx, &resp.Diagnostics, state)

	// Extract resource ID from current state
	id := ops.IdFromPlan(state)

	// Invoke API to delete resource by ID
	if err := operator.InvokeDelete(id); err != nil {
		// Report API error
		resp.Diagnostics.AddError("API error", err.Error())

		// Exit on API error
		return
	}

	// Remove resource from Terraform state
	resp.State.RemoveResource(ctx)
}

// Builders //

// BuildQuery constructs a Query struct from the given Terraform state.
func BuildQuery(ctx context.Context, diags *diag.Diagnostics, state types.List) *client.Query {
	// Convert Terraform list into model blocks using generic helper
	blocks, ok := BuildList[models.QueryModel](ctx, diags, state)

	// Return nil if conversion failed or list was empty
	if !ok {
		return nil
	}

	// Extract the very first block from the blocks slice
	block := blocks[0]

	// Construct and return the result
	return &client.Query{
		Fetch: BuildString(block.Fetch),
		Spec:  buildQuerySpec(ctx, diags, block.Spec),
	}
}

// buildQuerySpec constructs a Spec struct from the given Terraform state.
func buildQuerySpec(ctx context.Context, diags *diag.Diagnostics, state types.List) *client.QuerySpec {
	// Convert Terraform list into model blocks using generic helper
	blocks, ok := BuildList[models.QuerySpecModel](ctx, diags, state)

	// Return nil if conversion failed or list was empty
	if !ok {
		return nil
	}

	// Extract the very first block from the blocks slice
	block := blocks[0]

	// Construct and return the result
	return &client.QuerySpec{
		Match: BuildString(block.Match),
		Terms: buildQuerySpecTerms(ctx, diags, block.Terms),
	}
}

// buildQuerySpecTerms constructs a []client.Term slice from the given Terraform state.
func buildQuerySpecTerms(ctx context.Context, diags *diag.Diagnostics, state types.List) *[]client.QueryTerm {
	// Convert Terraform list into model blocks using generic helper
	blocks, ok := BuildList[models.QuerySpecTermModel](ctx, diags, state)

	// Return nil if conversion failed or list was empty
	if !ok {
		return nil
	}

	// Declare the output slice
	output := []client.QueryTerm{}

	// Iterate over each block and construct an output item
	for _, block := range blocks {
		// Construct the item
		item := client.QueryTerm{
			Op:       BuildString(block.Op),
			Property: BuildString(block.Property),
			Rel:      BuildString(block.Rel),
			Tag:      BuildString(block.Tag),
			Value:    BuildString(block.Value),
		}

		// Add the item to the output slice
		output = append(output, item)
	}

	// Return a pointer to the output
	return &output
}

// BuildTracing constructs a client.Tracing from Terraform state.
func BuildTracing(ctx context.Context, diags *diag.Diagnostics, lightstepTracingState types.List, otelTracingState types.List, cplnTracingState types.List) *client.Tracing {
	// Retrieve Lightstep tracing block from Terraform state
	lightstepTracingBlock := getLightstepTracingBlock(ctx, diags, lightstepTracingState)

	// Retrieve OpenTelemetry tracing block from Terraform state
	otelTracingBlock := getOtelTracingBlock(ctx, diags, otelTracingState)

	// Retrieve Control Plane tracing block from Terraform state
	cplnTracingBlock := getCplnTracingBlock(ctx, diags, cplnTracingState)

	// If no tracing blocks are defined, skip tracing configuration and return nil
	if lightstepTracingBlock == nil && otelTracingBlock == nil && cplnTracingBlock == nil {
		return nil
	}

	// Initialize tracing provider variables
	var lightstepProvider *client.TracingProviderLightstep
	var otelProvider *client.TracingProviderOtel
	var cplnProvider *client.TracingProviderControlPlane

	// Variables to hold common tracing settings
	var sampling types.Float64
	var customTags types.Map

	// Populate tracing settings based on the defined block
	if lightstepTracingBlock != nil {
		sampling = lightstepTracingBlock.Sampling
		customTags = lightstepTracingBlock.CustomTags

		// Build LightstepTracing details
		lightstepProvider = &client.TracingProviderLightstep{
			Endpoint:    BuildString(lightstepTracingBlock.Endpoint),
			Credentials: BuildString(lightstepTracingBlock.Credentials),
		}
	} else if otelTracingBlock != nil {
		sampling = otelTracingBlock.Sampling
		customTags = otelTracingBlock.CustomTags

		// Build OtelTelemetry details
		otelProvider = &client.TracingProviderOtel{
			Endpoint: BuildString(otelTracingBlock.Endpoint),
		}
	} else {
		sampling = cplnTracingBlock.Sampling
		customTags = cplnTracingBlock.CustomTags

		// Initialize ControlPlaneTracing without extra details
		cplnProvider = &client.TracingProviderControlPlane{}
	}

	// Assemble and return the final Tracing configuration
	return &client.Tracing{
		Sampling: BuildFloat64(sampling),
		Provider: &client.TracingProvider{
			Lightstep:    lightstepProvider,
			Otel:         otelProvider,
			ControlPlane: cplnProvider,
		},
		CustomTags: buildCustomTags(ctx, diags, customTags),
	}
}

// buildCustomTags builds a map[string]client.CustomTag from Terraform state.
func buildCustomTags(ctx context.Context, diags *diag.Diagnostics, state types.Map) *map[string]client.TracingCustomTag {
	// Convert Terraform state map to a Go map[string]interface{}
	customTags := BuildMapString(ctx, diags, state)

	// If the map is nil, return nil
	if customTags == nil {
		return nil
	}

	// Initialize the output map for CustomTag values
	output := map[string]client.TracingCustomTag{}

	// Iterate over each entry in the state-derived map
	for key, value := range *customTags {
		// Create a CustomTag with a Literal wrapping the value pointer
		output[key] = client.TracingCustomTag{
			Literal: &client.TracingCustomTagValue{
				Value: StringPointerFromInterface(value),
			},
		}
	}

	// Return a pointer to the assembled map of CustomTags
	return &output
}

// Flatteners //

// FlattenQuery transforms client.Query into a Terraform types.List.
func FlattenQuery(ctx context.Context, diags *diag.Diagnostics, input *client.Query) types.List {
	// Get attribute types
	elementType := models.QueryModel{}.AttributeTypes()

	// Check if the input is nil
	if input == nil {
		// Return a null list
		return types.ListNull(elementType)
	}

	// Build a single block
	block := models.QueryModel{
		Fetch: types.StringPointerValue(input.Fetch),
		Spec:  flattenQuerySpec(ctx, diags, input.Spec),
	}

	// Return the successfully created types.List
	return FlattenList(ctx, diags, []models.QueryModel{block})
}

// flattenQuerySpec transforms client.Spec into a Terraform types.List.
func flattenQuerySpec(ctx context.Context, diags *diag.Diagnostics, input *client.QuerySpec) types.List {
	// Get attribute types
	elementType := models.QuerySpecModel{}.AttributeTypes()

	// Check if the input is nil
	if input == nil {
		// Return a null list
		return types.ListNull(elementType)
	}

	// Build a single block
	block := models.QuerySpecModel{
		Match: types.StringPointerValue(input.Match),
		Terms: flattenQuerySpecTerms(ctx, diags, input.Terms),
	}

	// Return the successfully created types.List
	return FlattenList(ctx, diags, []models.QuerySpecModel{block})
}

// flattenQuerySpecTerms transforms []client.Term into a Terraform types.List.
func flattenQuerySpecTerms(ctx context.Context, diags *diag.Diagnostics, input *[]client.QueryTerm) types.List {
	// Get attribute types
	elementType := models.QuerySpecTermModel{}.AttributeTypes()

	// Check if the input is nil
	if input == nil {
		// Return a null list
		return types.ListNull(elementType)
	}

	// Define the blocks slice
	blocks := []models.QuerySpecTermModel{}

	// Iterate over the slice and construct the blocks
	for _, item := range *input {
		// Construct a block
		block := models.QuerySpecTermModel{
			Op:       types.StringPointerValue(item.Op),
			Property: types.StringPointerValue(item.Property),
			Rel:      types.StringPointerValue(item.Rel),
			Tag:      types.StringPointerValue(item.Tag),
			Value:    types.StringPointerValue(item.Value),
		}

		// Append the constructed block to the blocks slice
		blocks = append(blocks, block)
	}

	// Return the successfully created types.List
	return FlattenList(ctx, diags, blocks)
}

// FlattenTracing transforms a client.Tracing object into separate Terraform types.List values for each provider.
func FlattenTracing(ctx context.Context, diags *diag.Diagnostics, input *client.Tracing) (types.List, types.List, types.List) {
	// Determine attribute types for each tracing provider
	lightstepTracingElementType := models.LightstepTracingModel{}.AttributeTypes()
	otelTracingElementType := models.OtelTracingModel{}.AttributeTypes()
	cplnTracingElementType := models.ControlPlaneTracingModel{}.AttributeTypes()

	// Return null lists when tracing or provider is not configured
	if input == nil || input.Provider == nil {
		return types.ListNull(lightstepTracingElementType), types.ListNull(otelTracingElementType), types.ListNull(cplnTracingElementType)
	}

	// Flatten Lightstep tracing block with sampling and custom tags
	lightstepTracingList := flattenLightstepTracing(ctx, diags, input.Provider.Lightstep, input.Sampling, input.CustomTags)

	// Flatten Otel tracing block with sampling and custom tags
	otelTracingList := flattenOtelTracing(ctx, diags, input.Provider.Otel, input.Sampling, input.CustomTags)

	// Flatten Control Plane tracing block with sampling and custom tags
	cplnTracingList := flattenCplnTracing(ctx, diags, input.Provider.ControlPlane, input.Sampling, input.CustomTags)

	// Return the lists for each tracing provider
	return lightstepTracingList, otelTracingList, cplnTracingList
}

// flattenLightstepTracing transforms client.LightstepTracing into a Terraform types.List.
func flattenLightstepTracing(ctx context.Context, diags *diag.Diagnostics, input *client.TracingProviderLightstep, sampling *float64, customTags *map[string]client.TracingCustomTag) types.List {
	// Get attribute types
	elementType := models.LightstepTracingModel{}.AttributeTypes()

	// Check if the input is nil
	if input == nil {
		// Return a null list
		return types.ListNull(elementType)
	}

	// Build a single block
	block := models.LightstepTracingModel{
		Sampling:    FlattenFloat64(sampling),
		Endpoint:    types.StringPointerValue(input.Endpoint),
		Credentials: types.StringPointerValue(input.Credentials),
		CustomTags:  flattenCustomTags(customTags),
	}

	// Return the successfully created types.List
	return FlattenList(ctx, diags, []models.LightstepTracingModel{block})
}

// flattenOtelTracing transforms client.OtelTelemetry into a Terraform types.List.
func flattenOtelTracing(ctx context.Context, diags *diag.Diagnostics, input *client.TracingProviderOtel, sampling *float64, customTags *map[string]client.TracingCustomTag) types.List {
	// Get attribute types
	elementType := models.OtelTracingModel{}.AttributeTypes()

	// Check if the input is nil
	if input == nil {
		// Return a null list
		return types.ListNull(elementType)
	}

	// Build a single block
	block := models.OtelTracingModel{
		Sampling:   FlattenFloat64(sampling),
		Endpoint:   types.StringPointerValue(input.Endpoint),
		CustomTags: flattenCustomTags(customTags),
	}

	// Return the successfully created types.List
	return FlattenList(ctx, diags, []models.OtelTracingModel{block})
}

// flattenCplnTracing transforms client.ControlPlaneTracing into a Terraform types.List.
func flattenCplnTracing(ctx context.Context, diags *diag.Diagnostics, input *client.TracingProviderControlPlane, sampling *float64, customTags *map[string]client.TracingCustomTag) types.List {
	// Get attribute types
	elementType := models.ControlPlaneTracingModel{}.AttributeTypes()

	// Check if the input is nil
	if input == nil {
		// Return a null list
		return types.ListNull(elementType)
	}

	// Build a single block
	block := models.ControlPlaneTracingModel{
		Sampling:   FlattenFloat64(sampling),
		CustomTags: flattenCustomTags(customTags),
	}

	// Return the successfully created types.List
	return FlattenList(ctx, diags, []models.ControlPlaneTracingModel{block})
}

// flattenCustomTags transforms a map[string]client.CustomTag into a Terraform types.Map.
func flattenCustomTags(input *map[string]client.TracingCustomTag) types.Map {
	// Return null map when input is nil
	if input == nil {
		return types.MapNull(types.StringType)
	}

	// Prepare a native map to accumulate literal values
	output := map[string]interface{}{}

	// Iterate through CustomTag entries and extract literal values
	for key, value := range *input {
		output[key] = *value.Literal.Value
	}

	// Convert the native map to a Terraform types.Map
	return FlattenMapString(&output)
}

// Blocks //

// getLightstepTracingBlock constructs a models.LightstepTracing struct from the given Terraform state.
func getLightstepTracingBlock(ctx context.Context, diags *diag.Diagnostics, state types.List) *models.LightstepTracingModel {
	// Convert Terraform list into model blocks using generic helper
	blocks, ok := BuildList[models.LightstepTracingModel](ctx, diags, state)

	// Return nil if conversion failed or list was empty
	if !ok {
		return nil
	}

	// Extract the very first block from the blocks slice
	block := blocks[0]

	// Return a pointer to the block
	return &block
}

// getOtelTracingBlock constructs a models.OtelTracing struct from the given Terraform state.
func getOtelTracingBlock(ctx context.Context, diags *diag.Diagnostics, state types.List) *models.OtelTracingModel {
	// Convert Terraform list into model blocks using generic helper
	blocks, ok := BuildList[models.OtelTracingModel](ctx, diags, state)

	// Return nil if conversion failed or list was empty
	if !ok {
		return nil
	}

	// Extract the very first block from the blocks slice
	block := blocks[0]

	// Return a pointer to the block
	return &block
}

// getCplnTracingBlock constructs a models.ControlPlaneTracing struct from the given Terraform state.
func getCplnTracingBlock(ctx context.Context, diags *diag.Diagnostics, state types.List) *models.ControlPlaneTracingModel {
	// Convert Terraform list into model blocks using generic helper
	blocks, ok := BuildList[models.ControlPlaneTracingModel](ctx, diags, state)

	// Return nil if conversion failed or list was empty
	if !ok {
		return nil
	}

	// Extract the very first block from the blocks slice
	block := blocks[0]

	// Return a pointer to the block
	return &block
}
