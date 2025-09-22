package cpln

import (
	"context"
	"fmt"
	"sort"
	"strings"

	client "github.com/controlplane-com/terraform-provider-cpln/internal/provider/client"
	models "github.com/controlplane-com/terraform-provider-cpln/internal/provider/models/policy"
	"github.com/controlplane-com/terraform-provider-cpln/internal/provider/validators"
	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/resourcevalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
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
	_ resource.Resource                = &PolicyResource{}
	_ resource.ResourceWithImportState = &PolicyResource{}
)

/*** Resource Model ***/

// PolicyResourceModel holds the Terraform state for the resource.
type PolicyResourceModel struct {
	EntityBaseModel
	TargetKind  types.String `tfsdk:"target_kind"`
	Gvc         types.String `tfsdk:"gvc"`
	TargetLinks types.Set    `tfsdk:"target_links"`
	TargetQuery types.List   `tfsdk:"target_query"`
	Target      types.String `tfsdk:"target"`
	Origin      types.String `tfsdk:"origin"`
	Binding     types.Set    `tfsdk:"binding"`
}

/*** Resource Configuration ***/

// PolicyResource is the resource implementation.
type PolicyResource struct {
	EntityBase
	Operations EntityOperations[PolicyResourceModel, client.Policy]
}

// NewPolicyResource returns a new instance of the resource implementation.
func NewPolicyResource() resource.Resource {
	return &PolicyResource{}
}

// Configure configures the resource before use.
func (pr *PolicyResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	pr.EntityBaseConfigure(ctx, req.ProviderData, &resp.Diagnostics)
	pr.Operations = NewEntityOperations(pr.client, &PolicyResourceOperator{})
}

// ImportState sets up the import operation to map the imported ID to the "id" attribute in the state.
func (pr *PolicyResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// Metadata provides the resource type name.
func (pr *PolicyResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = "cpln_policy"
}

// Schema defines the schema for the resource.
func (pr *PolicyResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: MergeAttributes(pr.EntityBaseAttributes("Policy"), map[string]schema.Attribute{
			"gvc": schema.StringAttribute{
				Description: "The GVC for `identity`, `workload` and `volumeset` target kinds only.",
				Optional:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
				Validators: []validator.String{
					validators.NameValidator{},
				},
			},
			"target_kind": schema.StringAttribute{
				Description: "The kind of resource to target (e.g., gvc, serviceaccount, etc.).",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"target_links": schema.SetAttribute{
				Description: "List of the targets this policy will be applied to. Not used if `target` is set to `all`.",
				ElementType: types.StringType,
				Optional:    true,
				Validators: []validator.Set{
					setvalidator.SizeAtMost(200),
				},
				PlanModifiers: []planmodifier.Set{
					pr.RequiresReplaceOnRemoval(),
				},
			},
			"target": schema.StringAttribute{
				Description: "Set this value of this attribute to `all` if this policy should target all objects of the given target_kind. Otherwise, do not include the attribute.",
				Optional:    true,
				Validators: []validator.String{
					stringvalidator.OneOf("all"),
				},
			},
			"origin": schema.StringAttribute{
				Description: "Origin of the Policy. Either `builtin` or `default`.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
		}),
		Blocks: map[string]schema.Block{
			"target_query": schema.ListNestedBlock{
				Description:  "A defined set of criteria or conditions used to identify the target entities or resources to which the policy applies.",
				NestedObject: pr.QuerySchema(),
				Validators: []validator.List{
					listvalidator.SizeAtMost(1),
				},
			},
			"binding": schema.SetNestedBlock{
				Description: "The association between a target kind and the bound permissions to service principals.",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"permissions": schema.SetAttribute{
							Description: "List of permissions to allow.",
							ElementType: types.StringType,
							Required:    true,
						},
						"principal_links": schema.SetAttribute{
							Description: "List of the principals this binding will be applied to. Principal links format: `group/GROUP_NAME`, `user/USER_EMAIL`, `gvc/GVC_NAME/identity/IDENTITY_NAME`, `serviceaccount/SERVICE_ACCOUNT_NAME`.",
							ElementType: types.StringType,
							Required:    true,
							Validators: []validator.Set{
								setvalidator.SizeAtLeast(1),
								setvalidator.SizeAtMost(200),
							},
						},
					},
				},
				Validators: []validator.Set{
					setvalidator.SizeAtMost(50),
				},
			},
		},
	}
}

// ConfigValidators enforces mutual exclusivity between attributes.
func (pr *PolicyResource) ConfigValidators(ctx context.Context) []resource.ConfigValidator {
	expressions := []path.Expression{
		path.MatchRoot("target"),
		path.MatchRoot("target_links"),
		path.MatchRoot("target_query"),
	}

	return []resource.ConfigValidator{
		resourcevalidator.AtLeastOneOf(expressions...),
	}
}

