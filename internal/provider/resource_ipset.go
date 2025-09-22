package cpln

import (
	"context"

	client "github.com/controlplane-com/terraform-provider-cpln/internal/provider/client"
	models "github.com/controlplane-com/terraform-provider-cpln/internal/provider/models/ipset"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure resource implements required interfaces.
var (
	_ resource.Resource                = &IpSetResource{}
	_ resource.ResourceWithImportState = &IpSetResource{}
)

/*** Resource Model ***/

// IpSetResourceModel holds the Terraform state for the resource.
type IpSetResourceModel struct {
	EntityBaseModel
	Link      types.String `tfsdk:"link"`
	Locations types.Set    `tfsdk:"location"`
	Status    types.List   `tfsdk:"status"`
}

/*** Resource Configuration ***/

// IpSetResource is the resource implementation.
type IpSetResource struct {
	EntityBase
	Operations EntityOperations[IpSetResourceModel, client.IpSet]
}

// NewIpSetResource returns a new instance of the resource implementation.
func NewIpSetResource() resource.Resource {
	return &IpSetResource{}
}

// Configure configures the resource before use.
func (isr *IpSetResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	isr.EntityBaseConfigure(ctx, req.ProviderData, &resp.Diagnostics)
	isr.Operations = NewEntityOperations(isr.client, &IpSetResourceOperator{})
}

// ImportState sets up the import operation to map the imported ID to the "id" attribute in the state.
func (isr *IpSetResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// Metadata provides the resource type name.
func (isr *IpSetResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = "cpln_ipset"
}

// Schema defines the schema for the resource.
func (isr *IpSetResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: MergeAttributes(isr.EntityBaseAttributes("IP Set"), map[string]schema.Attribute{
			"link": schema.StringAttribute{
				Description: "The self link of a workload or a GVC.",
				Optional:    true,
			},
			"status": schema.ListNestedAttribute{
				Description: "The status of the IP Set.",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"ip_address": schema.ListNestedAttribute{
							Description: "",
							Computed:    true,
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"name": schema.StringAttribute{
										Description: "",
										Computed:    true,
									},
									"ip": schema.StringAttribute{
										Description: "",
										Computed:    true,
									},
									"id": schema.StringAttribute{
										Description: "",
										Computed:    true,
									},
									"state": schema.StringAttribute{
										Description: "",
										Computed:    true,
									},
									"created": schema.StringAttribute{
										Description: "",
										Computed:    true,
									},
								},
							},
						},
						"error": schema.StringAttribute{
							Description: "",
							Computed:    true,
						},
						"warning": schema.StringAttribute{
							Description: "",
							Computed:    true,
						},
					},
				},
			},
		}),
		Blocks: map[string]schema.Block{
			"location": schema.SetNestedBlock{
				Description: "",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							Description: "The self link of a location.",
							Required:    true,
						},
						"retention_policy": schema.StringAttribute{
							Description: "",
							Required:    true,
							Validators: []validator.String{
								stringvalidator.OneOf("keep", "free"),
							},
						},
					},
				},
			},
		},
	}
}

// Create creates the resource.
func (isr *IpSetResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	CreateGeneric(ctx, req, resp, isr.Operations)
}

// Read fetches the current state of the resource.
func (isr *IpSetResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	ReadGeneric(ctx, req, resp, isr.Operations)
}

// Update modifies the resource.
func (isr *IpSetResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	UpdateGeneric(ctx, req, resp, isr.Operations)
}

// Delete removes the resource.
func (isr *IpSetResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	DeleteGeneric(ctx, req, resp, isr.Operations)
}

/*** Resource Operator ***/

// IpSetResourceOperator is the operator for managing the state.
type IpSetResourceOperator struct {
	EntityOperator[IpSetResourceModel]
}

// NewAPIRequest creates a request payload from a state model.
func (isro *IpSetResourceOperator) NewAPIRequest(isUpdate bool) client.IpSet {
	// Initialize a new request payload
	requestPayload := client.IpSet{
		Spec: &client.IpSetSpec{},
	}

	// Initialize the spec struct
	var spec *client.IpSetSpec = &client.IpSetSpec{}

	// Populate Base fields from state
	isro.Plan.Fill(&requestPayload.Base, isUpdate)

	// Map planned state attributes to the API struct
	if isUpdate {
		requestPayload.SpecReplace = spec
	} else {
		requestPayload.Spec = spec
	}

	// Set specific attributes
	spec.Link = BuildString(isro.Plan.Link)
	spec.Locations = isro.buildLocations(isro.Plan.Locations)

	// Return constructed request payload
	return requestPayload
}

