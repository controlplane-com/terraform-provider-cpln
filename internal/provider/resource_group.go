package cpln

import (
	"context"
	"fmt"
	"strings"

	client "github.com/controlplane-com/terraform-provider-cpln/internal/provider/client"
	models "github.com/controlplane-com/terraform-provider-cpln/internal/provider/models/group"
	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure resource implements required interfaces.
var (
	_ resource.Resource                = &GroupResource{}
	_ resource.ResourceWithImportState = &GroupResource{}
)

/*** Resource Model ***/

// GroupResourceModel holds the Terraform state for the resource.
type GroupResourceModel struct {
	EntityBaseModel
	UserIdsAndEmails types.Set    `tfsdk:"user_ids_and_emails"`
	ServiceAccounts  types.Set    `tfsdk:"service_accounts"`
	MemberQuery      types.List   `tfsdk:"member_query"`
	IdentityMatcher  types.List   `tfsdk:"identity_matcher"`
	Origin           types.String `tfsdk:"origin"`
}

/*** Resource Configuration ***/

// GroupResource is the resource implementation.
type GroupResource struct {
	EntityBase
	Operations EntityOperations[GroupResourceModel, client.Group]
}

// NewGroupResource returns a new instance of the resource implementation.
func NewGroupResource() resource.Resource {
	return &GroupResource{}
}

// Configure configures the resource before use.
func (gr *GroupResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	gr.EntityBaseConfigure(ctx, req.ProviderData, &resp.Diagnostics)
	gr.Operations = NewEntityOperations(gr.client, &GroupResourceOperator{})
}

// ImportState sets up the import operation to map the imported ID to the "id" attribute in the state.
func (gr *GroupResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// Metadata provides the resource type name.
func (gr *GroupResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = "cpln_group"
}

// Schema defines the schema for the resource.
func (gr *GroupResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: MergeAttributes(gr.EntityBaseAttributes("Group"), map[string]schema.Attribute{
			"user_ids_and_emails": schema.SetAttribute{
				Description: "List of either the user ID or email address for a user that exists within the configured org. Group membership will fail if the user ID / email does not exist within the org.",
				ElementType: types.StringType,
				Optional:    true,
			},
			"service_accounts": schema.SetAttribute{
				Description: "List of service accounts that exists within the configured org. Group membership will fail if the service account does not exits within the org.",
				ElementType: types.StringType,
				Optional:    true,
			},
			"origin": schema.StringAttribute{
				Description: "Origin of the service account. Either `builtin` or `default`.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
		}),
		Blocks: map[string]schema.Block{
			"member_query": schema.ListNestedBlock{
				Description:  "A predefined set of criteria or conditions used to query and retrieve members within the group.",
				NestedObject: gr.QuerySchema(),
				Validators: []validator.List{
					listvalidator.SizeAtMost(1),
				},
			},
			"identity_matcher": schema.ListNestedBlock{
				Description: "Executes the expression against the users' claims to decide whether a user belongs to this group. This method is useful for managing the grouping of users logged-in with SAML providers.",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"expression": schema.StringAttribute{
							Description: "Executes the expression against the users' claims to decide whether a user belongs to this group. This method is useful for managing the grouping of users logged in with SAML providers.",
							Required:    true,
						},
						"language": schema.StringAttribute{
							Description: "Language of the expression. Either `jmespath` or `javascript`. Default: `jmespath`.",
							Optional:    true,
							Computed:    true,
							Default:     stringdefault.StaticString("jmespath"),
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
func (gr *GroupResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	CreateGeneric(ctx, req, resp, gr.Operations)
}

// Read fetches the current state of the resource.
func (gr *GroupResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	ReadGeneric(ctx, req, resp, gr.Operations)
}

// Update modifies the resource.
func (gr *GroupResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	UpdateGeneric(ctx, req, resp, gr.Operations)
}

// Delete removes the resource.
func (gr *GroupResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	DeleteGeneric(ctx, req, resp, gr.Operations)
}

/*** Resource Operator ***/

// GroupResourceOperator is the operator for managing the state.
type GroupResourceOperator struct {
	EntityOperator[GroupResourceModel]
}

// NewAPIRequest creates a request payload from a state model.
func (gro *GroupResourceOperator) NewAPIRequest(isUpdate bool) client.Group {
	// Initialize a new request payload
	requestPayload := client.Group{}

	// Populate Base fields from state
	gro.Plan.Fill(&requestPayload.Base, isUpdate)

	// Set specific attributes
	requestPayload.MemberLinks = gro.buildMemberLinks(gro.Plan)
	requestPayload.MemberQuery = gro.BuildQuery(gro.Plan.MemberQuery)
	requestPayload.IdentityMatcher = gro.buildIdentityMatcher(gro.Plan.IdentityMatcher)

	// Return constructed request payload
	return requestPayload
}

// MapResponseToState creates a state model from response payload.
func (gro *GroupResourceOperator) MapResponseToState(group *client.Group, isCreate bool) GroupResourceModel {
	// Initialize empty state model
	state := GroupResourceModel{}

	// Populate common fields from base resource data
	state.From(group.Base)

	// Set specific attributes
	gro.flattenMemberLinks(group.MemberLinks, &state)
	state.MemberQuery = gro.FlattenQuery(group.MemberQuery)
	state.IdentityMatcher = gro.flattenIdentityMatcher(group.IdentityMatcher)
	state.Origin = types.StringPointerValue(group.Origin)

	// Return completed state model
	return state
}

// InvokeCreate invokes the Create API to create a new resource.
func (gro *GroupResourceOperator) InvokeCreate(req client.Group) (*client.Group, int, error) {
	return gro.Client.CreateGroup(req)
}

// InvokeRead invokes the Get API to retrieve an existing resource by name.
func (gro *GroupResourceOperator) InvokeRead(name string) (*client.Group, int, error) {
	return gro.Client.GetGroup(name)
}

// InvokeUpdate invokes the Update API to update an existing resource.
func (gro *GroupResourceOperator) InvokeUpdate(req client.Group) (*client.Group, int, error) {
	return gro.Client.UpdateGroup(req)
}

// InvokeDelete invokes the Delete API to delete a resource by name.
func (gro *GroupResourceOperator) InvokeDelete(name string) error {
	return gro.Client.DeleteGroup(name)
}

// Builders //

// buildMemberLinks builds a list of member link strings combining user IDs/emails and service accounts.
func (gro *GroupResourceOperator) buildMemberLinks(state GroupResourceModel) *[]string {
	// Initialize an empty slice to hold member link strings
	memberLinks := []string{}

	// Retrieve user ID and email links from state
	userIdsAndEmails := BuildSetString(gro.Ctx, gro.Diags, state.UserIdsAndEmails)

	// Retrieve service account links from state
	serviceAccounts := BuildSetString(gro.Ctx, gro.Diags, state.ServiceAccounts)

	// If user ID and email links are present, append them to the result slice
	if userIdsAndEmails != nil {
		// Iterate over the users/emails, construct their self link and append it to the memberLinks slice
		for _, u := range *userIdsAndEmails {
			memberLinks = append(memberLinks, fmt.Sprintf("/org/%s/user/%s", gro.Client.Org, u))
		}
	}

	// If service account links are present, append them to the result slice
	if serviceAccounts != nil {
		// Iterate over the service accounts, construct their self link and append it to the memberLinks slice
		for _, s := range *serviceAccounts {
			memberLinks = append(memberLinks, fmt.Sprintf("/org/%s/serviceaccount/%s", gro.Client.Org, s))
		}
	}

	// Return a pointer to the slice of member links
	return &memberLinks
}

// buildIdentityMatcher constructs a IdentityMatcher struct from the given Terraform state.
func (gro *GroupResourceOperator) buildIdentityMatcher(state types.List) *client.GroupIdentityMatcher {
	// Convert Terraform list into model blocks using generic helper
	blocks, ok := BuildList[models.IdentityMatcherModel](gro.Ctx, gro.Diags, state)

	// Return nil if conversion failed or list was empty
	if !ok {
		return nil
	}

	// Extract the very first block from the blocks slice
	block := blocks[0]

	// Construct and return the result
	return &client.GroupIdentityMatcher{
		Expression: BuildString(block.Expression),
		Language:   BuildString(block.Language),
	}
}

// Flatteners //

// FlattenMemberLinks splits a flat list of member link URLs into separate Terraform sets for users and service accounts.
func (gro *GroupResourceOperator) flattenMemberLinks(input *[]string, state *GroupResourceModel) {
	// Handle missing attribute by setting both sets to null
	if input == nil {
		(*state).UserIdsAndEmails = types.SetNull(types.StringType)
		(*state).ServiceAccounts = types.SetNull(types.StringType)
		return
	}

	// Base URL prefix for all group members
	linksPrefix := fmt.Sprintf("/org/%s", gro.Client.Org)

	// Prefix identifying user links
	userIdPrefix := fmt.Sprintf("%s/user/", linksPrefix)

	// Prefix identifying service account links
	serviceAccountsPrefix := fmt.Sprintf("%s/serviceaccount/", linksPrefix)

	// Slice to collect user link URLs
	userIdsAndEmails := []string{}

	// Slice to collect service account link URLs
	serviceAccounts := []string{}

	// Iterate over each provided link
	for _, item := range *input {
		// Classify as user link if prefix matches
		if strings.HasPrefix(item, userIdPrefix) {
			userIdsAndEmails = append(userIdsAndEmails, strings.TrimPrefix(item, userIdPrefix))
		} else if strings.HasPrefix(item, serviceAccountsPrefix) {
			// Classify as service account link if prefix matches
			serviceAccounts = append(serviceAccounts, strings.TrimPrefix(item, serviceAccountsPrefix))
		}
	}

	// Convert found user links to Terraform set or set null if none
	if len(userIdsAndEmails) != 0 {
		(*state).UserIdsAndEmails = FlattenSetString(&userIdsAndEmails)
	} else {
		(*state).UserIdsAndEmails = types.SetNull(types.StringType)
	}

	// Convert found service account links to Terraform set or set null if none
	if len(serviceAccounts) != 0 {
		(*state).ServiceAccounts = FlattenSetString(&serviceAccounts)
	} else {
		(*state).ServiceAccounts = types.SetNull(types.StringType)
	}
}

// flattenIdentityMatcher transforms client.IdentityMatcher into a Terraform types.List.
func (gro *GroupResourceOperator) flattenIdentityMatcher(input *client.GroupIdentityMatcher) types.List {
	// Get attribute types
	elementType := models.IdentityMatcherModel{}.AttributeTypes()

	// Check if the input is nil
	if input == nil {
		// Return a null list
		return types.ListNull(elementType)
	}

	// Build a single block
	block := models.IdentityMatcherModel{
		Expression: types.StringPointerValue(input.Expression),
		Language:   types.StringPointerValue(input.Language),
	}

	// Return the successfully created types.List
	return FlattenList(gro.Ctx, gro.Diags, []models.IdentityMatcherModel{block})
}