// Create creates the resource.
func (pr *PolicyResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	CreateGeneric(ctx, req, resp, pr.Operations)
}

// Read fetches the current state of the resource.
func (pr *PolicyResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	ReadGeneric(ctx, req, resp, pr.Operations)
}

// Update modifies the resource.
func (pr *PolicyResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	UpdateGeneric(ctx, req, resp, pr.Operations)
}

// Delete removes the resource.
func (pr *PolicyResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	DeleteGeneric(ctx, req, resp, pr.Operations)
}

/*** Plan Modifiers ***/

// RequiresReplaceOnRemoval returns a plan modifier that forces a replace if the prior state was non-empty and the new plan is null/empty.
func (pr *PolicyResource) RequiresReplaceOnRemoval() planmodifier.Set {
	return setplanmodifier.RequiresReplaceIf(
		func(ctx context.Context,
			req planmodifier.SetRequest,
			resp *setplanmodifier.RequiresReplaceIfFuncResponse,
		) {
			// Extract the old values
			old := BuildSetString(ctx, &resp.Diagnostics, req.StateValue)

			// Skip if old is nil, there is nothing to do here
			if old == nil {
				return
			}

			// Extract the new values
			new := BuildSetString(ctx, &resp.Diagnostics, req.PlanValue)

			// Just in case the user wants to set target_links to null, we don't allow them
			if new == nil {
				resp.RequiresReplace = true
			}
		},
		"Recreate when removing existing target_links",
		"If this attribute previously had items, removing them will recreate the resource.",
	)
}

/*** Resource Operator ***/

// PolicyResourceOperator is the operator for managing the state.
type PolicyResourceOperator struct {
	EntityOperator[PolicyResourceModel]
}

// NewAPIRequest creates a request payload from a state model.
func (pro *PolicyResourceOperator) NewAPIRequest(isUpdate bool) client.Policy {
	// Initialize a new request payload
	requestPayload := client.Policy{}

	// Populate Base fields from state
	pro.Plan.Fill(&requestPayload.Base, isUpdate)

	// Extract the target kind from the plan
	targetKind := BuildString(pro.Plan.TargetKind)

	// Set specific attributes
	requestPayload.TargetKind = targetKind
	requestPayload.TargetLinks = pro.buildTargetLinks(pro.Plan.TargetLinks, *targetKind, BuildString(pro.Plan.Gvc))
	requestPayload.TargetQuery = pro.BuildQuery(pro.Plan.TargetQuery)
	requestPayload.Target = BuildString(pro.Plan.Target)
	requestPayload.Bindings = pro.buildBinding(pro.Plan.Binding)

	// Return constructed request payload
	return requestPayload
}

// MapResponseToState constructs the Terraform state model from the API response payload.
func (pro *PolicyResourceOperator) MapResponseToState(apiResp *client.Policy, isCreate bool) PolicyResourceModel {
	// Initialize empty state model
	state := PolicyResourceModel{}

	// Populate common fields from base resource data
	state.From(apiResp.Base)

	// Attempt to extract the GVC from the plan
	if pro.Plan.Gvc.IsNull() || pro.Plan.Gvc.IsUnknown() {
		state.Gvc = types.StringNull()
	} else {
		state.Gvc = pro.Plan.Gvc
	}

	// Set specific attributes
	state.TargetKind = types.StringPointerValue(apiResp.TargetKind)
	state.TargetLinks = pro.flattenTargetLinks(apiResp.TargetLinks, pro.Plan.TargetLinks)
	state.TargetQuery = pro.FlattenQuery(apiResp.TargetQuery)
	state.Target = types.StringPointerValue(apiResp.Target)
	state.Origin = types.StringPointerValue(apiResp.Origin)
	state.Binding = pro.flattenBinding(apiResp.Bindings)

	// Return completed state model
	return state
}

// InvokeCreate invokes the Create API to create a new resource.
func (pro *PolicyResourceOperator) InvokeCreate(req client.Policy) (*client.Policy, int, error) {
	return pro.Client.CreatePolicy(req)
}

// InvokeRead invokes the Get API to retrieve an existing resource by name.
func (pro *PolicyResourceOperator) InvokeRead(name string) (*client.Policy, int, error) {
	return pro.Client.GetPolicy(name)
}

