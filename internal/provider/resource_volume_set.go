package cpln

import (
	"context"
	"encoding/json"
	"fmt"
	"regexp"
	"strings"

	client "github.com/controlplane-com/terraform-provider-cpln/internal/provider/client"
	models "github.com/controlplane-com/terraform-provider-cpln/internal/provider/models/volume_set"
	"github.com/controlplane-com/terraform-provider-cpln/internal/provider/validators"
	"github.com/hashicorp/terraform-plugin-framework-validators/float64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/int32validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure resource implements required interfaces.
var (
	_ resource.Resource                   = &VolumeSetResource{}
	_ resource.ResourceWithImportState    = &VolumeSetResource{}
	_ resource.ResourceWithValidateConfig = &VolumeSetResource{}
)

/*** Resource Model ***/

// VolumeSetResourceModel holds the Terraform state for the resource.
type VolumeSetResourceModel struct {
	EntityBaseModel
	Gvc                types.String               `tfsdk:"gvc"`
	Status             types.List                 `tfsdk:"status"`
	InitialCapacity    types.Int32                `tfsdk:"initial_capacity"`
	PerformanceClass   types.String               `tfsdk:"performance_class"`
	StorageClassSuffix types.String               `tfsdk:"storage_class_suffix"`
	FileSystemType     types.String               `tfsdk:"file_system_type"`
	CustomEncryption   types.List                 `tfsdk:"custom_encryption"`
	Snapshots          []models.SnapshotsModel    `tfsdk:"snapshots"`
	Autoscaling        []models.AutoscalingModel  `tfsdk:"autoscaling"`
	MountOptions       []models.MountOptionsModel `tfsdk:"mount_options"`
	VolumesetLink      types.String               `tfsdk:"volumeset_link"`
}

/*** Resource Configuration ***/

// VolumeSetResource is the resource implementation.
type VolumeSetResource struct {
	EntityBase
	Operations EntityOperations[VolumeSetResourceModel, client.VolumeSet]
}

// NewVolumeSetResource returns a new instance of the resource implementation.
func NewVolumeSetResource() resource.Resource {
	return &VolumeSetResource{}
}

// Configure configures the resource before use.
func (vsr *VolumeSetResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	vsr.EntityBaseConfigure(ctx, req.ProviderData, &resp.Diagnostics)
	vsr.Operations = NewEntityOperations(vsr.client, &VolumeSetResourceOperator{})
}

// ImportState sets up the import operation to map the imported ID to the "id" attribute in the state.
func (vsr *VolumeSetResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Split the import ID
	parts := strings.SplitN(req.ID, ":", 2)

	// Validate that ID has exactly three non-empty segments
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		// Report error when import identifier format is unexpected
		resp.Diagnostics.AddError(
			"Unexpected Import Identifier",
			fmt.Sprintf(
				"Expected import identifier with format: "+
					"'gvc:volume_set_name'. Got: %q", req.ID,
			),
		)

		// Abort import operation on error
		return
	}

	// Extract gvc and volumeSetName from parts
	gvc, volumeSetName := parts[0], parts[1]

	// Set the generated ID attribute in the Terraform state
	resp.Diagnostics.Append(
		resp.State.SetAttribute(ctx, path.Root("id"), types.StringValue(volumeSetName))...,
	)

	// Set the GVC attribute in the Terraform state
	resp.Diagnostics.Append(
		resp.State.SetAttribute(ctx, path.Root("gvc"), types.StringValue(gvc))...,
	)
}

// Metadata provides the resource type name.
func (vsr *VolumeSetResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = "cpln_volume_set"
}