// MapResponseToState constructs the Terraform state model from the API response payload.
func (isro *IpSetResourceOperator) MapResponseToState(apiResp *client.IpSet, isCreate bool) IpSetResourceModel {
	// Initialize empty state model
	state := IpSetResourceModel{}

	// Populate common fields from base resource data
	state.From(apiResp.Base)

	// Just in case GVC spec is nil
	if apiResp.Spec != nil {
		state.Link = types.StringPointerValue(apiResp.Spec.Link)
		state.Locations = isro.flattenLocations(apiResp.Spec.Locations)
	} else {
		state.Link = types.StringNull()
		state.Locations = types.SetNull(models.LocationModel{}.AttributeTypes())
	}

	// Set specific attributes
	state.Status = isro.flattenStatus(apiResp.Status)

	// Return completed state model
	return state
}

// InvokeCreate invokes the Create API to create a new resource.
func (isro *IpSetResourceOperator) InvokeCreate(req client.IpSet) (*client.IpSet, int, error) {
	return isro.Client.CreateIpSet(req)
}

// InvokeRead invokes the Get API to retrieve an existing resource by name.
func (isro *IpSetResourceOperator) InvokeRead(name string) (*client.IpSet, int, error) {
	return isro.Client.GetIpSet(name)
}

// InvokeUpdate invokes the Update API to update an existing resource.
func (isro *IpSetResourceOperator) InvokeUpdate(req client.IpSet) (*client.IpSet, int, error) {
	return isro.Client.UpdateIpSet(req)
}

// InvokeDelete invokes the Delete API to delete a resource by name.
func (isro *IpSetResourceOperator) InvokeDelete(name string) error {
	return isro.Client.DeleteIpSet(name)
}

// Builders //

// buildLocations constructs a []client.IpSetLocation from the given Terraform state.
func (isro *IpSetResourceOperator) buildLocations(state types.Set) *[]client.IpSetLocation {
	// Convert Terraform set into model blocks using generic helper
	blocks, ok := BuildSet[models.LocationModel](isro.Ctx, isro.Diags, state)

	// Return nil if conversion failed or set was empty
	if !ok {
		return nil
	}

	// Prepare the output slice
	output := []client.IpSetLocation{}

	// Iterate over each block and construct an output item
	for _, block := range blocks {
		// Construct the item
		item := client.IpSetLocation{
			Name:            BuildString(block.Name),
			RetentionPolicy: BuildString(block.RetentionPolicy),
		}

		// Add the item to the output slice
		output = append(output, item)
	}

	// Return a pointer to the output
	return &output
}

// Flatteners //

// flattenLocations transforms *[]client.IpSetLocation into a Terraform types.Set.
func (isro *IpSetResourceOperator) flattenLocations(input *[]client.IpSetLocation) types.Set {
	// Get attribute types
	elementType := models.LocationModel{}.AttributeTypes()

	// Check if the input is nil
	if input == nil {
		// Return a null set
		return types.SetNull(elementType)
	}

	// Define the blocks slice
	var blocks []models.LocationModel

	// Iterate over the slice and construct the blocks
	for _, item := range *input {
		// Construct a block
		block := models.LocationModel{
			Name:            types.StringPointerValue(item.Name),
			RetentionPolicy: types.StringPointerValue(item.RetentionPolicy),
		}

		// Append the constructed block to the blocks slice
		blocks = append(blocks, block)
	}

	// Return the successfully created types.Set
	return FlattenSet(isro.Ctx, isro.Diags, blocks)
}

// flattenStatus transforms *client.IpSetStatus into a Terraform types.List.
func (isro *IpSetResourceOperator) flattenStatus(input *client.IpSetStatus) types.List {
	// Get attribute types
	elementType := models.StatusModel{}.AttributeTypes()

	// Check if the input is nil
	if input == nil {
		// Return a null list
		return types.ListNull(elementType)
	}

	// Build a single block
	block := models.StatusModel{
		IpAddresses: isro.flattenStatusIpAddresses(input.IpAddresses),
		Error:       types.StringPointerValue(input.Error),
		Warning:     types.StringPointerValue(input.Warning),
	}

	// Return the successfully created types.List
	return FlattenList(isro.Ctx, isro.Diags, []models.StatusModel{block})
}

// flattenStatusIpAddresses transforms *[]client.IpSetIpAddress into a Terraform types.List.
func (isro *IpSetResourceOperator) flattenStatusIpAddresses(input *[]client.IpSetIpAddress) types.List {
	// Get attribute types
	elementType := models.IpAddressModel{}.AttributeTypes()

	// Check if the input is nil
	if input == nil {
		// Return a null list
		return types.ListNull(elementType)
	}

	// Define the blocks slice
	var blocks []models.IpAddressModel

	// Iterate over the slice and construct the blocks
	for _, item := range *input {
		// Construct a block
		block := models.IpAddressModel{
			Name:    types.StringPointerValue(item.Name),
			IP:      types.StringPointerValue(item.IP),
			ID:      types.StringPointerValue(item.ID),
			State:   types.StringPointerValue(item.State),
			Created: types.StringPointerValue(item.Created),
		}

		// Append the constructed block to the blocks slice
		blocks = append(blocks, block)
	}

	// Return the successfully created types.List
	return FlattenList(isro.Ctx, isro.Diags, blocks)
}