// InvokeUpdate invokes the Update API to update an existing resource.
func (pro *PolicyResourceOperator) InvokeUpdate(req client.Policy) (*client.Policy, int, error) {
	return pro.Client.UpdatePolicy(client.PolicyUpdate{
		Base:        req.Base,
		TargetKind:  req.TargetKind,
		TargetLinks: req.TargetLinks,
		TargetQuery: req.TargetQuery,
		Target:      req.Target,
		Bindings:    req.Bindings,
	})
}

// InvokeDelete invokes the Delete API to delete a resource by name.
func (pro *PolicyResourceOperator) InvokeDelete(name string) error {
	return pro.Client.DeletePolicy(name)
}

// Builders //

// buildTargetLinks constructs full API links for each target based on user input and GVC context.
func (pro *PolicyResourceOperator) buildTargetLinks(state types.Set, kind string, gvc *string) *[]string {
	// Initialize the output slice
	output := []string{}

	// Build the planned target links
	targetLinks := pro.BuildSetString(state)

	// Return nil if the planned target links is nil
	if targetLinks == nil {
		return nil
	}

	// Iterate over each target link and process it
	for _, targetLink := range *targetLinks {
		// Typically, the target link that is specified by the user is just the name of the resource or but it could also be a full link
		// So, if the target link is a full path, use it as is
		if strings.HasPrefix(targetLink, "/org/") {
			output = append(output, targetLink)
			continue
		}

		// Otherwise, construct the self link from the specified target link resource name
		finalTargetLink := fmt.Sprintf("/org/%s/%s/%s", pro.Client.Org, kind, targetLink)

		// If the target kind is a GVC scoped resource then let's modify the final target link
		if gvc != nil && IsGvcScopedResource(kind) {
			finalTargetLink = fmt.Sprintf("/org/%s/gvc/%s/%s/%s", pro.Client.Org, *gvc, kind, targetLink)
		}

		// Add the final target link to the output
		output = append(output, finalTargetLink)
	}

	// Return the constructed output
	return &output
}

// buildBinding constructs a []client.Binding from the given Terraform state.
func (pro *PolicyResourceOperator) buildBinding(state types.Set) *[]client.Binding {
	// Convert Terraform set into model blocks using generic helper
	blocks, ok := BuildSet[models.BindingModel](pro.Ctx, pro.Diags, state)

	// Return an empty slice if conversion failed or set was empty
	if !ok {
		return &[]client.Binding{}
	}

	// Prepare the output slice
	output := []client.Binding{}

	// Iterate over each block and construct an output item
	for _, block := range blocks {
		// Construct the item
		item := client.Binding{
			Permissions:    pro.BuildSetString(block.Permissions),
			PrincipalLinks: pro.buildPrincipalLinks(block.PrincipalLinks),
		}

		// Add the item to the output slice
		output = append(output, item)
	}

	// Return a pointer to the output
	return &output
}

// buildPrincipalLinks constructs a []string from the given Terraform state.
func (pro *PolicyResourceOperator) buildPrincipalLinks(state types.Set) *[]string {
	// Build the string slice
	p := pro.BuildSetString(state)

	// Return nil if the slice is nil
	if p == nil {
		return nil
	}

	// Build the org link prefix string
	orgPrefix := fmt.Sprintf(`/org/%s`, pro.Client.Org)

	// Prepare the output slice
	output := []string{}

	// Iterate over each block and construct an output item
	for _, item := range *p {
		// Only add org prefix if it was not provided by the user
		if !strings.HasPrefix(item, orgPrefix) {
			item = fmt.Sprintf("%s/%s", orgPrefix, item)
		}

		// Add the principal link to the output
		output = append(output, item)
	}

	// Return a pointer to the output
	return &output
}

// Flatteners //

// flattenTargetLinks converts API returned links into a Terraform set, preserving full paths when user originally specified them.
func (pro *PolicyResourceOperator) flattenTargetLinks(input *[]string, state types.Set) types.Set {
	// Return a null set if input list is nil
	if input == nil {
		return types.SetNull(types.StringType)
	}

	// Retrieve the user's original planned links from the Terraform state
	plannedLinks := pro.BuildSetString(state)

	// Initialize a set for full path links
	fullPathSet := make(map[string]struct{})

	// Categorize each planned link by full path if necessary
	if plannedLinks != nil {
		for _, plannedLink := range *plannedLinks {
			if strings.HasPrefix(plannedLink, "/org/") {
				fullPathSet[plannedLink] = struct{}{}
			}
		}
	}

	// Prepare output slice with capacity matching number of API links
	output := make([]string, 0, len(*input))

	// Process each link returned by the API
	for _, apiLink := range *input {
		// Preserve full path if user originally specified it
		if _, ok := fullPathSet[apiLink]; ok {
			output = append(output, apiLink)
			continue
		}

		// Otherwise strip the link to its resource name
		output = append(output, apiLink[strings.LastIndexAny(apiLink, "/")+1:])
	}

	// Convert the processed links into a Terraform set
	return FlattenSetString(&output)
}