// Schema defines the schema for the resource.
func (vsr *VolumeSetResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: MergeAttributes(vsr.EntityBaseAttributes("volume set"), map[string]schema.Attribute{
			"gvc": schema.StringAttribute{
				Description: "Name of the associated GVC.",
				Required:    true,
				Validators: []validator.String{
					validators.NameValidator{},
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"status": schema.ListNestedAttribute{
				Description: "Status of the Volume Set.",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"parent_id": schema.StringAttribute{
							Description: "The GVC ID.",
							Computed:    true,
						},
						"used_by_workload": schema.StringAttribute{
							Description: "The url of the workload currently using this volume set (if any).",
							Computed:    true,
						},
						"workload_links": schema.SetAttribute{
							Description: "Contains a list of workload links that are using this volume set.",
							ElementType: types.StringType,
							Computed:    true,
						},
						"binding_id": schema.StringAttribute{
							Description: "Uniquely identifies the connection between the volume set and its workload. Every time a new connection is made, a new id is generated (e.g., If a workload is updated to remove the volume set, then updated again to reattach it, the volume set will have a new binding id).",
							Computed:    true,
						},
						"locations": schema.SetAttribute{
							Description: "Contains a list of actual volumes grouped by location.",
							ElementType: types.StringType,
							Computed:    true,
						},
					},
				},
			},
			"initial_capacity": schema.Int32Attribute{
				Description: "The initial volume size in this set, specified in GB. The minimum size for the performance class `general-purpose-ssd` is `10 GB`, while `high-throughput-ssd` requires at least `200 GB`.",
				Required:    true,
			},
			"performance_class": schema.StringAttribute{
				Description: "Each volume set has a single, immutable, performance class. Valid classes: `general-purpose-ssd` or `high-throughput-ssd`.",
				Required:    true,
				Validators: []validator.String{
					stringvalidator.OneOf("general-purpose-ssd", "high-throughput-ssd", "shared"),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"storage_class_suffix": schema.StringAttribute{
				Description: "For self-hosted locations only. The storage class used for volumes in this set will be {performanceClass}-{fileSystemType}-{storageClassSuffix} if it exists, otherwise it will be {performanceClass}-{fileSystemType}",
				Optional:    true,
				Validators: []validator.String{
					stringvalidator.RegexMatches(
						regexp.MustCompile(`^[a-zA-Z][0-9a-zA-Z\-_]*$`),
						"must be a valid storage class suffix",
					),
				},
			},
			"file_system_type": schema.StringAttribute{
				Description: "Each volume set has a single, immutable file system. Valid types: `xfs` or `ext4`.",
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString("ext4"),
				Validators: []validator.String{
					stringvalidator.OneOf("xfs", "ext4", "shared"),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"volumeset_link": schema.StringAttribute{
				Description: "Output used when linking a volume set to a workload.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
		}),
		Blocks: map[string]schema.Block{
			"custom_encryption": schema.ListNestedBlock{
				Description: "Configuration for customer-managed encryption keys, keyed by region.",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"regions": schema.MapAttribute{
							Description: "Map of region identifiers to encryption key configuration.",
							ElementType: models.CustomEncryptionRegionModel{}.AttributeTypes(),
							Required:    true,
						},
					},
				},
				Validators: []validator.List{
					listvalidator.SizeAtMost(1),
				},
			},
			"snapshots": schema.ListNestedBlock{
				Description: "Point-in-time copies of data stored within the volume set, capturing the state of the data at a specific moment.",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"create_final_snapshot": schema.BoolAttribute{
							Description: "If true, a volume snapshot will be created immediately before deletion of any volume in this set. Default: `true`",
							Optional:    true,
							Computed:    true,
							Default:     booldefault.StaticBool(true),
						},
						"retention_duration": schema.StringAttribute{
							Description: "The default retention period for volume snapshots. This string should contain a floating point number followed by either d, h, or m. For example, \"10d\" would retain snapshots for 10 days.",
							Optional:    true,
							Validators: []validator.String{
								stringvalidator.RegexMatches(
									regexp.MustCompile(`^([0-9]+(\.[0-9]+)?[dhm])$`),
									"must be a valid retention duration",
								),
							},
						},
						"schedule": schema.StringAttribute{
							Description: "A standard cron schedule expression used to determine when a snapshot will be taken. (i.e., `0 * * * *` Every hour). Note: snapshots cannot be scheduled more often than once per hour.",
							Optional:    true,
						},
					},
				},
				Validators: []validator.List{
					listvalidator.SizeAtMost(1),
				},
			},
			"autoscaling": schema.ListNestedBlock{
				Description: "Automated adjustment of the volume set's capacity based on predefined metrics or conditions.",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"max_capacity": schema.Int32Attribute{
							Description: "The maximum size in GB for a volume in this set. A volume cannot grow to be bigger than this value. Minimum value: `10`.",
							Optional:    true,
							Validators: []validator.Int32{
								int32validator.AtLeast(10),
							},
						},
						"min_free_percentage": schema.Int32Attribute{
							Description: "The guaranteed free space on the volume as a percentage of the volume's total size. Control Plane will try to maintain at least that many percent free by scaling up the total size. Minimum percentage: `1`. Maximum Percentage: `100`.",
							Optional:    true,
							Validators: []validator.Int32{
								int32validator.Between(1, 100),
							},
						},
						"scaling_factor": schema.Float64Attribute{
							Description: "When scaling is necessary, then `new_capacity = current_capacity * storageScalingFactor`. Minimum value: `1.1`.",
							Optional:    true,
							Validators: []validator.Float64{
								float64validator.AtLeast(1.1),
							},
						},
					},
				},
				Validators: []validator.List{
					listvalidator.SizeAtMost(1),
				},
			},
			"mount_options": schema.ListNestedBlock{
				Description: "A list of mount options to use when mounting volumes in this set.",
				NestedObject: schema.NestedBlockObject{
					Blocks: map[string]schema.Block{
						"resources": schema.ListNestedBlock{
							Description: "For volume sets using the shared file system, this object specifies the CPU and memory resources allotted to each mount point.",
							NestedObject: schema.NestedBlockObject{
								Attributes: map[string]schema.Attribute{
									"max_cpu": schema.StringAttribute{
										Description: "",
										Optional:    true,
										Computed:    true,
										Default:     stringdefault.StaticString("2000m"),
									},
									"min_cpu": schema.StringAttribute{
										Description: "",
										Optional:    true,
										Computed:    true,
										Default:     stringdefault.StaticString("500m"),
									},
									"min_memory": schema.StringAttribute{
										Description: "",
										Optional:    true,
										Computed:    true,
										Default:     stringdefault.StaticString("1Gi"),
									},
									"max_memory": schema.StringAttribute{
										Description: "",
										Optional:    true,
										Computed:    true,
										Default:     stringdefault.StaticString("2Gi"),
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

// ValidateConfig validates the configuration of the resource.
func (vsr *VolumeSetResource) ValidateConfig(ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse) {
	// Declare variable to store desired resource plan
	var plan VolumeSetResourceModel

	// Populate plan variable from config and capture diagnostics
	resp.Diagnostics.Append(req.Config.Get(ctx, &plan)...)

	// Halt further processing if plan retrieval failed
	if resp.Diagnostics.HasError() {
		return
	}

	// Extract the desired initial capacity and performance class from the plan
	initialCapacity := BuildInt(plan.InitialCapacity)
	performanceClass := BuildString(plan.PerformanceClass)

	// Skip validation if one of the following is nil
	if initialCapacity == nil || performanceClass == nil {
		return
	}

	// Initialize variable to hold minimum capacity threshold
	var minCapacity int

	// Assign minimum capacity based on performance class
	switch *performanceClass {
	case "general-purpose-ssd":
		// General-Purpose SSD requires 10 GB minimum
		minCapacity = 10
	case "high-throughput-ssd":
		// High-Throughput SSD requires 200 GB minimum
		minCapacity = 200
	default:
		// No validation needed for other performance classes
		return
	}

	// Add an attribute error if the planned capacity is below the minimum threshold
	if *initialCapacity < minCapacity {
		resp.Diagnostics.AddAttributeError(
			path.Root("initial_capacity"),
			"initial_capacity too small",
			fmt.Sprintf(
				"For performance_class %q, initial_capacity must be at least %d GB (you provided %d GB)",
				*performanceClass, minCapacity, *initialCapacity,
			),
		)
	}
}

// Create creates the resource.
func (vsr *VolumeSetResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	CreateGeneric(ctx, req, resp, vsr.Operations)
}

// Read fetches the current state of the resource.
func (vsr *VolumeSetResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	ReadGeneric(ctx, req, resp, vsr.Operations)
}

// Update modifies the resource.
func (vsr *VolumeSetResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	UpdateGeneric(ctx, req, resp, vsr.Operations)
}

// Delete removes the resource.
func (vsr *VolumeSetResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	DeleteGeneric(ctx, req, resp, vsr.Operations)
}

/*** Resource Operator ***/

// VolumeSetResourceOperator is the operator for managing the state.
type VolumeSetResourceOperator struct {
	EntityOperator[VolumeSetResourceModel]
}

// NewAPIRequest creates a request payload from a state model.
func (vsro *VolumeSetResourceOperator) NewAPIRequest(isUpdate bool) client.VolumeSet {
	// Initialize a new request payload
	requestPayload := client.VolumeSet{}

	// Initialize the GVC spec struct
	var spec *client.VolumeSetSpec = &client.VolumeSetSpec{}

	// Populate Base fields from state
	vsro.Plan.Fill(&requestPayload.Base, isUpdate)

	// Map planned state attributes to the API struct
	if isUpdate {
		requestPayload.SpecReplace = spec
	} else {
		requestPayload.Spec = spec
	}

	// Set specific attributes
	spec.InitialCapacity = BuildInt(vsro.Plan.InitialCapacity)
	spec.PerformanceClass = BuildString(vsro.Plan.PerformanceClass)
	spec.StorageClassSuffix = BuildString(vsro.Plan.StorageClassSuffix)
	spec.FileSystemType = BuildString(vsro.Plan.FileSystemType)
	spec.CustomEncryption = vsro.buildCustomEncryption(vsro.Plan.CustomEncryption)
	spec.Snapshots = vsro.buildSnapshots(vsro.Plan.Snapshots)
	spec.AutoScaling = vsro.buildAutoscaling(vsro.Plan.Autoscaling)
	spec.MountOptions = vsro.buildMountOptions(vsro.Plan.MountOptions)

	// Return constructed request payload
	return requestPayload
}

// MapResponseToState constructs the Terraform state model from the API response payload.
func (vsro *VolumeSetResourceOperator) MapResponseToState(apiResp *client.VolumeSet, isCreate bool) VolumeSetResourceModel {
	// Initialize empty state model
	state := VolumeSetResourceModel{}

	// Populate common fields from base resource data
	state.From(apiResp.Base)

	// Set specific attributes
	state.Gvc = types.StringPointerValue(BuildString(vsro.Plan.Gvc))
	state.Status = vsro.flattenStatus(apiResp.Status)
	state.VolumesetLink = types.StringValue(fmt.Sprintf("cpln://volumeset/%s", *apiResp.Name))

	// Just in case the spec is nil
	if apiResp.Spec == nil {
		state.InitialCapacity = types.Int32Null()
		state.PerformanceClass = types.StringNull()
		state.StorageClassSuffix = types.StringNull()
		state.FileSystemType = types.StringNull()
		state.CustomEncryption = types.ListNull(models.CustomEncryptionModel{}.AttributeTypes())
		state.Snapshots = nil
		state.Autoscaling = nil
		state.MountOptions = nil
	} else {
		state.InitialCapacity = FlattenInt(apiResp.Spec.InitialCapacity)
		state.PerformanceClass = types.StringPointerValue(apiResp.Spec.PerformanceClass)
		state.StorageClassSuffix = types.StringPointerValue(apiResp.Spec.StorageClassSuffix)
		state.FileSystemType = types.StringPointerValue(apiResp.Spec.FileSystemType)
		state.CustomEncryption = vsro.flattenCustomEncryption(apiResp.Spec.CustomEncryption)
		state.Snapshots = vsro.flattenSnapshots(apiResp.Spec.Snapshots)
		state.Autoscaling = vsro.flattenAutoscaling(apiResp.Spec.AutoScaling)
		state.MountOptions = vsro.flattenMountOptions(apiResp.Spec.MountOptions)
	}

	// Return completed state model
	return state
}

// InvokeCreate invokes the Create API to create a new resource.
func (vsro *VolumeSetResourceOperator) InvokeCreate(req client.VolumeSet) (*client.VolumeSet, int, error) {
	return vsro.Client.CreateVolumeSet(req, vsro.Plan.Gvc.ValueString())
}

// InvokeRead invokes the Get API to retrieve an existing resource by name.
func (vsro *VolumeSetResourceOperator) InvokeRead(name string) (*client.VolumeSet, int, error) {
	return vsro.Client.GetVolumeSet(name, vsro.Plan.Gvc.ValueString())
}

// InvokeUpdate invokes the Update API to update an existing resource.
func (vsro *VolumeSetResourceOperator) InvokeUpdate(req client.VolumeSet) (*client.VolumeSet, int, error) {
	return vsro.Client.UpdateVolumeSet(req, vsro.Plan.Gvc.ValueString())
}

// InvokeDelete invokes the Delete API to delete a resource by name.
func (vsro *VolumeSetResourceOperator) InvokeDelete(name string) error {
	return vsro.Client.DeleteVolumeSet(name, vsro.Plan.Gvc.ValueString())
}

// Builders //

// buildCustomEncryption constructs a VolumeSetCustomEncryption from the given Terraform state list.
func (vsro *VolumeSetResourceOperator) buildCustomEncryption(state types.List) *client.VolumeSetCustomEncryption {
	// Convert Terraform list into model blocks using generic helper
	blocks, ok := BuildList[models.CustomEncryptionModel](vsro.Ctx, vsro.Diags, state)

	// Return nil if conversion failed or list was empty
	if !ok {
		return nil
	}

	// Take the first (and only) block
	block := blocks[0]

	// Construct and return the output
	return &client.VolumeSetCustomEncryption{
		Regions: vsro.buildCustomEncryptionRegions(block.Regions),
	}
}

// buildCustomEncryptionRegions constructs a map of VolumeSetCustomEncryptionRegion from the given Terraform map.
func (vsro *VolumeSetResourceOperator) buildCustomEncryptionRegions(state types.Map) *map[string]*client.VolumeSetCustomEncryptionRegion {
	// Return nil if state is null or unknown
	if state.IsNull() || state.IsUnknown() {
		return nil
	}

	// Convert Terraform map into Go map
	var regions map[string]models.CustomEncryptionRegionModel
	vsro.Diags.Append(state.ElementsAs(vsro.Ctx, &regions, false)...)

	// Return nil if conversion failed or regions were nil
	if vsro.Diags.HasError() || regions == nil {
		return nil
	}

	// Construct output map
	output := make(map[string]*client.VolumeSetCustomEncryptionRegion, len(regions))

	// Iterate over regions and populate output map
	for key, value := range regions {
		keyID := BuildString(value.KeyId)
		output[key] = &client.VolumeSetCustomEncryptionRegion{
			KeyId: keyID,
		}
	}

	// Return constructed map
	return &output
}

// buildSnapshots constructs a VolumeSetSnapshots from the given Terraform state.
func (vsro *VolumeSetResourceOperator) buildSnapshots(state []models.SnapshotsModel) *client.VolumeSetSnapshots {
	// Return nil if state is not specified
	if len(state) == 0 {
		return nil
	}

	// Take the first (and only) block
	block := state[0]

	// Construct and return the output
	return &client.VolumeSetSnapshots{
		CreateFinalSnapshot: BuildBool(block.CreateFinalSnapshot),
		RetentionDuration:   BuildString(block.RetentionDuration),
		Schedule:            BuildString(block.Schedule),
	}
}

// buildAutoscaling constructs a VolumeSetScaling from the given Terraform state.
func (vsro *VolumeSetResourceOperator) buildAutoscaling(state []models.AutoscalingModel) *client.VolumeSetScaling {
	// Return nil if state is not specified
	if len(state) == 0 {
		return nil
	}

	// Take the first (and only) block
	block := state[0]

	// Construct and return the output
	return &client.VolumeSetScaling{
		MaxCapacity:       BuildInt(block.MaxCapacity),
		MinFreePercentage: BuildInt(block.MinFreePercentage),
		ScalingFactor:     BuildFloat64(block.ScalingFactor),
	}
}

// buildMountOptions constructs a VolumeSetMountOptions from the given Terraform state.
func (vsro *VolumeSetResourceOperator) buildMountOptions(state []models.MountOptionsModel) *client.VolumeSetMountOptions {
	// Return nil if state is not specified
	if len(state) == 0 {
		return nil
	}

	// Take the first (and only) block
	block := state[0]

	// Construct and return the output
	return &client.VolumeSetMountOptions{
		Resources: vsro.buildMountOptionsResources(block.Resources),
	}
}

// buildMountOptionsResources constructs a VolumeSetMountOptionsResources from the given Terraform state.
func (vsro *VolumeSetResourceOperator) buildMountOptionsResources(state []models.MountOptionsResourcesModel) *client.VolumeSetMountOptionsResources {
	// Return nil if state is not specified
	if len(state) == 0 {
		return nil
	}

	// Take the first (and only) block
	block := state[0]

	// Construct and return the output
	return &client.VolumeSetMountOptionsResources{
		MaxCpu:    BuildString(block.MaxCpu),
		MinCpu:    BuildString(block.MinCpu),
		MinMemory: BuildString(block.MinMemory),
		MaxMemory: BuildString(block.MaxMemory),
	}
}

// Flatteners //

// flattenCustomEncryption transforms *client.VolumeSetCustomEncryption into a Terraform types.List.
func (vsro *VolumeSetResourceOperator) flattenCustomEncryption(input *client.VolumeSetCustomEncryption) types.List {
	// Get attribute types
	elementType := models.CustomEncryptionModel{}.AttributeTypes()

	// Check if the input is nil
	if input == nil {
		// Return a null list
		return types.ListNull(elementType)
	}

	// Build a single block
	block := models.CustomEncryptionModel{
		Regions: vsro.flattenCustomEncryptionRegions(input.Regions),
	}

	// Return the successfully created types.List
	return FlattenList(vsro.Ctx, vsro.Diags, []models.CustomEncryptionModel{block})
}

// flattenCustomEncryptionRegions transforms *map[string]*client.VolumeSetCustomEncryptionRegion into a Terraform types.Map.
func (vsro *VolumeSetResourceOperator) flattenCustomEncryptionRegions(input *map[string]*client.VolumeSetCustomEncryptionRegion) types.Map {
	// Get attribute types
	elementType := models.CustomEncryptionRegionModel{}.AttributeTypes()

	// Check if the input is nil
	if input == nil {
		// Return a null list
		return types.MapNull(elementType)
	}

	// Build regions map
	regions := make(map[string]models.CustomEncryptionRegionModel, len(*input))

	// Iterate over each region in the input map
	for key, value := range *input {
		// Initialize a new region model
		region := models.CustomEncryptionRegionModel{}

		// Set KeyId attribute
		if value != nil {
			region.KeyId = types.StringPointerValue(value.KeyId)
		} else {
			region.KeyId = types.StringNull()
		}

		// Add region to regions map
		regions[key] = region
	}

	// Convert the regions map into a Terraform types.Map
	result, diags := types.MapValueFrom(vsro.Ctx, elementType, regions)
	vsro.Diags.Append(diags...)

	// Check for errors during conversion and return null map if any
	if vsro.Diags.HasError() {
		return types.MapNull(elementType)
	}

	// Return the successfully created types.Map
	return result
}

// flattenStatus transforms *client.VolumeSetStatus into a Terraform types.List.
func (vsro *VolumeSetResourceOperator) flattenStatus(input *client.VolumeSetStatus) types.List {
	// Get attribute types
	elementType := models.StatusModel{}.AttributeTypes()

	// Check if the input is nil
	if input == nil {
		// Return a null list
		return types.ListNull(elementType)
	}

	// Build a single block
	block := models.StatusModel{
		ParentId:       types.StringPointerValue(input.ParentID),
		UsedByWorkload: types.StringPointerValue(input.UsedByWorkload),
		WorkloadLinks:  FlattenSetString(input.WorkloadLinks),
		BindingId:      types.StringPointerValue(input.BindingID),
		Locations:      vsro.flattenStatusLocations(input.Locations),
	}

	// Return the successfully created types.List
	return FlattenList(vsro.Ctx, vsro.Diags, []models.StatusModel{block})
}

// flattenStatusLocations flattens the provided list of location objects into a Terraform set of JSON-encoded strings.
func (vsro *VolumeSetResourceOperator) flattenStatusLocations(locations *[]interface{}) types.Set {
	// If locations pointer is nil return a null Terraform set
	if locations == nil {
		return types.SetNull(types.StringType)
	}

	// Create a slice to hold JSON strings with length equal to number of locations
	result := make([]string, len(*locations))

	// Iterate over each location interface value
	for i, location := range *locations {
		// Marshal the location into JSON bytes
		jsonData, err := json.Marshal(location)

		// If marshaling fails record an error message string
		if err != nil {
			result[i] = fmt.Sprintf("Error serializing to JSON: %s", err)
		} else {
			// If marshaling succeeds convert bytes to string
			result[i] = string(jsonData)
		}
	}

	// Convert the slice of strings into a Terraform set of strings
	return FlattenSetString(&result)
}

// flattenSnapshots transforms *client.VolumeSetSnapshots into a []models.SnapshotsModel.
func (vsro *VolumeSetResourceOperator) flattenSnapshots(input *client.VolumeSetSnapshots) []models.SnapshotsModel {
	// Check if the input is nil
	if input == nil {
		return nil
	}

	// Build a single block
	block := models.SnapshotsModel{
		CreateFinalSnapshot: types.BoolPointerValue(input.CreateFinalSnapshot),
		RetentionDuration:   types.StringPointerValue(input.RetentionDuration),
		Schedule:            types.StringPointerValue(input.Schedule),
	}

	// Return a slice containing the single block
	return []models.SnapshotsModel{block}
}

// flattenAutoscaling transforms *client.VolumeSetScaling into a []models.AutoscalingModel.
func (vsro *VolumeSetResourceOperator) flattenAutoscaling(input *client.VolumeSetScaling) []models.AutoscalingModel {
	// Check if the input is nil
	if input == nil {
		return nil
	}

	// Build a single block
	block := models.AutoscalingModel{
		MaxCapacity:       FlattenInt(input.MaxCapacity),
		MinFreePercentage: FlattenInt(input.MinFreePercentage),
		ScalingFactor:     FlattenFloat64(input.ScalingFactor),
	}

	// Return a slice containing the single block
	return []models.AutoscalingModel{block}
}

// flattenMountOptions transforms *client.VolumeSetMountOptions into a []models.MountOptionsResourcesModel.
func (vsro *VolumeSetResourceOperator) flattenMountOptions(input *client.VolumeSetMountOptions) []models.MountOptionsModel {
	// Check if the input is nil
	if input == nil {
		return nil
	}

	// Build a single block
	block := models.MountOptionsModel{
		Resources: vsro.flattenMountOptionsResources(input.Resources),
	}

	// Return a slice containing the single block
	return []models.MountOptionsModel{block}
}

// flattenMountOptionsResources transforms *client.VolumeSetMountOptionsResources into a []models.MountOptionsResourcesModel.
func (vsro *VolumeSetResourceOperator) flattenMountOptionsResources(input *client.VolumeSetMountOptionsResources) []models.MountOptionsResourcesModel {
	// Check if the input is nil
	if input == nil {
		return nil
	}

	// Build a single block
	block := models.MountOptionsResourcesModel{
		MaxCpu:    types.StringPointerValue(input.MaxCpu),
		MinCpu:    types.StringPointerValue(input.MinCpu),
		MinMemory: types.StringPointerValue(input.MinMemory),
		MaxMemory: types.StringPointerValue(input.MaxMemory),
	}

	// Return a slice containing the single block
	return []models.MountOptionsResourcesModel{block}
}
