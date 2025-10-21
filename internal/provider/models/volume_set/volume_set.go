package volume_set

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Status //

type StatusModel struct {
	ParentId       types.String `tfsdk:"parent_id"`
	UsedByWorkload types.String `tfsdk:"used_by_workload"`
	WorkloadLinks  types.Set    `tfsdk:"workload_links"`
	BindingId      types.String `tfsdk:"binding_id"`
	Locations      types.Set    `tfsdk:"locations"`
}

func (s StatusModel) AttributeTypes() attr.Type {
	return types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"parent_id":        types.StringType,
			"used_by_workload": types.StringType,
			"workload_links":   types.SetType{ElemType: types.StringType},
			"binding_id":       types.StringType,
			"locations":        types.SetType{ElemType: types.StringType},
		},
	}
}

// Custom Encryption //

type CustomEncryptionModel struct {
	Regions types.Map `tfsdk:"regions"`
}

func (c CustomEncryptionModel) AttributeTypes() attr.Type {
	return types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"regions": types.MapType{
				ElemType: CustomEncryptionRegionModel{}.AttributeTypes(),
			},
		},
	}
}

// Custom Encryption -> Region //

type CustomEncryptionRegionModel struct {
	KeyId types.String `tfsdk:"key_id"`
}

func (c CustomEncryptionRegionModel) AttributeTypes() attr.Type {
	return types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"key_id": types.StringType,
		},
	}
}

// Snapshots //

type SnapshotsModel struct {
	CreateFinalSnapshot types.Bool   `tfsdk:"create_final_snapshot"`
	RetentionDuration   types.String `tfsdk:"retention_duration"`
	Schedule            types.String `tfsdk:"schedule"`
}

// Autoscaling //

type AutoscalingModel struct {
	MaxCapacity       types.Int32   `tfsdk:"max_capacity"`
	MinFreePercentage types.Int32   `tfsdk:"min_free_percentage"`
	ScalingFactor     types.Float64 `tfsdk:"scaling_factor"`
}

// Mount Options //

type MountOptionsModel struct {
	Resources []MountOptionsResourcesModel `tfsdk:"resources"`
}

// Mount Options -> Resources //

type MountOptionsResourcesModel struct {
	MaxCpu    types.String `tfsdk:"max_cpu"`
	MinCpu    types.String `tfsdk:"min_cpu"`
	MinMemory types.String `tfsdk:"min_memory"`
	MaxMemory types.String `tfsdk:"max_memory"`
}