// flattenBinding transforms *[]client.Binding into a Terraform types.Set.
func (pro *PolicyResourceOperator) flattenBinding(input *[]client.Binding) types.Set {
	// Get attribute types
	elementType := models.BindingModel{}.AttributeTypes()

	// Check if the input is nil
	if input == nil {
		// Return a null set
		return types.SetNull(elementType)
	}

	// Define the blocks slice
	var blocks []models.BindingModel

	// Iterate over the slice and construct the blocks
	for _, item := range *input {
		// Find the matching planned binding
		plannedPrincipalLinks := pro.findMatchingPlannedPrincipalLinks(item)

		// Construct a block
		block := models.BindingModel{
			Permissions:    FlattenSetString(item.Permissions),
			PrincipalLinks: pro.flattenPrincipalLinks(item.PrincipalLinks, plannedPrincipalLinks),
		}

		// Append the constructed block to the blocks slice
		blocks = append(blocks, block)
	}

	// Return the successfully created types.Set
	return FlattenSet(pro.Ctx, pro.Diags, blocks)
}

// flattenPrincipalLinks transforms *[]string into a types.Set.
func (pro *PolicyResourceOperator) flattenPrincipalLinks(input *[]string, plannedLinks *[]string) types.Set {
	// Check if the input is nil
	if input == nil {
		// Return a null set
		return types.SetNull(types.StringType)
	}

	// Build the org link prefix string
	orgPrefix := fmt.Sprintf(`/org/%s/`, pro.Client.Org)

	// Initialize a set for full path links
	fullPathPrincipalLinks := make(map[string]struct{})

	// Categorize each planned link by full path if necessary
	if plannedLinks != nil {
		for _, plannedLink := range *plannedLinks {
			if strings.HasPrefix(plannedLink, orgPrefix) {
				fullPathPrincipalLinks[plannedLink] = struct{}{}
			}
		}
	}

	// Define the items slice
	var result []string

	// Iterate over the slice and construct the blocks
	for _, l := range *input {
		// Preserve full path if user originally specified it
		if _, ok := fullPathPrincipalLinks[l]; ok {
			result = append(result, l)
			continue
		}

		// Remove the org link prefix from the principal link and add it to the result
		result = append(result, strings.TrimPrefix(l, orgPrefix))
	}

	// Return the successfully accumulated blocks
	return FlattenSetString(&result)
}

// Helpers //

// findMatchingPlannedPrincipalLinks returns the pointer to the planned principal links slice
func (pro *PolicyResourceOperator) findMatchingPlannedPrincipalLinks(input client.Binding) *[]string {
	// Extract input permissions slice (or empty if nil)
	var inPerms []string
	if input.Permissions != nil {
		inPerms = *input.Permissions
	}

	// Extract input principal-links length
	inLinksLen := 0
	if input.PrincipalLinks != nil {
		inLinksLen = len(*input.PrincipalLinks)
	}

	// Sort the input permissions for reliable comparison
	sort.Strings(inPerms)

	// Build the binding set from the plan
	if pro.Plan.Binding.IsNull() || pro.Plan.Binding.IsUnknown() {
		return nil
	}

	// Convert Terraform set into model blocks using generic helper
	proposedBindings, ok := BuildSet[models.BindingModel](pro.Ctx, pro.Diags, pro.Plan.Binding)

	// Return nil if conversion failed
	if !ok {
		return nil
	}

	// Scan through each planned binding
	for _, pb := range proposedBindings {
		// Pull out the planned permissions into a []string
		pbPerms := pro.BuildSetString(pb.Permissions)

		// Skip this one if extraction failed
		if pbPerms == nil {
			continue
		}

		// Quick length check
		if len(*pbPerms) != len(inPerms) {
			continue
		}

		// Sort planned perms and compare element-by-element
		sort.Strings(*pbPerms)

		// Assume that there is a match until proven otherwise
		match := true

		// Iterate over permissions and compare
		for i := range *pbPerms {
			if (*pbPerms)[i] != inPerms[i] {
				match = false
				break
			}
		}

		// If there is no match, then skip this binding and move on
		if !match {
			continue
		}

		// Now extract planned principal-links and compare lengths
		pbLinks := pro.BuildSetString(pb.PrincipalLinks)

		// Skip this one if extraction failed
		if pbLinks == nil {
			continue
		}

		// Quick length check
		if len(*pbLinks) != inLinksLen {
			continue
		}

		// Both criteria match: this is our planned binding
		return pbLinks
	}

	// Return nil for the binding that was not found
	return nil
}
